package webhooks

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"text/template"

	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
	"github.com/rs/zerolog"
)

type Service struct {
	repo DefinitionRepository
}

type TemplateContext struct {
	Environment   environments.Environment
	Configuration []configvalues.ConfigValue
}

func (s *Service) CreateWebhookDefinition(ctx context.Context, webhookDefinition Definition, envIDs ...int) (Definition, error) {
	wh, err := s.repo.CreateWebhookDefinition(ctx, webhookDefinition)
	if err != nil {
		return wh, err
	}

	if envIDs != nil {
		err = s.repo.AssociateWebhookToEnvironments(ctx, wh.ID, envIDs)
	}

	return wh, err
}

func RunWebhook(
	ctx context.Context,
	logger zerolog.Logger,
	client *http.Client,
	webhook Definition,
	env environments.Environment,
	configValues []configvalues.ConfigValue,
) error {
	renderContext := TemplateContext{
		Environment:   env,
		Configuration: configValues,
	}

	template, err := template.New("webhookTemplate").Parse(webhook.Template)
	if err != nil {
		return err
	}

	renderedContent := bytes.NewBuffer([]byte{})
	if err := template.Execute(renderedContent, renderContext); err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, webhook.URL, renderedContent)
	if err != nil {
		return err
	}

	if webhook.AuthzHeader != "" {
		req.Header.Set("Authorization", webhook.AuthzHeader)
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if logger.GetLevel() == zerolog.DebugLevel {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		logger.Debug().
			Bytes("webhookResponseBody", bodyBytes).
			Str("webhookUrl", webhook.URL).
			Int("webhookResponseStatusCode", resp.StatusCode).
			Int("webhookId", webhook.ID).
			Msg("webhook response")
	}

	return nil
}
