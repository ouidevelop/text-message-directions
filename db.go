package main

//import (
//	"database/sql"
//	"fmt"
//	"log"
//)
//
//func startDB(mysqlPassword string) *sql.DB {
//	db, err := sql.Open("mysql", mysqlPassword)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	err = db.Ping()
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	createTableCommand := `CREATE TABLE IF NOT EXISTS users(
//				   ID INT NOT NULL AUTO_INCREMENT,
//				   PHONE_NUMBER CHAR(20) NOT NULL,
//				   DIRECTIONS_USED INT NOT NULL,
//				   DIRECTIONS_LEFT VARCHAR(100) NOT NULL,
//				   CREATED DATETIME NOT NULL,
//				   PRIMARY KEY  (ID)
//				)`
//	_, err = db.Exec(createTableCommand)
//
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	return db
//}
//
//func addUser(phoneNumber string) error {
//	stmt, err := DB.Prepare("DELETE FROM alerts WHERE PHONE_NUMBER = ?;")
//	if err != nil {
//		return err
//	}
//
//	res, err := stmt.Exec(alert.PhoneNumber)
//	if err != nil {
//		return err
//	}
//	affected, err := res.RowsAffected()
//	if err != nil {
//		return err
//	}
//	fmt.Println("rows affected: ", affected)
//	return nil
//}


