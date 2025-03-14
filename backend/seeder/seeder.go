package seeder

import (
	"backend/models"
	"log"
	"math/rand/v2"

	"github.com/beego/beego/v2/client/orm"
	"github.com/bxcodec/faker/v3"
	"golang.org/x/crypto/bcrypt"
)

func SeedDatabase() {
	SeedUsersWithDataAndSkills()
}

var workType = map[int]string{
	1: "on-site",
	2: "remote",
	3: "hybrid",
}

func getSkills() []string {
	return []string{
		"JavaScript", "Python", "Java", "C++", "C#", "Go", "Ruby", "Swift", "Kotlin", "PHP",
		"TypeScript", "Rust", "HTML", "CSS", "SQL", "NoSQL", "MongoDB", "PostgreSQL", "MySQL", "Django",
		"Flask", "Spring", "React", "Angular", "Vue.js", "Node.js", "Express.js", "GraphQL", "Docker", "Kubernetes",
		"Terraform", "AWS", "Azure", "GCP", "Firebase", "Linux", "Bash", "Git", "Jenkins", "CI/CD",
		"Machine Learning", "Deep Learning", "NLP", "Computer Vision", "TensorFlow", "PyTorch", "Pandas", "NumPy", "Scikit-Learn", "Data Analysis",
		"Cybersecurity", "Penetration Testing", "Cryptography", "Blockchain", "Solidity", "Smart Contracts", "DevOps", "Site Reliability Engineering", "Networking", "Cloud Computing",
		"Agile", "Scrum", "Project Management", "UI/UX Design", "Figma", "Adobe XD", "Photoshop", "Illustrator", "3D Modeling", "Animation",
		"Game Development", "Unity", "Unreal Engine", "Cocos2d", "Mobile Development", "iOS", "Android", "React Native", "Flutter", "Xamarin",
		"AR/VR", "IoT", "Embedded Systems", "Robotics", "Automation", "Big Data", "Hadoop", "Spark", "Kafka", "Elasticsearch",
		"Functional Programming", "Haskell", "Scala", "Lisp", "Erlang", "Dart", "Perl", "R", "MATLAB", "COBOL",
	}
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
			"client":     "Client222",
			"freelancer": "Freelancer222",
			"admin":      "Admin222",
		}

		for _, role := range roles {
			for range 10 {
				user := models.User{}

				user.Email = faker.Username() + "@gmail.com"
				user.Role = role
				user.Ban = false
				user.Name = faker.Name()
				user.Surname = faker.LastName()

				password := passwords[role]
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
				if err != nil {
					log.Printf("Error hashing password for %s: %v", role, err)
					continue
				}
				user.Password = string(hashedPassword)

				userID, err := o.Insert(&user)
				if err != nil {
					log.Printf("Error inserting user into database: %v", err)
				}

				switch role {
				case "freelancer":
					freelancerData := models.FreelancerData{}
					freelancerData.User = &models.User{Id: int(userID)}
					freelancerData.Description = faker.Paragraph()
					freelancerData.HourlyRate = float64(rand.IntN(50) + 10)
					freelancerData.HoursPerWeek = rand.IntN(100) + 10
					freelancerData.Title = faker.Word() + " " + faker.Word()
					freelancerData.WorkType = workType[rand.IntN(3)+1]

					if _, err := o.Insert(&freelancerData); err != nil {
						log.Printf("Error inserting freelancer data: %v", err)
						continue
					}

					var skills []models.Skill
					_, err := o.Raw("SELECT * FROM skills ORDER BY RANDOM() LIMIT ?", rand.IntN(3)+1).QueryRows(&skills)
					if err != nil {
						log.Printf("Error fetching random skills: %v", err)
						continue
					}

					m2m := o.QueryM2M(&freelancerData, "Skills")
					for _, skill := range skills {
						if _, err := m2m.Add(&skill); err != nil {
							log.Printf("Error adding skill %s to freelancer %d: %v", skill.Name, userID, err)
						}
					}

				case "client":
					clientData := models.ClientData{}
					clientData.User = &models.User{Id: int(userID)}
					clientData.Description = faker.Paragraph()
					clientData.CompanyName = faker.Word()
					clientData.Location = faker.Word()
					clientData.Industry = faker.Word()

					if _, err := o.Insert(&clientData); err != nil {
						log.Printf("Error inserting client data: %v", err)
						continue
					}
				}

			}
		}
		log.Println("User table seeded successfully")
	} else {
		log.Println("User table already contains data.")
	}
}

// Add 100 skills
func SeedSkillsTable() {
	o := orm.NewOrm()

	count, err := o.QueryTable(new(models.Skill)).Count()
	if err != nil {
		log.Fatalf("Error checking skills table: %v", err)
		return
	}

	if count > 0 {
		log.Println(" Skills table already contains data.")
		return
	}

	skills := getSkills()
	for _, skillName := range skills {
		skill := models.Skill{Name: skillName}
		if _, err := o.Insert(&skill); err != nil {
			log.Printf("Error inserting skill %s: %v", skillName, err)
		}
	}

}
