// +build !jsoniter

package json

import (
	"github.com/pquerna/ffjson/ffjson"
)

var (
	Marshal   = ffjson.Marshal
	Unmarshal = ffjson.Unmarshal
)
