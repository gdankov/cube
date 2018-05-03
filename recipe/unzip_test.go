package main_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/julz/cube/recipe"
)

var _ = Describe("Unzip function", func() {

	var (
		targetDir string
		srcZip    string
		err       error
	)

	JustBeforeEach(func() {
		err = Unzip(srcZip, targetDir)
	})

	Context("Unzip succeeds", func() {

		zippedFiles := map[string]string{
			"file1":                       "this is the content of test file 1",
			"innerDir/file2":              "this is the content of test file 2",
			"innerDir/innermostDir/file3": "this is the content of test file 3",
		}

		getRoot := func(path string) string {
			pathParts := strings.Split(path, "/")
			return pathParts[0]
		}

		removeFile := func(file string) {
			if err := os.RemoveAll(file); err != nil {
				panic(err)
			}
		}

		cleanUpFiles := func() {
			for filePath := range zippedFiles {
				rootDir := getRoot(filePath)
				removeFile(rootDir)
			}
		}

		assertFileContents := func(file string, expectedContent string) {
			path := filepath.Join(targetDir, file)
			content, err := ioutil.ReadFile(path)
			if err != nil {
				panic(err)
			}
			Expect(content).To(Equal([]byte(expectedContent)))
		}

		assertFilesUnzippedSuccessfully := func() {
			It("should not fail", func() {
				Expect(err).ToNot(HaveOccurred())
			})

			It("should unzip the files in the target directory", func() {
				for fileName := range zippedFiles {
					path := filepath.Join(targetDir, fileName)
					Expect(path).To(BeAnExistingFile())
				}
			})

			It("should not change file contents", func() {
				for file, expectedContent := range zippedFiles {
					assertFileContents(file, expectedContent)
				}
			})
		}

		Context("When target directory is not specified", func() {

			BeforeEach(func() {
				srcZip = "testdata/unzip_me.zip"
				targetDir = ""
			})

			AfterEach(func() {
				cleanUpFiles()
			})

			assertFilesUnzippedSuccessfully()
		})

		Context("When target directory is not empty string", func() {
			BeforeEach(func() {
				srcZip = "testdata/unzip_me.zip"
				targetDir = "testdata/tmp"

				err := os.Mkdir(targetDir, 0755)
				if err != nil {
					panic(err)
				}
			})

			AfterEach(func() {
				err := os.RemoveAll(targetDir)
				if err != nil {
					panic(err)
				}
			})

			assertFilesUnzippedSuccessfully()
		})
	})

	Context("Unzip fails", func() {
		Context("When target dir does not exist", func() {

			BeforeEach(func() {
				targetDir = "non-existent"
				srcZip = "testdata/unzip_me.zip"
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
			})

		})

		Context("When target dir is not a directory", func() {

			BeforeEach(func() {
				targetDir = "testdata/unzip_me.zip"
				srcZip = "testdata/unzip_me.zip"
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
			})

		})

		Context("When source zip archive does not exist", func() {

			BeforeEach(func() {
				targetDir = "testdata"
				srcZip = "non-existent"
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
			})

		})

		Context("When source is not a zip archive", func() {

			BeforeEach(func() {
				targetDir = "testdata"
				srcZip = "testdata/file.notzip"
			})

			It("should fail", func() {
				Expect(err).To(HaveOccurred())
			})

		})
	})

})
