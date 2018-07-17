package ssb

import (
	"encoding/base64"

	"github.com/pkg/errors"

	"go.cryptoscope.co/secretstream"
	"go.cryptoscope.co/secretstream/secrethandshake"
)

var SHSClient *secretstream.Client

func InitClient(fname string) error {
	ssbAppKey, err := base64.StdEncoding.DecodeString("1KHLiKZvAvjbY1ziZEHMXawbCEIM6qwjCDm3VYRan/s=") // default appKey
	if err != nil {
		return errors.Wrap(err, "ssb appkey?!")
	}

	kp, err := secrethandshake.LoadSSBKeyPair(fname)
	if err != nil {
		return errors.Wrap(err, "load ssb keyPair")
	}

	SHSClient, err = secretstream.NewClient(*kp, ssbAppKey)
	return errors.Wrap(err, "ssb NewClient")
}
