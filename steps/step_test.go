package steps_test

import (
	"bytes"
	"context"
	"errors"
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
		logger  *log.Logger
		buffer  *bytes.Buffer
		steps   []Step
		attempt int
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = log.New(buffer, "root", 0)
		attempt = 1
	})

	Context("Given multiple steps", func() {
		BeforeEach(func() {
			steps = []Step{
				{
					Name:      StepFoo,
					DependsOn: []string{StepBar},
					Do: func(ctx context.Context) error {
						ContextLogger(ctx, logger, "[Foo]").Println("hello foo")
						return nil
					}},
				{
					Name: StepBar,
					Do: func(ctx context.Context) error {
						ContextLogger(ctx, logger, "[Bar]").Println("hello bar")
						return nil
					},
				},
				{
					Name:      "NullStep",
					DependsOn: []string{StepBar},
				},
				{
					Name:      "TrySomething",
					DependsOn: []string{StepFoo},
					Do: func(ctx context.Context) error {
						ContextLogger(ctx, logger, "[OM]").Println("Not today")
						if attempt < 2 {
							attempt++
							return errors.New("Failed")
						}
						return nil
					},
					Retry: 5,
				},
			}
		})

		It("Prefixes logs", func() {
			err := Run(context.Background(), steps, logger)
			Expect(err).ToNot(HaveOccurred())
			lines := strings.Split(string(buffer.Bytes()), "\n")
			Expect(lines[0]).To(Equal("root [Bar] BarStep hello bar"))
			Expect(lines[1]).To(Equal("root [Foo] FooStep hello foo"))
			Expect(lines[2]).To(Equal("root [OM] TrySomething Not today"))
			Expect(lines[3]).To(Equal("root [Steps] TrySomething Attempt 1 retrying error: Failed"))
			Expect(lines[4]).To(Equal("root [OM] TrySomething Not today"))
		})
	})

})
