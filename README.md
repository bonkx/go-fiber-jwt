# Golang Fiber with Docker

Golang Rest API with [Fiber](https://github.com/gofiber/fiber) and [GORM](https://github.com/go-gorm/gorm) and Docker

With pattern Handler, Usecase, Repository

---

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
  - [ ] List
  - [ ] Get
  - [ ] Create
  - [ ] Update (Upload Photo)
  - [ ] Delete
  - [ ] List User Deleted
  - [ ] Hard Delete (constraint:OnDelete:CASCADE)
  - [ ] Activated User
  - [ ] De Activated User
  - [ ] List User Activity
  - [ ] Get List User Activity by ID
  - [ ] Get List User Wishlist by ID
- [ ] API Key (Admin)
  - [ ] Original API Key
  - [ ] List with encode
  - [ ] Get with encode
  - [ ] Create
  - [ ] Revoke
  - [ ] Middleware API Key
- [ ] User
  - [ ] Verify Email
  - [ ] Profile
  - [ ] Update Profile
  - [ ] Upload Photo Profile
  - [ ] Request Forgot Password and send link change password to email
  - [ ] Change Password
  - [ ] Deletion Account
  - [ ] List Wishlist
  - [ ] Post Wishlist
  - [ ] Post UnWishlist
  - [ ] Save User Activity (last login at, ip address)
- [ ] Fact
  - [ ] List
  - [ ] Get
  - [ ] Create
  - [ ] Update
  - [ ] Delete
- [ ] Products
  - [ ] Populate
  - [ ] List
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
[Vinika Anthwal - Medium](https://medium.com/@22vinikaanthwal/register-login-api-with-jwt-authentication-in-golang-gin-740633e5707b) - [Github](https://github.com/VinikaAnthwal/go-jwt)  
[codevoweb.com](https://codevoweb.com/how-to-properly-use-jwt-for-authentication-in-golang/) - [Github](https://github.com/wpcodevo/golang-fiber-jwt-rs256)
