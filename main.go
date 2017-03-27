package main

import (
	"context"
	"fmt"
	"os"
	"reflect"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/deoxxa/aws_signing_client"
)

// 必要最低限のフィールドだけでも大丈夫っぽい
type Tweet struct {
	User string `json:"user"`
	Text string `json:"text"`
}

// AWS_ES_ENDPOINT=https://search-XXX.ap-northeast-1.es.amazonaws.com AWS_ACCESS_KEY_ID=AXXXA AWS_SECRET_ACCESS_KEY=aXXX go run main.go

func main() {
	creds := credentials.NewEnvCredentials()

	signer := v4.NewSigner(creds)
	awsClient, _ := aws_signing_client.New(signer, nil, "es", "ap-northeast-1")

	client, _ := elastic.NewClient(
		elastic.SetURL(os.Getenv("AWS_ES_ENDPOINT")),
		elastic.SetScheme("https"),
		elastic.SetHttpClient(awsClient),
		elastic.SetSniff(false),
	)

	// https://gist.github.com/olivere/114347ff9d9cfdca7bdc0ecea8b82263
	query := elastic.NewTermQuery("text", "mackerel")
	searchResult, _ := client.Search().Index("twitter_public_timeline").Query(query).Do(context.Background())

	var ttyp Tweet
	for _, item := range searchResult.Each(reflect.TypeOf(ttyp)) {
		if t, ok := item.(Tweet); ok {
			fmt.Printf("Tweet by %s: %s\n", t.User, t.Text)
		}
	}
}
