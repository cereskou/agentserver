package cx

import jsoniter "github.com/json-iterator/go"

//File -
type File struct {
	S3     string `json:"s3"`
	Local  string `json:"local"`
	Size   uint64 `json:"size"`
	Num    uint32 `json:"num"`
	Offset uint64 `json:"offset"`
	Length uint64 `json:"length"`
}

//Result -
type Result struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
	Error      error  `json:"error"`
}

//ToJSON Convert to JSOn
func ToJSON(files []*File) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	bytes, _ := json.Marshal(&files)

	return string(bytes)
}
