package appinfo

// AppName - экспортируемая константа.
// Начинается с заглавной буквы, поэтому видна из другого пакета.
const AppName = "Struct application demo"

// User - экспортируемый тип.
// Его можно использовать в main.go как appinfo.User.
type User struct {
	Name string
}

// GetMessage - экспортируемая функция.
// Ее можно вызвать из main.go как appinfo.GetMessage(...).
func GetMessage(user User) string {
	return "Hello, " + user.Name + "! " + secretMessage()
}

// secretMessage - неэкспортируемая функция.
// Начинается со строчной буквы, поэтому доступна только внутри package appinfo.
// В main.go нельзя написать appinfo.secretMessage().
func secretMessage() string {
	return "This message was created inside appinfo package."
}
