package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"strings"
)

//type userStorage struct {
//	users map[int]*user
//}

//func NewStorage() *userStorage {
//	uSt := userStorage{make(map[int]*user)}
//	return &uSt
//}

//func (us *userStorage) AddFriend(sourceId, targetId int) error {
//	err := us.users[sourceId].AddFriend(targetId)
//	if err != nil {
//		return err
//	}
//	err = us.users[targetId].AddFriend(sourceId)
//	return err
//}

//func (us *userStorage) GetUser(userId int) *user {
//	return us.users[userId]
//}

//func (us *userStorage) AddUser(name string, age int, friends []int) int {
//	newUser := user{name, age, friends}
//	us.users[len(us.users)] = &newUser
//	return len(us.users) - 1
//}
//
//func (us *userStorage) FriendsToStr(userId int) string {
//	result := strings.Builder{}
//	for _, friendId := range us.users[userId].friends {
//		result.WriteString(us.UserToStr(friendId))
//	}
//	return result.String()
//}
//func (us *userStorage) UserToStr(userId int) string {
//	user := us.GetUser(userId)
//	return fmt.Sprintf("\nuser_id:%v%v", userId, user.ToString())
//}
//
//func (us *userStorage) DeleteUser(userId int) error {
//	for _, friend := range us.users[userId].friends {
//		err := us.users[friend].DeleteFriend(userId)
//		if err != nil {
//			return err
//		}
//	}
//	delete(us.users, userId)
//	return nil
//}

type userStorage struct {
	dbStore *pgx.Conn
}

// NewStorage TODO: handle errors
func NewStorage() (*userStorage, error) {
	urlDB, err := getUrlDB()
	if err != nil {
		return nil, err
	}
	conn, err := pgx.Connect(context.Background(), urlDB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	return &userStorage{conn}, nil
}

// AddFriend TODO: implement re-entry verification
func (us *userStorage) AddFriend(sourceId, targetId int) error {

	if us.friendCheck(targetId, sourceId) {
		return errors.New("friend has already been added")
	}
	_, err := us.dbStore.Exec(context.Background(), ""+
		"UPDATE public.users "+
		"SET  \"friends\"= array_append(\"friends\", $1) "+
		"WHERE \"user_id\"=$2;", sourceId, targetId)
	if err != nil {
		return err
	}

	if us.friendCheck(sourceId, targetId) {
		return errors.New("friend has already been added")
	}
	_, err = us.dbStore.Exec(context.Background(), ""+
		"UPDATE public.users "+
		"SET  \"friends\"= array_append(\"friends\", $2) "+
		"WHERE \"user_id\"=$1;", sourceId, targetId)
	return err
}

// GetUser TODO: handle errors
func (us *userStorage) GetUser(userId int) (*user, error) {
	rows, err := us.dbStore.Query(context.Background(), ""+
		"SELECT \"user_id\", \"name\", \"age\", \"friends\" "+
		"FROM users "+
		"WHERE \"user_id\" = $1;", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var u user
	for rows.Next() {
		err = rows.Scan(&u.userID, &u.name, &u.age, &u.friends)
		if err != nil {
			return nil, err
		}
	}
	return &u, nil
}

// AddUser TODO: handle errors
func (us *userStorage) AddUser(name string, age int, friends []int) (int, error) {
	_, err := us.dbStore.Exec(context.Background(), "INSERT INTO public.users(\"name\", \"age\", \"friends\") "+
		"VALUES ($1,$2,$3);", name, age, friends)
	if err != nil {
		return -1, err
	}
	rows, err := us.dbStore.Query(context.Background(), ""+
		"select \"user_id\" "+
		"from users "+
		"order by \"user_id\" desc "+
		"limit 1")
	if err != nil {
		return -1, err
	}
	var resultID int
	for rows.Next() {
		err = rows.Scan(&resultID)
		if err != nil {
			return -1, err
		}
	}
	return resultID, nil
}

// FriendsToStr TODO: finish him
func (us *userStorage) FriendsToStr(userId int) (string, error) {
	result := strings.Builder{}
	u, err := us.GetUser(userId)
	if err != nil {
		return "", err
	}
	for _, id := range u.GetFriends() {
		friend, err := us.GetUser(id)
		if err != nil {
			return "", err
		}
		result.WriteString(friend.ToString())
	}
	return result.String(), nil
}

func (us *userStorage) DeleteUser(userId int) error {
	_, err := us.dbStore.Exec(context.Background(), ""+
		"UPDATE users "+
		"SET \"friends\"=array_remove(\"friends\",$1) "+
		"WHERE  $1 = Any(\"friends\"::int[]);", userId)
	if err != nil {
		return err
	}
	_, err = us.dbStore.Exec(context.Background(), "DELETE FROM public.users "+
		"WHERE \"user_id\"=$1;", userId)
	return err
}

func (us *userStorage) UpdateAge(userId, age int) error {
	_, err := us.dbStore.Exec(context.Background(), ""+
		"UPDATE users "+
		"SET \"age\"=$1 "+
		"WHERE \"user_id\"=$2;", age, userId)
	return err
}

func (us *userStorage) friendCheck(sourceId, targetId int) bool {
	u, _ := us.GetUser(sourceId)
	return contains(u.GetFriends(), targetId)
}

func getUrlDB() (string, error) {
	result := ""
	file, err := os.Open("linkFromDB.txt")
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		result += scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}
	return result, nil
}

func contains(arrayInt []int, desiredValue int) bool {
	for _, value := range arrayInt {
		if value == desiredValue {
			return true
		}
	}
	return false
}
