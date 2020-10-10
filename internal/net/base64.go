package net

import "encoding/base64"

func Base64EncodeToString(s string) string {
	return base64.URLEncoding.EncodeToString([]byte(s))
}
