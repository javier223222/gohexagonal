package admin

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("Error  en las credenciales")


type AdminService struct {
	repo AdminRepository
	
}


func NewAdminService(repo AdminRepository) *AdminService {
    return &AdminService{repo: repo}
}

func (s *AdminService) GenerateToken(admin *Admin) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "id":       admin.ID,
        "username": admin.Username,
        "role":    admin.IdRole,
        "exp":      time.Now().Add(time.Hour * 72).Unix(), // Token válido por 72 horas
    })

    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}



func (s *AdminService) Save(admin *Admin) error {
 
	
	return s.repo.Save(admin)
}

func (s *AdminService) Login(userNameOrEmail string, password string) (*Admin, error) {
	user, err := s.repo.GetByUsernameOrEmail(userNameOrEmail)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}
	return user, nil
}

func (s *AdminService) Get(page int64, limit int64) ([]Admin, int64, int64, error) {
	return s.repo.Get(page, limit)
}

func (s *AdminService) GetByID(id int64) (*Admin, error) {
	return s.repo.GetByID(id)
}

func (s *AdminService) Delete(id int64) error {
	return s.repo.Delete(id)
}


func (s *AdminService) UpdatePassword(id int64,oldpassword string, password string) error {
	user, err := s.repo.GetByIDWithPassword(id)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrInvalidCredentials
	}
	log.Println("Password:", user)
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldpassword)); err != nil {
		return err
	}
	
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    
	return s.repo.UpdatePassword(id, string(hashedPassword))
}


func (s *AdminService) UpdateUsername(id int64, username string) error {
	return s.repo.UpdateUsername(id, username)
}

func (s *AdminService) UpdateName(id int64, name string) error {
	return s.repo.UpdateName(id, name)
}

func (s *AdminService) UpdateLastName(id int64, lastName string) error {
	return s.repo.UpdateLastName(id, lastName)
}

func (s *AdminService) UpdateEmail(id int64, email string) error {
	return s.repo.UpdateEmail(id, email)
}

func (s *AdminService) UpdateNumber(id int64, number string) error {
	return s.repo.UpdateNumber(id, number)
}

func (s *AdminService) UpdateRole(id int64, role int64) error {
	return s.repo.UpdateRole(id, role)
}

