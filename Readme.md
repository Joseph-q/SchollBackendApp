# 🎓 School Manage Backend  

Welcome to the **School Manage Backend**! This project is built with **Go (Golang)** and serves as the backend system for managing courses and students.  

## 🛠️ Prerequisites  

Before getting started, ensure the following dependencies are installed on your machine:  

- **Go (Golang) v1.23.2**  
- **SQLite**  
- **Docker**  

## 🚀 Getting Started  

### 1️⃣ Clone the Repository  
```sh
git clone https://github.com/Joseph-q/SchollBackendApp
```

### 2️⃣ Navigate to the Project Directory  
```sh
cd SchoolBackendApp
```

### 3️⃣ Build the Docker Image  
```sh
docker build -t school_backend .
```

### 4️⃣ Run the Project in Detached Mode  
```sh
docker compose -f docker-compose.yml up -d
```

### ⚙️ Configuration  

#### 📂 Changing Data Storage Location  
If you need to modify where the data is stored, edit the **docker-compose.yml** file and update the database URL:  
```yaml
services:
  SCHOOL_API:
    .....
    volumes:
      - /home/joseph/registroAlumnos/SchoolPr/config:/root/config
```

#### 🛠️ General Configuration  
To change other settings, such as the database URL or additional configuration options, edit the **config/config_develop.yaml** file and update the necessary values:  
```yaml
database:
  url: "your_new_url"
```

## 🌐 Frontend Interface  

Looking for the **frontend interface**? You can find it here:  
🔗 [School Frontend App](https://github.com/Joseph-q/SchollFrontendApp)  