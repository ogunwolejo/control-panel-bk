package tiers

import (
	"control-panel-bk/utils"
	"encoding/json"
	"log"
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
		utils.ErrorException(w, err, statusCode)
		return
	}

	reads, e := json.Marshal(tier)
	if e != nil {
		utils.ErrorException(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Encoding", "application/json")
	w.WriteHeader(statusCode)
	_, writeErr := w.Write(reads)
	if writeErr != nil {
		utils.ErrorException(w, writeErr, http.StatusInternalServerError)
		return
	}

}

func HandleFetchTiers(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var ftr FetchTiersRequest

	err := json.NewDecoder(r.Body).Decode(&ftr)
	if err != nil {
		utils.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	tiers, err, statusCode := FetchTiers(ftr)
	if err != nil {
		utils.ErrorException(w, err, statusCode)
	}

	w.Header().Set("Content-Encoding", "application/json")
	w.WriteHeader(statusCode)
	e := json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"data":   *tiers,
	})

	if e != nil {
		utils.ErrorException(w, err, http.StatusInternalServerError)
	}
}

func HandleFetchTier(w http.ResponseWriter, r *http.Request) {
	log.Println("r: ", r.URL.Fragment)
	log.Println(r)
}
