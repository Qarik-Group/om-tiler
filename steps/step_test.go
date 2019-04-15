package steps_test

import (
	"bytes"
	"context"
	"log"
	"strings"

	. "github.com/starkandwayne/om-tiler/steps"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	StepFoo string = "FooStep"
	StepBar        = "BarStep"
)

var _ = Describe("Step", func() {
	var (
		logger *log.Logger
		buffer *bytes.Buffer
		steps  []Step
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = log.New(buffer, "root", 0)

	})

	Context("Given multiple steps", func() {
		BeforeEach(func() {
			steps = []Step{
				Step{
					Name:      StepFoo,
					DependsOn: []string{StepBar},
					Do: func(ctx context.Context) error {
						ContextLogger(ctx, logger, "[Foo]").Println("hello foo")
						return nil
					}},
				Step{
					Name: StepBar,
					Do: func(ctx context.Context) error {
						ContextLogger(ctx, logger, "[Bar]").Println("hello bar")
						return nil
					},
				},
				Step{
					Name:      "NullStep",
					DependsOn: []string{StepBar},
				},
			}
		})

		It("Prefixes logs", func() {
			err := Run(context.Background(), steps)
			Expect(err).ToNot(HaveOccurred())
			lines := strings.Split(string(buffer.Bytes()), "\n")
			Expect(lines[0]).To(Equal("root [Bar] BarStep hello bar"))
			Expect(lines[1]).To(Equal("root [Foo] FooStep hello foo"))
		})
	})

})
