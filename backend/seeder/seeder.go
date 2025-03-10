package seeder

import (
	"backend/models"
	"log"

	"github.com/beego/beego/v2/client/orm"
	"github.com/bxcodec/faker/v3"
	"golang.org/x/crypto/bcrypt"
)

func SeedDatabase() {
	SeedUsersWithDataAndSkills()
}

// TODO: Add client and freelancer data tables creation!
func SeedUsersWithDataAndSkills() {
	SeedSkillsTable()
	o := orm.NewOrm()

	count, err := o.QueryTable(new(models.User)).Count()
	if err != nil {
		log.Fatalf("Error checking count in User table: %v", err)
		return
	}
	if count == 0 {
		roles := []string{"client", "freelancer", "admin"}
		passwords := map[string]string{
			"client":     "client222",
			"freelancer": "freelancer222",
			"admin":      "admin222",
		}

		for _, role := range roles {
			for range 1 {
				user := models.User{}

				user.Email = faker.Username() + "@gmail.com"
				user.Role = role
				user.Ban = false
				user.Name = faker.Name()
				user.Surname = faker.LastName()
				user.Description = faker.Sentence()

				password := passwords[role]
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					log.Printf("Error hashing password for %s: %v", role, err)
					continue
				}
				user.Password = string(hashedPassword)

				_, err = o.Insert(&user)
				if err != nil {
					log.Printf("Error inserting user into database: %v", err)
				}
				// Create data tables for users based on their roles
				// Assign 1-3 random skills to freelancers
			}
		}
		log.Println("User table seeded successfully")
	} else {
		log.Println("User table already contains data.")
	}
}

// Add 100 skills
func SeedSkillsTable() {

}
