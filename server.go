package main

import (
	"lab4/repo"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

type User struct {
	Name  string `json:"name" xml:"name" form:"name" query:"name" db:"name"`
	Email string `json:"email" xml:"email" form:"email" query:"email" db:"email"`
}

var (
	db, err = repo.NewPostgresDB(repo.Config{
		Host:     "localhost",
		Port:     "5432",
		Username: "postgres",
		Password: "1234",
		DBName:   "onelab",
		SSLMode:  "disable",
	})
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	e.POST("/users", saveUserBody)
	e.GET("/users/:id", getUser)
	e.PUT("/users", updateUser)
	e.DELETE("/users", deleteUserQuery)
	e.Logger.Fatal(e.Start(":8080"))

}

func getUser(c echo.Context) error {
	// User ID from path `users/:id`
	user := User{}
	id := c.Param("id")
	println(id)
	db.Get(&user, "SELECT * FROM users WHERE name = $1", id)
	return c.JSON(http.StatusOK, user)
}

func deleteUserQuery(c echo.Context) error {
	id := c.QueryParam("id")
	deleteUser := "DELETE FROM users WHERE name = $1"

	db.MustExec(deleteUser, id)
	return c.String(http.StatusOK, id)
}
func saveUserBody(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	// this will pull the first place directly into p
	addUser := "INSERT INTO users (email, name) VALUES($1 , $2)"

	db.MustExec(addUser, u.Email, u.Name)
	return c.JSON(http.StatusCreated, u)
}
func updateUser(c echo.Context) error {

	name := c.QueryParam("oldName")
	newName := c.QueryParam("newName")

	updateUser := "UPDATE users SET name = $1 WHERE name = $2"

	db.MustExec(updateUser, newName, name)

	return c.JSON(http.StatusCreated, newName)
}
