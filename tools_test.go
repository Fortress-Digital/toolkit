package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	s := testTools.RandomString(10)

	if len(s) != 10 {
		t.Errorf("TestTools.RandomString returned wrong length: got %v want %v", len(s), 10)
	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{
		name:          "allowed no rename",
		allowedTypes:  []string{"image/jpeg", "image/png"},
		renameFile:    false,
		errorExpected: false,
	},
	{
		name:          "allowed rename",
		allowedTypes:  []string{"image/jpeg", "image/png"},
		renameFile:    true,
		errorExpected: false,
	},
	{
		name:          "not allowed",
		allowedTypes:  []string{"image/jpeg"},
		renameFile:    false,
		errorExpected: true,
	},
}

func TestTools_UploadFiles(t *testing.T) {
	for _, test := range uploadTests {
		// Simulate multipart upload
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			// Create form data field
			part, err := writer.CreateFormFile("file", "./test-data/img.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./test-data/img.png")
			if err != nil {
				t.Error(err)
			}

			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("error decoding image", err)
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error("error encoding image", err)
			}
		}()

		// Read from pipe
		request := httptest.NewRequest(http.MethodPost, "/upload", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools

		testTools.AllowedFileTypes = test.allowedTypes

		uploadedFiles, err := testTools.UploadFiles(request, "./test-data/uploads/", test.renameFile)
		if err != nil && !test.errorExpected {
			t.Error(err)
		}

		if !test.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./test-data/uploads/%s", uploadedFiles[0].NewFileName)); os.IsExist(err) {
				t.Errorf("%s: expected file to exist %s", test.name, err.Error())
			}

			_ = os.Remove(fmt.Sprintf("./test-data/uploads/%s", uploadedFiles[0].NewFileName))
		}

		if !test.errorExpected && err != nil {
			t.Errorf("%s: error expected but none received", test.name)
		}

		wg.Wait()
	}
}

func TestTools_UploadFile(t *testing.T) {
	for _, test := range uploadTests {
		// Simulate multipart upload
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)

		go func() {
			defer writer.Close()

			// Create form data field
			part, err := writer.CreateFormFile("file", "./test-data/img.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./test-data/img.png")
			if err != nil {
				t.Error(err)
			}

			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("error decoding image", err)
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error("error encoding image", err)
			}
		}()

		// Read from pipe
		request := httptest.NewRequest(http.MethodPost, "/upload", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools

		uploadedFiles, err := testTools.UploadFile(request, "./test-data/uploads/", true)
		if err != nil {
			t.Error(err)
		}

		if _, err := os.Stat(fmt.Sprintf("./test-data/uploads/%s", uploadedFiles.NewFileName)); os.IsExist(err) {
			t.Errorf("%s: expected file to exist %s", test.name, err.Error())
		}

		_ = os.Remove(fmt.Sprintf("./test-data/uploads/%s", uploadedFiles.NewFileName))
	}
}
