package client

import (
	"context"

	"github.com/config-source/cdb/auth"
)

func (ec *Client) Login(ctx context.Context, email, password string) (auth.TokenSet, error) {
	var data auth.TokenSet
	_, err := ec.Do(ctx, requestSpec{
		url:    "/api/v1/auth/login",
		method: "POST",
		body: struct {
			Email    string
			Password string
		}{
			Email:    email,
			Password: password,
		},
	}, &data)
	ec.token = data.IDToken
	return data, err
}

func (ec *Client) IssueAPIToken(ctx context.Context, email, password string) (auth.APIToken, error) {
	var data auth.APIToken

	_, err := ec.Do(ctx, requestSpec{
		method: "POST",
		url:    "/api/v1/auth/api-tokens",
	}, &data)

	return data, err
}
