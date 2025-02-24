package seeder

import (
	"backend/models"
	"log"
	"math/rand"

	"github.com/beego/beego/v2/client/orm"
)

func SeedDatabase() {
	SeedNumberTable()
}

func SeedNumberTable() {
	o := orm.NewOrm()

	// Check if table has data
	count, err := o.QueryTable(new(models.Number)).Count()
	if err != nil {
		log.Fatalf("Error checking count in Number table: %v", err)
		return
	}
	if count == 0 {
		for i := 0; i < 10; i++ {
			num := models.Number{
				Value: rand.Intn(100), // Generate a random number between 0 and 99
			}

			// Insert the new number into the table
			_, err := o.Insert(&num)
			if err != nil {
				log.Printf("Error inserting number: %v", err)
			}
		}
		log.Println("Number table seeded successfully.")
	} else {
		log.Println("Number table already contains data.")
	}
}
