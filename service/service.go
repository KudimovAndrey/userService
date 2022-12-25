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
	storage, _ := NewStorage()
	srv := Service{*storage}
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
	var crtUsr createUser
	err := decoder.Decode(&crtUsr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	userId, _ := s.storage.AddUser(crtUsr.Name, crtUsr.Age, crtUsr.Friends)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User was created"))
	w.Write([]byte(fmt.Sprintf("\nuser_id:%v", userId)))
}

func (s *Service) MakeFriends(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var mF makeFriends
	err := decoder.Decode(&mF)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = s.storage.AddFriend(mF.SourceID, mF.TargetID)
	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte(err.Error()))
		return
	}
	firstFriend, _ := s.storage.GetUser(mF.SourceID)
	secondFriend, _ := s.storage.GetUser(mF.TargetID)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("%v и %v теперь друзья", firstFriend.GetName(), secondFriend.GetName())))
}

func (s *Service) Get(w http.ResponseWriter, r *http.Request) {
	body := trimFirstRune(r.URL.Path)
	id, err := strconv.Atoi(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	friends, err := s.storage.FriendsToStr(id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(friends))
}

func (s *Service) Delete(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var dU deleteUser
	err := decoder.Decode(&dU)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	nameRemote, _ := s.storage.GetUser(dU.TargetID)
	err = s.storage.DeleteUser(dU.TargetID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("A user with the name was deleted:%v\nUser_id:%v\n", nameRemote.GetName(), dU.TargetID)))
}

func (s *Service) Put(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(trimFirstRune(r.URL.Path))
	decoder := json.NewDecoder(r.Body)
	var updateAge updateUserAge
	err := decoder.Decode(&updateAge)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	s.storage.UpdateAge(id, updateAge.Age)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("user's age has been successfully updated"))
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
