package db

import(
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gin-todo/models"
)

func Init() *gorm.DB{
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil{
		panic("データベースへの接続に失敗しました")
	}
	db.AutoMigrate(&models.Todo{})
	//db.Migrator().DropTable(&models.Todo{})
	return db
}

func UserInit() *gorm.DB{
	db,err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil{
		panic("データベースの接続に失敗しました")
	}
	db.AutoMigrate(&models.User{})
	//db.Migrator().DropTable(&models.User{})
	return db
}