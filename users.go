package rentals

type User struct {
	// Primary key
	ID uid `gorm:"primary_key" json:"id"`

	// Username
	Username string `gorm:"unique" json:"username"`

	// Password hash. Not included in json responses
	PasswordHash string `json:"-"`

	// Role
	Role string `json:"role"`
}

type UserSession struct {
	// Primary key
	ID uint `gorm:"primary_key"`

	// Generated token
	Token string

	// User associated to this session
	UserID uint
	User   User
}

type UserCreateInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserCreateOutput struct {
	User
}

type UserReadInput struct {
	// ID to lookup the apartment
	Id string
}

type UserReadOutput struct {
	User
}

type UserAllInput struct {
}

type UserAllOutput struct {
	Users []User
}

func (o *UserAllOutput) Public() interface{} {
	return o.Users
}

type UserUpdateInput struct {
	Id       string `json:"-"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UserUpdateOutput struct {
	User
}

type UserDeleteInput struct {
	Id string
}

type UserDeleteOutput struct {
	Message string `json:"message"`
}

type UserService interface {
	Create(UserCreateInput) (*UserCreateOutput, error)
	Read(UserReadInput) (*UserReadOutput, error)
	All(UserAllInput) (*UserAllOutput, error)
	Update(UserUpdateInput) (*UserUpdateOutput, error)
	Delete(UserDeleteInput) (*UserDeleteOutput, error)
}
