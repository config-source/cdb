package webhooks

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/config-source/cdb/pkg/configvalues"
	"github.com/config-source/cdb/pkg/environments"
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

func RunWebhook(ctx context.Context, webhook Definition, env environments.Environment, configValues []configvalues.ConfigValue) error {
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

	fmt.Println(renderedContent.String())

	return nil
}
