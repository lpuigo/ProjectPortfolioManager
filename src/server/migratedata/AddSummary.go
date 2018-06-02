package migratedata

import (
	"archive/zip"
	"fmt"
	ris "github.com/lpuig/prjptf/src/server/manager/recordindexedset"
	"io"
	"io/ioutil"
	"strings"
)

func NewSetFrom(r io.Reader) (*ris.RecordLinkedIndexedSet, error) {
	cs := ris.NewRecordLinkedIndexedSet(
		ris.NewIndexDesc("Issue", "ISSUE"),
	)
	//cs.AddLink(ris.NewLinkDesc("Issue-Summary", "Issue", "Summary"))
	err := cs.AddCSVDataFrom(r)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func NewTargetSetFrom(r io.Reader) (*ris.RecordLinkedIndexedSet, error) {
	cs := ris.NewRecordLinkedIndexedSet(
		ris.NewIndexDesc("Issue", "ISSUE"),
		ris.NewIndexDesc("Summary", "SUMMARY"),
	)
	cs.AddLink(ris.NewLinkDesc("Issue-Summary", "Issue", "Summary"))
	err := cs.AddCSVDataFrom(r)
	if err != nil {
		return nil, err
	}
	return cs, nil
}

func MigrateSet(r io.Reader, targetModel *ris.RecordLinkedIndexedSet) (*ris.RecordLinkedIndexedSet, error) {
	ss, err := NewSetFrom(r)
	if err != nil {
		return nil, err
	}
	// check if Set from r is already target
	if _, err := ss.GetRecordColNumByName("SUMMARY"); err == nil {
		return ss, nil
	}
	// too bad, its not ... let's migrate it
	a := targetModel.GetLink("Issue-Summary")

	nt, err := targetModel.CreateSubSet([]ris.IndexDesc{ris.NewIndexDesc("Issue", "ISSUE")}, nil)

	for _, r := range ss.GetRecords() {
		summary := a.Get(ss.GetRecordKeyByIndex("Issue", r), "")
		nr := append(r, strings.TrimLeft(summary, "!"))
		nt.AddRecord(nr)
	}
	return nt, nil
}

func ProcessDir(dir string, action func(file string) error) error {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("input dir: %s", err.Error())
	}
	for _, file := range files {
		err := action(dir + file.Name())
		if err != nil {
			return err
		}
	}
	return nil
}

func ExtractAndProcess(archiveFileName string, process func(r io.Reader, fileName string) error) error {
	unzipArchiveFile := func(zfile *zip.File) error {
		fileReader, err := zfile.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()
		err = process(fileReader, zfile.Name)
		if err != nil {
			return err
		}
		return nil
	}

	unzipArchive := func(afile string) error {
		zipReader, err := zip.OpenReader(afile)
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

	err := unzipArchive(archiveFileName)
	if err != nil {
		return err
	}
	return nil
}
