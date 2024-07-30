package controller

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prajjwal-w/cetec_golang_practical/database"
	"github.com/prajjwal-w/cetec_golang_practical/models"
)

// initilazing the pseudo database conn
var db = database.DatabaseConn()

// Get person handler
func GetPerson() gin.HandlerFunc {
	return func(c *gin.Context) {
		//getting the person id from url
		person_id := c.Param("person_id")

		//query for getting the person by the person id
		query := `Select p.name, ph.number, a.city, a.state, a.street1, a.street2, a.zip_code From person p
		          JOIN phone ph on p.id = ph.person_id
				  JOIN AddressJoin aj on p.id = aj.person_id
				  JOIN Address a on aj.address_id = a.id 
				  Where p.id = $1`

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		var person models.Person

		err := db.QueryRowContext(ctx, query, person_id).Scan(
			&person.Name,
			&person.PhoneNumber,
			&person.City,
			&person.State,
			&person.Street1,
			&person.Street2,
			&person.Zip_code,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				c.JSON(http.StatusNotFound, gin.H{"error": "No Person Found"})
				return
			}
		}

		c.JSON(http.StatusOK, person)

	}
}

func CreatePerson() gin.HandlerFunc {
	return func(c *gin.Context) {
		var person models.Person

		if err := c.BindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating user"})
			return
		}
		defer tx.Rollback()

		var person_id int
		query := `Insert into person (name, age) Values($1,$2) Returning id`
		err = tx.QueryRow(query, person.Name).Scan(&person_id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Insert user"})
			return
		}

		var phone_id int
		query = `Insert into phone (number, person_id) values ($1, $2) Returning id`
		err = tx.QueryRow(query, person.PhoneNumber, person_id).Scan(&phone_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Insert user"})
			return
		}

		var address_id int
		query = `Insert into address (city, state, street1, street2, zip_code) Returning id`
		err = tx.QueryRow(query, person.City, person.State, person.Street1, person.Street2, person.Zip_code).Scan(&address_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Insert user"})
			return
		}

		var address_join int

		query = `insert into address_join (person_id, address_id) Values($1,$2) Returning id`
		err = tx.QueryRow(query, person_id, address_id).Scan(&address_join)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Insert user"})
			return
		}

		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Insert user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Person Created Successfully"})

	}
}
