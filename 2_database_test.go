package main

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func TestEmpty(t *testing.T) {
	fmt.Println(true)
}

func TestOpenConnection(t *testing.T) {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/rnd")
	if err != nil {
		panic(err)
	}
	defer db.Close()
}

func TestRead(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	rows, err := db.Query("select name, gender from persons order by id desc")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var name, gender string
		err = rows.Scan(&name, &gender)
		if err != nil {
			panic(err)
		}

		fmt.Println(name, gender)
	}
}

func TestCreate(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()
	defer db.Close()

	query := "insert into customers(id, name, email, phone) values('ricid', 'ricid', 'ricid@gmail.com', '083815530634')"
	_, err := db.ExecContext(ctx, query)

	if err != nil {
		panic(err)
	}

	fmt.Println("Done")
}

func TestQuerySql(t *testing.T) {
	db := GetConnection()
	ctx := context.Background()
	defer db.Close()

	query := "select id, name from customers"
	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		if err != nil {
			panic(err)
		}
		fmt.Println("ID :", id)
		fmt.Println("Name :", name)
	}

	defer rows.Close()
}

func TestQuerySqlComplex(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	script := "SELECT id, name, email, balance, rating, birth_date, married, created_at FROM customer"
	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id, name string
		var email sql.NullString
		var balance int32
		var rating float64
		var birthDate sql.NullTime
		var createdAt time.Time
		var married bool

		err = rows.Scan(&id, &name, &email, &balance, &rating, &birthDate, &married, &createdAt)
		if err != nil {
			panic(err)
		}

		fmt.Println("================")
		fmt.Println("Id:", id)
		fmt.Println("Name:", name)

		if email.Valid {
			fmt.Println("Email:", email.String)
		}

		fmt.Println("Balance:", balance)
		fmt.Println("Rating:", rating)

		if birthDate.Valid {
			fmt.Println("Birth Date:", birthDate.Time)
		}

		fmt.Println("Married:", married)
		fmt.Println("Created At:", createdAt)
	}
}

func TestSqlInjection(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "ricidkumbara@gmail.com'; #"
	password := "salah"

	script := "select email from users where email = '" + email + "' and password = '" + password + "'"
	fmt.Println(script)

	rows, err := db.QueryContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			panic(err)
		}

		fmt.Println("Sukses Login", email)
	} else {
		fmt.Println("Gagal Login")
	}
}

func TestSqlInjectionSafe(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	email := "ricidkumbara@gmail.com'; #"
	password := "salah"

	script := "select email from users where email = ? and password = ? limit 1"
	fmt.Println(script)

	rows, err := db.QueryContext(ctx, script, email, password)
	if err != nil {
		panic(err)
	}

	if rows.Next() {
		var email string
		err := rows.Scan(&email)
		if err != nil {
			panic(err)
		}

		fmt.Println("Sukses Login", email)
	} else {
		fmt.Println("Gagal Login")
	}
}

func TestExecSqlParameter(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()

	name := "Ricid"
	gender := "M"
	age := 25

	script := "insert into persons (name, gender, age) values (?, ?, ?)"
	result, err := db.ExecContext(ctx, script, name, gender, age)
	if err != nil {
		panic(err)
	}

	insertId, err := result.LastInsertId()
	if err != nil {
		panic(err)
	}

	fmt.Println("Success insert new comment with id", insertId)
}

func TestPreparedStatement(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	script := "insert into persons (name, gender, age) values (?, ?, ?)"
	statement, err := db.PrepareContext(ctx, script)
	if err != nil {
		panic(err)
	}
	defer statement.Close()

	for i := 0; i < 10; i++ {
		name := "Ricid-" + strconv.Itoa(i)
		gender := "M"
		age := 25

		result, err := statement.ExecContext(ctx, name, gender, age)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id ", id)
	}
}

func TestSqlTransaction(t *testing.T) {
	db := GetConnection()
	defer db.Close()

	ctx := context.Background()
	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	script := "insert into persons (name, gender, age) values (?, ?, ?)"
	for i := 0; i < 10; i++ {
		name := "Ricid-" + strconv.Itoa(i)
		gender := "M"
		age := 25

		result, err := tx.ExecContext(ctx, script, name, gender, age)
		if err != nil {
			panic(err)
		}

		id, err := result.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println("Comment Id ", id)
	}

	err = tx.Commit()
	// err = tx.Rollback()

	if err != nil {
		panic(err)
	}
}
