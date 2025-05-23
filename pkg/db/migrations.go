package db

import (
	models2 "Ecadr/internal/app/models"
	"Ecadr/internal/app/models/seeds"
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
		&models2.Company{},
		&models2.Achievement{},
		&models2.SchoolGrade{},
		&models2.KindergartenNote{},
		&models2.Vacancy{},
		&models2.Criteria{},
		&models2.Course{},
		&models2.Recommend{},
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

	err = seeds.SeedRoles(dbConn)
	if err != nil {
		return err
	}

	err = seeds.SeedUsers(dbConn)
	if err != nil {
		return err
	}

	err = seeds.SeedCompany(dbConn)
	if err != nil {
		return err
	}

	err = seeds.SeedKinderGarten(dbConn)
	if err != nil {
		return err
	}

	err = seeds.SeedGrades(dbConn)
	if err != nil {
		return err
	}

	err = seeds.SeedDastovar(dbConn)
	if err != nil {
		return err
	}

	err = seeds.SeedVacancy(dbConn)
	if err != nil {
		return err
	}

	err = seeds.SeedCourse(dbConn)
	if err != nil {
		return err
	}

	err = seeds.SeedCriteria(dbConn)
	if err != nil {
		return err
	}

	err = seeds.SeedTests(dbConn)
	if err != nil {
		return err
	}

	err = seeds.SeedTestSubmissions(dbConn)
	if err != nil {
		return err
	}

	return nil
}
