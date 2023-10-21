package helper

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gosimple/slug"
)

func FileUploaded(r *http.Request) (bool, int, string) {
	alias := r.FormValue("alias")
	alias = slug.Make(alias)

	uploadedFile, handler, err := r.FormFile("file")
	if err != nil {
		return false, http.StatusInternalServerError, err.Error()
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return false, http.StatusInternalServerError, err.Error()
	}

	filename := handler.Filename
	if alias != "" {
		filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
	}

	fileLocation := filepath.Join(dir, "files", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return false, http.StatusInternalServerError, err.Error()
	}

	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return false, http.StatusInternalServerError, err.Error()
	}

	return true, 0, filename
}