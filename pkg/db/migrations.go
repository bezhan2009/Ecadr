package db

import (
	models2 "Ecadr/internal/app/models"
	"Ecadr/pkg/logger"
	"errors"
)

func Migrate() error {
	if dbConn == nil {
		logger.Error.Printf("[db.Migrate] Error because database connection is nil")

		return errors.New("database connection is not initialized")
	}

	//if userDBConn == nil {
	//	logger.Error.Printf("[db.Migrate] Error because users database connection is nil")
	//
	//	return errors.New("users database connection is not initialized")
	//}
	//
	//err := userDBConn.AutoMigrate(
	//	&models2.User{},
	//	&models2.Admin{},
	//)
	//if err != nil {
	//	logger.Error.Printf("[db.Migrate] Error migrating users tables: %v", err)
	//
	//	return err
	//}

	err := dbConn.AutoMigrate(
		&models2.Role{},
		&models2.User{},
		&models2.Achievement{},
		&models2.SchoolGrade{},
		&models2.KindergartenNote{},
		&models2.Vacancy{},
		&models2.Course{},
		&models2.Test{},
		&models2.Question{},
		&models2.Choice{},
		&models2.TestSubmission{},
		&models2.Answer{},
		&models2.Message{},
	)
	if err != nil {
		logger.Error.Printf("[db.Migrate] Error migrating tables: %v", err)

		return err
	}
	//
	//err = seeds.SeedRoles(dbConn)
	//if err != nil {
	//	return err
	//}
	//
	//err = seeds.SeedUsers(dbConn)
	//if err != nil {
	//	return err
	//}
	//
	//err = seeds.SeedKinderGarten(dbConn)
	//if err != nil {
	//	return err
	//}
	//
	//err = seeds.SeedGrades(dbConn)
	//if err != nil {
	//	return err
	//}
	//
	//err = seeds.SeedDastovar(dbConn)
	//if err != nil {
	//	return err
	//}
	//
	//err = seeds.SeedVacansy(dbConn)
	//if err != nil {
	//	return err
	//}
	//
	//err = seeds.SeedCourse(dbConn)
	//if err != nil {
	//	return err
	//}

	return nil
}
