package panelAdmins

import (
	"control-panel-bk/util"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func HandleCreateRole(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body CRole
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Println("Error 1", err)
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	output, outputErr := CreateRole(body, r.Context())
	if outputErr != nil {
		log.Println("Error 2", outputErr)
		util.ErrorException(w, outputErr, http.StatusInternalServerError)
		return
	}

	respBytes, e := json.Marshal(output)
	if e != nil {
		log.Println("Error 3", e)
		util.ErrorException(w, e, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

}

func HandleFetchRole(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	roleId := chi.URLParam(r, "id")
	err, cde := FetchRole(roleId, r.Context())

	if err != nil {
		util.ErrorException(w, err, cde)
		return
	}
	c := "done"
	w.Write([]byte(c))

}
