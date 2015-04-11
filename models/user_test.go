package models

import (
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"testing"
	"time"
)

func TestFindOrCreateUserByExternalId(t *testing.T) {
	user := makeFakeUser()
	user.Id = "TestFindOrCreateUserByExternalId-1"
	user.ExternalId = "TestFindOrCreateUserByExternalId-ex1"
	result, err := FindOrCreateUserByExternalId(user)
	if err != nil {
		t.Fatal(err)
	}
	if *user != *result {
		t.Errorf("FindOrCreateUserByExternalId: making user, got (%v), want (%v)", result, user)
	}

	user1 := &User{
		ExternalId: user.ExternalId,
	}
	result, err = FindOrCreateUserByExternalId(user1)
	if err != nil {
		t.Fatal(err)
	}

	result.setTime(user.CreatedAt)
	if *user != *result {
		t.Errorf("FindOrCreateUserByExternalId: finding user, got (%v), want (%v)", result, user)
	}
}

func TestFindUser(t *testing.T) {
	id := "TestFindUser-1"
	user := makeFakeUser()
	user.Id = id
	err := Db.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	result, err := FindUser(id)
	if err != nil {
		t.Fatal(err)
	}
	result.setTime(user.CreatedAt)
	if *user != *result {
		t.Errorf("FindUser: got (%v), want (%v)", result, user)
	}

	id2 := "TestFindUser-2"
	_, err = FindUser(id2)
	if err == nil {
		t.Fatalf("FindUser: should have no such id: %s", id2)
	}
}

func TestFindUserByExternalIdAndTeam(t *testing.T) {
	id := "TestFindUserByExternalIdAndTeam-1"
	extId := "ExternalId-1"
	team := "Team1"
	user := makeFakeUser()
	user.Id = id
	user.ExternalId = extId
	user.TeamId = team
	err := Db.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	result, err := FindUserByExternalIDAndTeam(extId, team)
	if err != nil {
		t.Fatal(err)
	}
	result.setTime(user.CreatedAt)
	if *user != *result {
		t.Errorf("FindUserByExternalIDAndTeam: got (%v), want (%v)", result, user)
	}

	extId2 := "TestFindUserByExternalIdAndTeam-2"
	_, err = FindUserByExternalIDAndTeam(extId2, team)
	if err == nil {
		t.Fatalf("FindUserByExternalIDAndTeam: should have no such id: %s", extId2)
	}
}

func TestFindUserByUsernameAndTeam(t *testing.T) {
	id := "TestFindUserByUsernameAndTeam-1"
	username := "username"
	team := "Team1"
	user := makeFakeUser()
	user.Id = id
	user.TeamId = team
	user.Username = username
	err := Db.Insert(user)
	if err != nil {
		t.Fatal(err)
	}

	result, err := FindUserByUsernameAndTeam(username, team)
	if err != nil {
		t.Fatal(err)
	}
	result.setTime(user.CreatedAt)
	if *user != *result {
		t.Errorf("FindUserByUsernameAndTeam: got (%v), want (%v)", result, user)
	}

	username2 := "TestFindUserByUsernameAndTeam-2"
	_, err = FindUserByExternalIDAndTeam(username2, team)
	if err == nil {
		t.Fatalf("FindUserByUsernameAndTeam: should have no such id: %s", username2)
	}
}

func TestFindUsers(t *testing.T) {
	team := "TestFindUsers-team"
	users := []*User{makeFakeUser(), makeFakeUser(), makeFakeUser()}
	for i, user := range users {
		user.Id = fmt.Sprintf("TestFindUsers-%d", i)
		user.TeamId = team
		err := Db.Insert(user)
		if err != nil {
			t.Fatal(err)
		}
	}

	result, err := FindUsers(team)
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != len(result) {
		t.Fatalf("FindUsers: got %d results, want %d results", len(result), len(users))
	}
	for i, newUser := range result {
		origUser := users[i]
		newUser.setTime(origUser.CreatedAt)
		if *origUser != newUser {
			t.Errorf("FindUsers: got (%v), want (%v)", newUser, origUser)
		}
	}
}

// This method can't be fully tested due to preInsert/preUpdate hook
func TestFindRecentlyOnlineUsers(t *testing.T) {
	team := "TestFindRecentlyOnlineUsers-team"
	users := []*User{makeFakeUser(), makeFakeUser(), makeFakeUser()}
	for i, user := range users {
		user.Id = fmt.Sprintf("TestFindRecentlyOnlineUsers-%d", i)
		user.TeamId = team
		err := Db.Insert(user)
		if err != nil {
			t.Fatal(err)
		}
	}
	expiredUser := makeFakeUser()
	expiredUser.Id = "TestFindRecentlyOnlineUsers-3"
	expiredUser.TeamId = team
	expiredUser.LastOnlineAt = time.Now().Add(-time.Duration(2) * time.Hour)
	err := Db.Insert(expiredUser)
	if err != nil {
		t.Fatal(err)
	}

	result, err := FindUsers(team)
	if err != nil {
		t.Fatal(err)
	}
	/*if len(users) != len(result) {
	    t.Fatalf("FindRecentlyOnlineUsers: got %d results, want %d results", len(result), len(users))
	}*/
	for i, origUser := range users {
		newUser := result[i]
		newUser.setTime(origUser.CreatedAt)
		if *origUser != newUser {
			t.Errorf("FindRecentlyOnlineUsers: got (%v), want (%v)", newUser, origUser)
		}
	}
}

// used to make timestamps uniform since they are modified before insertion
func (o *User) setTime(t time.Time) {
	o.CreatedAt = t
	o.UpdatedAt = t
	o.LastOnlineAt = t
}

func makeFakeUser() *User {
	t := time.Now()
	return &User{
		Id:           "1",
		CreatedAt:    t,
		UpdatedAt:    t,
		LastOnlineAt: t,
		TeamId:       "team1",
		AvatarUrl:    "http://imgur.com/1.jpg",
		Email:        "user@host.com",
		ExternalId:   "ex1",
		ProfileUrl:   "http://imgur.com/2.jpg",
		RealName:     "Max User",
		Username:     "maxuser",
	}
}
