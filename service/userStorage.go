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

func (u *userStorage) AddFriend(sourceId, targetId int) {
	u.users[sourceId] = u.users[sourceId].AddFriend(targetId)
	u.users[targetId] = u.users[targetId].AddFriend(sourceId)
}

// getuser
func (u *userStorage) GetUser(userId int) string {
	res := u.users[userId]
	return fmt.Sprintf("\nUser_id:%d\nName:%s\nAge:%d\nFriends:%v\n", userId, res.GetName(), res.GetAge(), res.GetFriends())
}

func (u *userStorage) AddUser(name string, age int, friends []int) int {
	newUser := user{name, age, friends}
	u.users[len(u.users)] = &newUser
	return len(u.users) - 1
}

func (u *userStorage) GetFriendsToStr(userId int) string {
	result := strings.Builder{}
	for _, friendId := range u.users[userId].friends {
		result.WriteString(u.GetUser(friendId))
	}
	return result.String()
}

func (u *userStorage) DeleteUser(userId int) {
	delete(u.users, userId)
}
