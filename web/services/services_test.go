package services

import (
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDatabase() *gorm.DB {
	// TODO: refactor this in a common infrastructure init package
	viper.AutomaticEnv()
	viper.SetDefault("dbhost", "localhost")
	viper.SetDefault("dbport", "32432")
	viper.SetDefault("dbuser", "postgres")
	viper.SetDefault("dbpassword", "postgres")
	viper.SetDefault("dbname", "trento_test")

	host := viper.GetString("dbhost")
	port := viper.GetString("dbport")
	user := viper.GetString("dbuser")
	password := viper.GetString("dbpassword")
	dbname := viper.GetString("dbname")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}
