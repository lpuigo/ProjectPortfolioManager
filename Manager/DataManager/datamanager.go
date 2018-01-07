package DataManager

import (
	"log"
	"os"
	"sync"
	"time"
)

const (
	JSONFileLock = `./Ressources/nosave.json`
	SaveDelay    = 2
)

type DataManager struct {
	mutex       sync.RWMutex
	saveTimer   *time.Timer
	noSave      bool
	persistFunc func() error
}

func NewDataManager(persistFunc func() error) *DataManager {
	m := &DataManager{}
	if lf, err := os.Open(JSONFileLock); err == nil {
		lf.Close()
		log.Printf("NoSave file found : changes won't be persisted")
		m.noSave = true
	}
	m.persistFunc = persistFunc
	return m
}

// WLockPtf locks the DataManager for Concurentsafe Modifications
func (dm *DataManager) WLock() {
	dm.mutex.Lock()
}

// RLockPtf locks the DataManager for Read operations
func (dm *DataManager) RLock() {
	dm.mutex.RLock()
}

// WUnlockPtf unlocks the DataManager after Modifications and triggers JSON file persist mechanism with delay (if not already armed)
func (dm *DataManager) WUnlockWithPersist() {
	if dm.saveTimer == nil {
		dm.saveTimer = time.AfterFunc(SaveDelay*time.Second, dm.persist)
	}
	dm.mutex.Unlock()
}

// WUnlockPtfNoPersist unlocks the DataManager after Modifications without triggering persist mechanism
func (dm *DataManager) WUnlock() {
	dm.mutex.Unlock()
}

// RUnlockPtf unlocks the DataManager after Read operations
func (dm *DataManager) RUnlock() {
	dm.mutex.RUnlock()
}

func (dm *DataManager) persist() {
	if dm.noSave {
		log.Println("NoSave mode ON")
		return
	}
	dm.RLock()
	if err := dm.persistFunc(); err != nil {
		log.Println("Unable to persist Data :", err.Error())
	} else {
		log.Println("Data persisted successfully")
	}
	dm.saveTimer = nil
	dm.RUnlock()
}
