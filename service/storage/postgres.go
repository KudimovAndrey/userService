package storage

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"os"
	"strings"
)

type PostgresStorage struct {
	dbStore *pgx.Conn
}

func NewPostgres() (*PostgresStorage, error) {
	urlDB, err := getUrlDB()
	if err != nil {
		return nil, err
	}
	conn, err := pgx.Connect(context.Background(), urlDB)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return nil, err
	}
	return &PostgresStorage{conn}, nil
}

func (ps *PostgresStorage) AddFriend(sourceId, targetId int) error {
	isFriends, err := ps.friendCheck(targetId, sourceId)
	if err != nil {
		return err
	}
	if isFriends {
		return errors.New("friend has already been added")
	}
	_, err = ps.dbStore.Exec(context.Background(), ""+
		"UPDATE public.users "+
		"SET  \"friends\"= array_append(\"friends\", $1) "+
		"WHERE \"user_id\"=$2;", sourceId, targetId)
	if err != nil {
		return err
	}
	isFriends, err = ps.friendCheck(sourceId, targetId)
	if err != nil {
		return err
	}
	if isFriends {
		return errors.New("friend has already been added")
	}
	_, err = ps.dbStore.Exec(context.Background(), ""+
		"UPDATE public.users "+
		"SET  \"friends\"= array_append(\"friends\", $2) "+
		"WHERE \"user_id\"=$1;", sourceId, targetId)
	return err
}

func (ps *PostgresStorage) GetUser(userId int) (*User, error) {
	rows, err := ps.dbStore.Query(context.Background(), ""+
		"SELECT \"user_id\", \"name\", \"age\", \"friends\" "+
		"FROM users "+
		"WHERE \"user_id\" = $1;", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var u User
	for rows.Next() {
		err = rows.Scan(&u.userID, &u.name, &u.age, &u.friends)
		if err != nil {
			return nil, err
		}
	}
	return &u, nil
}

func (ps *PostgresStorage) AddUser(name string, age int, friends []int) (int, error) {
	_, err := ps.dbStore.Exec(context.Background(), "INSERT INTO public.users(\"name\", \"age\", \"friends\") "+
		"VALUES ($1,$2,$3);", name, age, friends)
	if err != nil {
		return -1, err
	}
	rows, err := ps.dbStore.Query(context.Background(), ""+
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

func (ps *PostgresStorage) GetUsers(usersID ...int) (map[int]User, error) {
	result := make(map[int]User)
	rows, err := ps.dbStore.Query(context.Background(), ""+
		"SELECT *FROM users "+
		"WHERE \"user_id\" = Any($1::integer[]) "+
		"order BY \"user_id\", \"name\", \"age\", \"friends\"", usersID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var u User
	for i := 0; rows.Next(); i++ {
		err = rows.Scan(&u.userID, &u.name, &u.age, &u.friends)
		if err != nil {
			return nil, err
		}
		result[i] = u
	}
	return result, nil
}

func (ps *PostgresStorage) FriendsToStr(userId int) (string, error) {
	result := strings.Builder{}
	u, err := ps.GetUser(userId)
	if err != nil {
		return "", err
	}
	friends, err := ps.GetUsers(u.GetFriends()...)
	for _, user := range friends {
		_, err = result.WriteString(user.ToString())
		if err != nil {
			return "", nil
		}
	}
	return result.String(), nil
}

func (ps *PostgresStorage) DeleteUser(userId int) error {
	_, err := ps.dbStore.Exec(context.Background(), ""+
		"UPDATE users "+
		"SET \"friends\"=array_remove(\"friends\",$1) "+
		"WHERE  $1 = Any(\"friends\"::int[]);", userId)
	if err != nil {
		return err
	}
	_, err = ps.dbStore.Exec(context.Background(), "DELETE FROM public.users "+
		"WHERE \"user_id\"=$1;", userId)
	return err
}

func (ps *PostgresStorage) UpdateAge(userId, age int) error {
	_, err := ps.dbStore.Exec(context.Background(), ""+
		"UPDATE users "+
		"SET \"age\"=$1 "+
		"WHERE \"user_id\"=$2;", age, userId)
	return err
}

func (ps *PostgresStorage) friendCheck(sourceId, targetId int) (bool, error) {
	u, err := ps.GetUser(sourceId)
	if err != nil {
		return false, err
	}
	return contains(u.GetFriends(), targetId), nil
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
