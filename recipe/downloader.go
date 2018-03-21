package main

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/julz/cube"
	"github.com/pkg/errors"
)

type Downloader struct {
	Cfclient cube.CfClient
}

func (d *Downloader) Download(appId string, filepath string) error {
	if appId == "" {
		return errors.New("empty appId provided")
	}

	if filepath == "" {
		return errors.New("empty filepath provided")
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}

	resp, err := d.Cfclient.GetAppBitsByAppGuid(appId)
	if err != nil {
		return errors.Wrap(err, "failed to perform request")
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Download failed. Status Code %s", resp.StatusCode))
	}

	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to copy content to file")
	}

	return nil
}

//func (d *Downloader) Download(url string, filename string) error {
//if url == "" {
//return errors.New("empty url provided")
//}

//if filename == "" {
//return errors.New("empty filename provided")
//}

//file, err := os.Create(filename)
//if err != nil {
//return err
//}

//req, err := http.NewRequest("GET", url, nil)
//if err != nil {
//return errors.Wrap(err, "failed to create http request")
//}

//resp, err := d.Client.Do(req)
//if err != nil {
//return errors.Wrap(err, "failed to perform request")
//}

//if resp.StatusCode != http.StatusOK {
//return errors.New(fmt.Sprintf("Download failed. Status Code %s", resp.StatusCode))
//}

//defer resp.Body.Close()

//_, err = io.Copy(file, resp.Body)
//if err != nil {
//return errors.Wrap(err, "failed to copy content to file")
//}

//return nil
//}
