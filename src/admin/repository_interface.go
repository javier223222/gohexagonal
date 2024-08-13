package admin

type AdminRepository interface {
	Save(admin *Admin) error
	GetByUsernameOrEmail(field string) (*Admin, error)
	Get(page int64,limit int64) ([]Admin, int64, int64, error)
	GetByID(id int64) (*Admin,error)
	GetByIDWithPassword(id int64) (*Admin,error)
	Delete(id int64) error
	UpdatePassword(id int64,password string) error
	UpdateUsername(id int64,username string) error
	UpdateName(id int64,name string) error
	UpdateLastName(id int64,lastName string) error
	UpdateEmail(id int64,email string) error
	UpdateNumber(id int64,number string) error
	UpdateRole(id int64,role int64) error
	
}