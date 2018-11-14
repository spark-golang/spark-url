// +build jsoniter

package json

import (
	"github.com/json-iterator/go"
)

var (
	Marshal   = jsoniter.Marshal
	Unmarshal = jsoniter.Unmarshal
)
