package streetsweeping

import (
	"database/sql"
	"log"
)

// FindReadyAlerts finds all alerts that are ready to be sent (that is, that has a "next call" that is before now),
// and sends a text message reminder to those alerts.
func FindReadyAlerts(sender smsMessager) {
	findReadyAlertStmt, err := DB.Prepare("select ID, PHONE_NUMBER, NTH_DAY, TIMEZONE, WEEKDAY from alerts where NEXT_CALL < ?")
	if err != nil {
		log.Println("In FindReadyAlerts, problem preparing database statement: ", err)
	}
	defer findReadyAlertStmt.Close()

	updateStmt, err := DB.Prepare("UPDATE alerts SET NEXT_CALL = ? WHERE ID = ?")
	if err != nil {
		log.Println("In FindReadyAlerts, problem preparing database update statement: ", err)
	}
	defer updateStmt.Close()

	nowUTC := Now().Unix()
	rows, err := findReadyAlertStmt.Query(nowUTC)
	if err != nil {
		log.Println("In FindReadyAlerts, problem executing statement: err", err)
	}

	defer rows.Close()
	for rows.Next() {
		var id int
		alert := alert{}
		day := day{}

		err := rows.Scan(&id, &alert.PhoneNumber, &day.NthWeek, &alert.Timezone, &day.Weekday)
		if err != nil {
			log.Println("problem scanning rows: err", err)
		}

		nextCall, err := CalculateNextCall(day.NthWeek, day.Weekday, alert.Timezone)
		if err != nil {
			log.Println("error calculating next call: err", err)
		}

		_, err = updateStmt.Exec(nextCall, id)
		if err != nil {
			log.Println("error exicuting update statement: err", err)
		}

		remind(alert.PhoneNumber, sender, id)
	}
	if err = rows.Err(); err != nil {
		log.Println("problem iterating through the rows: ", err)
	}
}

func save(alert alert) error {
	tx, err := DB.Begin()
	stmt, err := tx.Prepare("INSERT INTO alerts (PHONE_NUMBER, NTH_DAY, TIMEZONE, WEEKDAY, NEXT_CALL, COUNTRY_CODE) VALUES (?,?,?,?,?,1)")
	if err != nil {
		log.Println("problem preparing transaction", err)
		err := tx.Rollback()
		return err
	}

	for _, t := range alert.Times {
		nextCall, err := CalculateNextCall(t.NthWeek, t.Weekday, alert.Timezone)
		log.Println("in save ..., next call: ", nextCall)

		if err != nil {
			log.Println("problem calculating next call: ", err)
			err := tx.Rollback()
			return err
		}
		result, err := stmt.Exec(alert.PhoneNumber, t.NthWeek, alert.Timezone, t.Weekday, nextCall)
		rowsAffected, _ := result.RowsAffected()
		lastInsertID, _ := result.LastInsertId()
		log.Println("new alert created: ", rowsAffected, lastInsertID)
		if err != nil {
			log.Println("problem exicuting statement: ", err)
			err := tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func StartDB(mysqlPassword string) *sql.DB {
	db, err := sql.Open("mysql", mysqlPassword)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	createTableCommand := `CREATE TABLE IF NOT EXISTS alerts(
				   ID INT NOT NULL AUTO_INCREMENT,
				   PHONE_NUMBER CHAR(10) NOT NULL,
				   NTH_DAY INT NOT NULL,
				   TIMEZONE VARCHAR(100) NOT NULL,
				   WEEKDAY VARCHAR(20) NOT NULL,
				   NEXT_CALL BIGINT NOT NULL,
				   PRIMARY KEY  (ID)
				)`
	_, err = db.Exec(createTableCommand)

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func removeAlerts(alert removeAlert) error {
	stmt, err := DB.Prepare("DELETE FROM alerts WHERE PHONE_NUMBER = ?;")
	if err != nil {
		return err
	}

	res, err := stmt.Exec(alert.PhoneNumber)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	log.Println("rows affected: ", affected)
	return nil
}
