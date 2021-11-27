package common

import "encoding/json"

func SetupJSONEncoder(enc *json.Encoder) {
	enc.SetIndent("", "\t")
}
