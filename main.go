package main

import (
	"database/sql"
	"fmt"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/render"
	"net/http"
)

type User struct {
	Name        string
	Email       string
	Description string
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}

func SetupDB() *sql.DB {
	db, err := sql.Open("mysql", "root@/go_dev")
	//db, err := sql.Open("mysql", "USER:PASSWORD@tcp(VM IP:VM PORT/go_dev)")
	PanicIf(err)
	return db
}

func main() {
	m := martini.Classic()
	m.Map(SetupDB())

	// reads "templates" directory by default
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))

	m.Post("/users", func(ren render.Render, r *http.Request, db *sql.DB) {

		fmt.Println(r.FormValue("name"))
		fmt.Println(r.FormValue("email"))
		fmt.Println(r.FormValue("description"))

		_, err := db.Query("INSERT INTO users (name, email, description) VALUES (?, ?, ?)",
			r.FormValue("name"),
			r.FormValue("email"),
			r.FormValue("description"))

		PanicIf(err)

		rows, err := db.Query("SELECT * FROM users")
		PanicIf(err)
		defer rows.Close()

		users := []User{}
		for rows.Next() {
			user := User{}
			err := rows.Scan(&user.Name, &user.Email, &user.Description)
			PanicIf(err)
			users = append(users, user)

		}
		fmt.Println(users)
		ren.HTML(200, "users", users)
	})

	m.Run()

}
