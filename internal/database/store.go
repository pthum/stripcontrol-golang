package database

import (
	"fmt"
	"log"
	"reflect"

	"github.com/pthum/stripcontrol-golang/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBReader interface {
	GetAll(dest interface{}) error
	Get(id string, obj interface{}) error
	Close()
}

type DBWriter interface {
	Save(input interface{}) (err error)
	Update(dbObject interface{}, input interface{}) (err error)
	Create(input interface{}) (err error)
	Delete(input interface{}) (err error)
	Close()
}

//go:generate mockery --name=DBHandler --with-expecter=true
type DBHandler interface {
	DBReader
	DBWriter
}

type GeneralDbHandler struct {
	db *gorm.DB
}

// DB database object

func New(cfg config.DatabaseConfig) DBHandler {
	dbh := &GeneralDbHandler{}
	dbh.connect(cfg)
	return dbh
}

// Connect set up the connection to the database
func (d *GeneralDbHandler) connect(cfg config.DatabaseConfig) {
	var conn gorm.Dialector
	configString := cfg.Host
	log.Printf("Setup database with %s", configString)
	conn = sqlite.Open(configString)

	db, err := gorm.Open(conn, &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	d.db = db
}

// Close closes the db connection
func (d *GeneralDbHandler) Close() {
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
func (d *GeneralDbHandler) GetAll(targets interface{}) (err error) {
	return d.db.Find(targets).Error
}

// Get loads an object from the database
func (d *GeneralDbHandler) Get(ID string, target interface{}) (err error) {
	return d.db.Where("id = ?", ID).First(&target).Error
}

// Update updates the object
func (d *GeneralDbHandler) Update(dbObject interface{}, input interface{}) (err error) {
	// calculate the difference, as gorm seem to update too much fields
	fields := findPartialUpdateFields(dbObject, input)
	return d.db.Model(dbObject).Debug().Select(fields).Updates(input).Error
}

func (d *GeneralDbHandler) Create(input interface{}) (err error) {
	return d.db.Create(input).Error
}

func (d *GeneralDbHandler) Save(input interface{}) (err error) {
	return d.db.Debug().Save(input).Error
}

func (d *GeneralDbHandler) Delete(input interface{}) (err error) {
	return d.db.Delete(input).Error
}

// findPartialUpdateFields find the fields that need to be updated
func findPartialUpdateFields(dbObject interface{}, input interface{}) (fields []string) {
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
