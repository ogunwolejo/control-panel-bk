package pkg

import (
	"control-panel-bk/config"
	"control-panel-bk/internal/aws"
	"control-panel-bk/util"
	"encoding/json"
	"net/http"
	"os"
	"time"
)

var (
	ClientID = os.Getenv("AWS_CLIENT_ID")
)

type Credential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Username struct {
	Username string `json:"username"`
}

type ChangePassword struct {
	NewPassword string `json:"new_password"`
	OldPassword string `json:"old_password"`
}

type ForgetPasswordCred struct {
	Credential
	OtpCode string `json:"otp_code"`
}

func RefreshTokenAuth(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")

	if err != nil {
		util.ErrorException(w, err, http.StatusBadRequest)
		return
	}

	output, e := aws.AuthViaRefreshToken(config.AwsConfig, ClientID, cookie.Value)
	if e != nil {
		util.ErrorException(w, e, http.StatusNotImplemented)
		return
	}

	rt := output.AuthenticationResult.RefreshToken
	rtCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    *rt,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Second * time.Duration(output.AuthenticationResult.ExpiresIn)),
	}

	util.SetHttpOnlyCookie(w, rtCookie)

	data := map[string]string{
		"AccessToken": *output.AuthenticationResult.AccessToken,
		"IdToken":     *output.AuthenticationResult.IdToken,
	}

	if respBytes, respErr := util.GetBytesResponse(http.StatusOK, data); respErr != nil {
		util.ErrorException(w, respErr, http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var cred Credential
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	output, err := aws.LogInUser(config.AwsConfig, ClientID, cred.Username, cred.Password)
	if err != nil {
		util.ErrorException(w, err, http.StatusNotImplemented)
		return
	}

	rt := output.AuthenticationResult.RefreshToken
	rtCookie := &http.Cookie{
		Name:     "refresh_token",
		Value:    *rt,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(time.Second * time.Duration(output.AuthenticationResult.ExpiresIn)),
	}

	util.SetHttpOnlyCookie(w, rtCookie)

	data := map[string]string{
		"AccessToken": *output.AuthenticationResult.AccessToken,
		"IdToken":     *output.AuthenticationResult.IdToken,
	}

	respBytes, respErr := util.GetBytesResponse(http.StatusOK, data)

	if respErr != nil {
		util.ErrorException(w, respErr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("access_token").(string)

	if err := aws.LogOutUser(config.AwsConfig, token); err != nil {
		util.ErrorException(w, err, http.StatusNotImplemented)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("logout"))
}

func ChangePasswordHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	token := r.Context().Value("access_token").(string)

	var body ChangePassword
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	output, e := aws.ChangeUserPassword(config.AwsConfig, token, body.NewPassword, body.OldPassword)
	if e != nil {
		util.ErrorException(w, e, http.StatusNotImplemented)
		return
	}

	respBytes, respErr := util.GetBytesResponse(http.StatusOK, output.ResultMetadata)
	if respErr != nil {
		util.ErrorException(w, respErr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func ForgetPasswordOtpHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var cred Username
	if err := json.NewDecoder(r.Body).Decode(&cred); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	_, err := aws.ForgetPasswordOtp(config.AwsConfig, cred.Username)
	if err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	respBytes, respErr := util.GetBytesResponse(http.StatusOK, "confirmation code sent")
	if respErr != nil {
		util.ErrorException(w, respErr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}

func ForgetPasswordHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var fCred ForgetPasswordCred
	if err := json.NewDecoder(r.Body).Decode(&fCred); err != nil {
		util.ErrorException(w, err, http.StatusInternalServerError)
		return
	}

	_, err := aws.ForgetPassword(config.AwsConfig, fCred.Username, fCred.OtpCode, fCred.Password)

	if err != nil {
		util.ErrorException(w, err, http.StatusNotImplemented)
		return
	}

	respBytes, respErr := util.GetBytesResponse(http.StatusOK, "password has be changed")
	if respErr != nil {
		util.ErrorException(w, respErr, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
