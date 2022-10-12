package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"
)

type Service struct {
	storage userStorage
}

func NewService() *Service {
	srv := Service{*NewStorage()}
	return &srv
}

func (s *Service) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.Post(w, r)
	case http.MethodGet:
		s.Get(w, r)
	case http.MethodDelete:
		s.Delete(w, r)
	case http.MethodPut:
		s.Put(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Service) Post(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var u newUser
	err := decoder.Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	userId := s.storage.AddUser(u.Name, u.Age, u.Friends)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User was created"))
	w.Write([]byte(fmt.Sprintf("\nuser_id:%v", userId)))
}

// MakeFriends TODO: i don't like this
func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mF makeFriends
	err := decoder.Decode(&mF)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	s.storage.AddFriend(mF.SourceID, mF.TargetID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v и %v теперь друзья", s.storage.users[mF.SourceID].name, s.storage.users[mF.TargetID].name)))
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	body := trimFirstRune(r.URL.Path)
	id, err := strconv.Atoi(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte(s.storage.GetUser(id)))
	friends := s.storage.GetFriendsToStr(id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(friends))
}

// Delete TODO:implement the removal of this user from all friends
func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var dU deleteUser
	err := decoder.Decode(&dU)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	nameRemote := s.storage.users[dU.TargetID].name
	s.storage.DeleteUser(dU.TargetID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("A user with the name was deleted:%v\nUser_id:%v\n", nameRemote, dU.TargetID)))
}

// Put TODO:looks like doesn't good
func (s *Service) Put(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(trimFirstRune(r.URL.Path))
	decoder := json.NewDecoder(r.Body)
	var nA newAge
	err := decoder.Decode(&nA)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	s.storage.users[id] = s.storage.users[id].NewAge(nA.Age)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user's age has been successfully updated"))
}

// TODO:she doesn't belong here
func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}