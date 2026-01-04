package lambda_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/petergtz/go-alexa/lambda"
)

var _ = Describe("Lambda", func() {
	It("WithSnippets truncates long values in map", func() {
		snippets := lambda.WithSnippets(map[string]interface{}{
			"k1": "short_value",
			"k2": "this is a very long value that should be truncated by the WithSnippets function because it exceeds the maximum length allowed",
		}, 20)
		Expect(snippets).To(HaveKeyWithValue("k1", "short_value"))
		Expect(snippets).To(HaveKeyWithValue("k2", `"this is a very long... (107 B truncated)`))
	})

	It("WithSnippets properly truncates maps", func() {
		snippets := lambda.WithSnippets(map[string]interface{}{
			"key": map[string]interface{}{
				"nested_key": "this is a very long nested value that should be truncated by the WithSnippets function because it exceeds the maximum length allowed",
			},
		}, 30)
		Expect(snippets).To(HaveKeyWithValue("key", `{"nested_key":"this is a very ... (119 B truncated)`))
	})
})
