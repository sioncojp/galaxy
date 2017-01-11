package galaxy

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Commits ...Commits table.
type Commits struct {
	ID        int64     `gorm:"column:id"`
	Number    string    `gorm:"column:number"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

// getConnectionString ...Output connection string for gorm.
func (config *Config) getConnectionString() string {
	user := config.Database.User
	password := config.Database.Password
	host := config.Database.Host
	port := config.Database.Port
	dbname := config.Database.DBname

	return fmt.Sprintf("%s:%s@([%s]:%d)/%s?charset=utf8&parseTime=True", user, password, host, port, dbname)
}

// DBConnect ...Connection DB.
func (config *Config) DBConnect() (*gorm.DB, error) {
	db, err := gorm.Open(
		config.Database.Driver,
		config.getConnectionString(),
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}
