# b2

`b2` is a Go wrapper for the Backblaze B2 API.

It provides convenient access to buckets and files within B2.

## Installation

`go get github.com/euantorano/b2`

## Usage

```go
package main

import (
	"fmt"
	"github.com/euantorano/b2"
)

func main() {
	api, err := b2.NewClient(b2.Credentials{"YOUR_ACCOUNT_ID", "YOUR_APPLICATION_KEY"})

	if err != nil {
		panic(err)
	}

	buckets, err := api.ListBuckets()

	for _, bucket := range buckets {
		fmt.Println("Bucket name: %s", bucket.BucketName)
	}
}
```

### Uploading files

In order to upload files, you mst know the ID of the bucket you want to upload to, and you must get an upload URL first. An example can be seen below:

```go
package main

import (
	"fmt"
	"github.com/euantorano/b2"
)

func main() {
	api, err := b2.NewClient(b2.Credentials{"YOUR_ACCOUNT_ID", "YOUR_APPLICATION_KEY"})

	if err != nil {
		panic(err)
	}

	uploadUrl, err := api.GetUploadUrl("BUCKET_ID")

	if err != nil {
		panic(err)
	}

	uploaded, err := api.UploadFile(uploadUrl.Url, uploadUrl.AuthToken, "/Users/euan/Desktop/hello_world.txt")

	if err != nil {
		panic(err)
	}

	/*
	uploaded now contains details about the uploaded file, such as:

	- fileId: The unique identifier for this version of this file. Used with b2_get_file_info, b2_download_file_by_id, and b2_delete_file_version.

	- fileName: The name of this file, which can be used with b2_download_file_by_name.

	- accountId: Your account ID.

	- bucketId: The bucket that the file is in.

	- contentLength: The number of bytes stored in the file.

	- contentSha1: The SHA1 of the bytes stored in the file.

	- contentType: The MIME type of the file.

	- fileInfo: The custom information that was uploaded with the file. This is a JSON object, holding the name/value pairs that were uploaded with the file.
	*/
}
```

## Licence

```
The MIT License (MIT)

Copyright (c) 2015 Euan T

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
