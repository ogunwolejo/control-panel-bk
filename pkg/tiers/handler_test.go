package tiers

import (
	"bytes"
	"context"
	"control-panel-bk/utils"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var mockCreateTierRequest = CreateTierRequest{
	Name:     "Test Tier",
	Amount:   1000,
	Interval: IntervalMonthly,
	Currency: CurrencyNGN,
}

// Mock the CreateTier function for testing
var mockCreateTier = func(w http.ResponseWriter, r *http.Request) {
	resp, err, code := func() (*TierResponse, error, int) {
		return &TierResponse{
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
				Name:     mockCreateTierRequest.Name,
				Amount:   int(mockCreateTierRequest.Amount),
				Interval: string(mockCreateTierRequest.Interval),
				Currency: string(mockCreateTierRequest.Currency),
			},
		}, nil, http.StatusCreated
	}()

	respBytes, _ := json.Marshal(resp)

	if err != nil {
		utils.ErrorException(w, err, code)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(respBytes)
}

func TestHandleTierCreation(t *testing.T) {
	requestBody := mockCreateTierRequest

	body, _ := json.Marshal(requestBody)
	// Create a request
	req := httptest.NewRequest(http.MethodPost, "/tiers", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockCreateTier(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, resp.StatusCode)
	}

	var response TierResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Message != "Tier created successfully" {
		t.Errorf("Expected message 'Tier created successfully', got '%s'", response.Message)
	}
}

// Mock the FetchTiers function for testing
var mockFetchTiers = func(w http.ResponseWriter, r *http.Request) {
	resp, err, cde := func() (*FetchTiersResponse, error, int) {
		return &FetchTiersResponse{
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
		}, nil, http.StatusOK
	}()

	if err != nil {
		utils.ErrorException(w, err, cde)
		return
	}

	respBytes, _ := json.Marshal(resp)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(cde)
	w.Write(respBytes)
}

func TestHandleFetchTiers(t *testing.T) {
	// Create a sample request body
	requestBody := FetchTiersRequest{
		PerPage: 10,
		Page:    1,
	}
	body, _ := json.Marshal(requestBody)

	// Create a request
	req := httptest.NewRequest(http.MethodGet, "/tiers/all", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mockFetchTiers(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response FetchTiersResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Message != "Tiers retrieved successfully" {
		t.Errorf("Expected message 'Tiers retrieved successfully', got '%s'", response.Message)
	}
}

// Mock the GetTier function for testing
var mockGetTier = func(planCode string, ctx context.Context) (*FetchTierResponse, error, int) {
	return &FetchTierResponse{
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
	}, nil, http.StatusOK
}

func TestHandleFetchTier(t *testing.T) {
	// Create a request with a URL parameter
	r := chi.NewRouter()
	r.Get("/tiers/{id}", func(writer http.ResponseWriter, request *http.Request) {
		phoneCode := chi.URLParam(request, "id")

		response, err, code := mockGetTier(phoneCode, request.Context())
		if err != nil {
			utils.ErrorException(writer, err, code)
			return
		}

		responseByte, _ := json.Marshal(response)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(code)
		writer.Write(responseByte)
	})

	req := httptest.NewRequest(http.MethodGet, "/tiers/PLAN123", nil)
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response FetchTierResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Message != "Tier retrieved successfully" {
		t.Errorf("Expected message 'Tier retrieved successfully', got '%s'", response.Message)
	}
}

// Mock the UpdateTier function for testing
var mockUpdateTier = func(planCode string, updateOption UpdateTierRequest, ctx context.Context) (*UpdateTierResponse, error, int) {
	return &UpdateTierResponse{
		Status:  true,
		Message: "Tier updated successfully",
	}, nil, http.StatusOK
}

func TestHandleUpdateTier(t *testing.T) {
	// Create a sample request body
	requestBody := UpdateTierRequest{
		CreateTierRequest: CreateTierRequest{
			Name:     "Updated Tier",
			Amount:   2000,
			Interval: IntervalMonthly,
			Currency: CurrencyUSD,
		},
		UpdateExistingSubscriptions: true,
	}
	body, _ := json.Marshal(requestBody)

	// Create a request with a URL parameter
	r := chi.NewRouter()
	r.Put("/tiers/{id}", func(writer http.ResponseWriter, request *http.Request) {
		planCode := chi.URLParam(request, "id")
		resp, err, code := mockUpdateTier(planCode, requestBody, request.Context())

		if err != nil {
			utils.ErrorException(writer, err, code)
			return
		}

		respoBytes, _ := json.Marshal(resp)

		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(code)
		writer.Write(respoBytes)
	})

	req := httptest.NewRequest(http.MethodPut, "/tiers/PLAN123", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Serve the request
	r.ServeHTTP(w, req)

	// Check the response
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var response UpdateTierResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.Message != "Tier updated successfully" {
		t.Errorf("Expected message 'Tier updated successfully', got '%s'", response.Message)
	}
}
