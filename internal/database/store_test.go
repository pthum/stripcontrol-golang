package database

import (
	"testing"

	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestDBCreateAndRead(t *testing.T) {
	dbh := initHandler()
	dbh.db.AutoMigrate(&model.ColorProfile{})
	defer dbh.Close()

	testProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 123},
		Blue:       null.IntFrom(1),
		Brightness: null.IntFrom(2),
		Red:        null.IntFrom(3),
		Green:      null.IntFrom(4),
	}

	dbh.Create(&testProfile)
	var result model.ColorProfile
	dbh.Get("123", &result)
	assert.Equal(t, testProfile, result)
}

func TestDBSaveGetAllAndDelete(t *testing.T) {
	dbh := initHandler()
	dbh.db.AutoMigrate(&model.ColorProfile{})
	defer dbh.Close()

	testProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 234},
		Blue:       null.IntFrom(1),
		Brightness: null.IntFrom(2),
		Red:        null.IntFrom(3),
		Green:      null.IntFrom(4),
	}

	dbh.Create(&testProfile)
	testProfile.Blue = null.IntFrom(42)
	dbh.Save(&testProfile)
	var result []model.ColorProfile
	dbh.GetAll(&result)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, testProfile, result[0])

	dbh.Delete(&testProfile)

	var resultAfterDelete []model.ColorProfile
	dbh.GetAll(&resultAfterDelete)
	assert.Equal(t, 0, len(resultAfterDelete))
}
func TestUpdate(t *testing.T) {
	dbh := initHandler()
	dbh.db.AutoMigrate(&model.ColorProfile{})
	defer dbh.Close()
	testProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 235},
		Blue:       null.IntFrom(1),
		Brightness: null.IntFrom(2),
		Red:        null.IntFrom(3),
		Green:      null.IntFrom(4),
	}
	dbh.Create(&testProfile)
	// copy the profile
	otherProfile := testProfile
	otherProfile.Blue = null.IntFrom(42)
	// test with changes
	dbh.Update(testProfile, otherProfile)
	var result model.ColorProfile
	dbh.Get("235", &result)
	assert.Equal(t, otherProfile, result)
}
func TestUpdateFields_Empty(t *testing.T) {
	testProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 234},
		Blue:       null.IntFrom(1),
		Brightness: null.IntFrom(2),
		Red:        null.IntFrom(3),
		Green:      null.IntFrom(4),
	}
	// copy the profile
	otherProfile := testProfile
	// test with no changes
	changedFields := findPartialUpdateFields(testProfile, otherProfile)
	assert.Equal(t, 0, len(changedFields))
}

func TestUpdateFields_WithChanges(t *testing.T) {
	testProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 235},
		Blue:       null.IntFrom(1),
		Brightness: null.IntFrom(2),
		Red:        null.IntFrom(3),
		Green:      null.IntFrom(4),
	}
	// copy the profile
	otherProfile := testProfile
	otherProfile.Blue = null.IntFrom(42)
	// test with changes
	changedFields := findPartialUpdateFields(testProfile, otherProfile)
	assert.Equal(t, 1, len(changedFields))
	assertContains(t, changedFields, "Blue")
}

func TestUpdateFields_DifferentTypes(t *testing.T) {
	testProfile := model.ColorProfile{
		BaseModel:  model.BaseModel{ID: 235},
		Blue:       null.IntFrom(1),
		Brightness: null.IntFrom(2),
		Red:        null.IntFrom(3),
		Green:      null.IntFrom(4),
	}
	// have a different type
	otherProfile := model.LedStrip{}
	// there should be no changes
	changedFields := findPartialUpdateFields(testProfile, otherProfile)
	assert.Equal(t, 0, len(changedFields))
}

func assertContains(t *testing.T, s []string, e string) {
	for _, a := range s {
		if a == e {
			return
		}
	}
	t.Errorf("expected %v to contain %v", s, e)
}

func initHandler() *GeneralDbHandler {
	dbConf := config.DatabaseConfig{
		Host: ":memory:",
	}
	return New(dbConf).(*GeneralDbHandler)
}
