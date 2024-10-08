package client

import (
	"context"
	"fmt"

	"github.com/config-source/cdb/pkg/environments"
)

var baseEnvURL = "/api/v1/environments"

func (ec *Client) GetEnvironmentByName(ctx context.Context, serviceName, name string) (environments.Environment, error) {
	var data environments.Environment

	_, err := ec.Do(ctx, requestSpec{
		method: "GET",
		url:    fmt.Sprintf("%s/%s/by-name/%s", baseEnvURL, serviceName, name),
	}, &data)

	return data, err
}

func (ec *Client) GetEnvironment(ctx context.Context, id int) (environments.Environment, error) {
	var data environments.Environment

	_, err := ec.Do(ctx, requestSpec{
		method: "GET",
		url:    fmt.Sprintf("%s/by-id/%d", baseEnvURL, id),
	}, &data)

	return data, err
}

func (ec *Client) CreateEnvironment(ctx context.Context, env environments.Environment) (environments.Environment, error) {
	var data environments.Environment

	_, err := ec.Do(ctx, requestSpec{
		method: "POST",
		url:    baseEnvURL,
		body:   env,
	}, &data)

	return data, err
}

func (ec *Client) ListEnvironments(ctx context.Context) ([]environments.Environment, error) {
	var data []environments.Environment

	_, err := ec.Do(ctx, requestSpec{
		method: "GET",
		url:    baseEnvURL,
	}, &data)

	return data, err
}

func (ec *Client) GetEnvironmentTree(ctx context.Context) ([]environments.Tree, error) {
	var data []environments.Tree

	_, err := ec.Do(ctx, requestSpec{
		method: "GET",
		url:    fmt.Sprintf("%s/tree", baseEnvURL),
	}, &data)

	return data, err
}
