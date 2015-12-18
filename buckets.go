package b2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type BucketType int

const (
	BUCKET_TYPE_ALL_PRIVATE BucketType = iota
	BUCKET_TYPE_ALL_PUBLIC
)

const (
	listBucketsURL  string = "/b2api/v1/b2_list_buckets"
	createBucketURL string = "/b2api/v1/b2_create_bucket"
	deleteBucketURL string = "/b2api/v1/b2_delete_bucket"
	updateBucketURL string = "/b2api/v1/b2_update_bucket"
)

type Bucket struct {
	BucketId   string `json:"bucketId"`
	AccountId  string `json:"accountId"`
	BucketName string `json:"bucketName"`
	BucketType string `json:"bucketType"`
}

type listBucketsResult struct {
	Buckets []Bucket `json:"buckets"`
}

func (bt BucketType) String() string {
	switch bt {
	case BUCKET_TYPE_ALL_PUBLIC:
		return "allPublic"
	default:
		return "allPrivate"
	}
}

func (c *Client) ListBuckets() ([]Bucket, error) {
	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"accountId": "%s"}`, c.AccountId))
	if req, err := http.NewRequest("POST", c.buildRequestUrl(listBucketsURL), reqBody); err != nil {
		return nil, err
	} else {
		c.setHeaders(req)

		var result listBucketsResult
		err = c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return result.Buckets, nil
	}
}

func (c *Client) CreateBucket(bucketName string, bucketType BucketType) (*Bucket, error) {
	requestBodyData := struct {
		AccountId  string `json:"accountId"`
		BucketName string `json:"bucketName"`
		BucketType string `json:"bucketType"`
	}{
		c.AccountId,
		bucketName,
		bucketType.String(),
	}

	requestBody, err := json.Marshal(requestBodyData)

	if err != nil {
		return nil, err
	}

	if req, err := http.NewRequest("POST", c.buildRequestUrl(createBucketURL), bytes.NewBuffer(requestBody)); err != nil {
		return nil, err
	} else {
		c.setHeaders(req)

		var result Bucket
		err := c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}
}

func (c *Client) DeleteBucket(bucketId string) (*Bucket, error) {
	requestBodyData := struct {
		AccountId string `json:"accountId"`
		BucketId  string `json:"bucketId"`
	}{
		c.AccountId,
		bucketId,
	}

	requestBody, err := json.Marshal(requestBodyData)

	if err != nil {
		return nil, err
	}

	if req, err := http.NewRequest("POST", c.buildRequestUrl(deleteBucketURL), bytes.NewBuffer(requestBody)); err != nil {
		return nil, err
	} else {
		c.setHeaders(req)

		var result Bucket
		err := c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}
}

func (c *Client) UpdateBucket(bucketId string, bucketType BucketType) (*Bucket, error) {
	requestBodyData := struct {
		AccountId  string `json:"accountId"`
		BucketId   string `json:"bucketId"`
		BucketType string `json:"bucketType"`
	}{
		c.AccountId,
		bucketId,
		bucketType.String(),
	}

	requestBody, err := json.Marshal(requestBodyData)

	if err != nil {
		return nil, err
	}

	if req, err := http.NewRequest("POST", c.buildRequestUrl(updateBucketURL), bytes.NewBuffer(requestBody)); err != nil {
		return nil, err
	} else {
		c.setHeaders(req)

		var result Bucket
		err := c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}
}
