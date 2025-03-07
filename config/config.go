package config

import (
	"fmt"
	"log"
	"os"
)

type Headers struct {
	Authorization string
	ContentType   string
}

type PayStack struct {
	Host        string
	PublicKey   string
	SecretKey   string
	CallBackUrl string
	WebHookUrl  string
	Port        int
	Headers     Headers
}

var PayStackConfig PayStack

func DefaultPayStackConfiguration() *PayStack {
	PayStackConfig = PayStack{
		os.Getenv("PAY_STACK_HOME_NAME"),
		os.Getenv("PAY_STACK_PUBLIC_KEY"),
		os.Getenv("PAY_STACK_SECRET_KEY"),
		os.Getenv("PAY_STACK_CALLBACK_URL"),
		os.Getenv("PAY_STACK_WEBHOOK_URL"),
		443,
		Headers{
			fmt.Sprintf("Bearer %s", os.Getenv("PAY_STACK_SECRET_KEY")),
			"application/json",
		},
	}

	log.Println("Result Pay stack config ", PayStackConfig)
	return &PayStackConfig
}
