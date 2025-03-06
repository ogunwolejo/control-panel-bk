package tiers

import (
	"bytes"
	"context"
	cfg "control-panel-bk/config"
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"log"
	"net/http"
	burl "net/url"
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

var psConfig cfg.PayStack = cfg.PayStackConfig

type TierData struct {
	Name         string   `json:"name"`
	Amount       int      `json:"amount"`
	Interval     Interval `json:"interval"`
	Integration  int      `json:"integration"`
	Domain       string   `json:"domain"`
	PlanCode     string   `json:"plan_code"`
	SendInvoices bool     `json:"send_invoices"`
	SendSMS      bool     `json:"send_sms"`
	HostedPage   bool     `json:"hosted_page"`
	Currency     Currency `json:"currency"`
	ID           int      `json:"id"`
	CreatedAt    string   `json:"createdAt"`
	UpdatedAt    string   `json:"updatedAt"`
}

type TierResponse struct {
	Status  bool     `json:"status"`
	message string   `json:"message"`
	data    TierData `json:"data"`
}

type CreateTierRequest struct {
	Name         string   `json:"name"`
	Amount       int64    `json:"amount"`
	Interval     Interval `json:"interval"`
	Description  string   `json:"description,omitempty"`
	SendInvoices bool     `json:"send_invoices,omitempty"`
	SendSMS      bool     `json:"send_sms,omitempty"`
	Currency     Currency `json:"currency, omitempty"`
	InvoiceLimit int      `json:"invoice_limit,omitempty"`
}

type FetchTiersRequest struct {
	PerPage  int      `json:"perPage"`
	Page     int      `json:"page"`
	Status   int      `json:"status,omitempty"`
	Interval Interval `json:"interval,omitempty"`
	Amount   int64    `json:"amount,omitempty"`
}

type UpdateTierRequest struct {
	CreateTierRequest,
	UpdateExistingSubscriptions bool `json:"update_existing_subscriptions,omitempty"`
}

var client = &http.Client{}
var url = fmt.Sprintf("%s:%d/plan", psConfig.Host, psConfig.Port)

// CreateTier creates a new tier on the PayStack API
func CreateTier(tier CreateTierRequest) (*TierResponse, error) {
	ctx := context.Background()
	body, err := json.Marshal(tier)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", psConfig.Headers.Authorization)
	req.Header.Set("Content-Type", psConfig.Headers.ContentType)
	req.Header.Set("Cache-Control", "no-cache")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to create tier: %d", resp.StatusCode)
	}

	var createdTier TierResponse
	if err := json.NewDecoder(resp.Body).Decode(&createdTier); err != nil {
		return nil, err
	}

	return &createdTier, nil
}

// GetTier retrieves a tier from the PayStack API by the plan code
func GetTier(planCode string) (*TierResponse, error) {
	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/%s", url, planCode), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", psConfig.Headers.Authorization)
	req.Header.Set("Content-Type", psConfig.Headers.ContentType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get tier: %d", resp.StatusCode)
	}

	var tier TierResponse
	if err := json.NewDecoder(resp.Body).Decode(&tier); err != nil {
		return nil, err
	}

	return &tier, nil
}

// FetchTiers retrieves all tiers from the PayStack API
func FetchTiers(arg FetchTiersRequest) (*TierResponse, error) {
	ctx := context.Background()
	body, err := json.Marshal(arg)
	if err != nil {
		return nil, err
	}

	v, err := query.Values(body)
	if err != nil {
		log.Fatal("Query values error: ", err)
		return nil, err
	}

	baseUrl, err := burl.Parse(url)
	if err != nil {
		return nil, err
	}

	baseUrl.RawQuery = v.Encode()
	baseUrlStr := baseUrl.String()

	req, err := http.NewRequestWithContext(ctx, "GET", baseUrlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", psConfig.Headers.Authorization)
	req.Header.Set("Content-Type", psConfig.Headers.ContentType)
	req.Header.Set("Cache-Control", "cache")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get tiers: %d", resp.StatusCode)
	}

	var fetchTiers TierResponse
	if err := json.NewDecoder(resp.Body).Decode(&fetchTiers); err != nil {
		return nil, err
	}

	return &fetchTiers, nil
}



// UpdateTier updates a tier via the PayStack API
func UpdateTier(planCode string, updateOption UpdateTierRequest) (*TierResponse, error) {
	ctx := context.Background()
	body, err := json.Marshal(updateOption)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/%s", url, planCode), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", psConfig.Headers.Authorization)
	req.Header.Set("Content-Type", psConfig.Headers.ContentType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to update tier: %d", resp.StatusCode)
	}

	var updatedTier TierResponse
	if err := json.NewDecoder(resp.Body).Decode(&updatedTier); err != nil {
		return nil, err
	}

	return &updatedTier, nil
}
