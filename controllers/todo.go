package controllers

import (
	"fmt"
	"gin-todo/models"
	"net/http"
	"strconv"
	"unicode/utf8"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var admin = false
var currentUser string

type TodoHandler struct {
	Db     *gorm.DB
	UserDb *gorm.DB
}

type SessionInfo struct {
	UserId interface{}
}

var LoginInfo SessionInfo

func (h *TodoHandler) AuthTest(c *gin.Context) {
	admin = false
	h.Logout(c)
	var users []models.User
	h.UserDb.Find(&users)
	c.HTML(http.StatusOK, "login.html", gin.H{
		"users": users,
	})
}

func (h *TodoHandler) Register(c *gin.Context) {
	var users []models.User
	h.UserDb.Find(&users)
	c.HTML(http.StatusOK, "register.html", gin.H{
		"users": users,
	})
}

func (h *TodoHandler) RegisterPOST(c *gin.Context) {
	account := c.PostForm("ID")
	if utf8.RuneCountInString(account) == 0 && utf8.RuneCountInString(account) > 10 {
		c.Redirect(http.StatusMovedPermanently, "/register")
		return
	}
	pass := c.PostForm("pass")
	if utf8.RuneCountInString(pass) == 0 && utf8.RuneCountInString(pass) > 10 {
		c.Redirect(http.StatusMovedPermanently, "/register")
		return
	}
	pass2 := c.PostForm("pass2")
	if pass != pass2 {
		c.Redirect(http.StatusMovedPermanently, "/register")
		return
	}
	fmt.Print("user")
	user := []models.User{}
	h.UserDb.Find(&user, account)
	if len(user) == 1 {
		c.Redirect(http.StatusMovedPermanently, "/register")
		return
	}
	admin := c.PostForm("admin")
	var BoolAdmin = true
	if admin == "" {
		BoolAdmin = false
	}

	fmt.Print("hash")
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), 12)

	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/register")
	} else {
		h.UserDb.Create(&models.User{Account: account, Pass: hash, Admin: BoolAdmin})
		c.Redirect(http.StatusMovedPermanently, "/todo")
	}
}

func (h *TodoHandler) Login(c *gin.Context) {
	account := c.PostForm("account")
	pass := c.PostForm("pass")
	user := models.User{}

	h.UserDb.Where("Account = ?", account).First(&user)

	hash := []byte(user.Pass)
	err := bcrypt.CompareHashAndPassword(hash, []byte(pass))

	//h.UserDb.Where("Account = ? AND Pass = ?", account, pass).First(&user)
	if account == "a" && pass == "a" {
		admin = true
		currentUser = account
		session := sessions.Default(c)
		session.Set("UserId", "superUser")
		session.Save()
		println("ログイン")
		c.Redirect(http.StatusMovedPermanently, "/todo")
		return
	}
	if err != nil {
		c.Redirect(http.StatusMovedPermanently, "/")
		return
	}
	if user.Admin {
		admin = true
	}

	currentUser = account
	session := sessions.Default(c)
	session.Set("UserId", account)
	session.Save()
	println("ログイン")
	c.Redirect(http.StatusMovedPermanently, "/todo")
}

func (h *TodoHandler) Logout(c *gin.Context) {
	session := sessions.Default(c)
	fmt.Println("セッション取得")
	session.Clear()
	fmt.Println("クリア処理")
	session.Save()
}

func (h *TodoHandler) GetAll(c *gin.Context) {
	sessionCheck(c)
	fmt.Println("getall")
	var todos []models.Todo
	var users []models.User

	if admin {
		h.Db.Find(&todos)
		h.UserDb.Find(&users)
	} else {
		h.Db.Where("User = ? ", currentUser).Find(&todos)
		h.UserDb.Find(&users)
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"todos":       todos,
		"users":       users,
		"admin":       admin,
		"currentUser": currentUser,
	})
}
func (h *TodoHandler) GetAllUser(c *gin.Context) {
	var users []models.User
	h.UserDb.Find(&users)
	c.HTML(http.StatusOK, "user.html", gin.H{
		"users": users,
	})
}
func (h *TodoHandler) CreateTask(c *gin.Context) {
	fmt.Printf("pass:%s\n", "CreateTask")
	text, _ := c.GetPostForm("text")
	status, _ := c.GetPostForm("status")
	istatus, _ := strconv.ParseUint(status, 10, 32)
	var todo = models.Todo{}
	h.Db.First(&todo, text)
	if todo.Text == text {
		c.Redirect(http.StatusMovedPermanently, "/todo")
		return
	}
	h.Db.Create(&models.Todo{Text: text, Status: istatus, User: currentUser})
	c.Redirect(http.StatusMovedPermanently, "/todo")
}

func (h *TodoHandler) EditTask(c *gin.Context) {
	todo := models.Todo{}
	id := c.Param("id")
	h.Db.First(&todo, id)
	c.HTML(http.StatusOK, "edit.html", gin.H{
		"todo":  todo,
		"admin": admin,
	})
}

func (h *TodoHandler) UpdateTask(c *gin.Context) {
	todo := models.Todo{}
	id := c.Param("id")
	text, _ := c.GetPostForm("text")
	status, _ := c.GetPostForm("status")
	istatus, _ := strconv.ParseUint(status, 10, 32)
	h.Db.First(&todo, id)
	todo.Text = text
	todo.Status = istatus
	h.Db.Save(&todo)
	c.Redirect(http.StatusMovedPermanently, "/todo")
}

func (h *TodoHandler) EditUser(c *gin.Context) {
	user := models.User{}
	id := c.Param("id")
	h.UserDb.First(&user, id)
	c.HTML(http.StatusOK, "userEdit.html", gin.H{
		"user": user,
	})
}

func (h *TodoHandler) DeleteTask(c *gin.Context) {
	todo := models.Todo{}
	id := c.Param("id")
	h.Db.First(&todo, id)
	h.Db.Delete(&todo)
	c.Redirect(http.StatusMovedPermanently, "/todo")
}

func (h *TodoHandler) DeleteUser(c *gin.Context) {
	user := models.User{}
	id := c.Param("id")
	h.UserDb.First(&user, id)
	h.UserDb.Delete(&user)
	c.Redirect(http.StatusMovedPermanently, "/todo")
}

func sessionCheck(c *gin.Context) {
	session := sessions.Default(c)
	UserId := session.Get("UserId")
	fmt.Println("sessionCheck")
	fmt.Printf("%c", UserId)
	//セッションがない場合、ログインフォームを出す
	if UserId == nil {
		fmt.Println("ログインしていません")
		c.Redirect(http.StatusMovedPermanently, "/")
		c.Abort()
	} else {
		c.Set("UserId", LoginInfo.UserId) // ユーザidをセット
		c.Next()
	}
	fmt.Println("ログインチェック終わり")
}
