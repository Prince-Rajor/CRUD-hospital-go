# 🏥 Hospital Management API

A RESTful CRUD API for managing doctors and patients, built with Go, Gin, GORM, and MySQL.

---

## 🛠️ Tech Stack

| Technology | Purpose                     |
| ---------- | --------------------------- |
| **Go**     | Programming language        |
| **Gin**    | Web framework               |
| **GORM**   | ORM for database operations |
| **MySQL**  | Database                    |

---

## 📁 Project Structure

```
CRUD-hospital-go/
├── main.go                    # Application entry point
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependency lock file
├── config/
│   └── config.go              # Database connection configuration
├── database/
│   └── database.go            # Database initialization
├── models/
│   ├── doctor.go              # Doctor model
│   └── patient.go             # Patient model
├── controllers/
│   ├── doctor_controller.go   # Doctor CRUD handlers
│   └── patient_controller.go  # Patient CRUD handlers
└── routers/
    └── router.go              # API route definitions
```

---

## 📋 Prerequisites

- Go 1.21+ installed
- MySQL 8.0+ installed and running
- Postman (for API testing)

---

## 🚀 Installation & Setup

### 1. Clone the repository

```bash
git clone https://github.com/razorpay/CRUD-hospital-go.git
cd CRUD-hospital-go
```

### 2. Create MySQL database

```sql
CREATE DATABASE hospital_db;
```

### 3. Configure database connection

Update `config/config.go` with your MySQL credentials:

```go
dsn := "root:your_password@tcp(127.0.0.1:3306)/hospital_db?charset=utf8mb4&parseTime=True&loc=Local"
```

### 4. Install dependencies

```bash
go mod tidy
```

### 5. Run the application

```bash
go run main.go
```

The server will start at `http://localhost:8080`

---

## 📊 Database Models

### Doctor

| Field     | Type   | Description                  |
| --------- | ------ | ---------------------------- |
| ID        | uint   | Primary key (auto-increment) |
| Name      | string | Doctor's name                |
| ContactNo | string | Contact number               |
| Address   | string | Address                      |
| CreatedAt | time   | Record creation timestamp    |
| UpdatedAt | time   | Last update timestamp        |
| DeletedAt | time   | Soft delete timestamp        |

### Patient

| Field     | Type   | Description                  |
| --------- | ------ | ---------------------------- |
| ID        | uint   | Primary key (auto-increment) |
| Name      | string | Patient's name               |
| ContactNo | string | Contact number               |
| Address   | string | Address                      |
| DoctorID  | uint   | Foreign key to Doctor        |
| CreatedAt | time   | Record creation timestamp    |
| UpdatedAt | time   | Last update timestamp        |
| DeletedAt | time   | Soft delete timestamp        |

---

## 🔌 API Endpoints

### Base URL: `http://localhost:8080`

### Welcome

| Method | Endpoint | Description     |
| ------ | -------- | --------------- |
| GET    | `/`      | Welcome message |

---

### 🩺 Doctor Endpoints

| Method | Endpoint                       | Description                 |
| ------ | ------------------------------ | --------------------------- |
| POST   | `/doctor/`                     | Create a new doctor         |
| GET    | `/doctor/:id`                  | Get doctor by ID            |
| PATCH  | `/doctor/:id`                  | Update doctor (partial)     |
| DELETE | `/doctor/:id`                  | Delete doctor (soft delete) |
| GET    | `/searchDoctorByName?name=xxx` | Search doctors by name      |

---

### 🏥 Patient Endpoints

| Method | Endpoint                             | Description                  |
| ------ | ------------------------------------ | ---------------------------- |
| POST   | `/patient/`                          | Create a new patient         |
| GET    | `/patient/:id`                       | Get patient by ID            |
| PATCH  | `/patient/:id`                       | Update patient (partial)     |
| DELETE | `/patient/:id`                       | Delete patient (soft delete) |
| GET    | `/fetchPatientByDoctorId/:doctor_id` | Get patients by doctor ID    |
| GET    | `/searchPatientByName?name=xxx`      | Search patients by name      |

---

## 📬 API Usage Examples

### Create Doctor

**Request:**

```http
POST /doctor/
Content-Type: application/json

{
    "name": "Dr. John Smith",
    "contact_no": "9876543210",
    "address": "123 Medical Street, Mumbai"
}
```

**Response:**

```json
{
  "ID": 1,
  "CreatedAt": "2025-12-28T10:30:00Z",
  "UpdatedAt": "2025-12-28T10:30:00Z",
  "DeletedAt": null,
  "name": "Dr. John Smith",
  "contact_no": "9876543210",
  "address": "123 Medical Street, Mumbai"
}
```

---

### Get Doctor by ID

**Request:**

```http
GET /doctor/1
```

**Response:**

```json
{
  "ID": 1,
  "CreatedAt": "2025-12-28T10:30:00Z",
  "UpdatedAt": "2025-12-28T10:30:00Z",
  "DeletedAt": null,
  "name": "Dr. John Smith",
  "contact_no": "9876543210",
  "address": "123 Medical Street, Mumbai"
}
```

---

### Update Doctor (Partial Update)

**Request:**

```http
PATCH /doctor/1
Content-Type: application/json

{
    "name": "Dr. John Smith Jr."
}
```

Only the fields provided will be updated. Other fields remain unchanged.

---

### Search Doctor by Name

**Request:**

```http
GET /searchDoctorByName?name=john
```

**Response:**

```json
[
    {
        "ID": 1,
        "name": "Dr. John Smith",
        ...
    },
    {
        "ID": 2,
        "name": "Dr. Johnny Doe",
        ...
    }
]
```

---

### Create Patient

**Request:**

```http
POST /patient/
Content-Type: application/json

{
    "name": "Rahul Kumar",
    "contact_no": "9123456789",
    "address": "456 Patient Lane, Delhi",
    "doctor_id": 1
}
```

---

### Get Patients by Doctor ID

**Request:**

```http
GET /fetchPatientByDoctorId/1
```

**Response:**

```json
[
    {
        "ID": 1,
        "name": "Rahul Kumar",
        "doctor_id": 1,
        ...
    },
    {
        "ID": 2,
        "name": "Priya Sharma",
        "doctor_id": 1,
        ...
    }
]
```

---

### Delete Doctor/Patient

**Request:**

```http
DELETE /doctor/1
```

**Response:**

```json
{
  "message": "Doctor deleted successfully"
}
```

> **Note:** This is a soft delete. The record is not permanently removed but marked with a `DeletedAt` timestamp.

---

## 🧪 Testing with Postman

1. Import the endpoints into Postman
2. Set the base URL to `http://localhost:8080`
3. For POST/PATCH requests:
   - Go to **Body** tab
   - Select **raw**
   - Choose **JSON** from dropdown
   - Enter your JSON payload

---

## 📝 Features

- ✅ Full CRUD operations for Doctors and Patients
- ✅ Partial updates using pointer fields
- ✅ Soft delete (records are not permanently deleted)
- ✅ Search by name with LIKE query
- ✅ Doctor-Patient relationship
- ✅ Auto-migration of database tables
- ✅ JSON API responses

---

## 🔧 Configuration

### Database Configuration

Edit `config/config.go`:

```go
dsn := "username:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True&loc=Local"
```

### Server Port

Edit `main.go`:

```go
router.Run(":8080")  // Change 8080 to your desired port
```

---

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

## 👤 Author

Prince Rajor

---

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

# CRUD-hospital-go


