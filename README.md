# Golang Fiber with Docker

Golang Rest API with [Fiber](https://github.com/gofiber/fiber) and [GORM](https://github.com/go-gorm/gorm) and Docker

With pattern Handler, Usecase, Repository

---

## How to Run

```bash
# clone the repo
$ git clone repo

# go into repo's directory
$ cd repo

# copy and edit env file
$ cp .env.example .env

# seeds data to database like status
$ db migrate, edit pkg/configs/connect_db.go

# build docker
$ make build

# start docker
$ make run
```

### Todo List

- [x] Fiber Log file, Favicon
- [x] Fiber Monitor (metrics)
- [x] Golang Architecture Pattern
  - [x] Handler (delivery/controller)
  - [x] Usecase (bridge - logic process)
  - [x] Repository (process to DB)
  - [x] Error Handling with *fiber.Error
- [x] Auth
  - [x] Register, Send verification Email
  - [x] Send Email with Goroutines
  - [x] Open Link Verification Email
  - [x] Resend Verification Email Code
  - [x] Login
  - [x] JWT Auth Middleware + Redis
  - [x] Refresh Token
  - [x] Forgot Password, send email OTP
  - [x] Forgot Password Verify OTP
  - [x] Reset Password
  - [x] Logout
- [x] Account
  - [x] Get Profile
  - [x] Update Profile
  - [x] Update Photo Profile + thumbnail
  - [x] Upload File, upload image(compressed)
  - [x] Change Password
  - [x] Deletion Account with OTP
  - [x] Recover deleted account (Admin role)
  - [x] User Activity with interval (last login at, ip address in middleware)
- [x] Golang Swagger
- [x] CRUD
  - [x] Pagination with custom Paginate [pagination-using-gorm-scopes](https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f)
  - [x] Sort + Search function in List Data
  - [x] Create Data
  - [x] Edit Data
  - [x] Delete Data
- [x] Preload Model (Associations Struct)
- [x] Struct MarshalJSON (Custom representation)
- [ ] Open API with API KEY middleware
- [x] Upload Files
- [x] Remove Files
- [x] Upload Videos
- [x] Create thumbnail from videos with ffmpeg
- [x] Upload Images and Compress Image with libvips
- [x] Create thumbnail from image
- [x] Image Processing with [libvips](https://www.libvips.org/)

---

### Instalation LibVips

go to [https://www.libvips.org/](https://www.libvips.org/)

---


### Credit

[tutorial-go-fiber-rest-api](https://github.com/koddr/tutorial-go-fiber-rest-api)  
[Twilio](https://www.twilio.com/blog/build-restful-api-using-golang-and-gin)  
[Vinika Anthwal - Medium](https://medium.com/@22vinikaanthwal/register-login-api-with-jwt-authentication-in-golang-gin-740633e5707b) - [Github](https://github.com/VinikaAnthwal/go-jwt)  
[codevoweb.com](https://codevoweb.com/how-to-properly-use-jwt-for-authentication-in-golang/) - [Github](https://github.com/wpcodevo/golang-fiber-jwt-rs256)
