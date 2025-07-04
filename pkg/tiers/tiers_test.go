package tiers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// Mock HTTP server for testing
func mockServer(statusCode int, response interface{}) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
	}))
}

func TestCreateTier(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   TierResponse
		mockStatusCode int
		request        CreateTierRequest
		expectedError  error
		expectedCode   int
	}{
		{
			name: "Success - Tier Created",
			mockResponse: TierResponse{
				Status:  true,
				Message: "Tier created successfully",
				Data: struct {
					Name         string    `json:"name"`
					Amount       int       `json:"amount"`
					Interval     string    `json:"interval"`
					Integration  int       `json:"integration"`
					Domain       string    `json:"domain"`
					PlanCode     string    `json:"plan_code"`
					SendInvoices bool      `json:"send_invoices"`
					SendSms      bool      `json:"send_sms"`
					HostedPage   bool      `json:"hosted_page"`
					Currency     string    `json:"currency"`
					Id           int       `json:"id"`
					CreatedAt    time.Time `json:"createdAt"`
					UpdatedAt    time.Time `json:"updatedAt"`
				}{
					Name:     "Test Tier",
					Amount:   1000,
					Interval: string(IntervalMonthly),
					Currency: string(CurrencyUSD),
				},
			},
			mockStatusCode: http.StatusCreated,
			request: CreateTierRequest{
				Name:     "Test Tier",
				Amount:   1000,
				Interval: IntervalMonthly,
				Currency: CurrencyNGN,
			},
			expectedError: nil,
			expectedCode:  http.StatusCreated,
		},
		{
			name:           "Error - Duplicate Tier",
			mockResponse:   TierResponse{},
			mockStatusCode: http.StatusConflict,
			request: CreateTierRequest{
				Name:     "Test Tier",
				Amount:   1000,
				Interval: IntervalMonthly,
				Currency: CurrencyNGN,
			},
			expectedError: errors.New("an active tier with the same data already exists"),
			expectedCode:  http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mockServer(tt.mockStatusCode, tt.mockResponse)
			defer server.Close()

			// Override the global URL with the mock server URL
			url = server.URL

			response, err, code := CreateTier(tt.request, context.Background())
			if err != nil && response == nil {
				t.Logf("Except error code %d not equals to 201", code)
			}

			if response != nil {
				t.Logf("Except response status %d to be equals to %d", tt.expectedCode, code)
				t.Logf("Except response data name %s equals to %s", tt.name, response.Data.Name)
			}
		})
	}
}

