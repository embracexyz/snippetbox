package models

import (
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
	GetUser(id int) (*User, error)
}
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) GetUser(id int) (*User, error) {
	stmt := `SELECT id, name, email, created FROM users WHERE id = ?`

	u := &User{}
	err := m.DB.QueryRow(stmt, id).Scan(&u.ID, &u.Name, &u.Email, &u.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

func (m *UserModel) Insert(name, email, password string) error {
	HashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (name, email, hashed_password, created)
	VALUES(?, ?, ?, UTC_TIMESTAMP())`

	_, err = m.DB.Exec(stmt, name, email, string(HashedPassword))
	// 尝试match mysql 特定error，来判断是否出现了duplicate recrod的错误
	if err != nil {
		var myerr *mysql.MySQLError
		if errors.As(err, &myerr) && myerr.Number == 1062 {
			return ErrDuplicateEmail
		}
		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`
	var id int
	var hashedPassword []byte
	err := m.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exist bool
	stmt := `select exists(select true from users where id = ?)`
	err := m.DB.QueryRow(stmt, id).Scan(&exist)
	return exist, err
}
