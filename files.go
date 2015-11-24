package b2

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type HideFileAction string

const (
	HIDE_FILE_ACTION_HIDE   HideFileAction = "hide"
	HIDE_FILE_ACTION_UPLOAD                = "upload"
)

const (
	DELETE_FILE_VERSION_URL string = "/b2api/v1/b2_delete_file_version"
	DOWNLOAD_FILE_BY_ID_URL string = "/b2api/v1/b2_download_file_by_id"
	GET_FILE_INFO_URL       string = "/b2api/v1/b2_get_file_info"
	GET_UPLOAD_URL_URL      string = "/b2api/v1/b2_get_upload_url"
	HIDE_FILE_URL           string = "/b2api/v1/b2_hide_file"
	LIST_FILE_NAMES_URL     string = "/b2api/v1/b2_list_file_names"
)

type FileVersion struct {
	FileId   string `json:"fileId"`
	FileName string `json:"fileName"`
}

type FileInfo struct {
	ContentLength int               `json:"contentLength"`
	ContentType   string            `json:"contentType"`
	FileId        string            `json:"fileId"`
	FileName      string            `json:"fileName"`
	ContentSha1   string            `json:"contentSha1"`
	BucketId      string            `json:"bucketId"`
	AccountId     string            `json:"accountId"`
	Info          map[string]string `json:"fileInfo"`
}

type UploadUrlDetails struct {
	BucketId  string `json:"bucketId"`
	Url       string `json:"uploadUrl"`
	AuthToken string `json:"authorizationToken"`
}

type File struct {
	FileId          string         `json:"fileId"`
	FileName        string         `json:"fileName"`
	UploadTimestamp int            `json:"uploadTimestamp"`
	Action          HideFileAction `json:"action"`
	Size            int64          `json:"size"`
}

type FileCollection struct {
	Files        []File `json:"files"`
	NextFileName string `json:"nextFileName"`
}

func (c *Client) DeleteFileVersion(fileName string, fileId string) (*FileVersion, error) {
	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"fileName": "%s", "fileId": "%s"}`, fileName, fileId))
	if req, err := http.NewRequest("POST", c.buildRequestUrl(DELETE_FILE_VERSION_URL), reqBody); err != nil {
		return nil, err
	} else {
		c.setHeaders(req)

		var result FileVersion
		err = c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}
}

func (c *Client) DownloadFileById(fileId string) ([]byte, *FileInfo, error) {
	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"fileId": "%s"}`, fileId))
	if req, err := http.NewRequest("POST", c.buildFileRequestUrl(DOWNLOAD_FILE_BY_ID_URL), reqBody); err != nil {
		return nil, nil, err
	} else {
		c.setHeaders(req)

		data, header, err := c.requestBytes(req)

		if err != nil {
			return nil, nil, err
		}

		contentLen, err := strconv.Atoi(header.Get("Content-Length"))

		info := &FileInfo{
			ContentLength: contentLen,
			ContentType:   header.Get("Content-Type"),
			FileId:        header.Get("X-Bz-File-Id"),
			FileName:      header.Get("X-Bz-File-Name"),
			ContentSha1:   header.Get("X-Bz-Content-Sha1"),
			Info:          make(map[string]string),
		}

		for key, _ := range header {
			// make the key lowercase as some of the headers B2 returns are lowercase, some uppercase...
			if strings.HasPrefix(strings.ToLower(key), "x-bz-info-") {
				info.Info[key] = header.Get(key)
			}
		}

		return data, info, nil
	}
}

func (c *Client) DownloadFileByName(bucketName string, fileName string) ([]byte, *FileInfo, error) {
	requestPath := fmt.Sprintf("/file/%s/%s", bucketName, fileName)
	if req, err := http.NewRequest("GET", c.buildFileRequestUrl(requestPath), nil); err != nil {
		return nil, nil, err
	} else {
		req.Header.Set("Authorization", c.AuthToken)

		data, header, err := c.requestBytes(req)

		if err != nil {
			return nil, nil, err
		}

		contentLen, err := strconv.Atoi(header.Get("Content-Length"))

		info := &FileInfo{
			ContentLength: contentLen,
			ContentType:   header.Get("Content-Type"),
			FileId:        header.Get("X-Bz-File-Id"),
			FileName:      header.Get("X-Bz-File-Name"),
			ContentSha1:   header.Get("X-Bz-Content-Sha1"),
			Info:          make(map[string]string),
		}

		for key, _ := range header {
			// make the key lowercase as some of the headers B2 returns are lowercase, some uppercase...
			if strings.HasPrefix(strings.ToLower(key), "x-bz-info-") {
				info.Info[key] = header.Get(key)
			}
		}

		return data, info, nil
	}
}

func (c *Client) GetFileInfo(fileId string) (*FileInfo, error) {
	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"fileId": "%s"}`, fileId))
	if req, err := http.NewRequest("POST", c.buildRequestUrl(GET_FILE_INFO_URL), reqBody); err != nil {
		return nil, err
	} else {
		c.setHeaders(req)

		var result FileInfo
		err = c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}
}

func (c *Client) GetUploadUrl(bucketId string) (*UploadUrlDetails, error) {
	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"bucketId": "%s"}`, bucketId))
	if req, err := http.NewRequest("POST", c.buildRequestUrl(GET_UPLOAD_URL_URL), reqBody); err != nil {
		return nil, err
	} else {
		c.setHeaders(req)

		var result UploadUrlDetails
		err = c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}
}

func (c *Client) HideFile(bucketId, fileName string) (*File, error) {
	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"bucketId": "%s", "fileName": "%s"}`, bucketId, fileName))
	if req, err := http.NewRequest("POST", c.buildRequestUrl(HIDE_FILE_URL), reqBody); err != nil {
		return nil, err
	} else {
		c.setHeaders(req)

		var result File
		err = c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}
}

func (c *Client) ListFileNamesWithCountAndOffset(bucketId, startFileName string, maxFileCount int) (*FileCollection, error) {
	if !validateMaxFileCount(maxFileCount) {
		return nil, errors.New("maxFileCount must be between 1 and 1000")
	}

	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"bucketId": "%s", "startFileName": "%s", "maxFileCount": %d}`, bucketId, startFileName, maxFileCount))
	if req, err := http.NewRequest("POST", c.buildRequestUrl(LIST_FILE_NAMES_URL), reqBody); err != nil {
		return nil, err
	} else {
		c.setHeaders(req)

		var result FileCollection
		err = c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}
}

func (c *Client) ListFileNames(bucketId string) (*FileCollection, error) {
	return c.ListFileNamesWithCountAndOffset(bucketId, "", 100)
}

func validateMaxFileCount(count int) bool {
	return count > 0 && count <= 1000
}
