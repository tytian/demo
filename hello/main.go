package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

var (
	Db *sql.DB
	err error
)

func init() {
	Db, err = sql.Open("mysql", "root:abc123@tcp(127.0.0.1:3306)/test?parseTime=true")
	if err != nil{
		fmt.Println(err)
	}
	//defer Db.Close()
	Db.SetMaxOpenConns(20) //最大链接数的设置
	Db.SetMaxIdleConns(20) //设置最大空闲链接数
	if err := Db.Ping(); err != nil{
		fmt.Println("-----ping fail----",err)
		return
	}
	fmt.Println("pong")
}

func ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

func putting(c *gin.Context) {
	name := c.Query("name")
	rs, err := Db.Exec("INSERT INTO person (name) VALUES (?)", name)
	if err != nil {
		fmt.Println("----db insert fail!------", err)
		return
	}

	id, err := rs.LastInsertId()
	if err != nil {
		fmt.Println("----get id fail!------", err)
		return
	}
	fmt.Println("insert person Id {}", id)
	msg := fmt.Sprintf("insert successful %d", id)
	c.JSON(http.StatusOK, gin.H{
		"msg": msg,
	})
}

func getting(c *gin.Context) {
	var name string
	id := c.Query("id")
	row := Db.QueryRow("select name from person where id=?", id)
	err := row.Scan(&name)
	if err != nil {
		fmt.Println("------query fail-----", err)
		c.JSON(http.StatusOK, gin.H{
			"person": nil,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"person": name,
	})
}

func main() {
	router := gin.Default()
	router.GET("/ping", ping)
	router.GET("/get", getting)
	router.PUT("/put", putting)
	router.Run()
}

