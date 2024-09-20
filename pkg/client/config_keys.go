package client

import (
	"context"
	"fmt"

	"github.com/config-source/cdb/pkg/configkeys"
)

var baseConfigKeyURL = "/api/v1/config-keys"

func (c *Client) GetConfigKey(ctx context.Context, id int) (configkeys.ConfigKey, error) {
	var data configkeys.ConfigKey

	_, err := c.Do(ctx, requestSpec{
		method: "GET",
		url:    fmt.Sprintf("%s/by-id/%d", baseConfigKeyURL, id),
	}, &data)

	return data, err
}

func (c *Client) GetConfigKeyByName(ctx context.Context, serviceName, name string) (configkeys.ConfigKey, error) {
	var data configkeys.ConfigKey

	_, err := c.Do(ctx, requestSpec{
		method: "GET",
		url:    fmt.Sprintf("%s/%s/by-name/%s", baseConfigKeyURL, serviceName, name),
	}, &data)

	return data, err
}
