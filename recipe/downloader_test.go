package main_test

import (
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/julz/cube/recipe"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Downloader", func() {
	var (
		downloader *Downloader
		fakeServer *ghttp.Server
	)

	BeforeEach(func() {
		downloader = &Downloader{&http.Client{}}
		fakeServer = ghttp.NewServer()
	})

	It("contains a http.Client field", func() {
		Expect(downloader.Client).Should(BeAssignableToTypeOf(&http.Client{}))
	})

	Context("Download", func() {

		It("should return an error if an empty url is provided", func() {
			err := downloader.Download("", "")
			Expect(err).To(HaveOccurred())

			Expect(err).To(MatchError(ContainSubstring("empty url provided")))
		})

		It("should return an error if an empty file name is provided", func() {
			err := downloader.Download("http://download-me.com", "")
			Expect(err).To(HaveOccurred())

			Expect(err).To(MatchError(ContainSubstring("empty filename provided")))
		})

		It("performs an request against the provided url", func() {
			fakeServer.AppendHandlers(
				ghttp.VerifyRequest("GET", "/download-me"),
			)

			err := downloader.Download(fakeServer.URL()+"/download-me", "file")
			Expect(err).ToNot(HaveOccurred())
			Expect(fakeServer.ReceivedRequests()).Should(HaveLen(1))
		})

		Context("When the downlad request is successful", func() {

			BeforeEach(func() {
				fakeServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/download-me"),
						ghttp.RespondWith(http.StatusOK, `appbits`),
					),
				)
			})

			AfterEach(func() {
				err := os.Remove("test/file")
				Expect(err).ToNot(HaveOccurred())
			})

			It("writes the downloaded content to the given file", func() {
				err := downloader.Download(fakeServer.URL()+"/download-me", "test/file")
				Expect(err).ToNot(HaveOccurred())
				Expect("test/file").Should(BeAnExistingFile())

				file, err := ioutil.ReadFile("test/file")
				Expect(err).ToNot(HaveOccurred())
				Expect(string(file)).To(Equal("appbits"))
			})
		})

		Context("When the download request fails", func() {
			BeforeEach(func() {
				fakeServer.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", "/download-me"),
						ghttp.RespondWith(http.StatusInternalServerError, nil),
					),
				)
			})

			It("should error with an corresponding error message", func() {
				err := downloader.Download(fakeServer.URL()+"/download-me", "test/file")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError(ContainSubstring("Download failed")))
			})
		})
	})
})
