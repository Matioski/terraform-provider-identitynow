package custom

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	sailpoint "github.com/sailpoint-oss/golang-sdk/v2"
	"golang.org/x/oauth2"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func NewAPIClient(spApiClient *sailpoint.APIClient, config *sailpoint.Configuration) *APIClient {
	return &APIClient{
		ApiClient: spApiClient,
		config:    config,
	}
}

type APIClient struct {
	ApiClient *sailpoint.APIClient
	config    *sailpoint.Configuration
	token     *oauth2.Token
}

func (c *APIClient) doCall(ctx context.Context, method, uri string, body *string, headers map[string]string) (*http.Response, error) {
	fullUrl, err := url.Parse(c.config.ClientConfiguration.BaseURL + uri)
	if err != nil {
		return nil, err
	}
	var request *http.Request
	if body != nil {
		request, err = http.NewRequestWithContext(ctx, method, fullUrl.String(), strings.NewReader(*body))
	} else {
		request, err = http.NewRequestWithContext(ctx, method, fullUrl.String(), nil)
	}
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		request.Header.Add(k, v)
	}
	token, err := c.getAuthToken(ctx)
	if err != nil {
		return nil, err
	}
	token.SetAuthHeader(request)
	response, err := c.config.HTTPClient.StandardClient().Do(request)
	if err != nil || response == nil {
		return nil, err
	}
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusBadRequest {
		return response, fmt.Errorf("error calling %s: %v", fullUrl.String(), response.Status)
	}
	return response, nil
}

func (c *APIClient) unmarshalBody(response *http.Response, v interface{}) error {
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	response.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if err != nil {
		return err
	}
	return json.Unmarshal(bodyBytes, v)
}
func (c *APIClient) getAuthToken(ctx context.Context) (token *oauth2.Token, err error) {
	if c.token != nil && c.token.Valid() {
		return c.token, nil
	}
	config := c.ApiClient.Beta.GetConfig()
	token, err = c.getAccessToken(ctx, config.ClientId, config.ClientSecret, config.TokenURL)
	if err != nil {
		return nil, err
	}
	c.token = token
	return token, nil
}

func (c *APIClient) getAccessToken(ctx context.Context, clientId string, clientSecret string, tokenURL string) (*oauth2.Token, error) {
	tflog.Info(ctx, "Requesting Access Token from "+tokenURL)
	requestUrl := tokenURL
	method := "POST"
	client := &http.Client{}
	form := url.Values{
		"grant_type":    {"client_credentials"},
		"client_id":     {clientId},
		"client_secret": {clientSecret},
	}
	req, err := http.NewRequestWithContext(ctx, method, requestUrl, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var token oauth2.Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
