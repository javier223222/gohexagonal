package admin
import (
	"net/http"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"strconv"
)

type AdminHandler struct {
	service *AdminService
}

func NewAdminHandler(service *AdminService) *AdminHandler {
    return &AdminHandler{service: service}
}

func (h *AdminHandler) Save(c *gin.Context) {
	var req struct{
		Username string `json:"username"`
		Name string `json:"name"`
		LastName string `json:"lastName"`
		FullName string `json:"fullName"`
		Email string `json:"email"`
		Number string `json:"number"`
		IdRole int64 `json:"idRol"`
		CreatedBy int64 `json:"createdBy"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashedPassword,_:= bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	
	admin := &Admin{
		ID: 0,
	Username: req.Username,
	Name: req.Name,
	LastName: req.LastName,
	FullName: req.FullName,
	Email: req.Email,
	Number:req.Number,
	IdRole: req.IdRole,
	CreateAt: "",
	IsDeleted: false,
	CreatedBy: req.CreatedBy,
	Password: string(hashedPassword),
		
	}
	err := h.service.repo.Save(admin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "registered"})
}

func (h *AdminHandler)Login(c *gin.Context){
	var req struct{
		UserNameOrEmail string `json:"usernameOrEmail"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.service.Login(req.UserNameOrEmail, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	token, err := h.service.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "token": token})		
}

func (h *AdminHandler) Get(c *gin.Context) {
	idQuery := c.Query("id")
	page := c.Query("page")
	limit := c.Query("limit")

	if idQuery == "" {
		if page == "" || limit == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "page and limit query params are required"})
			return
		}

		page, _ := strconv.ParseInt(page, 10, 64)
		limit, _ := strconv.ParseInt(limit, 10, 64)
		admins, totalRecords, totalPages, err := h.service.Get(page, limit)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": admins, "totalRecords": totalRecords, "totalPages": totalPages})
		return

	}

	id, err := strconv.ParseInt(idQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	admin, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": admin})
}


func (h *AdminHandler) Delete(c *gin.Context) {
	idQuery := c.Param("id")
	id, err := strconv.ParseInt(idQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "admin deleted"})
}

func (h *AdminHandler) DeleteMYOnCount(c *gin.Context){
	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user id not found"})
		return
	}
	var userId64 int64
	switch v := userId.(type) {
	case int64:
		userId64 = v
	case float64:
		userId64 = int64(v)
	}
	err := h.service.Delete(userId64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "admin deleted"})


}


func (h *AdminHandler) UpdatePassword(c *gin.Context){

	userId, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "user id not found"})
		return
	}
	var userId64 int64
	switch v := userId.(type) {
	case int64:
		userId64 = v
	case float64:
		userId64 = int64(v)
	}
	var req struct{
		OldPassword string `json:"oldPassword"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.UpdatePassword(userId64, req.OldPassword, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "password updated"})
}

func (h *AdminHandler) UpdatePasswordByID(c *gin.Context){
	idQuery := c.Param("id")
	id, err := strconv.ParseInt(idQuery, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct{
		OldPassword string `json:"oldPassword"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.UpdatePassword(id, req.Password, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "password updated"})
}
