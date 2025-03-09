package tiers

import (
	"bytes"
	"context"
	cfg "control-panel-bk/config"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-querystring/query"
	"io"
	"net/http"
	burl "net/url"
	"time"
)

type Interval string
type Currency string

const (
	IntervalDaily      Interval = "daily"
	IntervalWeekly     Interval = "weekly"
	IntervalMonthly    Interval = "monthly"
	IntervalAnnually   Interval = "annually"
	IntervalBiannually Interval = "biannually"
	IntervalQuarterly  Interval = "quarterly"
)

const (
	CurrencyUSD Currency = "USD"
	CurrencyNGN Currency = "NGN"
	CurrencyGHS Currency = "GHS"
	CurrencyZAR Currency = "ZAR"
)

// TierResponse Create Tier Response
type TierResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
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
	} `json:"data"`
}

type FetchTiersResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    []struct {
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
	} `json:"data"`
	Meta struct {
		Total     int `json:"total"`
		Skipped   int `json:"skipped"`
		PerPage   int `json:"perPage"`
		Page      int `json:"page"`
		PageCount int `json:"pageCount"`
	} `json:"meta"`
}

type FetchTierResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
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
	} `json:"data"`
	Meta struct {
		Total     int `json:"total"`
		Skipped   int `json:"skipped"`
		PerPage   int `json:"perPage"`
		Page      int `json:"page"`
		PageCount int `json:"pageCount"`
	} `json:"meta"`
}

type UpdateTierResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

type CreateTierRequest struct {
	Name         string   `json:"name"`
	Amount       int64    `json:"amount"`
	Interval     Interval `json:"interval"`
	Description  string   `json:"description,omitempty"`
	SendInvoices bool     `json:"send_invoices,omitempty"`
	SendSMS      bool     `json:"send_sms,omitempty"`
	Currency     Currency `json:"currency,omitempty"`
	InvoiceLimit int      `json:"invoice_limit,omitempty"`
}

type FetchTiersRequest struct {
	PerPage  int      `json:"perPage"`
	Page     int      `json:"page"`
	Status   string   `json:"status,omitempty"`
	Interval Interval `json:"interval,omitempty"`
	Amount   int64    `json:"amount,omitempty"`
}

type UpdateTierRequest struct {
	CreateTierRequest,
	UpdateExistingSubscriptions bool `json:"update_existing_subscriptions,omitempty"`
}

type APIError struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Meta    map[string]interface{} `json:"meta"`
	Status  bool                   `json:"status"`
	Type    string                 `json:"type"`
}

var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	},
}

var payStackConfig = cfg.DefaultPayStackConfiguration()
var url = payStackConfig.PlanUrl()

// CreateTier creates a new tier on the PayStack API
func CreateTier(tier CreateTierRequest, ctx context.Context) (*TierResponse, error, int) {
	body, err := json.Marshal(tier)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	// We need to check if the Tiers exist already
	fetchTiers := make(chan *FetchTiersResponse)
	fetchErr := make(chan error)
	fetchErrCode := make(chan int)
	go func() {
		par := FetchTiersRequest{
			PerPage:  0,
			Amount:   tier.Amount,
			Interval: tier.Interval,
			Page:     1,
			Status:   "active",
		}

		t, e, scd := FetchTiers(par, ctx)
		fetchTiers <- t
		fetchErr <- e
		fetchErrCode <- scd

		close(fetchTiers)
		close(fetchErr)
		close(fetchErrCode)
	}()

	select {
	case er := <-fetchErr:
		if er != nil {
			return nil, er, <-fetchErrCode
		}

	case list := <-fetchTiers:
		if len(list.Data) > 0 {
			return nil, errors.New("an active tier with the same data already exists"), http.StatusOK
		}
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", cfg.PayStackConfig.Headers.Authorization)
	req.Header.Set("Content-Type", cfg.PayStackConfig.Headers.ContentType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	defer resp.Body.Close()

	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if resp.StatusCode != http.StatusCreated {
		var errResponse APIError
		if err := json.Unmarshal(responseBytes, &errResponse); err != nil {
			return nil, err, http.StatusInternalServerError
		}

		return nil, errors.New(errResponse.Message), resp.StatusCode
	}

	var createdTier TierResponse
	if err := json.Unmarshal(responseBytes, &createdTier); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return &createdTier, nil, http.StatusCreated
}

// GetTier retrieves a tier from the PayStack API by the plan code
func GetTier(planCode string, ctx context.Context) (*FetchTierResponse, error, int) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", url, planCode), nil)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", cfg.PayStackConfig.Headers.Authorization)
	req.Header.Set("Content-Type", cfg.PayStackConfig.Headers.ContentType)

	resp, respErr := client.Do(req)
	defer resp.Body.Close()

	if respErr != nil {
		return nil, err, http.StatusInternalServerError
	}

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(respBytes, &apiErr); err != nil {
			return nil, errors.New(apiErr.Message), resp.StatusCode
		}
		return nil, fmt.Errorf("failed to get tier: %d", resp.StatusCode), resp.StatusCode
	}

	var tier FetchTierResponse
	if err := json.Unmarshal(respBytes, &tier); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return &tier, nil, http.StatusOK
}

