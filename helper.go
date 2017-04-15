package mastodon

import (
	"encoding/base64"
	"net/http"
	"os"
)

func Base64EncodeFileName(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	return Base64Encode(file)
}

func Base64Encode(file *os.File) (string, error) {
	fi, err := file.Stat()
	if err != nil {
		return "", err
	}

	d := make([]byte, fi.Size())
	_, err = file.Read(d)
	if err != nil {
		return "", err
	}

	return "data:" + http.DetectContentType(d) +
		";base64," + base64.StdEncoding.EncodeToString(d), nil
}

// String is a helper function to get the pointer value of a string.
func String(v string) *string { return &v }
