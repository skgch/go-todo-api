package infrastructure

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var Db *gorm.DB

func Connect() {
	DBMS := "mysql"
	USER := "root"
	PROTOCOL := "tcp(127.0.0.1:3306)"
	DBNAME := "mydb"
	CONNECT := USER + "@" + PROTOCOL + "/" + DBNAME

	var err error
	Db, err = gorm.Open(DBMS, CONNECT)
	if err != nil {
		panic(err.Error())
	}
}

func Close() {
	Db.Close()
}
