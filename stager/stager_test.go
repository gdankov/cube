package stager_test

import (
	"context"

	"github.com/julz/cube"
	"github.com/julz/cube/opi"
	. "github.com/julz/cube/stager"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stager", func() {
	Context("Run", func() {

		var (
			task   opi.Task
			stager cube.Stager
		)

		BeforeEach(func() {
			task = opi.Task{}
			stager = Stager{
				Desirer: opi.DesireTaskFunc(func(_ context.Context, tasks []opi.Task) error {
					return nil
				}),
			}
		})

		It("converts and desires a staging request to a Task", func() {
			err := stager.Run(task)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
