package db

import "sync"

type Table struct {
	mutex sync.RWMutex
	maxId uint64
	data  map[uint64]DBObject
}

func NewTable() *Table {
	table := Table{
		sync.RWMutex{},
		0,
		make(map[uint64]DBObject),
	}
	return &table
}

func (table *Table) Store(obj DBObject) {
	table.mutex.Lock()
	defer table.mutex.Unlock()

	table.maxId++ // 0 is an invalid id: first increment, then add
	obj.setPK(table.maxId)
	table.data[table.maxId] = obj
}

func (table *Table) Load(pk uint64) (DBObject, error) {
	table.mutex.RLock()
	defer table.mutex.RUnlock()

	if val, ok := table.data[pk]; ok {
		return val, nil
	} else {
		return nil, &ObjectNotFoundError{pk}
	}
}

func (table *Table) Update(obj DBObject) error {
	// write-lock is only necessary to prevent access to the same object
	// i'd write a more fine-grained mutex but it's only a mock so who gives a shit

	pk := obj.GetPK()

	table.mutex.Lock()
	defer table.mutex.Unlock()

	if _, ok := table.data[pk]; ok {
		table.data[pk] = obj
		return nil
	} else {
		return &ObjectNotFoundError{pk}
	}
}

func (table *Table) Delete(pk uint64) error {
	// id remains untouched
	table.mutex.Lock()
	defer table.mutex.Unlock()

	if _, ok := table.data[pk]; ok {
		delete(table.data, pk)
		return nil
	} else {
		return &ObjectNotFoundError{pk}
	}
}
