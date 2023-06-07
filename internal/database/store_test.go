package database

import (
	"testing"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestDBCreateAndRead(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	testProfile := createTestProfile(123)

	dbh.Create(&testProfile)
	result, err := dbh.Get("123")
	assert.Equal(t, testProfile, *result)
	assert.NoError(t, err)
}

func TestDBSaveGetAllAndDelete(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	testProfile := createTestProfile(234)

	dbh.Create(&testProfile)
	testProfile.Blue = null.IntFrom(42)
	dbh.Save(&testProfile)
	result, err := dbh.GetAll()
	assert.Equal(t, 1, len(result))
	assert.Equal(t, testProfile, result[0])
	assert.NoError(t, err)

	dbh.Delete(&testProfile)

	resultAfterDelete, err := dbh.GetAll()
	assert.Equal(t, 0, len(resultAfterDelete))
	assert.NoError(t, err)
}
func TestUpdate(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	testProfile := createTestProfile(235)
	dbh.Create(&testProfile)
	// copy the profile
	otherProfile := testProfile
	otherProfile.Blue = null.IntFrom(42)
	// test with changes
	dbh.Update(testProfile, otherProfile)
	result, err := dbh.Get("235")
	assert.Equal(t, otherProfile, *result)
	assert.NoError(t, err)
}
func TestUpdateFields_Empty(t *testing.T) {
	testProfile := createTestProfile(234)
	// copy the profile
	otherProfile := testProfile
	// test with no changes
	changedFields := findPartialUpdateFields(testProfile, otherProfile)
	assert.Equal(t, 0, len(changedFields))
}

func TestUpdateFields_WithChanges(t *testing.T) {
	testProfile := createTestProfile(235)
	// copy the profile
	otherProfile := testProfile
	otherProfile.Blue = null.IntFrom(42)
	// test with changes
	changedFields := findPartialUpdateFields(testProfile, otherProfile)
	assert.Equal(t, 1, len(changedFields))
	assertContains(t, changedFields, "Blue")
}

func assertContains(t *testing.T, s []string, e string) {
	for _, a := range s {
		if a == e {
			return
		}
	}
	t.Errorf("expected %v to contain %v", s, e)
}

func initHandler[T any](t *testing.T) *GeneralDbHandler[T] {
	dbConf := config.DatabaseConfig{
		Host: ":memory:",
	}
	dbh := New[T](dbConf).(*GeneralDbHandler[T])
	dbh.db.AutoMigrate(&model.ColorProfile{})
	t.Cleanup(func() {
		all, _ := dbh.GetAll()
		for _, a := range all {
			dbh.Delete(&a)
		}
	})
	return dbh
}

func createTestProfile(id int64) model.ColorProfile {
	return model.ColorProfile{
		BaseModel:  model.BaseModel{ID: id},
		Blue:       null.IntFrom(1),
		Brightness: null.IntFrom(2),
		Red:        null.IntFrom(3),
		Green:      null.IntFrom(4),
	}
}
