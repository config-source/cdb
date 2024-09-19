package webhooks_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/config-source/cdb/pkg/configkeys"
	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/config-source/cdb/pkg/webhooks"
	"github.com/rs/zerolog"
)

func TestRunWebhook(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	called := false
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			called = true
		},
	))
	defer server.Close()

	webhook := webhooks.Definition{
		ID:       1,
		Template: "Test",
		URL:      server.URL,
	}

	env := environments.Environment{}
	configValues := []configvalues.ConfigValue{}

	err := webhooks.RunWebhook(context.Background(), logger, server.Client(), webhook, env, configValues)
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Error("Expected the server to be called")
	}
}

func TestRunWebhookRendersTemplate(t *testing.T) {
	logger := zerolog.New(nil).Level(zerolog.Disabled)
	expectedBody := `
	{
		"enviromnent": {{ .Environment | toJSON }},
		"configValues": {
			{{ range .Configuration }}
			"{{ .Name }}": {{ .Value }},
			{{ end }}
		}
	}
	`
	receivedBody := ""
	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			reqBody, err := io.ReadAll(r.Body)
			defer r.Body.Close()

			if err != nil {
				receivedBody = err.Error()
			} else {
				receivedBody = string(reqBody)
			}
		},
	))
	defer server.Close()

	webhook := webhooks.Definition{
		ID: 1,
		Template: `
		{
			"enviromnent": {{ .Environment | toJSON }},
			"configValues": {
				{{ range .Configuration }}
				"{{ .Name }}": {{ .Value }},
				{{ end }}
			}
		}`,
		URL: server.URL,
	}

	env := environments.Environment{
		ID:        1,
		Name:      "test",
		Service:   "testSvc",
		ServiceID: 100,
	}

	val1 := 75
	val2 := "seventy-five"
	configValues := []configvalues.ConfigValue{
		{
			ID:            5,
			ConfigKeyID:   10,
			EnvironmentID: 1,
			Name:          "testIntValue",
			ValueType:     configkeys.TypeInteger,
			IntValue:      &val1,
		},
		{
			ID:            6,
			ConfigKeyID:   11,
			EnvironmentID: 1,
			Name:          "testStringValue",
			ValueType:     configkeys.TypeString,
			StrValue:      &val2,
		},
	}

	err := webhooks.RunWebhook(context.Background(), logger, server.Client(), webhook, env, configValues)
	if err != nil {
		t.Fatal(err)
	}

	if receivedBody != expectedBody {
		t.Errorf("Expected %s, got %s", expectedBody, receivedBody)
	}
}