func TestGetTier(t *testing.T) {
	tests := []struct {
		name           string
		planCode       string
		mockResponse   FetchTierResponse
		mockStatusCode int
		expectedError  error
		expectedCode   int
	}{
		{
			name:     "Success - Get Tier",
			planCode: "PLAN123",
			mockResponse: FetchTierResponse{
				Status:  true,
				Message: "Tier retrieved successfully",
				Data: struct {
					Subscriptions []struct {
						Customer         int    `json:"customer"`
						Plan             int    `json:"plan"`
						Integration      int    `json:"integration"`
						Domain           string `json:"domain"`
						Start            int    `json:"start"`
						Status           string `json:"status"`
						Quantity         int    `json:"quantity"`
						Amount           int    `json:"amount"`
						SubscriptionCode string `json:"subscription_code"`
						EmailToken       string `json:"email_token"`
						Authorization    struct {
							AuthorizationCode string `json:"authorization_code"`
							Bin               string `json:"bin"`
							Last4             string `json:"last4"`
							ExpMonth          string `json:"exp_month"`
							ExpYear           string `json:"exp_year"`
							Channel           string `json:"channel"`
							CardType          string `json:"card_type"`
							Bank              string `json:"bank"`
							CountryCode       string `json:"country_code"`
							Brand             string `json:"brand"`
							Reusable          bool   `json:"reusable"`
							Signature         string `json:"signature"`
							AccountName       string `json:"account_name"`
						} `json:"authorization"`
						EasyCronId      interface{} `json:"easy_cron_id"`
						CronExpression  string      `json:"cron_expression"`
						NextPaymentDate time.Time   `json:"next_payment_date"`
						OpenInvoice     interface{} `json:"open_invoice"`
						Id              int         `json:"id"`
						CreatedAt       time.Time   `json:"createdAt"`
						UpdatedAt       time.Time   `json:"updatedAt"`
					} `json:"subscriptions"`
					Integration       int         `json:"integration"`
					Domain            string      `json:"domain"`
					Name              string      `json:"name"`
					PlanCode          string      `json:"plan_code"`
					Description       interface{} `json:"description"`
					Amount            int         `json:"amount"`
					Interval          string      `json:"interval"`
					SendInvoices      bool        `json:"send_invoices"`
					SendSms           bool        `json:"send_sms"`
					HostedPage        bool        `json:"hosted_page"`
					HostedPageUrl     interface{} `json:"hosted_page_url"`
					HostedPageSummary interface{} `json:"hosted_page_summary"`
					Currency          string      `json:"currency"`
					Id                int         `json:"id"`
					CreatedAt         time.Time   `json:"createdAt"`
					UpdatedAt         time.Time   `json:"updatedAt"`
				}{
					Name:     "Test Tier",
					PlanCode: "PLAN123",
					Amount:   1000,
					Interval: string(IntervalMonthly),
					Currency: string(CurrencyUSD),
				},
			},
			mockStatusCode: http.StatusOK,
			expectedError:  nil,
			expectedCode:   http.StatusOK,
		},
		{
			name:           "Error - Tier Not Found",
			planCode:       "INVALID",
			mockResponse:   FetchTierResponse{},
			mockStatusCode: http.StatusNotFound,
			expectedError:  errors.New("failed to get tier: 404"),
			expectedCode:   http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mockServer(tt.mockStatusCode, tt.mockResponse)
			defer server.Close()

			// Override the global URL with the mock server URL
			url = server.URL

			response, err, code := GetTier(tt.planCode, context.Background())
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}
			if code != tt.expectedCode {
				t.Errorf("Expected status code: %d, got: %d", tt.expectedCode, code)
			}
			if response != nil && response.Message != tt.mockResponse.Message {
				t.Errorf("Expected message: %s, got: %s", tt.mockResponse.Message, response.Message)
			}
		})
	}
}

