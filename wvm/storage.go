package wvm

import (
	"fmt"

	"github.com/alphabill-org/alphabill/keyvaluedb"
	"github.com/alphabill-org/alphabill/keyvaluedb/memorydb"
)

type MemoryStorage struct {
	db keyvaluedb.KeyValueDB
}

func NewMemoryStorage() *MemoryStorage {
	memDB, _ := memorydb.New()
	return &MemoryStorage{
		db: memDB,
	}
}

func NewAlwaysFailsMemoryStorage() *MemoryStorage {
	mdb, _ := memorydb.New()
	mdb.MockWriteError(fmt.Errorf("mock db write failed"))
	return &MemoryStorage{
		db: mdb,
	}
}

func (s *MemoryStorage) Get(key []byte) ([]byte, error) {
	var stateFile = make([]byte, 0)
	found, err := s.db.Read(key, &stateFile)
	if err != nil {
		return nil, fmt.Errorf("read program state failed, %v", err)
	}
	if !found {
		return nil, fmt.Errorf("state file not found")
	}
	return stateFile, nil
}

func (s *MemoryStorage) Put(key []byte, file []byte) error {
	if err := s.db.Write(key, file); err != nil {
		return fmt.Errorf("storage write failed, %w", err)
	}
	return nil
}
