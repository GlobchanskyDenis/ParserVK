package main

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"strings"
	"strconv"
)

func ConnectDatabase(conf Sql) error {
	var dsn string
	var err error

	dsn = "user=" + conf.User + " password=" + conf.Pass + " dbname=" + conf.DBName + " host=" + conf.Host + " sslmode=disable"
	gDbConn, err = sql.Open(conf.DBType, dsn)
	if err != nil {
		return err
	}
	return gDbConn.Ping()
}

func (user *User) SaveInDatabase() error {
	var interests = "{" + strings.Join(user.Interests, ",") + "}"
	/*
	**	Transaction start
	 */
	tx, err := gDbConn.Begin()
	if err != nil {
		return errors.New("Database transaction error: " + err.Error())
	}
	defer tx.Rollback()
	/*
	**	Set new user
	 */
	stmt, err := tx.Prepare(`INSERT INTO users (mail, encryptedPass, fname, lname,
		birth, gender, orientation, bio, latitude, longitude, interests, status, search_visibility, rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING uid`)
	if err != nil {
		return errors.New("Database preparing error: " + err.Error())
	}
	defer stmt.Close()
	err = stmt.QueryRow(user.Mail, user.EncryptedPass, user.Fname, user.Lname,
		user.Birth, user.Gender, user.Orientation, user.Bio, user.Latitude, user.Longitude,
		interests, user.Status, true, user.Rating).Scan(&user.Uid)
	if err != nil {
		return errors.New("Database executing error: " + err.Error())
	}
	/*
	**	Set user should ignore himself in search queries
	 */
	stmt, err = tx.Prepare("INSERT INTO ignores (uidSender, uidReceiver) VALUES ($1, $1)")
	if err != nil {
		return errors.New("Database preparing error: " + err.Error())
	}
	result, err := stmt.Exec(user.Uid)
	if err != nil {
		return errors.New("Database executing error: " + err.Error())
	}
	// handle results
	nbr64, err := result.RowsAffected()
	if err != nil {
		return errors.New("Database executing error: " + err.Error())
	}
	if int(nbr64) == 0 {
		return errors.New("Database executing error: record not created")
	}
	if int(nbr64) != 1 {
		return errors.New("Database executing error: too many records was created")
	}
	/*
	**	Set new photo
	 */
	stmt, err = tx.Prepare("INSERT INTO photos (uid, src) VALUES ($1, $2) RETURNING pid")
	if err != nil {
		return errors.New("Database preparing error: " + err.Error())
	}
	err = stmt.QueryRow(user.Uid, user.Avatar).Scan(&user.AvaID)
	if err != nil {
		return errors.New("Database executing error: " + err.Error())
	}
	/*
	**	Update users ava ID
	 */
	stmt, err = tx.Prepare("UPDATE users SET avaId = $1 WHERE uid = $2")
	if err != nil {
		return errors.New("Database preparing error: " + err.Error())
	}
	result, err = stmt.Exec(user.AvaID, user.Uid)
	if err != nil {
		return errors.New("Database executing error: " + err.Error())
	}
	// handle results
	nbr64, err = result.RowsAffected()
	if err != nil {
		return errors.New("Database executing error: " + err.Error())
	}
	if int(nbr64) == 0 {
		return errors.New("Database executing error: record not updated")
	}
	if int(nbr64) != 1 {
		return errors.New("Database executing error: too many records was updated")
	}
	/*
	**	Close transaction
	 */
	err = tx.Commit()
	if err != nil {
		return errors.New("Database transaction error: " + err.Error())
	}
	return nil
}

func SaveInterestsToDB(unknownInterests []string) error {
	var query = "INSERT INTO interests (name) VALUES "
	var nameArr = []interface{}{}
	if len(unknownInterests) == 0 {
		return nil
	}
	for nbr, interest := range unknownInterests { ////// ПОХОЖЕ НА ГОВНОКОД. УЗНАТЬ ПОДРОБНЕЕ
		query += "($" + strconv.Itoa(nbr+1) + "), "
		nameArr = append(nameArr, interest) /// УБРАТЬ АЛЛОЦИРОВАНИЕ СЛАЙСА - ПРИНИМАТЬ СЛАЙС ИНТЕРФЕЙСОВ
	}
	query = string(query[:len(query)-2])
	stmt, err := gDbConn.Prepare(query)
	if err != nil {
		return errors.New("Database preparing error: " + err.Error())
	}
	defer stmt.Close()
	_, err = stmt.Exec(nameArr...)
	if err != nil {
		return errors.New("Database executing error: " + err.Error())
	}
	return nil
}
