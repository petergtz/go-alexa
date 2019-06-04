package dynamodb_test

import (
	"fmt"
	"time"

	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/petergtz/go-alexa"
	"github.com/petergtz/go-alexa/dynamodb"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

var _ = Describe("Persistence", func() {
	It("works", func() {
		logger, e := zap.NewDevelopment()
		Expect(e).NotTo(HaveOccurred())
		if os.Getenv("ACCESS_KEY_ID") == "" {
			logger.Fatal("env var ACCESS_KEY_ID not provided.")
		}

		if os.Getenv("SECRET_ACCESS_KEY") == "" {
			logger.Fatal("env var SECRET_ACCESS_KEY not provided.")
		}

		persistence := dynamodb.NewInteractionLogger(os.Getenv("ACCESS_KEY_ID"), os.Getenv("SECRET_ACCESS_KEY"), logger.Sugar(), "TestAlexaWikipediaRequests")
		persistence.Log(&alexa.Interaction{
			RequestID:     uuid.NewV4().String(),
			RequestType:   "RequestTestType",
			UnixTimestamp: time.Now().Unix(),
			Timestamp:     time.Now(),
			Locale:        "de-DE",
			UserID:        "userid1",
			SessionID:     "sessionid1",
			Attributes: map[string]interface{}{
				"SearchQuery": "Bla",
				"ActualTitle": "blub",
			},
		})

		persistence.Log(&alexa.Interaction{
			RequestID:     uuid.NewV4().String(),
			RequestType:   "RequestTestType2",
			UnixTimestamp: time.Now().Unix(),
			Timestamp:     time.Now(),
			Locale:        "de-DE",
			UserID:        "userid1",
			SessionID:     "sessionid2",
			Attributes: map[string]interface{}{
				"SearchQuery": "Bla2",
				"ActualTitle": "blub2",
			},
		})

		interactions := persistence.GetInteractionsByUser("userid1")
		for _, interaction := range interactions {
			fmt.Printf("%#v\n", *interaction)
		}
	})
})
