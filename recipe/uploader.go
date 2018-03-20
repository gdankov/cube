package main

import (
	"net/http"

	"github.com/JulzDiverse/cfclient"
	"github.com/pkg/errors"
)

type Uploader struct {
	Client *http.Client
}

func (u *Uploader) UploadWithCfClient(cfclient *cfclient.Client, guid string, path string) error {
	if guid == "" {
		return errors.New("empty-guid-provided")
	}

	if path == "" {
		return errors.New("empty-filepath-provided")
	}

	err := cfclient.PushDroplet("droplet", guid)
	if err != nil {
		return errors.Wrap(err, "perform-request-failed")
	}

	return nil
}
