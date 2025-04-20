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
	SeedSkills()
	SeedUsers()
	SeedJobs()
	SeedApplications()
}

var projectTypes = []string{"ongoing", "one-time"}
var rateTypes = []string{"hourly", "fixed"}
var lengthOptions = []string{"<1", "1-3", "3-6", "6-12", "12+"}
var hoursPerWeekOptions = []string{"<20", "20-40", "40-60", "60-80", "80+"}

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

func SeedUsers() {
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
		roleCounts := map[string]int{
			"client":     45,
			"freelancer": 45,
			"admin":      10,
		}

		for _, role := range roles {
			for i := 0; i < roleCounts[role]; i++ {
				user := models.User{}

				if i == 0 {
					user.Email = role + "@gmail.com"
				} else {
					user.Email = faker.Username() + "@gmail.com"
				}

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
					freelancerData.HoursPerWeek = hoursPerWeekOptions[rand.IntN(len(hoursPerWeekOptions))]
					freelancerData.Title = faker.Word() + " " + faker.Word()

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
					clientData := models.ClientData{}
					clientData.User = &models.User{Id: int(userID)}

					if _, err := o.Insert(&clientData); err != nil {
						log.Printf("Error inserting client data: %v", err)
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

					freelancerData := models.FreelancerData{}
					freelancerData.User = &models.User{Id: int(userID)}

					if _, err := o.Insert(&freelancerData); err != nil {
						log.Printf("Error inserting freelancer data: %v", err)
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
func SeedSkills() {
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

func SeedJobs() {
	o := orm.NewOrm()

	// Check if jobs table already has data
	count, err := o.QueryTable(new(models.Job)).Count()
	if err != nil {
		log.Fatalf("Error checking count in Jobs table: %v", err)
		return
	}
	if count > 0 {
		log.Println("Jobs table already contains data.")
		return
	}

	// Job type and rate options
	projectTypes := []string{"ongoing", "one-time"}
	rateTypes := []string{"hourly", "fixed"}
	lengthOptions := []string{"<1", "1-3", "3-6", "6-12", "12+"}
	hoursPerWeekOptions := []string{"<20", "20-40", "40-60", "60-80", "80+"}

	// Get all clients
	var clients []models.User
	_, err = o.QueryTable(new(models.User)).Filter("role", "client").All(&clients)
	if err != nil {
		log.Printf("Error fetching clients: %v", err)
		return
	}

	// Create 0-5 jobs for each client
	for _, client := range clients {
		numJobs := rand.IntN(6)

		for range numJobs {
			rateType := rateTypes[rand.IntN(len(rateTypes))]

			// Set amount based on rate type
			var amount int
			if rateType == "hourly" {
				amount = rand.IntN(81) + 20
			} else {
				amount = rand.IntN(101) + 10
				amount *= 1000
			}

			job := models.Job{
				Client:       &client,
				Title:        faker.Sentence()[:30],
				Description:  faker.Paragraph(),
				Type:         projectTypes[rand.IntN(len(projectTypes))],
				Rate:         rateType,
				Amount:       amount,
				Length:       lengthOptions[rand.IntN(len(lengthOptions))],
				HoursPerWeek: hoursPerWeekOptions[rand.IntN(len(hoursPerWeekOptions))],
				Status:       "open",
			}

			jobID, err := o.Insert(&job)
			if err != nil {
				log.Printf("Error inserting job: %v", err)
				continue
			}

			// Add random skills (1-3) to the job
			var skills []models.Skill
			_, err = o.Raw("SELECT * FROM skills ORDER BY RANDOM() LIMIT ?", rand.IntN(3)+1).QueryRows(&skills)
			if err != nil {
				log.Printf("Error fetching random skills: %v", err)
				continue
			}

			m2m := o.QueryM2M(&job, "Skills")
			for _, skill := range skills {
				if _, err := m2m.Add(&skill); err != nil {
					log.Printf("Error adding skill %s to job %d: %v", skill.Name, jobID, err)
				}
			}
		}
	}

	log.Println("Jobs table seeded successfully")
}

func SeedApplications() {
	o := orm.NewOrm()

	count, err := o.QueryTable(new(models.Application)).Count()
	if err != nil {
		log.Fatalf("Error checking applications table: %v", err)
		return
	}
	if count > 0 {
		log.Println("Applications table already contains data.")
		return
	}

	var jobs []models.Job
	_, err = o.QueryTable(new(models.Job)).All(&jobs)
	if err != nil {
		log.Printf("Error fetching jobs: %v", err)
		return
	}

	var freelancers []models.User
	_, err = o.QueryTable(new(models.User)).Filter("role", "freelancer").All(&freelancers)
	if err != nil {
		log.Printf("Error fetching freelancers: %v", err)
		return
	}

	// Seed applications: 3 for each job
	for _, job := range jobs {
		used := make(map[int]bool) // track used freelancer IDs

		appCount := 0
		for appCount < 3 {
			randIndex := rand.IntN(len(freelancers))
			freelancer := freelancers[randIndex]

			if used[freelancer.Id] {
				continue
			}

			used[freelancer.Id] = true

			app := models.Application{
				User:        &freelancer,
				Job:         &job,
				Description: faker.Paragraph(),
				Status:      "pending",
			}

			if _, err := o.Insert(&app); err != nil {
				log.Printf("Error inserting application for freelancer %d on job %d: %v", freelancer.Id, job.Id, err)
			} else {
				appCount++
			}
		}
	}

	log.Println("Applications seeded successfully.")
}
