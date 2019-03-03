package db

import (
	"strings"
	"testing"
)

// disclaimer: tests are only for single-table because this is mock and by the
// time i'll use 2 or more i'll probably switch to an actual database

// helpers
type testObject struct {
	id  uint64
	val string
}

func (obj *testObject) GetPK() uint64 {
	return obj.id
}

func (obj *testObject) setPK(pk uint64) {
	obj.id = pk
}

func (obj *testObject) GetTableName() string {
	return "test"
}

func NewTestObject(str string) *testObject {
	return &testObject{0, str}
}

func getTestData() []*testObject {
	return []*testObject{
		NewTestObject("foo"),
		NewTestObject("bar"),
		NewTestObject("foobar"),
		NewTestObject("barfoo"),
		NewTestObject("ham"),
		NewTestObject("spam"),
	}
}

func TestNewMockDB(t *testing.T) {
	testTables := []string{"test", "foo", "bar"}
	db := NewMockDB(testTables...)
	for _, tableName := range testTables {
		if _, ok := (*db)[tableName]; !ok {
			t.Errorf("table excpected but absent: %s", tableName)
		}
	}
}

func TestMockDb_StoreLoad(t *testing.T) {
	objects := getTestData()
	db := NewMockDB("test")

	for _, obj := range objects {
		if err := db.Store(obj); err != nil {
			t.Error(err)
		}
	}

	for _, obj := range objects {
		if loaded, err := db.Load(obj.GetTableName(), obj.GetPK()); err != nil {
			t.Error(err)
		} else if loaded.(*testObject).val != obj.val {
			t.Errorf("Loaded object differs from the stored one: %s != %s",
				obj.val, loaded.(*testObject).val)
		}
	}
}

func TestMockDb_StoreUpdateLoad(t *testing.T) {
	objects := getTestData()
	db := NewMockDB("test")

	for _, obj := range objects {
		if err := db.Store(obj); err != nil {
			t.Error(err)
		}
	}

	for _, obj := range objects {
		obj.val = obj.val + obj.val
		if err := db.Update(obj); err != nil {
			t.Error(err)
		}
	}

	for _, obj := range objects {
		if loaded, err := db.Load(obj.GetTableName(), obj.GetPK()); err != nil {
			t.Error(err)
		} else if loaded.(*testObject).val != obj.val {
			t.Errorf("Loaded object differs from the updated one: %s != %s",
				obj.val, loaded.(*testObject).val)
		}
	}
}

func TestMockDb_StoreDeleteLoad(t *testing.T) {
	objects := getTestData()
	db := NewMockDB("test")

	deleteFilter := func(obj *testObject) bool {
		return strings.Contains(obj.val, "foo")
	}

	for _, obj := range objects {
		if err := db.Store(obj); err != nil {
			t.Error(err)
		}
	}

	for _, obj := range objects {
		if deleteFilter(obj) {
			if err := db.Delete(obj.GetTableName(), obj.GetPK()); err != nil {
				t.Error(err)
			}
		}
	}

	for _, obj := range objects {
		loaded, err := db.Load(obj.GetTableName(), obj.GetPK());
		if err == nil {
			if deleteFilter(obj) {
				t.Errorf("object not deleted: %s", obj.val)
			} else if obj.val != loaded.(*testObject).val {
				t.Errorf("Loaded object differs from the updated one: %s != %s",
					obj.val, loaded.(*testObject).val)
			}
		} else {
			if !deleteFilter(obj) {
				t.Errorf("object not found: %s: %s", obj.val, err)
			} else {
				switch err.(type) {
				case *ObjectNotFoundError:
				default:
					t.Errorf("unexptected error when requesting a deleted object: %s", err)
				}
			}
		}
	}
}
