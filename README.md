# go-fiber jwt docker

Golang Rest API with [Fiber](https://github.com/gofiber/fiber) and [GORM](https://github.com/go-gorm/gorm) and Docker

With pattern Handler, Usecase, Repository

---

### List of packages that we will need to install for this project

```bash
// gin framework
go get -u github.com/gin-gonic/gin
// ORM library
go get -u github.com/jinzhu/gorm
// package that we will be used to authenticate and generate our JWT
go get -u github.com/dgrijalva/jwt-go
// to help manage our environment variables
go get -u github.com/joho/godotenv
// to encrypt our users password
go get -u golang.org/x/crypto
```

### Todo List

- [x] Fiber Log file
- [x] Fiber Monitor
- [x] CRUD
  - [x] Pagination with custom Paginate [pagination-using-gorm-scopes](https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f)
  - [ ] Redis List Pagination
  - [x] Sort + Search function in List Data
  - [x] Create Data
  - [x] Payload Sanitize
  - [x] Edit Data
  - [x] Delete Data
- [x] Preload Model (Associations in Go/Serializer in Django)
- [x] Auth (Register, Login)
  - [x] Register, Send verification Email
  - [x] Send Email with Goroutines (background task)
  - [x] Login
  - [x] Refresh Token
  - [ ] Logout
- [x] JWT Auth
- [x] Protected API (Auth Token)
- [x] Routes
  - [x] Delivery (handler/controller)
  - [x] Usecase (bridge - logic process)
  - [x] Repository (process to DB)
- [ ] Open API with API KEY middleware
- [x] Upload Files
- [ ] Upload Videos
- [ ] Create thumbnail from videos
- [x] Upload Images and Compress Image with libvips
- [x] Create thumbnail from image
- [x] Image Processing with [libvips](https://www.libvips.org/)

---

### Instalation LibVips

go to [https://www.libvips.org/](https://www.libvips.org/)

---

### Data

- [ ] User (Admin)
  - [x] List
  - [x] Get
  - [ ] Create
  - [ ] Update (Upload Photo)
  - [x] Delete
  - [x] List User Deleted
  - [ ] Hard Delete (constraint:OnDelete:CASCADE)
  - [x] Activated User
  - [x] De Activated User
  - [x] List User Activity
  - [x] Get List User Activity by ID
  - [x] Get List User Wishlist by ID
- [ ] API Key (Admin)
  - [x] Original API Key
  - [x] List with encode
  - [x] Get with encode
  - [x] Create
  - [x] Revoke
  - [x] Middleware API Key
- [ ] User
  - [x] Verify Email
  - [x] Profile
  - [x] Update Profile
  - [ ] Upload Photo Profile
  - [ ] Request Forgot Password and send link change password to email
  - [ ] Change Password
  - [x] Deletion Account
  - [x] List Wishlist
  - [x] Post Wishlist
  - [x] Post UnWishlist
  - [x] Save User Activity (last login at, ip address)
- [ ] Fact
  - [ ] List
  - [ ] Get
  - [ ] Create
  - [ ] Update
  - [ ] Delete
- [ ] Products
  - [x] Populate
  - [x] List
  - [ ] Get
  - [ ] Create
  - [ ] Update
  - [ ] Delete
- [ ] Merchants
  - [ ] List
  - [ ] Get
  - [ ] Create
  - [ ] Update
  - [ ] Delete (De Activated)

### Credit
[tutorial-go-fiber-rest-api](https://github.com/koddr/tutorial-go-fiber-rest-api)
[Twilio](https://www.twilio.com/blog/build-restful-api-using-golang-and-gin)  
[Vinika Anthwal - Medium](https://medium.com/@22vinikaanthwal/register-login-api-with-jwt-authentication-in-golang-gin-740633e5707b)  
[codevoweb.com](https://codevoweb.com/how-to-properly-use-jwt-for-authentication-in-golang/)