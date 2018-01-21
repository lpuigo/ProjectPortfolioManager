package FileProcesser

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type FileProcesser struct {
	inputDir   string
	archiveDir string
}

func NewFileManager(inputDir, archiveDir string) (*FileProcesser, error) {
	_, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(inputDir, `\`) {
		inputDir += `\`
	}
	if !strings.HasSuffix(archiveDir, `\`) {
		archiveDir += `\`
	}
	_, err = ioutil.ReadDir(archiveDir)
	if err != nil {
		return nil, err
	}
	return &FileProcesser{inputDir: inputDir, archiveDir: archiveDir}, nil
}

// Process undertakes given action function on each file found in InputDir. Once file is processed ok, it is compressed and move to ArchiveDir
func (fm *FileProcesser) Process(action func(file string) error) error {
	files, err := ioutil.ReadDir(fm.inputDir)
	if err != nil {
		return fmt.Errorf("input dir: %s", err.Error())
	}
	for _, file := range files {
		err := action(fm.inputDir + file.Name())
		if err != nil {
			return err
		}
		go fm.achiveFile(file.Name())
	}
	return nil
}

func (fm *FileProcesser) achiveFile(file string) error {
	archiveFile := file + ".zip"

	// Create new archive file
	zipf, err := os.Create(fm.archiveDir + archiveFile)
	if err != nil {
		return err
	}
	defer zipf.Close()

	zipWriter := zip.NewWriter(zipf)
	defer zipWriter.Close()

	// Add file to archive file
	archivedFile := fm.inputDir + file
	archivedf, err := os.Open(archivedFile)
	if err != nil {
		return err
	}
	defer archivedf.Close()

	info, err := archivedf.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	header.Method = zip.Deflate
	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, archivedf)
	if err != nil {
		return err
	}

	// and finnally delete zipped file
	archivedf.Close()
	return os.Remove(archivedFile)
}
