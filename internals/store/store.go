package store

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"sync"
)

type StoreManager struct {
	storeIds map[string]struct{}
}

var (
	once          sync.Once
	storeinstance *StoreManager
)

func NewStoreManager() (*StoreManager, error) {
	var err error
	once.Do(func() {
		storeinstance = &StoreManager{
			storeIds: make(map[string]struct{}),
		}
		err = storeinstance.LoadStoreIds()
	})
	return storeinstance, err
}

func (sm *StoreManager) LoadStoreIds() error {
	file, err := os.Open(os.Getenv("CSVFILEPATH"))
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Check if there are records and skip the header
	if len(records) > 0 {
		for _, record := range records[1:] {
			if len(record) > 2 {
				storeID := strings.TrimSpace(record[2])
				sm.storeIds[storeID] = struct{}{}
			}
		}
	}
	log.Println("All records read and stored in StoreManager")
	return nil
}

func (sm *StoreManager) CheckStoreIDExist(store_id string) bool {
	_, exists := sm.storeIds[store_id]
	return exists
}
