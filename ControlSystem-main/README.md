# ControlSystem

A web-based control system application built with **Go (Gin & Gorm)** on the backend and **React** on the frontend.  
It uses **PostgreSQL** as the database and **MinIO** for object storage.  
The project is fully containerized with **Docker**.

---

##  Tech Stack

- **Backend:** [Go](https://go.dev/) with [Gin](https://gin-gonic.com/) & [Gorm](https://gorm.io/)  
- **Frontend:** [React](https://react.dev/)  
- **Database:** [PostgreSQL](https://www.postgresql.org/)  
- **Storage:** [MinIO](https://min.io/)  
- **Containerization:** [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)  

---

## Getting Started

### Prerequisites
- [Docker](https://www.docker.com/get-started) installed  
- [Docker Compose](https://docs.docker.com/compose/) installed
- [MinIo](https://docs.min.io/community/minio-object-store/operations/deployments/installation.html) installed

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/SpiritFoxo/ControlSystem
   cd ControlSystem
   ```
2. **Set up environment variables**
Copy the example environment file and configure it with your own values:
   ```bash
   cp .env.example .env
   ```
3. **Build and run the application**
   ```bash
   docker-compose up --build
   ```

4. **Stop the application**
   ```bash
   docker-compose down
   ```

## Notes
Update the .env file before running the containers (default values are placeholders).

Containers include backend, frontend, database, and storage services.
