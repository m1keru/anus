package models

import "database/sql"

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCollection struct {
	Users []User `json:"items"`
}

func GetUsers(db *sql.DB) (UserCollection, error) {
	sql := "SELECT * from users"
	rows, err := db.Query(sql)
	if err != nil {
		return UserCollection{}, err
	}
	defer rows.Close()
	result := UserCollection{}
	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.ID, &user.Login, &user.Email, &user.Password)
		if err != nil {
			return UserCollection{}, err
		}
		result.Users = append(result.Users, user)
	}
	return result, nil
}

func GetUser(username string, db *sql.DB) (User, error) {
	sql := "SELECT * FROM users WHERE login=?"
	stmt, err := db.Prepare(sql)
	if err != nil {
		return User{}, err
	}
	result, err := stmt.Query(username)
	if err != nil {
		return User{}, err
	}
	user := User{}
	for result.Next() {
		err = result.Scan(&user.ID, &user.Login, &user.Email, &user.Password)
		if err != nil {
			return User{}, err
		}
	}
	return user, nil
}

func PutUsers(db *sql.DB, user User) error {
	//Пока нет мыслей зачем это понадобится
	return nil
}

//DeleteAnsibleScripts -- delete scripts from db
func DeleteUsers(db *sql.DB, user User) error {
	//Пока нет мыслей зачем это понадобится
	return nil
}
