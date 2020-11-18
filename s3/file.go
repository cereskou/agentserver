package s3

import "time"

//File -
type File struct {
	Err     error
	Bucket  string
	Key     string
	SaveTo  string
	Size    int64
	Md5     string
	S3      string
	ModTime time.Time
	Num     uint32
	Offset  uint64
	Length  uint64
}

//URI -
type URI struct {
	Bucket string
	Prefix string
}
