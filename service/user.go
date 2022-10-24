package service

import (
	"errors"
	"fmt"
)

type user struct {
	name    string
	age     int
	friends []int
}

func (u *user) GetName() string {
	return u.name
}

func (u *user) GetAge() int {
	return u.age
}

func (u *user) GetFriends() []int {
	return u.friends
}

func (u *user) NewAge(age int) {
	u.age = age
}

func (u *user) AddFriend(friendId int) error {
	if contains(u.friends, friendId) {
		return errors.New("friend has already been added")
	}
	u.friends = append(u.friends, friendId)
	return nil
}

func (u *user) DeleteFriend(friendId int) error {
	index := getIndexElement(u.friends, friendId)
	if index == -1 {
		return errors.New("not found")
	}
	u.friends = append(u.friends[:index], u.friends[index+1:]...)
	return nil
}

func (u *user) ToString() string {
	return fmt.Sprintf("\nName:%s\nAge:%d\nFriends:%v\n", u.GetName(), u.GetAge(), u.GetFriends())
}

func contains(arrayInt []int, desiredValue int) bool {
	for _, value := range arrayInt {
		if value == desiredValue {
			return true
		}
	}
	return false
}

func getIndexElement(arrayInt []int, desiredValue int) int {
	for i, value := range arrayInt {
		if value == desiredValue {
			return i
		}
	}
	return -1
}
