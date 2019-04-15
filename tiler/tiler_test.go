package tiler_test

import (
	"context"
	"log"

	"github.com/starkandwayne/om-tiler/steps"
	. "github.com/starkandwayne/om-tiler/tiler"
	"github.com/starkandwayne/om-tiler/tiler/tilerfakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegisterStep", func() {
	var (
		tiler *Tiler
		step  steps.Step
	)

	JustBeforeEach(func() {
		fakeOpsman := &tilerfakes.FakeOpsmanClient{}
		fakeMover := &tilerfakes.FakeMover{}
		logger := log.New(GinkgoWriter, "", 0)
		tiler = NewTiler(fakeOpsman, fakeMover, logger)
	})

	Context("given a step with a disallowed dependency", func() {
		BeforeEach(func() {
			step = steps.Step{
				Name:      "disallowedStep",
				DependsOn: []string{"nonExistingStep"},
			}
		})

		It("Does not register", func() {
			err := tiler.RegisterStep(BuildCallback, step)
			Expect(err).To(MatchError(
				"BuildCallback: disallowedStep may not DependOn: nonExistingStep"))
		})

	})

	Context("given a valid step with callback", func() {
		var callbackHasBeenCalled bool

		BeforeEach(func() {
			step = steps.Step{
				Name:      "fooStep",
				DependsOn: []string{StepWaitOpsmanOnline},
				Do: func(ctx context.Context) error {
					callbackHasBeenCalled = true
					return nil
				},
			}
		})

		It("Does not register", func() {
			err := tiler.RegisterStep(DeleteCallback, step)
			Expect(err).ToNot(HaveOccurred())
			err = tiler.Delete(context.Background())
			Expect(err).ToNot(HaveOccurred())
			Expect(callbackHasBeenCalled).To(Equal(true))
		})

	})

})
