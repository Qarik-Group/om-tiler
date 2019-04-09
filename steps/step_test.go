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

var _ = Describe("Step", func() {
	var (
		logger  *log.Logger
		buffer  *bytes.Buffer
		stepFoo func(r map[string]interface{}) (interface{}, error)
		stepBar func(r map[string]interface{}) (interface{}, error)
	)

	BeforeEach(func() {
		buffer = bytes.NewBuffer([]byte{})
		logger = log.New(buffer, "root", 0)

	})

	Context("Given multiple steps", func() {
		BeforeEach(func() {
			ctx := context.Background()
			stepFoo = Step(ctx, "stepFoo", func(ctx context.Context) error {
				ContextLogger(ctx, logger, "[Foo]").Println("hello foo")
				return nil
			})
			stepBar = Step(ctx, "stepBar", func(ctx context.Context) error {
				ContextLogger(ctx, logger, "[Bar]").Println("hello bar")
				return nil
			})
		})

		It("Prefixes logs", func() {
			_, err := stepFoo(map[string]interface{}{})
			Expect(err).ToNot(HaveOccurred())
			_, err = stepBar(map[string]interface{}{})
			Expect(err).ToNot(HaveOccurred())
			lines := strings.Split(string(buffer.Bytes()), "\n")
			Expect(lines[0]).To(Equal("root [Foo] stepFoo hello foo"))
			Expect(lines[1]).To(Equal("root [Bar] stepBar hello bar"))
		})
	})

})
