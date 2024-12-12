package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func initDB() {
	var err error
	db, err = sqlx.Connect("postgres", "user=postgres password=nasru dbname=Students sslmode=disable")
	if err != nil {
		log.Fatalln("Failed to connect to database:", err)
	}
	log.Println("Database connected successfully!")
}

type Person struct {
	Student_id int    `db:"student_id" json:"student_id"`
	First_name string `db:"firstname" json:"firstname"`
	Last_name  string `db:"lastname" json:"lastname"`
	Email      string `db:"email" json:"email"`
	Contact    string `db:"contact" json:"contact"`
}

func addStudents(c *gin.Context) {
	var students Person
	if err := c.ShouldBindJSON(&students); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	query := `INSERT INTO students (firstname, lastname, email, contact) 
              VALUES ($1, $2, $3, $4) RETURNING student_id`
	err := db.QueryRow(query, students.First_name, students.Last_name, students.Email, students.Contact).
		Scan(&students.Student_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add student"})
		return
	}

	c.JSON(http.StatusCreated, students)
}

func getStudents(c *gin.Context) {
	id := c.Param("id")
	var student Person
	err := db.Get(&student, "SELECT * FROM students WHERE student_id=$1", id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
		return
	}
	c.JSON(http.StatusOK, student)
}

func getAllStudents(c *gin.Context) {
	var students []Person
	err := db.Select(&students, "SELECT * FROM students")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch students"})
		return
	}
	c.JSON(http.StatusOK, students)
}

func updateStudents(c *gin.Context) {
	id := c.Param("id")
	var students Person
	if err := c.ShouldBindJSON(&students); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := `UPDATE students SET firstname=$1, lastname=$2, email=$3, contact=$4 WHERE student_id=$5`
	_, err := db.Exec(query, students.First_name, students.Last_name, students.Email, students.Contact, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update student"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Student updated successfully"})
}

func deleteStudents(c *gin.Context) {
	id := c.Param("id")
	_, err := db.Exec("DELETE FROM students WHERE student_id=$1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete student"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
}

func main() {
	initDB()

	r := gin.Default()

	r.POST("/students", addStudents)
	r.GET("/students/:id", getStudents)
	r.GET("/students", getAllStudents)
	r.PUT("/students/:id", updateStudents)
	r.DELETE("/students/:id", deleteStudents)

	log.Println("Server is running on port 8080")
	r.Run(":8080")
}
