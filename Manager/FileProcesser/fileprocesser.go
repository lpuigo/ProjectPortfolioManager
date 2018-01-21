package FileProcesser

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

type FileProcesser struct {
	inputDir   string
	archiveDir string
	wg         sync.WaitGroup
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

// ProcessAndArchive undertakes given action function on each file found in InputDir. Once file is processed ok, it is compressed and move to ArchiveDir
func (fm *FileProcesser) ProcessAndArchive(action func(file string) error) error {
	files, err := ioutil.ReadDir(fm.inputDir)
	if err != nil {
		return fmt.Errorf("input dir: %s", err.Error())
	}
	for _, file := range files {
		err := action(fm.inputDir + file.Name())
		if err != nil {
			return err
		}
		fm.wg.Add(1)
		go fm.achiveFile(file.Name())
	}
	fm.wg.Wait()
	return nil
}

func (fm *FileProcesser) RestoreArchives() error {
	archiveFiles, err := ioutil.ReadDir(fm.archiveDir)
	if err != nil {
		return fmt.Errorf("archive dir: %s", err.Error())
	}

	unzipArchiveFile := func(zfile *zip.File) error {
		path := fm.inputDir + zfile.Name
		fileReader, err := zfile.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()
		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, zfile.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()
		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
		return nil
	}

	unzipArchive := func(afi os.FileInfo) error {
		zipReader, err := zip.OpenReader(fm.archiveDir + afi.Name())
		if err != nil {
			return err
		}
		defer zipReader.Close()
		for _, file := range zipReader.File {
			err := unzipArchiveFile(file)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, archiveFileInfo := range archiveFiles {
		// unzip archive file and restore contained file(s) to inputDir
		err := unzipArchive(archiveFileInfo)
		if err != nil {
			return err
		}
		// remove archivefile
		err = os.Remove(fm.archiveDir + archiveFileInfo.Name())
	}
	return nil
}

func (fm *FileProcesser) achiveFile(file string) error {
	defer fm.wg.Done()
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
