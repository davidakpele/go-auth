# 🛡️ Go Authentication API

This is a secure and extensible Go-based authentication service built using the [Gin Web Framework](https://github.com/gin-gonic/gin). It provides user registration, login, email verification (via OTP), and role-based authorization for protected API endpoints.

---

## 📦 Features

- 🔐 **User Authentication** (`register`, `login`, `logout`)
- 📧 **Email Verification** using OTP
- ♻️ **Resend OTP**
- 🔒 **JWT-based Authorization Middleware**
- 🧑‍💼 **Role-based Access Control** (`USER`, `ADMIN`, `SUPER_USER`)
- ❌ **Custom Error Handling** for Unauthorized Access and Unknown Routes

---

## 🚀 Getting Started

### 1. Clone the Repository

```bash
git clone https://github.com/davidakpele/go-auth
cd go-auth-api
```
2. Install Dependencies
- Make sure you have Go installed (>=1.18), then run:
``` go mod tidy ```
3. Environment Variables
  - Created a .env file and configure database and secret keys:
```
PORT=8080
DATABASE_URL=postgres://user:password@localhost:5432/yourdb
JWT_SECRET=your_jwt_secret_key
EMAIL_SMTP_HOST=smtp.yourprovider.com
EMAIL_SMTP_PORT=587
EMAIL_USERNAME=your_email@example.com
EMAIL_PASSWORD=your_email_password
```
##📌 API Endpoints
 -  🔓 Public Routes (/auth)
  - Protected by AuthenticationMiddleware and RoleMiddleware.

| Method | Endpoint               | Description                       |
| ------ | ---------------------- | --------------------------------- |
| POST   | `/auth/register`       | Register new user                 |
| POST   | `/auth/login`          | Authenticate user                 |
| GET    | `/auth/logout`         | Logout user (client-side)         |
| POST   | `/auth/verify-account` | Verify OTP for account activation |
| POST   | `/auth/resend-otp`     | Resend OTP to user email          |

##🔐 Protected Routes (/api)
  - Protected by AuthenticationMiddleware and RoleMiddleware.

| Method | Endpoint        | Description       | Role Required                 |
| ------ | --------------- | ----------------- | ----------------------------- |
| GET    | `/api/user/:id` | Fetch user by ID  | `USER`, `ADMIN`, `SUPER_USER` |
| DELETE | `/api/user/:id` | Delete user by ID | `ADMIN`, `SUPER_USER`         |


##🔧 Middleware
**🛡️ AuthenticationMiddleware**
  - Validates JWT tokens.
  - Attaches the user to the request context.

## 🧑‍💼 RoleMiddleware(...roles)
  - Verifies the current user’s role against allowed roles.
  - Denies access with a 401 Unauthorized response if not permitted.

## 📬 Email Service
  - On registration, a 6-digit OTP is sent to the user's email.
  - The user must verify the OTP to activate their account.
  - A resend OTP endpoint is available for unverified users.


## SMTP Configuration
Ensure your ```.env``` file contains the correct email service credentials to send OTPs.

## ❗ Error Handling
> Custom error handling is in place for unknown routes and unauthorized access. The base route / and undefined routes return:
```
{
  "error": "Access Denied",
  "status": "error",
  "title": "Authentication Error",
  "message": "Authorization Access",
  "details": "Something went wrong with authentication",
  "code": "generic_authentication_error"
}
```
- Similarly, all undefined routes are handled with the following response:
```
{
  "error": "EndPoint Not Found",
  "message": "Access Denied",
  "status": "error",
  "title": "Authentication Error",
  "details": "Something went wrong with authentication",
  "code": "generic_authentication_error"
}
```

## 🧪 Run Locally
   - To run the service locally:
       - Ensure that your database is set up and running.
       - Configure environment variables (e.g., database connection, SMTP settings, JWT secret).
       - Start the application:

> Direct with Go
```go run main.go```

> If you want auto restart when there is changes in your code, run the following command in yor terminar 

```air``` 
> This will run the API locally at http://localhost:7099.

-  The **7099** is the port number i choose to use in this project, you can change it from **_main.go_** file
  
