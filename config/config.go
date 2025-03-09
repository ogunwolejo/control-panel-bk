package config

import (
	"fmt"
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

	return &PayStackConfig
}

func (p *PayStack) PlanUrl() string {
	if p == nil {
		panic("didn't initialized paystack")
	}

	s := fmt.Sprintf("https://%s:%d/plan", p.Host, p.Port)
	return s
}
