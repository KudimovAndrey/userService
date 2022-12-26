package service

import (
	"fmt"
)

type User struct {
	userID  int
	name    string
	age     int
	friends []int
}

func (u *User) SetID(id int) {
	u.userID = id
}
func (u *User) SetName(name string) {
	u.name = name
}
func (u *User) SetAge(age int) {
	u.age = age
}
func (u *User) SetFriends(friends []int) {
	u.friends = friends
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
	return fmt.Sprintf("\nUser_id:%d\nName:%s\nAge:%d\nFriends:%v\n", u.GetID(), u.GetName(), u.GetAge(), u.GetFriends())
}
