package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/JulzDiverse/cfclient"
	"github.com/julz/cube"
	"github.com/pkg/errors"
	"github.com/starkandwayne/goutils/ansi"

	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/runtimeschema/cc_messages"
)

func main() {
	downloadUrl := os.Getenv(cube.EnvDownloadUrl)
	uploadUrl := os.Getenv(cube.EnvUploadUrl)
	appId := os.Getenv(cube.EnvAppId)
	stagingGuid := os.Getenv(cube.EnvStagingGuid)
	completionCallback := os.Getenv(cube.EnvCompletionCallback)

	username := os.Getenv(cube.EnvCfUsername)
	password := os.Getenv(cube.EnvCfPassword)
	apiAddress := os.Getenv(cube.EnvApiAddress)
	cubeAddress := os.Getenv(cube.EnvCubeAddress)

	fmt.Println("STARTING WITH:", downloadUrl, uploadUrl, appId, stagingGuid, completionCallback)

	cfclient, err := cfclient.NewClient(&cfclient.Config{
		SkipSslValidation: true,
		Username:          username,
		Password:          password,
		ApiAddress:        apiAddress,
	})

	downloader := Downloader{cfclient}
	uploader := Uploader{cfclient}

	exitWithError(err)

	err = downloader.Download(appId, "/workspace/appbits")
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
	err = uploader.Upload(appId, "/out/droplet.tgz")
	exitWithError(err)

	fmt.Println("Upload successful!")
	result, err := readResultJson("/out/result.json")
	exitWithError(err)

	annotation := cc_messages.StagingTaskAnnotation{
		CompletionCallback: completionCallback,
	}

	annotationJson, err := json.Marshal(annotation)
	exitWithError(err)

	stagingCompleteResponse(
		cubeAddress,
		stagingGuid,
		string(annotationJson[:len(annotationJson)]),
		string(result[:len(result)]),
	)
	fmt.Println("Staging completed")
}

func readResultJson(path string) ([]byte, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to read result.json")
	}
	return file, nil
}

func stagingCompleteResponse(cubeAddress, stagingGuid, annotation string, result string) error {

	callbackResponse := models.TaskCallbackResponse{
		TaskGuid:   stagingGuid,
		Result:     result,
		Failed:     false,
		Annotation: annotation,
	}

	jsonBytes := new(bytes.Buffer)
	json.NewEncoder(jsonBytes).Encode(callbackResponse)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v1/staging/%s/completed", cubeAddress, stagingGuid), jsonBytes)
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}

	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "request failed")
	}

	if resp.StatusCode >= 400 {
		return errors.New("Request not successful")
	}

	return nil
}

//TODO: Don't use this unzipper it does weired stuff.
//Needs to be reimplemented/improved.
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
