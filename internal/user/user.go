package user

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	UserIDCounter = 1
	rootUser      = User{ID: 0, Name: "root", Email: "root@root.com", Password: "root"}
)

func GenId() (id int) {
	id = UserIDCounter
	UserIDCounter++
	return id
}
