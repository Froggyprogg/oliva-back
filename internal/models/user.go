package models

type User struct {
	Id          uint32 `json:"id" sql:"id"`
	Surname     string `json:"surname" validate:"required" sql:"surname"`
	Name        string `json:"name" validate:"required" sql:"name"`
	Middlename  string `json:"middlename" sql:"middlename"`
	PhoneNumber string `json:"phonenumber" validate:"required" sql:"phonenumber"`
	Email       string `json:"email" validate:"required" sql:"email"`
	Sex         string `json:"sex" sql:"sex"`
	Password    string `json:"password" validate:"required" sql:"password"`
}

func NewUser(Surname string, Name string, Middlename string, PhoneNumber string, Email string, Sex string, Password string) *User {
	return &User{Surname: Surname, Name: Name, Middlename: Middlename, PhoneNumber: PhoneNumber, Email: Email, Sex: Sex, Password: Password}
}

//func (u *User) AfterDelete(tx *gorm.DB) (err error) {
//	tx.Clauses(clause.Returning{}).Where("user_id = ?", u.ID).Delete(&User{})
//	return
//}
