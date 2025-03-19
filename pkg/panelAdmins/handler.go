package panelAdmins

import (
	"control-panel-bk/util"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"strconv"
)

func HandleCreateRole(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body CRole
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	output, outputErr := CreateRole(body, r.Context())
	if outputErr != nil {
		if errors.Is(outputErr, errors.New("a role having the same name already exists")) {
			util.ErrorException(w, outputErr, http.StatusOK)
			return
		}

		util.ErrorException(w, outputErr, http.StatusInternalServerError)
		return
	}

	response := Response{
		Status: http.StatusCreated,
		error:  nil,
		Data:   output,
	}

	respBytes, e := json.Marshal(&response)
	if e != nil {
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

func HandleFetchRoleByName(w http.ResponseWriter, r *http.Request) {
	roleName := r.URL.Query().Get("name")
	roles, err := FetchRoleByName(roleName, r.Context())

	if err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	response := Response{
		Status: http.StatusOK,
		error:  nil,
		Data:   roles,
	}

	respBytes, e := json.Marshal(&response)
	if e != nil {
		util.ErrorException(w, e, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respBytes); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

}

func HandleFetchRoleById(w http.ResponseWriter, r *http.Request) {
	roleId := chi.URLParam(r, "id")

	result, err, cde := FetchRoleById(roleId, r.Context())

	if err != nil {
		util.ErrorException(w, err, cde)
		return
	}

	log.Println("RESULT GOTTEN FROM FETCHING ROLE BY ID: ", result)

	resp, respErr := json.Marshal(&result)
	if respErr != nil {
		util.ErrorException(w, respErr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
	}
}

func HandleFetchRoles(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// Get all the query params i.e PAGE, LIMIT,
	page := query.Get("page")
	limit := query.Get("limit")

	var pgQuery, ltQuery int

	log.Printf("QUERY: %s - %s", page, limit)

	pg, pgErr := strconv.Atoi(page)
	if pgErr != nil {
		util.ErrorException(w, pgErr, http.StatusInternalServerError)
		return
	}

	lt, ltErr := strconv.Atoi(limit)
	if ltErr != nil {
		util.ErrorException(w, ltErr, http.StatusInternalServerError)
		return
	}

	pgQuery = pg
	ltQuery = lt

	result, err, code := FetchRoles(pgQuery, ltQuery, r.Context())

	if err != nil {
		util.ErrorException(w, err, code)
		return
	}

	response := Response{
		Status: code,
		error:  nil,
		Data:   result,
	}

	if respBytes, err := json.Marshal(&response); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(respBytes); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
		}
	}

}

func HandleHardDeleteOfRole(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var role Role

	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	if id, err, code := role.HardDeleteRole(r.Context()); err != nil {
		util.ErrorException(w, err, code)
		return
	} else {
		deleteRes := Response{
			Status: code,
			error: nil,
			Data: id,
		}

		deleteBytes, delErr := json.Marshal(&deleteRes)
		if delErr != nil {
			util.ErrorException(w, delErr, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		if _, err := w.Write(deleteBytes); err != nil {
			util.ErrorException(w, err, http.StatusInternalServerError)
		}
	}
}

func HandleGeneralUpdate(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var body Role

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	updateDoc, updateError, code := body.GeneralizedUpdate(r.Context())
	if updateError != nil {
		util.ErrorException(w, updateError, code)
		return
	}

	resp := Response{
		Status: http.StatusAccepted,
		error: nil,
		Data: updateDoc,
	}

	respBytes, respErr := json.Marshal(&resp)
	if respErr != nil {
		util.ErrorException(w, respErr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	if _, err := w.Write(respBytes); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
	}
}

func HandleArchiveRole(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var role Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	doc, docErr, code := role.ArchiveRole(r.Context())
	if docErr != nil {
		if errors.Is(docErr, errors.New("no document was found")) {
			util.ErrorException(w, docErr, code)
			return
		}

		util.ErrorException(w, docErr, code)
		return
	}

	archBytes, archErr := json.Marshal(&doc)
	if archErr != nil {
		util.ErrorException(w, archErr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(archBytes); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
	}

}

func HandleUnArchiveRole(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var role Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	doc, docErr, code := role.UnArchiveRole(r.Context())
	if docErr != nil {
		if errors.Is(docErr, errors.New("no document was found")) {
			util.ErrorException(w, docErr, code)
			return
		}

		util.ErrorException(w, docErr, code)
		return
	}

	response := Response{
		Status: code,
		error:  nil,
		Data:   doc,
	}

	unArchBytes, unArchErr := json.Marshal(&response)
	if unArchErr != nil {
		util.ErrorException(w, unArchErr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(unArchBytes); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
	}
}

func HandlePushRoleToBin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var role Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	bin, binErr, code := role.DeleteRole(r.Context())
	if binErr != nil {
		if errors.Is(binErr, errors.New("no document was found")) {
			util.ErrorException(w, binErr, code)
			return
		}

		util.ErrorException(w, binErr, code)
		return
	}

	response := Response{
		Status: code,
		error:  nil,
		Data:   bin,
	}

	binByte, bbErr := json.Marshal(&response)
	if bbErr != nil {
		util.ErrorException(w, bbErr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if _, err := w.Write(binByte); err != nil {
		util.ErrorException(w, bbErr, http.StatusInternalServerError)
	}
}
