package mysql

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

type demoEntity struct {
	ID       *int64 `column:"id"`
	AddTime  *time.Time
	EditTime *time.Time

	Name   string
	Status int
}

// The init() function in client_test.go will initialize ctx and client variables
// This file depends on client_test.go being included in the test run

func TestInsertEntities(t *testing.T) {
	cnt := 3
	entities := make([]interface{}, cnt)
	now := time.Now()
	for i := 0; i < cnt; i++ {
		entities[i] = &demoEntity{
			AddTime: &now,
			Name:    "demo" + strconv.Itoa(i),
			Status:  0,
		}
	}

	err := entityDao().InsertEntities(ctx, "demo", entities...)
	t.Log(err)
}

func TestSelectEntityByID(t *testing.T) {
	entity := new(demoEntity)
	err := entityDao().SelectEntityByID(ctx, "demo", 58, entity)
	t.Log(err, entity, NoRowsError(err))
	if err == nil {
		t.Log(*entity.ID, *entity.AddTime, *entity.EditTime, entity)
	}
}

func TestSimpleQueryEntityAnd(t *testing.T) {
	entity := new(demoEntity)
	condItems := []*SqlColQueryItem{
		{"name", SqlCondEqual, "demo", false},
	}
	err := entityDao().SimpleQueryEntityAnd(ctx, "demo", entity, condItems...)
	t.Log(err, NoRowsError(err))
	if err == nil {
		t.Log(*entity.ID, *entity.AddTime, *entity.EditTime, entity)
	}
}

func TestSimpleQueryEntitiesAnd(t *testing.T) {
	var entityList []*demoEntity
	condItems := []*SqlColQueryItem{
		{"name", SqlCondEqual, "demo", false},
	}
	params := &SqlQueryParams{
		CondItems: condItems,
		OrderBy:   "id desc",
		Offset:    0,
		Cnt:       10,
	}
	err := entityDao().SimpleQueryEntitiesAnd(ctx, "demo", params, &entityList)
	t.Log(err, NoRowsError(err))
	for i, entity := range entityList {
		t.Log(i, entity, *entity.ID, *entity.AddTime, *entity.EditTime)
	}
}

func TestColumnNameByFieldQuoting(t *testing.T) {
	type Reserved struct {
		Key  string `column:"key"`
		Data int
	}
	st := reflect.TypeOf(Reserved{})
	fKey, _ := st.FieldByName("Key")
	if got := ColumnNameByField(&fKey); got != "`key`" {
		t.Errorf("ColumnNameByField for tag, want \"`key`\", got %s", got)
	}
	fData, _ := st.FieldByName("Data")
	if got := ColumnNameByField(&fData); got != "`data`" {
		t.Errorf("ColumnNameByField default, want \"`data`\", got %s", got)
	}
}

func entityDao() *EntityDao {
	return &EntityDao{Dao{client}}
}
