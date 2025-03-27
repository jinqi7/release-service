package models

import "time"

type Model interface {
	Validate() error
}

func UnixMilli() int64 {
	return time.Now().UTC().UnixNano() / int64(time.Millisecond)
}
