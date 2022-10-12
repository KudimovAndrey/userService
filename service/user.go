package service

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

// NewAge TODO: pointer method
func (u user) NewAge(age int) *user {
	u.age = age
	return &u
}

func (u user) AddFriend(friendId int) *user {
	u.friends = append(u.friends, friendId)
	return &u
}

// gunc toString