// FetchTiers retrieves all tiers from the PayStack API
func FetchTiers(arg FetchTiersRequest, ctx context.Context) (*FetchTiersResponse, error, int) {
	v, queryErr := query.Values(arg)
	if queryErr != nil {
		return nil, queryErr, http.StatusInternalServerError
	}

	baseUrl, urlErr := burl.Parse(url)
	if urlErr != nil {
		return nil, urlErr, http.StatusInternalServerError
	}

	baseUrl.RawQuery = v.Encode()
	baseUrlStr := baseUrl.String()

	req, e := http.NewRequestWithContext(ctx, "GET", baseUrlStr, nil)
	if e != nil {
		return nil, e, http.StatusInternalServerError
	}

	req.Header.Set("Authorization", cfg.PayStackConfig.Headers.Authorization)
	req.Header.Set("Content-Type", cfg.PayStackConfig.Headers.ContentType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	defer resp.Body.Close()

	respBody, respBodyErr := io.ReadAll(resp.Body)
	if respBodyErr != nil {
		return nil, respBodyErr, http.StatusInternalServerError
	}

	if resp.StatusCode != http.StatusOK {
		var errorResponse APIError
		if e := json.Unmarshal(respBody, &errorResponse); e != nil {
			return nil, errors.New(errorResponse.Message), http.StatusInternalServerError
		}
		return nil, fmt.Errorf("failed to get tiers: %d", resp.StatusCode), resp.StatusCode
	}

	var fetchTiers FetchTiersResponse
	if err := json.Unmarshal(respBody, &fetchTiers); err != nil {
		return nil, err, http.StatusInternalServerError
	}
	return &fetchTiers, nil, http.StatusOK
}

// UpdateTier updates a tier via the PayStack API
func UpdateTier(planCode string, updateOption UpdateTierRequest, ctx context.Context) (*UpdateTierResponse, error, int) {
	body, err := json.Marshal(updateOption)
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/%s", url, planCode), bytes.NewReader(body))
	if err != nil {
		return nil, err, http.StatusInternalServerError
	}
	req.Header.Set("Authorization", cfg.PayStackConfig.Headers.Authorization)
	req.Header.Set("Content-Type", cfg.PayStackConfig.Headers.ContentType)

	resp, respErr := client.Do(req)
	if respErr != nil {
		return nil, respErr, resp.StatusCode
	}

	defer resp.Body.Close()

	respBytes, e := io.ReadAll(resp.Body)
	if e != nil {
		return nil, e, http.StatusInternalServerError
	}

	if resp.StatusCode != http.StatusOK {
		var apiErr APIError
		if err := json.Unmarshal(respBytes, &apiErr); err != nil {
			return nil, errors.New(apiErr.Message), resp.StatusCode
		}

		return nil, fmt.Errorf("failed to update tier"), resp.StatusCode
	}

	var updatedTier UpdateTierResponse
	if err := json.Unmarshal(respBytes, &updatedTier); err != nil {
		return nil, err, http.StatusInternalServerError
	}

	return &updatedTier, nil, http.StatusOK
}
