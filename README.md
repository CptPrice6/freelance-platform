# Freelance Platform

## Project Overview

This is a **Freelance Platform** developed as part of a university project. The platform focuses on software development tasks and allows **clients** to post coding related job listings and **freelancers** to apply for them. It includes user authentication, job management, and an application process.

## Author

**Name:** Domantas Petkevičius  
**University:** Vilniaus Universitetas, Matematikos ir Informatikos fakultetas (MIF)

## Tech Stack

- **Backend:** Beego (GoLang)
- **Frontend:** ReactJS, Bootstrap
- **Database:** PostgreSQL
- **Containerization:** Docker (with Docker Compose)

## Features

✔ **User Authentication:** Clients & Freelancers can register and log in.  
✔ **Job Management:** Clients can post, edit, and delete jobs.  
✔ **Applications:** Freelancers can browse and apply for jobs. 
✔ **Role-Based Access:** Different permissions for Clients, Freelancers (PrivateRoute) and Admins (AdminRoute).  
✔ **Database ORM:** Efficient data handling using Beego's ORM.  
✔ **File Uploads:** Attachments for applications.  
✔ **API Integration:** External APIs can be used for data retrieval.

## Project Structure

```
├── backend/                 # Beego backend (GoLang)
│   ├── controllers/         # Handles request logic
│   ├── models/              # Database models (Users, Jobs, Applications)
│   ├── routers/             # API routes
│   ├── middleware/          # Authentication middleware
│   ├── database/            # Database connections and configurations
│   ├── migrations/          # Database migration files
│   ├── seeder/              # Database seeding scripts
│   ├── types/               # Type definitions for requests and responses
│   ├── utils/               # Utility functions
│   ├── validators/          # Input validation logic
│   ├── Dockerfile           # Docker configuration for backend
│   ├── main.go              # Application entry point
│   └── conf/                # Configuration files
│
├── frontend/                # ReactJS Frontend
│   ├── public/
│   │   ├── index.html       # Main HTML file
│   ├── src/
│   │   ├── components/      # Reusable components
│   │   ├── pages/           # Page-level components
│   │   ├── styles/          # CSS styling files for pages
│   │   ├── utils/           # Utility functions
│   │   ├── App.js           # Main React component
│   │   └── index.js         # Entry point
│   ├── Dockerfile           # Docker configuration for frontend
│
├── docker-compose.yml       # Docker configuration
├── README.md                # Project documentation
└── .env                     # Environment variables
```

## Setup Instructions

### 1. Clone the Repository

```
git clone https://github.com/CptPrice6/freelance-platform.git
cd freelance-platform
```

### 2. Start Docker Containers

```
docker-compose up --build
```

### 3. Access the Application

- **Backend API:** http://localhost:8080
- **Frontend UI:** http://localhost:3000
