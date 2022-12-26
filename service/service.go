package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"unicode/utf8"
	"userService/service/storage"
)

type Service struct {
	store storage.PostgresStorage
}

func NewService() *Service {
	store, _ := storage.NewPostgres()
	srv := Service{*store}
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
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	userId, err := s.store.AddUser(crtUsr.Name, crtUsr.Age, crtUsr.Friends)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
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
	err = s.store.AddFriend(mF.SourceID, mF.TargetID)
	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte(err.Error()))
		return
	}
	firstFriend, err := s.store.GetUser(mF.SourceID)
	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte(err.Error()))
		return
	}
	secondFriend, err := s.store.GetUser(mF.TargetID)
	if err != nil {
		w.WriteHeader(http.StatusAlreadyReported)
		w.Write([]byte(err.Error()))
		return
	}
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
	friends, err := s.store.FriendsToStr(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
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
	nameRemote, err := s.store.GetUser(dU.TargetID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	err = s.store.DeleteUser(dU.TargetID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("A User with the name was deleted:%v\nUser_id:%v\n", nameRemote.GetName(), dU.TargetID)))
}

func (s *Service) Put(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(trimFirstRune(r.URL.Path))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	decoder := json.NewDecoder(r.Body)
	var updateAge updateUserAge
	err = decoder.Decode(&updateAge)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	s.store.UpdateAge(id, updateAge.Age)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User's age has been successfully updated"))
}

func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
