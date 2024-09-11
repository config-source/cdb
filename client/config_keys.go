package client

import (
	"context"
	"fmt"
	"strconv"

	"github.com/config-source/cdb/configkeys"
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

func (c *Client) GetConfigKeyByName(ctx context.Context, name string) (configkeys.ConfigKey, error) {
	var data configkeys.ConfigKey

	_, err := c.Do(ctx, requestSpec{
		method: "GET",
		url:    fmt.Sprintf("%s/by-name/%s", baseConfigKeyURL, name),
	}, &data)

	return data, err
}

func (c *Client) GetConfigKeyByNameOrID(nameOrID string) (configkeys.ConfigKey, error) {
	id, err := strconv.Atoi(nameOrID)
	if err == nil {
		return c.GetConfigKey(context.Background(), id)
	}

	return c.GetConfigKeyByName(context.Background(), nameOrID)
}
