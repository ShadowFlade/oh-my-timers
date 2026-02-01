package db

import (
	"fmt"
	"log"
	"strings"

	"shadowflade/timers/pkg/global"
	"shadowflade/timers/pkg/interfaces"
)

type User struct {
	TableName string
}

func (this *User) Create() *User {
	return &User{
		TableName: "users",
	}
}

func (this *User) CreateUser(user interfaces.User) int64 {
	db := Db{}
	db.Connect()
	defer db.Db.Close()
	tx := db.Db.MustBegin()
	query := `
	insert into users (name, password) values (:name, :password)
	`
	fmt.Printf("creating user %s with password %s", user.Name, user.Password)

	res, err := db.Db.NamedExec(query, user)
	if err != nil {
		log.Fatalf("Error on creating user. Query: %s. Err: %s", query, err.Error())
	}
	newId, err := res.LastInsertId()
	fmt.Println("new id: " + fmt.Sprint(newId))
	if err != nil {
		panic(err.Error())
	}
	tx.Commit()
	return newId
}

func (this *User) FindUserByHashedPassword(hashedPassword string) (*interfaces.User, error) {
	db := Db{}
	db.Connect()
	defer db.Db.Close()
	fmt.Printf("Raw bytes: %v\n", []byte(hashedPassword))
	fmt.Printf("Raw string: %v\n", hashedPassword)
	fmt.Printf("DBDB: %v\n", db.Db)

	query := fmt.Sprintf("SELECT * FROM users WHERE password = '%v' ;", hashedPassword)
	hashedPassword = strings.TrimSpace(hashedPassword)
	fmt.Printf("Hex: %x\n", []byte(hashedPassword))
	fmt.Printf("🔍 Hash as quoted string: %q\n", hashedPassword)
	fmt.Printf("Input hash (quoted): %q\n", hashedPassword)
	fmt.Printf("Input length: %d\n", len(hashedPassword))
	global.Logger.LogText("hashed: "+hashedPassword, "")
	global.Logger.LogText(query, "")

	var u interfaces.User
	err := db.Db.Get(&u, "SELECT * FROM users WHERE password = ?", hashedPassword)

	if err != nil {
		fmt.Printf("user: %+v", u)

		return nil, err
	}

	var user interfaces.User = u

	return &user, nil
}
