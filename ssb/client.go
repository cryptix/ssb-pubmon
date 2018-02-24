package ssb

import (
	"encoding/base64"

	"cryptoscope.co/go/secretstream"
	"cryptoscope.co/go/secretstream/secrethandshake"
)

var SHSClient *secretstream.Client

func init() {
	ssbAppKey, err := base64.StdEncoding.DecodeString("1KHLiKZvAvjbY1ziZEHMXawbCEIM6qwjCDm3VYRan/s=")
	if err != nil {
		panic(err)
	}
	kp, err := secrethandshake.GenEdKeyPair(nil)
	if err != nil {
		panic(err)
	}
	SHSClient, err = secretstream.NewClient(*kp, ssbAppKey)
	if err != nil {
		panic(err)
	}
}
