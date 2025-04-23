package amadeusapi

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/danielmesquitta/flight-api/internal/domain/errs"
)

const authHeaderKey = "Authorization"
const authExpiresAtKey = "Expires-At"

func (a *AmadeusAPI) refreshAccessToken(ctx context.Context) error {
	token := a.c.Header().Get(authHeaderKey)
	if token == "" {
		return a.authenticate(ctx)
	}

	expStr := a.c.Header().Get(authExpiresAtKey)
	if expStr == "" {
		return a.authenticate(ctx)
	}

	secs, err := strconv.ParseInt(expStr, 10, 64)
	if err != nil {
		return a.authenticate(ctx)
	}
	expiresAt := time.Unix(secs, 0)

	if time.Now().After(expiresAt) {
		return a.authenticate(ctx)
	}

	return nil
}

func (a *AmadeusAPI) authenticate(ctx context.Context) error {
	res, err := a.c.R().
		SetContext(ctx).
		SetFormData(map[string]string{
			"grant_type":    "client_credentials",
			"client_id":     a.e.AmadeusAPIKey,
			"client_secret": a.e.AmadeusAPISecret,
		}).
		Post("/v1/security/oauth2/token")
	if err != nil {
		return errs.New(err)
	}
	body := res.Bytes()
	if res.IsError() {
		return errs.New(string(body))
	}

	type AuthResponse struct {
		AccessToken string `json:"access_token" validate:"required"`
		ExpiresIn   int64  `json:"expires_in"   validate:"required,min=1"`
	}
	data := AuthResponse{}
	if err := json.Unmarshal(body, &data); err != nil {
		return errs.New(err)
	}

	if data.AccessToken == "" {
		return errs.New("access token is empty")
	}

	a.c.SetHeader(authHeaderKey, "Bearer "+data.AccessToken)

	expiresAt := time.Now().Add(time.Duration(data.ExpiresIn) * time.Second)
	a.c.SetHeader(authExpiresAtKey, fmt.Sprintf("%d", expiresAt.Unix()))

	return nil
}
