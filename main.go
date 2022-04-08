package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Employee struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Salary string `json:"salary"`
	Age    string `json:"age"`
}

func main() {
	user := os.Getenv("DBUSER")
	pwd := os.Getenv("DBPWD")
	fmt.Println("Env: ", user)
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=Employee sslmode=disable", user, pwd))
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		panic(err)
	} else {
		fmt.Println("DB connected...")
	}
	defer db.Close()
	e := echo.New()

	e.POST("/employee", func(c echo.Context) error {

		u := new(Employee)
		if err := c.Bind(u); err != nil {
			return err
		}
		sqlstatement := "INSERT INTO emp (id,name, salary,age)VALUES ($1,$2,$3,$4)"
		res, err := db.Query(sqlstatement, u.Id, u.Name, u.Salary, u.Age)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, "ok")
	})

	e.GET("/employee", func(c echo.Context) error {
		sqlstatement := "SELECT id, name, salary, age FROM emp order by id"
		rows, err := db.Query(sqlstatement)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		//creating new slice of type Employee.
		names := make([]Employee, 0)
		for rows.Next() {
			name := Employee{}
			err := rows.Scan(&name.Id, &name.Name, &name.Salary, &name.Age)
			if err != nil {
				return err
			}
			names = append(names, name)
		}
		fmt.Println(names)
		return c.JSON(http.StatusCreated, names)
	})

	//Get by Id.
	e.GET("/employee/:id", func(c echo.Context) error {

		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}
		sqlstatement := "SELECT id, name, salary, age FROM emp order by id"
		rows, err := db.Query(sqlstatement)
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		//creating new slice of type Employee.
		Emp_slice := make([]Employee, 0)
		for rows.Next() {
			Emp_struct := Employee{}
			err := rows.Scan(&Emp_struct.Id, &Emp_struct.Name, &Emp_struct.Salary, &Emp_struct.Age)
			if err != nil {
				return err
			}
			Emp_slice = append(Emp_slice, Emp_struct)
		}
		//Below code was used to loop through slice of names of type Employee
		for _, emp_struct := range Emp_slice {
			//if condition check id in db with passed id.
			//if ok, then it returns Json as response.
			if emp_struct.Id == id {
				fmt.Println(emp_struct)
				return c.JSON(http.StatusCreated, emp_struct)
			}
		}

		return c.NoContent(http.StatusNoContent)
	})
	//Update by id.
	e.PUT("/employee/:id", func(c echo.Context) error {
		u := new(Employee)
		if err := c.Bind(u); err != nil {
			return err
		}
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}
		sqlstatement := fmt.Sprintf("UPDATE emp SET name=$1,salary=$2, age=$3 where  id = %d ", id)
		res, err := db.Query(sqlstatement, u.Name, u.Salary, u.Age)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			fmt.Println("updated data:", u)
			return c.JSON(http.StatusCreated, u)
		}
		return c.String(http.StatusOK, "ok")
	})
	e.DELETE("/employee/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}
		sqlstatement := fmt.Sprintf("DELETE FROM emp WHERE id = %d", id)
		res, err := db.Query(sqlstatement)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(res)
			return c.JSON(http.StatusOK, "Deleted")
		}
		return c.String(http.StatusOK, "Deleted")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
