package tiers

import (
	"control-panel-bk/util"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func HandleTierCreation(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var ctr CreateTierRequest

	err := json.NewDecoder(r.Body).Decode(&ctr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctr.Amount = ctr.Amount * 100 // From the documentation whatever price is charge it must be by 100

	tier, err, statusCode := CreateTier(ctr, r.Context())
	if err != nil {
		util.ErrorException(w, err, statusCode)
		return
	}

	reads, e := json.Marshal(tier)
	if e != nil {
		util.ErrorException(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Encoding", "application/json")
	w.WriteHeader(statusCode)
	_, writeErr := w.Write(reads)
	if writeErr != nil {
		util.ErrorException(w, writeErr, http.StatusInternalServerError)
		return
	}

}

func HandleFetchTiers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var ftr FetchTiersRequest

	err := json.NewDecoder(r.Body).Decode(&ftr)
	if err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	tiers, err, statusCode := FetchTiers(ftr, r.Context())
	if err != nil {
		util.ErrorException(w, err, statusCode)
		return
	}

	reads, e := json.Marshal(tiers)
	if e != nil {
		util.ErrorException(w, err, statusCode)
		return
	}

	w.Header().Set("Content-Encoding", "application/json")
	w.WriteHeader(statusCode)
	_, writeErr := w.Write(reads)
	if writeErr != nil {
		util.ErrorException(w, writeErr, http.StatusInternalServerError)
		return
	}
}

func HandleFetchTier(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	planCode := chi.URLParam(r, "id")

	resp, err, statsCode := GetTier(planCode, r.Context())
	if err != nil {
		util.ErrorException(w, err, statsCode)
		return
	}

	respBytes, e := json.Marshal(resp)
	if e != nil {
		util.ErrorException(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Encoding", "application/json")
	w.WriteHeader(statsCode)
	_, writeErr := w.Write(respBytes)
	if writeErr != nil {
		util.ErrorException(w, writeErr, http.StatusInternalServerError)
		return
	}
}

func HandleUpdateTier(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	planCode := chi.URLParam(r, "id")

	var updateBody UpdateTierRequest
	if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	updateBody.Amount = updateBody.Amount * 100 // We have to multiply the amount by 100 - default paystack rule

	updated, updateError, updateStatCde := UpdateTier(planCode, updateBody, r.Context())
	if updateError != nil {
		util.ErrorException(w, updateError, updateStatCde)
		return
	}

	updatedBytes, e := json.Marshal(updated)
	if e != nil {
		util.ErrorException(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(updateStatCde)
	if _, err := w.Write(updatedBytes); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

}
