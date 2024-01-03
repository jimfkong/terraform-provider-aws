// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package amp

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/amp"
	"github.com/aws/aws-sdk-go-v2/service/amp/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/prometheusservice"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-provider-aws/internal/errs"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
)

func FindScraperByID(ctx context.Context, conn *amp.Client, id string) (*types.ScraperDescription, error) {
	input := &amp.DescribeScraperInput{
		ScraperId: aws.String(id),
	}

	output, err := conn.DescribeScraper(ctx, input)

	if errs.IsA[*types.ResourceNotFoundException](err) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Scraper == nil || output.Scraper.Status == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Scraper, nil
}

func FindWorkspaceByID(ctx context.Context, conn *prometheusservice.PrometheusService, id string) (*prometheusservice.WorkspaceDescription, error) {
	input := &prometheusservice.DescribeWorkspaceInput{
		WorkspaceId: aws.String(id),
	}

	output, err := conn.DescribeWorkspaceWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, prometheusservice.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.Workspace == nil || output.Workspace.Status == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.Workspace, nil
}

func FindLoggingConfigurationByWorkspaceID(ctx context.Context, conn *prometheusservice.PrometheusService, id string) (*prometheusservice.LoggingConfigurationMetadata, error) {
	input := &prometheusservice.DescribeLoggingConfigurationInput{
		WorkspaceId: aws.String(id),
	}

	output, err := conn.DescribeLoggingConfigurationWithContext(ctx, input)

	if tfawserr.ErrCodeEquals(err, prometheusservice.ErrCodeResourceNotFoundException) {
		return nil, &retry.NotFoundError{
			LastError:   err,
			LastRequest: input,
		}
	}

	if err != nil {
		return nil, err
	}

	if output == nil || output.LoggingConfiguration == nil || output.LoggingConfiguration.Status == nil {
		return nil, tfresource.NewEmptyResultError(input)
	}

	return output.LoggingConfiguration, nil
}

func FindWorkspaces(ctx context.Context, conn *prometheusservice.PrometheusService, alias string) ([]*prometheusservice.WorkspaceSummary, error) { // nosemgrep:ci.caps0-in-func-name
	input := &prometheusservice.ListWorkspacesInput{}
	if alias != "" {
		input.Alias = aws.String(alias)
	}
	var output []*prometheusservice.WorkspaceSummary

	err := conn.ListWorkspacesPagesWithContext(ctx, input, func(page *prometheusservice.ListWorkspacesOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, v := range page.Workspaces {
			if v == nil {
				continue
			}
			output = append(output, v)
		}

		return !lastPage
	})

	if err != nil {
		return nil, err
	}

	return output, nil
}
