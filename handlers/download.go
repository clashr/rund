package handlers

import (
	"net/http"

	"github.com/clashr/go-archive"
)

func download(src string, dest string) error {
	resp, err := http.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := archive.Untar(resp.Body, dest, nil); err != nil {
		return err
	}

	return nil
}
