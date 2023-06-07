package database

import (
	"fmt"
	"log"
	"reflect"
	"sync"

	"github.com/pthum/stripcontrol-golang/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// interface guard
var _ DBHandler[any] = (*GeneralDbHandler[any])(nil)

var once sync.Once
var adb *gorm.DB

type DBReader[T any] interface {
	GetAll() ([]T, error)
	Get(id string) (*T, error)
	Close()
}

type DBWriter[T any] interface {
	Save(input *T) (err error)
	Update(dbObject T, input T) (err error)
	Create(input *T) (err error)
	Delete(input *T) (err error)
	Close()
}

//go:generate mockery --name=DBHandler --with-expecter=true
type DBHandler[T any] interface {
	DBReader[T]
	DBWriter[T]
}

type GeneralDbHandler[T any] struct {
	db *gorm.DB
}

// DB database object

func New[T any](cfg config.DatabaseConfig) DBHandler[T] {
	dbh := &GeneralDbHandler[T]{}
	dbh.connect(cfg)
	return dbh
}

// Connect set up the connection to the database
func (d *GeneralDbHandler[T]) connect(cfg config.DatabaseConfig) {
	once.Do(func() {
		var conn gorm.Dialector
		var err error
		configString := cfg.Host
		log.Printf("Setup database with %s", configString)
		conn = sqlite.Open(configString)
		adb, err = gorm.Open(conn, &gorm.Config{})

		if err != nil {
			panic("Failed to connect to database!")
		}
	})

	d.db = adb
}

// Close closes the db connection
func (d *GeneralDbHandler[T]) Close() {
	// Get generic database object sql.DB to be able to close
	sqlDB, err := d.db.DB()
	if err != nil {
		log.Printf("err closing db connection: %s", err.Error())
	}

	if err = sqlDB.Close(); err != nil {
		log.Printf("err closing db connection: %s", err.Error())
	} else {
		log.Println("db connection gracefully closed")
	}
}

// GetAll get all objects
func (d *GeneralDbHandler[T]) GetAll() (trg []T, err error) {
	r := d.db.Find(&trg)
	return trg, r.Error
}

// Get loads an object from the database
func (d *GeneralDbHandler[T]) Get(ID string) (*T, error) {
	var trg T
	r := d.db.Where("id = ?", ID).First(&trg)
	return &trg, r.Error
}

// Update updates the object
func (d *GeneralDbHandler[T]) Update(dbObject T, input T) (err error) {
	// calculate the difference, as gorm seem to update too much fields
	fields := findPartialUpdateFields(dbObject, input)
	return d.db.Model(dbObject).Debug().Select(fields).Updates(input).Error
}

func (d *GeneralDbHandler[T]) Create(input *T) (err error) {
	return d.db.Create(input).Error
}

func (d *GeneralDbHandler[T]) Save(input *T) (err error) {
	return d.db.Debug().Save(input).Error
}

func (d *GeneralDbHandler[T]) Delete(input *T) (err error) {
	return d.db.Delete(input).Error
}

// findPartialUpdateFields find the fields that need to be updated
func findPartialUpdateFields[T any](dbObject T, input T) (fields []string) {
	tIn := reflect.TypeOf(input)
	tDb := reflect.TypeOf(dbObject)
	if tIn.Kind() != tDb.Kind() || tIn != tDb {
		log.Println("different kinds, no update")
		return
	}
	valIn := reflect.ValueOf(input)
	valDb := reflect.ValueOf(dbObject)

	for i := 0; i < valIn.NumField(); i++ {
		valueFieldIn := valIn.Field(i)
		valueFieldDb := valDb.Field(i)
		typeField := valIn.Type().Field(i)
		if valueFieldIn.Interface() != valueFieldDb.Interface() {
			fields = append(fields, typeField.Name)
		}
	}

	fmt.Printf("Fields to update: %v\n", fields)
	return
}
