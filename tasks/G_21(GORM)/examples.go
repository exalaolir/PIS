

import "gorm.io/gorm"

type User struct {
	gorm.Model        // Встраиваемая структура (ID, CreatedAt, UpdatedAt, DeletedAt)
	Name       string `gorm:"type:varchar(100);not null"`
	Email      string `gorm:"uniqueIndex"`
	Age        uint8  `gorm:"default:18"`
	Role       string `gorm:"-"` // Игнорировать это поле (не создавать колонку)
}


db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    Logger: logger.Default.LogMode(logger.Info),
})

db.AutoMigrate(&User{}, &Product{})

// 1. Создаем объект в памяти Go
newUser := User{Name: "Ivan", Age: 25}

// 2. Сохраняем в базу данных
db.Create(&newUser)



var user User
// Найти первого пользователя с ID = 10
db.First(&user, 10) 

// Найти пользователя по условию (Name = "Ivan")
db.Where("name = ?", "Ivan").First(&user)


var users []User // Слайс, куда сложим результат
// Найти всех пользователей, у которых Age > 18
db.Where("age > ?", 18).Find(&users)


db.Model(&User{}).Where("id = ?", 1).Update("name", "NewName")

updatedData := User{Name: "Max", Age: 0} 

db.Model(&User{}).Where("id = ?", 1).Updates(updatedData)





// Карта явно говорит: запиши в Age именно 0
db.Model(&User{}).Where("id = ?", 1).Updates(map[string]interface{}{
    "name": "Max",
    "age":  0,
})




var user User
db.First(&user, 1)

// Мягкое удаление
db.Delete(&user)

db.Unscoped().Delete(&user)


// Загрузить пользователя и все его заказы (Eager Loading)
db.Preload("Orders").Find(&users)


db.Transaction(func(tx *gorm.DB) error {
    if err := tx.Create(&User{Name: "Alice"}).Error; err != nil {
        return err // Автоматический Rollback
    }
    return nil // Автоматический Commit
})


stmt := db.Session(&gorm.Session{DryRun: true}).First(&user, 1).Statement
fmt.Println(stmt.SQL.String())