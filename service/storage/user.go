package storage

import (
	"fmt"
)

type User struct {
	userID  int
	name    string
	age     int
	friends []int
}

func (u *User) GetID() int {
	return u.userID
}

func (u *User) GetName() string {
	return u.name
}

func (u *User) GetAge() int {
	return u.age
}

func (u *User) GetFriends() []int {
	return u.friends
}

func (u *User) ToString() string {
	return fmt.Sprintf("\nUser_id:%d\nName:%s\nAge:%d\nFriends:%v\n", u.userID, u.name, u.age, u.friends)
}
