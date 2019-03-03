package db

type DBObject interface {
	GetPK() uint64
	setPK(uint64)

	GetTableName() string
}

// no createTable or anything - it's just a mock
type Database interface {
	Delete(string, uint64) error
	Load(string, uint64) (DBObject, error)
	Store(DBObject) error
	Update(DBObject) error
}

type MockDB map[string]*Table

func NewMockDB(tables ...string) *MockDB {
  db := make(MockDB, len(tables))
  for _, table := range tables {
  	db[table] = NewTable()
	}
  return &db
}

func (db *MockDB) getTable(tableName string) (*Table, error) {
	if table, ok := (*db)[tableName]; ok {
		return table, nil
	} else {
		return nil, &TableNotFoundError{tableName}
	}
}

func (db *MockDB) getTableFromObj(obj DBObject) (*Table, error) {
	return db.getTable(obj.GetTableName())
}

func (db *MockDB) Store(obj DBObject) (err error) {
	if table, err := db.getTableFromObj(obj); err == nil {
		table.Store(obj)
	}
	return
}

func (db *MockDB) Load(tableName string, pk uint64) (DBObject, error) {
	if table, err := db.getTable(tableName); err == nil {
		return table.Load(pk)
	} else {
		return nil, err
	}
}

func (db *MockDB) Update(obj DBObject) (err error) {
	if table, err := db.getTableFromObj(obj); err == nil {
		err = table.Update(obj)
	}
	return
}

func (db *MockDB) Delete(tableName string, pk uint64) (err error) {
	if table, err := db.getTable(tableName); err == nil {
		err = table.Delete(pk)
	}
	return
}
