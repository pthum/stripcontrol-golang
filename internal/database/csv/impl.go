package csv

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/go-co-op/gocron"
	"github.com/gocarina/gocsv"
	"github.com/pthum/stripcontrol-golang/internal/config"
	"github.com/pthum/stripcontrol-golang/internal/database"
	alog "github.com/pthum/stripcontrol-golang/internal/log"
	"github.com/pthum/stripcontrol-golang/internal/model"
	"github.com/samber/do"
)

// interface guard
var _ database.DBHandler[any] = (*CSVHandler[any])(nil)

type CSVHandler[T any] struct {
	cfg           *config.CSVConfig
	iMap          *SyncMap[string, T]
	lastCheckHash string
	l             alog.Logger
}

func NewHandler[T any](cfg *config.CSVConfig) *CSVHandler[T] {
	ch := &CSVHandler[T]{
		cfg:  cfg,
		iMap: NewSyncMap[string, T](),
		l:    alog.NewLogger("csvhandler"),
	}
	ch.load()
	return ch
}
func NewHandlerI[T any](inj *do.Injector) (database.DBHandler[T], error) {
	s := do.MustInvoke[*gocron.Scheduler](inj)
	acfg := do.MustInvoke[*config.Config](inj)
	ch := &CSVHandler[T]{
		cfg:  &acfg.CSV,
		iMap: NewSyncMap[string, T](),
		l:    alog.NewLogger("csvhandler"),
	}
	ch.load()
	ch.ScheduleJob(s)
	return ch, nil
}

func (c *CSVHandler[T]) GetAll() ([]T, error) {
	objs := c.iMap.LoadAll()

	sort.SliceStable(objs, func(i, j int) bool {
		a := c.asIDer(&objs[i])
		b := c.asIDer(&objs[j])
		if a == nil || b == nil {
			return false
		}
		return a.GetID() < b.GetID()
	})
	return objs, nil
}

func (c *CSVHandler[T]) Get(id string) (*T, error) {
	obj, ok := c.iMap.Load(id)
	if !ok {
		return nil, errors.New("object not found")
	}
	return &obj, nil
}

func (c *CSVHandler[T]) Save(input *T) (err error) {
	id := c.findId(input)
	c.iMap.Store(id, *input)
	return nil
}

func (c *CSVHandler[T]) Update(dbObject T, input T) (err error) {
	// only fullupdate atm
	return c.Save(&input)
}

func (c *CSVHandler[T]) Create(input *T) (err error) {
	return c.Save(input)
}

func (c *CSVHandler[T]) Delete(input *T) (err error) {
	id := c.findId(input)
	c.iMap.Delete(id)
	return nil
}

func (c *CSVHandler[T]) Close() {
	// nothing to close
}

func (c *CSVHandler[T]) Shutdown() error {
	return nil
}

func (c *CSVHandler[T]) ScheduleJob(s *gocron.Scheduler) {
	if c.cfg.DataDir == "" {
		// do nothing if no data dir given
		return
	}
	name := c.tableName()
	c.l.Info("Scheduling job for %v with interval of %v min", name, c.cfg.Interval)

	_, err := s.Every(c.cfg.Interval).Minutes().Tag(name).Do(c.persistIfNecessary)
	if err != nil {
		// handle the error related to setting up the job
		c.l.Error("error scheduling the %v job: %s", name, err.Error())
	}
}
func (c *CSVHandler[T]) findId(input any) string {
	ider := c.asIDer(input)
	if ider == nil {
		return ""
	}
	return strconv.FormatInt(ider.GetID(), 10)
}

func (c *CSVHandler[T]) tableName() string {
	var dummy T
	ider := c.asIDer(&dummy)
	if ider == nil {
		return ""
	}
	return ider.TableName()
}

func (c *CSVHandler[T]) asIDer(input any) model.IDer {
	if input == nil {
		return nil
	}
	ider, ok := input.(model.IDer)
	if !ok {
		return nil
	}
	return ider
}

func (c *CSVHandler[T]) openFile() (*os.File, error) {
	if c.cfg.DataDir == "" {
		return nil, errors.New("no datadir")
	}
	fName := c.tableName() + ".csv"
	return os.OpenFile(c.cfg.DataDir+fName, os.O_RDWR|os.O_CREATE, os.ModePerm)
}

func (c *CSVHandler[T]) load() {
	elems := []T{}
	if c.cfg.DataDir != "" {
		dataFile, err := c.openFile()
		if err != nil {
			panic(err)
		}
		defer dataFile.Close()

		if err := gocsv.UnmarshalFile(dataFile, &elems); err != nil { // Load elements from file
			if gocsv.ErrEmptyCSVFile != err {
				panic(err)
			}
		}
	} else {
		c.l.Warn("No data dir given, skip loading existing data")
	}

	for i := range elems {
		if err := c.Save(&elems[i]); err != nil {
			c.l.Error("error: %s\n", err.Error())
		}
	}
	var err error
	if c.lastCheckHash, err = c.hashEntries(); err != nil {
		panic(err)
	}
}

func (c *CSVHandler[T]) persistIfNecessary() {
	tName := c.tableName()
	c.l.Info("Running job for " + tName)
	currentHash, err := c.hashEntries()
	if err != nil {
		c.l.Error("error calculating the hash in job %v: %s\n", tName, err.Error())
		return
	}
	if strings.EqualFold(currentHash, c.lastCheckHash) {
		c.l.Info("Hashes are equal for job %v, skip writing", tName)
		return
	}
	err = c.persist()
	if err != nil {
		c.l.Error("error persisting updates for %v: %s\n", tName, err.Error())
		return
	}
	c.lastCheckHash = currentHash
}

func (c *CSVHandler[T]) hashEntries() (string, error) {
	models, err := c.GetAll()
	if err != nil {
		return "", err
	}

	content, err := gocsv.MarshalStringWithoutHeaders(&models)
	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(content))
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum, nil
}

func (c *CSVHandler[T]) persist() error {
	dataFile, err := c.openFile()
	if err != nil {
		return err
	}
	defer dataFile.Close()

	models, err := c.GetAll()
	if err != nil {
		return err
	}

	err = gocsv.MarshalFile(models, dataFile)
	if err != nil {
		return err
	}
	return nil
}
