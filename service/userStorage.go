package service

import (
	"fmt"
	"strings"
)

type userStorage struct {
	users map[int]*user
}

func NewStorage() *userStorage {
	uSt := userStorage{make(map[int]*user)}
	return &uSt
}

func (us *userStorage) AddFriend(sourceId, targetId int) error {
	err := us.users[sourceId].AddFriend(targetId)
	if err != nil {
		return err
	}
	err = us.users[targetId].AddFriend(sourceId)
	return err
}

func (us *userStorage) GetUser(userId int) *user {
	return us.users[userId]
}

func (us *userStorage) AddUser(name string, age int, friends []int) int {
	newUser := user{name, age, friends}
	us.users[len(us.users)] = &newUser
	return len(us.users) - 1
}

func (us *userStorage) FriendsToStr(userId int) string {
	result := strings.Builder{}
	for _, friendId := range us.users[userId].friends {
		result.WriteString(us.UserToStr(friendId))
	}
	return result.String()
}
func (us *userStorage) UserToStr(userId int) string {
	user := us.GetUser(userId)
	return fmt.Sprintf("\nuser_id:%v%v", userId, user.ToString())
}

func (us *userStorage) DeleteUser(userId int) error {
	for _, friend := range us.users[userId].friends {
		err := us.users[friend].DeleteFriend(userId)
		if err != nil {
			return err
		}
	}
	delete(us.users, userId)
	return nil
}
