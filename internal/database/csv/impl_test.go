package csv

import (
	"testing"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/pthum/null"
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestCreateAndRead(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	testProfile := createTestProfile(123)

	dbh.Create(&testProfile)
	result, err := dbh.Get("123")
	assert.Equal(t, testProfile, *result)
	assert.NoError(t, err)
}
func TestGetMissing(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	result, err := dbh.Get("123")
	assert.Nil(t, result)
	assert.Error(t, err)
}

func TestSaveGetAllAndDelete(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	testProfile := createTestProfile(242)
	dbh.Create(&testProfile)

	testProfile.Blue = null.IntFrom(42)
	dbh.Save(&testProfile)

	testProfile2 := createTestProfile(23)
	dbh.Save(&testProfile2)

	result, err := dbh.GetAll()

	assert.Equal(t, 2, len(result))
	// order of GetAll should be stable, by id
	assert.Equal(t, testProfile2, result[0])
	assert.Equal(t, testProfile, result[1])
	assert.NoError(t, err)

	dbh.Delete(&testProfile)

	resultAfterDelete, err := dbh.GetAll()
	assert.Equal(t, 1, len(resultAfterDelete))
	assert.Equal(t, testProfile2, resultAfterDelete[0])
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

func TestTableName(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	tn := dbh.tableName()
	assert.Equal(t, model.Table_ColorProfile, tn)
}

func TestOpenFile_MissingDataDir(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	dbh.cfg.DataDir = ""
	file, err := dbh.openFile()
	assert.Nil(t, file)
	assert.Error(t, err)
}

func TestOpenFile_Empty(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	dbh.cfg.DataDir = t.TempDir()
	file, err := dbh.openFile()
	assert.NotNil(t, file)
	assert.NoError(t, err)
}

func TestPersistAndLoad(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	dbh.cfg.DataDir = t.TempDir()

	// store a profile
	testProfile := createTestProfile(235)
	dbh.Create(&testProfile)
	initialHash := dbh.lastCheckHash

	// first persist
	dbh.persistIfNecessary()
	// hash should have changed
	hashAfterFirstSave := dbh.lastCheckHash
	assert.NotEqual(t, initialHash, hashAfterFirstSave)
	// second persist, without changes
	dbh.persistIfNecessary()
	// hashes should be equal
	hashAfterSecondSave := dbh.lastCheckHash
	assert.Equal(t, hashAfterFirstSave, hashAfterSecondSave)

	// cleanup map to load freshly from file
	all, _ := dbh.GetAll()
	for _, a := range all {
		dbh.Delete(&a)
	}
	// consistency check: map should be empty
	all, _ = dbh.GetAll()
	assert.Len(t, all, 0)

	dbh.load()
	// re-check that element has been loaded from file
	all, _ = dbh.GetAll()
	assert.Len(t, all, 1)
	assert.Equal(t, testProfile, all[0])
}

func TestLoadEmptyFile(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	dbh.cfg.DataDir = t.TempDir()
	// loading an empty/nonexisting file shouldn't panic
	dbh.load()
}

func TestScheduleJob_MissingDataDir(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	dbh.cfg.DataDir = ""
	s := gocron.NewScheduler(time.UTC)
	dbh.ScheduleJob(s)
	jobs := s.Jobs()
	// no jobs scheduled if dir is empty
	assert.Len(t, jobs, 0)
}
func TestScheduleJob_MissingInterval(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	dbh.cfg.DataDir = t.TempDir()
	dbh.cfg.Interval = 0
	s := gocron.NewScheduler(time.UTC)
	dbh.ScheduleJob(s)
	jobs := s.Jobs()
	// no jobs scheduled if interval is 0
	assert.Len(t, jobs, 0)
}

func TestScheduleJob(t *testing.T) {
	dbh := initHandler[model.ColorProfile](t)
	dbh.cfg.DataDir = t.TempDir()
	dbh.cfg.Interval = 10
	s := gocron.NewScheduler(time.UTC)
	dbh.ScheduleJob(s)
	jobs := s.Jobs()

	assert.Len(t, jobs, 1)
	assert.Contains(t, jobs[0].Tags(), dbh.tableName())
}

func initHandler[T any](t *testing.T) *CSVHandler[T] {
	dbConf := config.Config{}
	dbh := NewHandler[T](&dbConf.CSV)
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
