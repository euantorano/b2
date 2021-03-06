package b2

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type HideFileAction string

const (
	HIDE_FILE_ACTION_HIDE   HideFileAction = "hide"
	HIDE_FILE_ACTION_UPLOAD                = "upload"
)

const (
	deleteFileVersionURL string = "/b2api/v1/b2_delete_file_version"
	downloadFileByIdURL         = "/b2api/v1/b2_download_file_by_id"
	getFileInfoURL              = "/b2api/v1/b2_get_file_info"
	getUploadUrlURL             = "/b2api/v1/b2_get_upload_url"
	hideFileURL                 = "/b2api/v1/b2_hide_file"
	listFileNamesURL            = "/b2api/v1/b2_list_file_names"
	listFileVersionsURL         = "/b2api/v1/b2_list_file_versions"
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
	if req, err := http.NewRequest("POST", c.buildRequestUrl(deleteFileVersionURL), reqBody); err != nil {
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
	if req, err := http.NewRequest("POST", c.buildFileRequestUrl(downloadFileByIdURL), reqBody); err != nil {
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
	if req, err := http.NewRequest("POST", c.buildRequestUrl(getFileInfoURL), reqBody); err != nil {
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
	if req, err := http.NewRequest("POST", c.buildRequestUrl(getUploadUrlURL), reqBody); err != nil {
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
	if req, err := http.NewRequest("POST", c.buildRequestUrl(hideFileURL), reqBody); err != nil {
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
	if req, err := http.NewRequest("POST", c.buildRequestUrl(listFileNamesURL), reqBody); err != nil {
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

func (c *Client) ListFileVersionsWithCountAndOffset(bucketId, startFileName string, maxFileCount int) (*FileCollection, error) {
	if !validateMaxFileCount(maxFileCount) {
		return nil, errors.New("maxFileCount must be between 1 and 1000")
	}

	reqBody := bytes.NewBufferString(fmt.Sprintf(`{"bucketId": "%s", "startFileName": "%s", "maxFileCount": %d}`, bucketId, startFileName, maxFileCount))
	if req, err := http.NewRequest("POST", c.buildRequestUrl(listFileVersionsURL), reqBody); err != nil {
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

func (c *Client) ListFileVersions(bucketId string) (*FileCollection, error) {
	return c.ListFileVersionsWithCountAndOffset(bucketId, "", 100)
}

func (c *Client) UploadFile(uploadUrl, uploadAuthToken, filePath string) (*FileInfo, error) {
	fileName := filepath.Base(filePath)

	return c.UploadFileWithFileName(uploadUrl, uploadAuthToken, filePath, fileName)
}

func (c *Client) UploadFileWithFileName(uploadUrl, uploadAuthToken, filePath, fileName string) (*FileInfo, error) {
	filePath = filepath.Clean(filePath)

	fileDetails, err := os.Stat(filePath)

	if err != nil {
		return nil, err
	}

	fileLastModifiedMillis := fileDetails.ModTime().Unix()

	fileContent, err := ioutil.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	if req, err := http.NewRequest("POST", uploadUrl, bytes.NewBuffer(fileContent)); err != nil {
		return nil, err
	} else {
		req.Header.Set("Authorization", uploadAuthToken)
		req.Header.Set("X-Bz-File-Name", url.QueryEscape(fileName))
		req.Header.Set("Content-Type", "b2/x-auto")
		req.Header.Set("X-Bz-Content-Sha1", c.getHashForContent(fileContent))
		req.Header.Set("X-Bz-Info-src_last_modified_millis", strconv.FormatInt(fileLastModifiedMillis, 10))

		var result FileInfo
		err = c.requestJson(req, &result)

		if err != nil {
			return nil, err
		}

		return &result, nil
	}
}

func (c *Client) getHashForContent(content []byte) string {
	sum := sha1.Sum(content)

	return fmt.Sprintf("%x", sum)
}
