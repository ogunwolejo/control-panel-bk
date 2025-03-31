package pkg

import (
	"control-panel-bk/config"
	"control-panel-bk/internal/aws"
	"control-panel-bk/util"
	"net/http"
	"os"
)

var (
	ClientID = os.Getenv("AWS_CLIENT_ID")
)

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
		//Expires: time.Now().Add()
	}

	util.SetHttpOnlyCookie(w, rtCookie)

	resp := util.Response{
		Status: http.StatusOK,
		Data: map[string]string{
			"AccessToken": *output.AuthenticationResult.AccessToken,
			"IdToken":     *output.AuthenticationResult.IdToken,
		},
	}

	if respBytes, respErr := util.GetBytesResponse(http.StatusOK, resp); respErr != nil {
		util.ErrorException(w, respErr, http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBytes)
	}
}

func AccessTokenAuth(w http.ResponseWriter, r *http.Request) {}
