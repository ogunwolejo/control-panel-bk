package config

import (
	"os"
	"testing"
)

func TestDefaultPayStackConfiguration(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("PAY_STACK_HOME_NAME", "api.paystack.co")
	os.Setenv("PAY_STACK_PUBLIC_KEY", "test_public_key")
	os.Setenv("PAY_STACK_SECRET_KEY", "test_secret_key")
	os.Setenv("PAY_STACK_CALLBACK_URL", "http://localhost/callback")
	os.Setenv("PAY_STACK_WEBHOOK_URL", "http://localhost/webhook")

	// Call the function to get the default configuration
	config := DefaultPayStackConfiguration()

	// Validate the configuration
	if config.Host != "api.paystack.co" {
		t.Errorf("Expected Host 'api.paystack.co', got '%s'", config.Host)
	}
	if config.PublicKey != "test_public_key" {
		t.Errorf("Expected PublicKey 'test_public_key', got '%s'", config.PublicKey)
	}
	if config.SecretKey != "test_secret_key" {
		t.Errorf("Expected SecretKey 'test_secret_key', got '%s'", config.SecretKey)
	}
	if config.CallBackUrl != "http://localhost/callback" {
		t.Errorf("Expected CallBackUrl 'http://localhost/callback', got '%s'", config.CallBackUrl)
	}
	if config.WebHookUrl != "http://localhost/webhook" {
		t.Errorf("Expected WebHookUrl 'http://localhost/webhook', got '%s'", config.WebHookUrl)
	}
	if config.Port != 443 {
		t.Errorf("Expected Port 443, got %d", config.Port)
	}
	if config.Headers.Authorization != "Bearer test_secret_key" {
		t.Errorf("Expected Authorization 'Bearer test_secret_key', got '%s'", config.Headers.Authorization)
	}
	if config.Headers.ContentType != "application/json" {
		t.Errorf("Expected ContentType 'application/json', got '%s'", config.Headers.ContentType)
	}
}

func TestPlanUrl(t *testing.T) {
	// Set environment variables for testing
	os.Setenv("PAY_STACK_HOME_NAME", "api.paystack.co")
	os.Setenv("PAY_STACK_SECRET_KEY", "test_secret_key")

	// Get the default configuration
	config := DefaultPayStackConfiguration()

	// Call the PlanUrl method
	url := config.PlanUrl()

	// Validate the URL
	expectedUrl := "https://api.paystack.co:443/plan"
	if url != expectedUrl {
		t.Errorf("Expected URL '%s', got '%s'", expectedUrl, url)
	}
}

func TestPlanUrl_NilPayStack(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic when calling PlanUrl on nil PayStack, but did not panic")
		}
	}()

	var p *PayStack
	p.PlanUrl() // This should panic
}

func BenchmarkDefaultPayStackConfiguration(b *testing.B) {
	// Set environment variables for testing
	os.Setenv("PAY_STACK_HOME_NAME", "api.paystack.co")
	os.Setenv("PAY_STACK_PUBLIC_KEY", "test_public_key")
	os.Setenv("PAY_STACK_SECRET_KEY", "test_secret_key")
	os.Setenv("PAY_STACK_CALLBACK_URL", "http://localhost/callback")
	os.Setenv("PAY_STACK_WEBHOOK_URL", "http://localhost/webhook")

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		DefaultPayStackConfiguration()
	}
}

func BenchmarkPlanUrl(b *testing.B) {
	// Set environment variables for testing
	os.Setenv("PAY_STACK_HOME_NAME", "api.paystack.co")
	os.Setenv("PAY_STACK_SECRET_KEY", "test_secret_key")

	// Get the default configuration once
	config := DefaultPayStackConfiguration()

	// Run the benchmark for the PlanUrl method
	for i := 0; i < b.N; i++ {
		config.PlanUrl()
	}
}
