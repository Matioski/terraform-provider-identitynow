package custom

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

const defaultCronExpression = "0 0 0 * * ?"

// todo: getEntitlementAggregationSchedules

func (c *APIClient) ReadSourceAccountAggregationSchedule(ctx context.Context, sourceCloudId string) (*SourceAggregationSchedule, *http.Response, error) {
	uri := fmt.Sprintf("/cc/api/source/getAggregationSchedules/%s", sourceCloudId)

	return c.getScheduledAggregation(ctx, uri)
}

func (c *APIClient) ModifySourceAccountAggregationSchedule(ctx context.Context, sourceCloudId, cronExpression string) (*SourceAggregationSchedule, *http.Response, error) {
	uri := fmt.Sprintf("/cc/api/source/scheduleAggregation/%s", sourceCloudId)
	data := url.Values{}
	data.Set("enable", "true")
	data.Set("cronExp", cronExpression)

	return c.invokeScheduledAggregation(ctx, uri, data)
}

func (c *APIClient) DeleteSourceAccountAggregationSchedule(ctx context.Context, sourceCloudId string) (*http.Response, error) {
	uri := fmt.Sprintf("/cc/api/source/scheduleAggregation/%s", sourceCloudId)
	data := url.Values{}
	data.Set("enable", "false")
	data.Set("cronExp", defaultCronExpression)
	_, response, err := c.invokeScheduledAggregation(ctx, uri, data)
	return response, err
}

func (c *APIClient) ReadSourceEntitlementAggregationSchedule(ctx context.Context, sourceCloudId string) (*SourceAggregationSchedule, *http.Response, error) {
	uri := fmt.Sprintf("/cc/api/source/getEntitlementAggregationSchedules/%s", sourceCloudId)

	return c.getScheduledAggregation(ctx, uri)
}

func (c *APIClient) ModifySourceEntitlementAggregationSchedule(ctx context.Context, sourceCloudId, cronExpression string) (*SourceAggregationSchedule, *http.Response, error) {
	uri := fmt.Sprintf("/cc/api/source/scheduleEntitlementAggregation/%s", sourceCloudId)
	data := url.Values{}
	data.Set("enable", "true")
	data.Set("cronExp", cronExpression)

	return c.invokeScheduledAggregation(ctx, uri, data)
}

func (c *APIClient) DeleteSourceEntitlementAggregationSchedule(ctx context.Context, sourceCloudId string) (*http.Response, error) {
	uri := fmt.Sprintf("/cc/api/source/scheduleEntitlementAggregation/%s", sourceCloudId)
	data := url.Values{}
	data.Set("enable", "false")
	data.Set("cronExp", defaultCronExpression)
	_, response, err := c.invokeScheduledAggregation(ctx, uri, data)
	return response, err
}

func (c *APIClient) invokeScheduledAggregation(ctx context.Context, uri string, data url.Values) (*SourceAggregationSchedule, *http.Response, error) {
	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/x-www-form-urlencoded; charset=utf-8",
	}
	body := data.Encode()
	response, err := c.doCall(ctx, http.MethodPost, uri, &body, headers)
	if err != nil {
		return nil, nil, err
	}
	var config SourceAggregationSchedule
	if err = c.unmarshalBody(response, &config); err != nil {
		return nil, response, err
	}
	return &config, response, nil
}

func (c *APIClient) getScheduledAggregation(ctx context.Context, uri string) (*SourceAggregationSchedule, *http.Response, error) {
	headers := map[string]string{
		"Accept": "application/json",
	}
	response, err := c.doCall(ctx, http.MethodGet, uri, nil, headers)
	if err != nil {
		return nil, nil, err
	}
	var configs []SourceAggregationSchedule
	if err = c.unmarshalBody(response, &configs); err != nil {
		return nil, response, err
	}
	if len(configs) == 0 {
		return nil, response, nil
	}
	return &configs[0], response, nil
}

type SourceAggregationSchedule struct {
	CronExpressions []string `json:"cronExpressions"`
}
