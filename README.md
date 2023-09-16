# go-fiber jwt docker

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
- [ ] Auth
  - [x] Register, Send verification Email
  - [x] Send Email with Goroutines
  - [x] Open Link Verification Email
  - [x] Resend Verification Email Code
  - [x] Login
  - [x] JWT Auth Middleware + Redis
  - [x] Refresh Token
  - [ ] Forgot Password
  - [ ] Forgot Password Verify OTP
  - [ ] Change Password
  - [ ] Logout
- [ ] CRUD
  - [ ] Pagination with custom Paginate [pagination-using-gorm-scopes](https://dev.to/rafaelgfirmino/pagination-using-gorm-scopes-3k5f)
  - [ ] Redis List Pagination
  - [ ] Sort + Search function in List Data
  - [ ] Create Data
  - [ ] Payload Sanitize
  - [ ] Edit Data
  - [ ] Delete Data
- [ ] Preload Model (Associations in Go/Serializer in Django)
- [ ] Protected API (Auth Token)
- [ ] Open API with API KEY middleware
- [ ] Upload Files
- [ ] Upload Videos
- [ ] Create thumbnail from videos
- [ ] Upload Images and Compress Image with libvips
- [ ] Create thumbnail from image
- [ ] Image Processing with [libvips](https://www.libvips.org/)

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
