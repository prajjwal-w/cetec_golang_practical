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

		//creating the time out context for inturepting the long running query
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

// create person handler
func CreatePerson() gin.HandlerFunc {
	return func(c *gin.Context) {

		var person models.Person

		//binding the json to the person object

		if err := c.BindJSON(&person); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//creating the time out context for inturepting the long running query
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// as we have to insert the data into the multiple table for that we are using the transctions if one insert query fails then all fails
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while creating user"})
			return
		}
		//if the any query failes we are rolling back all our tranctions
		defer tx.Rollback()

		//inserting to the person table
		var person_id int
		query := `Insert into person (name, age) Values($1,$2) Returning id`
		err = tx.QueryRow(query, person.Name).Scan(&person_id)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Insert user"})
			return
		}

		//inserting to the phone table
		var phone_id int
		query = `Insert into phone (number, person_id) values ($1, $2) Returning id`
		err = tx.QueryRow(query, person.PhoneNumber, person_id).Scan(&phone_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Insert user"})
			return
		}

		//inserting to the address table
		var address_id int
		query = `Insert into address (city, state, street1, street2, zip_code) Returning id`
		err = tx.QueryRow(query, person.City, person.State, person.Street1, person.Street2, person.Zip_code).Scan(&address_id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Insert user"})
			return
		}

		//inserting to the address_join table
		var address_join int

		query = `insert into address_join (person_id, address_id) Values($1,$2) Returning id`
		err = tx.QueryRow(query, person_id, address_id).Scan(&address_join)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Insert user"})
			return
		}

		//commiting the transction
		if err = tx.Commit(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to Insert user"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Person Created Successfully"})

	}
}
