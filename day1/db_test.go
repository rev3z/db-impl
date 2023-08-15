package day1

import (
	"fmt"
	"testing"
)

func TestDb(t *testing.T) {
	db := CreateDb()
	db.Insert(Row{Name: "Jonh", Time: 22})
	db.Insert(Row{Name: "Mike", Time: 28})
	db.Insert(Row{Name: "Marry", Time: 13})
	result := db.QueryById("2")
	if len(result) != 1 {
		t.Errorf("QueryById return %d results which must be 1!", len(result))
	}
	if orig := (Row{"2", "Mike", 28}); result[0] != orig {
		t.Errorf("QueryById return the false record! Get %v expect %v.", result[0], orig)
	}
	db.Insert(Row{"2", "Leon", 27})
	result = db.QueryById("2")
	if len(result) != 1 {
		t.Errorf("QueryById return %d results which must be 1!", len(result))
	}
	if orig := (Row{"2", "Leon", 27}); result[0] != orig {
		t.Errorf("QueryById return the false record! Get %v expect %v.", result[0], orig)
	}
	fmt.Printf("QueryByTimeArea %v", db.QueryByTimeArea(13, 28))
	db.Remove("2")
	result = db.QueryById("2")
	if len(result) != 0 {
		t.Errorf("QueryById return %d results which must be 0!", len(result))
	}
	fmt.Printf("QueryByTimeArea %v", db.QueryByTimeArea(13, 29))
	result = db.QueryById("4")
	if len(result) != 0 {
		t.Errorf("QueryById return %d results which must be 0!", len(result))
	}
	db.Remove("4")
}
