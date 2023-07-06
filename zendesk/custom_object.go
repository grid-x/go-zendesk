package zendesk

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type CustomObjectRecord struct {
	Url                string                 `json:"url"`
	Name               string                 `json:"name"`
	ID                 string                 `json:"id"`
	CustomObjectKey    string                 `json:"custom_object_key"`
	CustomObjectFields map[string]interface{} `json:"custom_object_fields"`
	CreatedByUserID    string                 `json:"created_by_user_id"`
	UpdatedByUserID    string                 `json:"updated_by_user_id"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
	ExternalID         string                 `json:"external_id"`
}

// CustomObjectAPI an interface containing all custom object related methods
type CustomObjectAPI interface {
	CreateCustomObjectRecord(
		ctx context.Context, record CustomObjectRecord, customObjectKey string) (CustomObjectRecord, error)
	SearchCustomObjectRecords(
		ctx context.Context,
		customObjectKey string,
		opts *CustomObjectListOptions,
	) ([]CustomObjectRecord, Page, error)
}

// CustomObjectListOptions custom object search options
type CustomObjectListOptions struct {
	PageOptions
	Name string `url:"name"`
}

// CreateCustomObjectRecord CreateCustomObject create a custom object record
func (z *Client) CreateCustomObjectRecord(
	ctx context.Context, record CustomObjectRecord, customObjectKey string,
) (CustomObjectRecord, error) {

	var data, result struct {
		CustomObjectRecord CustomObjectRecord `json:"custom_object_record"`
	}
	data.CustomObjectRecord = record

	body, err := z.post(ctx, fmt.Sprintf("/custom_objects/%s/records.json", customObjectKey), data)
	if err != nil {
		return CustomObjectRecord{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return CustomObjectRecord{}, err
	}
	return result.CustomObjectRecord, nil
}

// SearchCustomObjectRecords search for a custom object record by the name field
// https://developer.zendesk.com/api-reference/custom-objects/custom_object_records/#search-custom-object-records
func (z *Client) SearchCustomObjectRecords(
	ctx context.Context, customObjectKey string, opts *CustomObjectListOptions) ([]CustomObjectRecord, Page, error) {
	var result struct {
		CustomObjectRecords []CustomObjectRecord `json:"custom_object_records"`
		Page
	}
	tmp := opts
	if tmp == nil {
		tmp = &CustomObjectListOptions{}
	}
	url := fmt.Sprintf("/custom_objects/%s/records/autocomplete", customObjectKey)
	urlWithOptions, err := addOptions(url, tmp)
	body, err := z.get(ctx, urlWithOptions)

	if err != nil {
		return nil, Page{}, err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, Page{}, err
	}
	return result.CustomObjectRecords, result.Page, nil
}
