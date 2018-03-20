package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/JulzDiverse/cfclient"
	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"
)

func main() {
	downloadUrl := os.Getenv("DOWNLOAD_URL")
	uploadUrl := os.Getenv("UPLOAD_URL")
	appId := os.Getenv("APP_ID")
	stagingGuid := os.Getenv("STAGING_GUID")
	completionCallback := os.Getenv("COMPLETION_CALLBACK") //TODO: implement callback path

	username := os.Getenv("CF_USERNAME")
	password := os.Getenv("CF_PASSWORD")
	apiAddress := os.Getenv("API_ADDRESS")

	fmt.Println("STARTING WITH:", downloadUrl, uploadUrl, appId, stagingGuid, completionCallback)

	downloader := Downloader{&http.Client{}}
	uploader := Uploader{&http.Client{}}
	cfclient, err := cfclient.NewClient(&cfclient.Config{
		SkipSslValidation: true,
		Username:          username,
		Password:          password,
		ApiAddress:        apiAddress,
	})

	exitWithError(err)

	err = downloader.DownloadWithCfClient(cfclient, appId, "/workspace/appbits")
	exitWithError(err)

	// TODO: Replace this with pure-go implementation of an unzipper
	err = execCmd(
		"unzip", []string{
			"/workspace/appbits",
		})
	exitWithError(err)

	err = os.Remove("/workspace/appbits")
	exitWithError(err)

	// TODO: Replace this with an inplace call of packs library
	err = execCmd(
		"/packs/builder", []string{
			"-buildpacksDir", "/var/lib/buildpacks",
			"-outputDroplet", "/out/droplet.tgz",
			"-outputBuildArtifactsCache", "/cache/cache.tgz",
			"-outputMetadata", "/out/result.json",
		})
	exitWithError(err)

	fmt.Println("Start Upload Process.")
	err = uploader.UploadWithCfClient(cfclient, appId, "/out/droplet.tgz")
	exitWithError(err)
}

//TODO: Don't use this unzipper it does weired stuff. Needs to be reimplemented.
func Unzip(archive, target string) error {
	if archive == "" || target == "" {
		return errors.New("source or destination path not defined")
	}

	reader, err := zip.OpenReader(archive)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
			continue
		}

		os.MkdirAll(filepath.Dir(path), file.Mode())

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

func exitWithError(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func execCmd(cmdname string, args []string) error {
	cmd := exec.Command(cmdname, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		return errors.Wrap(err, ansi.Sprintf("@R{Failed to run %s}", cmdname))
	}

	return nil
}
