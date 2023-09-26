package utils

import (
	"fmt"
	"vkpdeveloper/atlasian-cli/config"
	"vkpdeveloper/atlasian-cli/types"
)

type StatusPageClient struct {
	APIKey string
}

func NewStatusPageClient(config *AppConfig) *StatusPageClient {
	return &StatusPageClient{
		APIKey: config.APIKey,
	}
}

func (c *StatusPageClient) GetPages() (*[]types.Page, error) {
	pagesUrl := config.BaseUrl + "/pages"
	var pages = new([]types.Page)
	err := HttpGet(pagesUrl, c.APIKey, RequestQuery{}, pages)

	if err != nil {
		return nil, err
	}

	return pages, nil
}

func (c *StatusPageClient) GetPageComponents(pageId string) (*[]types.Component, error) {
	pageComponentedUrl := fmt.Sprintf("%s/pages/%s/components", config.BaseUrl, pageId)
	var components = new([]types.Component)
	err := HttpGet(pageComponentedUrl, c.APIKey, RequestQuery{
		PerPage: 100,
	}, components)

	if err != nil {
		return nil, err
	}

	return components, nil
}
