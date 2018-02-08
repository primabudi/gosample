package main

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"net/http"

	_ "github.com/lib/pq"
)

var dbHome *sql.DB
var err error
var counter = 0

type (
	User struct {
		ID          int
		Name        string
		Msisdn      string
		Email       string
		BirthDate   time.Time
		CreatedTime time.Time
		UpdateTime  time.Time
		Age         int
	}
)

var listUser = []User{}

func main() {
	dbHome, err = sql.Open("postgres", "postgres://st140804:apaajadeh@devel-postgre.tkpd/tokopedia-user?sslmode=disable")
	if err != nil {
		log.Print(err)
	}

	doQuery()

	http.HandleFunc("/html", handleHTML)
	http.ListenAndServe(":8080", nil)
}

func doQuery() {
	usr := &User{}
	stmt := "select user_id, full_name, msisdn, user_email, birth_date, create_time, update_time, DATE_PART('year', NOW()::date) - DATE_PART('year', birth_date::date) from ws_user ORDER BY user_id ASC limit 10"
	rows, err := dbHome.Query(stmt)

	for rows.Next() {
		_ = rows.Scan(&usr.ID, &usr.Name, &usr.Msisdn, &usr.Email, &usr.BirthDate, &usr.CreatedTime, &usr.UpdateTime, &usr.Age)
		listUser = append(listUser, *usr)
		//log.Print(usr.ID, " - ", usr.Name, " - ", usr.Msisdn, " - ", usr.BirthDate, " - ", usr.Email, " - ", usr.CreatedTime, " - ", usr.UpdateTime)
	}

	if err != nil {
		log.Print(err)
	}
}

func handleHTML(w http.ResponseWriter, r *http.Request) {
	counter = counter + 1
	strr :=
		"<html>" +
			"<body>" +

			"<table style>" +
			"<tr>" +
			"	<th>ID</th>" +
			"	<th>Name</th>" +
			"	<th>MSISDN</th>" +
			"	<th>Email</th>" +
			"	<th>Birthdate</th>" +
			"	<th>CreateTime</th>" +
			"	<th>UpdateTime</th>" +
			"	<th>Age</th>" +
			"</tr>"

	for i := 0; i < len(listUser); i++ {

		strr += "<tr>" +
			"<td>" + strconv.Itoa(listUser[i].ID) + "</td>" +
			"<td>" + listUser[i].Name + "</td>" +
			"<td>" + listUser[i].Msisdn + "</td>" +
			"<td>" + listUser[i].Email + "</td>" +
			"<td>" + listUser[i].BirthDate.String() + "</td>" +
			"<td>" + listUser[i].CreatedTime.String() + "</td>" +
			"<td>" + listUser[i].UpdateTime.String() + "</td>" +
			"<td>" + strconv.Itoa(listUser[i].Age) + "</td>" +
			"</tr>"
	}
	strr += "</table>" +
		"visit counter = " + strconv.Itoa(counter) +
		"</body>" +
		"</html>"
	byt := []byte(strr)
	w.Write(byt)
}
