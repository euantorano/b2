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
	api, err := b2.NewClient(Credentials{"YOUR_ACCOUNT_ID", "YOUR_APPLICATION_KEY"})

	if err != nil {
		panic(err)
	}

	buckets, err := api.ListBuckets()

	for _, bucket := range buckets {
		fmt.Println("Bucket name: %s", bucket.BucketName)
	}
}
```

## TODO

- [x] API Auth
- [x] Bucket handling
- [ ] File handling

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
