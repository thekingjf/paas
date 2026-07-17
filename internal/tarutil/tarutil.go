package tarutil

import (
	"archive/tar"
	"bytes"
	"os"
	"path/filepath"
)

func streamIn(folderLocation string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	folder, err := filepath.Abs(folderLocation)
	writer := tar.NewWriter(&buf)

	if err != nil {
		return nil, err
	}

	files, err := os.ReadDir(folder)

	if err != nil {
		return nil, err
	}

	for _, item := range files {
		err = streamInHelper(folderLocation, folder, "", item, writer)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()
	return &buf, nil
}

func streamInDive(folderLocation, src string, writer *tar.Writer) error {

	folder := filepath.Join(folderLocation, src)

	files, err := os.ReadDir(folder)

	if err != nil {
		return err
	}

	for _, item := range files {
		err = streamInHelper(folderLocation, folder, src, item, writer)
		if err != nil {
			return err
		}
	}

	return nil
}

func streamInHelper(folderLocation, folder, src string, item os.DirEntry, writer *tar.Writer) error {
	if item.IsDir() {
		err := streamInDive(folderLocation, filepath.Join(src, item.Name()), writer)

		if err != nil {
			return err
		}
	} else {
		file, err := item.Info()

		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(file, "")

		if err != nil {
			return err
		}

		header.Name = filepath.Join(src, file.Name())
		err = writer.WriteHeader(header)

		if err != nil {
			return err
		}

		book, err := os.ReadFile(filepath.Join(folder, file.Name()))

		if err != nil {
			return err
		}

		_, err = writer.Write(book)

		if err != nil {
			return err
		}
	}

	return nil
}
