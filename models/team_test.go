package models

import (
	"testing"
	"time"
)

func TestFindTeamById(t *testing.T) {
	team := insertFakeTeamWithId("TestFindTeamById-1", t)
	result := FindTeamById(team.Id)
	result.setTime(team.CreatedAt)
	if *team != *result {
		t.Errorf("TestFindTeamById: got (%v), want (%v)", result, team)
	}
}

func TestFindTeamBySlug(t *testing.T) {
	team := makeFakeTeam()
	team.Id = "TestFindTeamBySlug-1"
	team.Slug = "TestFindTeamBySlug-slug"
	_ = insertFakeTeam(team, t)
	result := FindTeamBySlug(team.Slug)
	result.setTime(team.CreatedAt)
	if *team != *result {
		t.Errorf("TestFindTeamById: got (%v), want (%v)", result, team)
	}
}

func TestFindOrCreateTeam(t *testing.T) {
	team := makeFakeTeam()
	team.Id = "TestFindOrCreateTeam-1"
	team.Slug = "TestFindOrCreateTeam-slug1"
	result, err := FindOrCreateTeam(team)
	if err != nil {
		t.Fatal("TestFindOrCreateTeam:", err)
	}
	result.setTime(team.CreatedAt)
	if *team != *result {
		t.Errorf("TestFindOrCreateTeam: got (%v), want (%v)", result, team)
	}
}

func TestFindTeamBySecret(t *testing.T) {
	team := makeFakeTeam()
	team.Id = "TestFindTeamBySecret-1"
	team.Slug = "TestFindTeamBySecret-slug1"
	team.SSOSecret = "TestFindTeamBySecret-sso1"
	_ = insertFakeTeam(team, t)
	result := FindTeamBySecret(team.Slug, team.SSOSecret)
	result.setTime(team.CreatedAt)
	if *team != *result {
		t.Errorf("TestFindTeamBySecret: got (%v), want (%v)", result, team)
	}
}

func TestUpdateTeam(t *testing.T) {
	team := makeFakeTeam()
	team.Id = "TestUpdateTeam-1"
	team.Slug = "TestUpdateTeam-slug1"
	_ = insertFakeTeam(team, t)

	fields := makeFakeTeam()
	fields.Email = "TestUpdateTeam-email2"
	fields.Slug = "TestUpdateTeam-slug2"
	fields.SSOUrl = "TestUpdateTeam-ssourl2"
	fields.SSOSecret = "TestUpdateTeam-secret2"

	result, err := UpdateTeam(team.Slug, fields)
	if err != nil {
		t.Fatal("TestUpdateTeam:", err)
	}
	team.Email = fields.Email
	team.Slug = fields.Slug
	team.SSOUrl = fields.SSOUrl
	team.SSOSecret = fields.SSOSecret
	result.setTime(team.CreatedAt)
	if *team != *result {
		t.Errorf("TestUpdateTeam: got (%v), want (%v)", result, team)
	}
}

func (o *Team) setTime(t time.Time) {
	o.CreatedAt = t
	o.UpdatedAt = t
}

func insertFakeTeam(team *Team, t *testing.T) *Team {
	err := Db.Insert(team)
	if err != nil {
		t.Fatal("Insert fake team error:", err)
	}
	return team
}

func insertFakeTeamWithId(id string, t *testing.T) *Team {
	team := makeFakeTeam()
	team.Id = id
	err := Db.Insert(team)
	if err != nil {
		t.Fatal("Insert fake team error:", err)
	}
	return team
}

func makeFakeTeam() *Team {
	t := time.Now()
	return &Team{
		Id:                "1",
		CreatedAt:         t,
		UpdatedAt:         t,
		Email:             "email@host.com",
		EncryptedPassword: "password",
		SSOSecret:         "secret",
		SSOUrl:            "www.sso.url",
		Slug:              "team1",
	}
}
