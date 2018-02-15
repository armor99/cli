package main

import "encoding/base64"

func encodeBasicAuth(user string, passwd string) string {
	data := []byte(user + ":" + passwd)
	str := base64.StdEncoding.EncodeToString(data)
	authHdr := "Basic " + str
	return authHdr
}
