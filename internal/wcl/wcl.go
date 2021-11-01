package wcl

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/semaphore"
)

var (
	clientID     string
	clientSecret string
	httpClient   = &http.Client{
		Transport: &customRoundTripper{},
	}
	token         string
	tokenMutex    sync.RWMutex
	callSemaphore = semaphore.NewWeighted(8)
	client        = graphql.NewClient(wclURL, httpClient)
)

const (
	wclURL     = "https://www.warcraftlogs.com/api/v2/client"
	wclAuthURL = "https://www.warcraftlogs.com/oauth/token"
)

func init() {
	clientID = os.Getenv("WCL_CLIENT_ID")
	clientSecret = os.Getenv("WCL_CLIENT_SECRET")
}

func authenticate(ctx context.Context) error {
	tokenMutex.RLock()
	if token != "" {
		tokenMutex.RUnlock()
		return nil // already acquired
	}
	tokenMutex.RUnlock()

	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	var buffer bytes.Buffer
	form := multipart.NewWriter(&buffer)
	if w, err := form.CreateFormField("grant_type"); err != nil {
		return err
	} else {
		w.Write([]byte(`client_credentials`)) //nolint:errcheck
	}
	form.Close() //nolint:errcheck
	req, err := http.NewRequestWithContext(ctx, "POST", wclAuthURL, &buffer)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", form.FormDataContentType())
	req.SetBasicAuth(clientID, clientSecret)
	var t struct {
		Token string `json:"access_token"`
	}
	err = doQuery(ctx, req, &t)
	if err != nil {
		return err
	}
	token = "Bearer " + t.Token
	return nil
}

func Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	if err := authenticate(ctx); err != nil {
		return err
	}
	if err := callSemaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	defer callSemaphore.Release(1)
	return client.Query(ctx, q, variables)
}

func QueryRaw(ctx context.Context, q interface{}, variables map[string]interface{}) (*json.RawMessage, error) {
	if err := authenticate(ctx); err != nil {
		return nil, err
	}
	if err := callSemaphore.Acquire(ctx, 1); err != nil {
		return nil, err
	}
	defer callSemaphore.Release(1)
	return client.QueryRaw(ctx, q, variables)
}

type Error struct {
	Status int    `json:"status"`
	Err    string `json:"error"`
}

func (e Error) Error() string {
	return e.Err
}

func doQuery(ctx context.Context, req *http.Request, dest interface{}) error {
	if err := callSemaphore.Acquire(ctx, 1); err != nil {
		return err
	}
	defer callSemaphore.Release(1)
	start := time.Now()
	resp, err := http.DefaultClient.Do(req)
	timing := time.Since(start)
	logrus.WithFields(logrus.Fields{
		"duration_ms": timing.Milliseconds(),
		"status_code": resp.StatusCode,
		"error":       err,
	}).Debug("WCL requested")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode > 299 {
		e := Error{}
		if err := json.Unmarshal(body, &e); err != nil {
			return err
		}
		logrus.WithError(err).Warn("error while querying WCL")
		return e
	}
	if err := json.Unmarshal(body, dest); err != nil {
		return err
	}
	return nil
}