func TestFetchTiers(t *testing.T) {
	tests := []struct {
		name           string
		request        FetchTiersRequest
		mockResponse   FetchTiersResponse
		mockStatusCode int
		expectedError  error
		expectedCode   int
	}{
		{
			name: "Success - Fetch Tiers",
			request: FetchTiersRequest{
				PerPage: 10,
				Page:    1,
			},
			mockResponse: FetchTiersResponse{
				Status:  true,
				Message: "Tiers retrieved successfully",
				Data: []struct {
					Subscriptions []struct {
						Customer         int    `json:"customer"`
						Plan             int    `json:"plan"`
						Integration      int    `json:"integration"`
						Domain           string `json:"domain"`
						Start            int    `json:"start"`
						Status           string `json:"status"`
						Quantity         int    `json:"quantity"`
						Amount           int    `json:"amount"`
						SubscriptionCode string `json:"subscription_code"`
						EmailToken       string `json:"email_token"`
						Authorization    struct {
							AuthorizationCode string `json:"authorization_code"`
							Bin               string `json:"bin"`
							Last4             string `json:"last4"`
							ExpMonth          string `json:"exp_month"`
							ExpYear           string `json:"exp_year"`
							Channel           string `json:"channel"`
							CardType          string `json:"card_type"`
							Bank              string `json:"bank"`
							CountryCode       string `json:"country_code"`
							Brand             string `json:"brand"`
							Reusable          bool   `json:"reusable"`
							Signature         string `json:"signature"`
							AccountName       string `json:"account_name"`
						} `json:"authorization"`
						EasyCronId      interface{} `json:"easy_cron_id"`
						CronExpression  string      `json:"cron_expression"`
						NextPaymentDate time.Time   `json:"next_payment_date"`
						OpenInvoice     interface{} `json:"open_invoice"`
						Id              int         `json:"id"`
						CreatedAt       time.Time   `json:"createdAt"`
						UpdatedAt       time.Time   `json:"updatedAt"`
					} `json:"subscriptions"`
					Integration       int         `json:"integration"`
					Domain            string      `json:"domain"`
					Name              string      `json:"name"`
					PlanCode          string      `json:"plan_code"`
					Description       interface{} `json:"description"`
					Amount            int         `json:"amount"`
					Interval          string      `json:"interval"`
					SendInvoices      bool        `json:"send_invoices"`
					SendSms           bool        `json:"send_sms"`
					HostedPage        bool        `json:"hosted_page"`
					HostedPageUrl     interface{} `json:"hosted_page_url"`
					HostedPageSummary interface{} `json:"hosted_page_summary"`
					Currency          string      `json:"currency"`
					Id                int         `json:"id"`
					CreatedAt         time.Time   `json:"createdAt"`
					UpdatedAt         time.Time   `json:"updatedAt"`
				}{
					{
						Name:     "Test Tier",
						PlanCode: "PLAN123",
						Amount:   1000,
						Interval: string(IntervalMonthly),
						Currency: string(CurrencyUSD),
					},
				},
			},
			mockStatusCode: http.StatusOK,
			expectedError:  nil,
			expectedCode:   http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mockServer(tt.mockStatusCode, tt.mockResponse)
			defer server.Close()

			// Override the global URL with the mock server URL
			url = server.URL

			response, err, code := FetchTiers(tt.request, context.Background())
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}
			if code != tt.expectedCode {
				t.Errorf("Expected status code: %d, got: %d", tt.expectedCode, code)
			}
			if response != nil && response.Message != tt.mockResponse.Message {
				t.Errorf("Expected message: %s, got: %s", tt.mockResponse.Message, response.Message)
			}
		})
	}
}

func TestUpdateTier(t *testing.T) {
	tests := []struct {
		name           string
		planCode       string
		request        UpdateTierRequest
		mockResponse   UpdateTierResponse
		mockStatusCode int
		expectedError  error
		expectedCode   int
	}{
		{
			name:     "Success - Update Tier",
			planCode: "PLAN123",
			request: UpdateTierRequest{
				CreateTierRequest: CreateTierRequest{
					Name:     "Updated Tier",
					Amount:   2000,
					Interval: IntervalMonthly,
					Currency: CurrencyUSD,
				},
				UpdateExistingSubscriptions: true,
			},
			mockResponse: UpdateTierResponse{
				Status:  true,
				Message: "Tier updated successfully",
			},
			mockStatusCode: http.StatusOK,
			expectedError:  nil,
			expectedCode:   http.StatusOK,
		},
		{
			name:     "Error - Tier Not Found",
			planCode: "INVALID",
			request: UpdateTierRequest{
				CreateTierRequest: CreateTierRequest{
					Name:     "Updated Tier",
					Amount:   2000,
					Interval: IntervalMonthly,
					Currency: CurrencyUSD,
				},
				UpdateExistingSubscriptions: true,
			},
			mockResponse:   UpdateTierResponse{},
			mockStatusCode: http.StatusNotFound,
			expectedError:  errors.New("failed to update tier"),
			expectedCode:   http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := mockServer(tt.mockStatusCode, tt.mockResponse)
			defer server.Close()

			// Override the global URL with the mock server URL
			url = server.URL

			response, err, code := UpdateTier(tt.planCode, tt.request, context.Background())
			if (err != nil && tt.expectedError == nil) || (err == nil && tt.expectedError != nil) {
				t.Errorf("Expected error: %v, got: %v", tt.expectedError, err)
			}
			if code != tt.expectedCode {
				t.Errorf("Expected status code: %d, got: %d", tt.expectedCode, code)
			}
			if response != nil && response.Message != tt.mockResponse.Message {
				t.Errorf("Expected message: %s, got: %s", tt.mockResponse.Message, response.Message)
			}
		})
	}
}
