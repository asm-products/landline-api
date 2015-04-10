package models

import (
	"fmt"
	"testing"
	"time"
)

func TestDeleteRoom(t *testing.T) {
	room := makeFakeRoom()
	room.Id = "TestDeleteRoom-1"
	room.Slug = "TestDeleteRoom-slug1"
	room.TeamId = "TestDeleteRoom-team1"
	_ = insertFakeRoom(room, t)

	err := DeleteRoom(room.Slug, room.TeamId)
	if err != nil {
		t.Error("TestDeleteRoom error:", err)
	}

	result := Room{}
	err = Db.SelectOne(&result, "select * from rooms where id=$1", room.Id)
	if err != nil {
		t.Error("TestDeleteRoom error:", err)
	}
	if result.DeletedAt == nil {
		t.Error("TestDeleteRoom: DeletedAt should be updated (%+v)", result)
	}
	result.setTime(room.CreatedAt)
	if *room != result {
		t.Errorf("TestDeleteRoom: got (%+v), wanted (%+v)", result, *room)
	}
}

func TestFindOrCreateRoom(t *testing.T) {
	room := makeFakeRoom()
	room.Id = "TestFindOrCreateRoom-1"
	room.Slug = "TestFindOrCreateRoom-slug1"
	room.TeamId = "TestFindOrCreateRoom-team1"

	_, err := FindOrCreateRoom(room)
	if err != nil {
		t.Error("TestFindOrCreateRoom error:", err)
	}

	result := Room{}
	err = Db.SelectOne(&result, "select * from rooms where id=$1", room.Id)
	if err != nil {
		t.Error("TestFindOrCreateRoom error:", err)
	}
	result.setTime(room.CreatedAt)
	if *room != result {
		t.Errorf("TestFindOrCreateRoom: got (%+v), wanted (%+v)", result, *room)
	}
}

func TestFindRoom(t *testing.T) {
	room := makeFakeRoom()
	room.Id = "TestFindRoom-1"
	room.Slug = "TestFindRoom-slug1"
	room.TeamId = "TestFindRoom-team1"
	_ = insertFakeRoom(room, t)

	result, err := FindRoom(room.Slug, room.TeamId)
	if err != nil {
		t.Error("TestFindRoom error:", err)
	}

	result.setTime(room.CreatedAt)
	if *room != *result {
		t.Errorf("TestFindRoom: got (%+v), wanted (%+v)", result, room)
	}
}

func TestFindRooms(t *testing.T) {
	teamId := "TestFindRooms-team"
	rooms := []*Room{makeFakeRoom(), makeFakeRoom(), makeFakeRoom()}
	for i, room := range rooms {
		room.Id = fmt.Sprintf("TestFindRooms-%d", i)
		room.Slug = fmt.Sprintf("TestFindRooms-slug%d", i)
		room.TeamId = teamId
		_ = insertFakeRoom(room, t)
	}

	result, err := FindRooms(teamId)
	if err != nil {
		t.Error("TestFindRooms error:", err)
	}
	if len(result) != len(rooms) {
		t.Fatalf("TestFindRooms result length: got %d, want %d", len(result), len(rooms))
	}
	for i, room := range result {
		room.setTime(rooms[i].CreatedAt)
		if *rooms[i] != room {
			t.Errorf("TestFindRooms: got (%+v), want (%+v)", room, *rooms[i])
		}
	}
}

func TestFindRoomById(t *testing.T) {
	room := makeFakeRoom()
	room.Id = "TestFindRoomById-1"
	_ = insertFakeRoom(room, t)

	result := FindRoomById(room.Id)

	result.setTime(room.CreatedAt)
	if *room != *result {
		t.Errorf("TestFindRoomById: got (%+v), wanted (%+v)", result, room)
	}
}

func TestUpdateRoom(t *testing.T) {
	room := makeFakeRoom()
	room.Id = "TestUpdateRoom-1"
	room.TeamId = "TestUpdateRoom-team1"
	room.Slug = "TestUpdateRoom-slug1"
	_ = insertFakeRoom(room, t)

	fields := makeFakeRoom()
	fields.Slug = "TestUpdateRoom-slug2"
	fields.Topic = "TestUpdateRoom-topic2"

	result, err := UpdateRoom(room.Slug, room.TeamId, fields)
	if err != nil {
		t.Fatal("TestUpdateRoom:", err)
	}
	room.Slug = fields.Slug
	room.Topic = fields.Topic
	result.setTime(room.CreatedAt)
	if *room != *result {
		t.Errorf("TestUpdateRoom: got (%v), want (%v)", result, room)
	}
}

func insertFakeRoom(room *Room, t *testing.T) *Room {
	err := Db.Insert(room)
	if err != nil {
		t.Fatal("Insert fake room error:", err)
	}
	return room
}

func (o *Room) setTime(t time.Time) {
	o.CreatedAt = t
	o.UpdatedAt = t
	o.DeletedAt = nil
}

func makeFakeRoom() *Room {
	t := time.Now()
	return &Room{
		Id:        "1",
		CreatedAt: t,
		UpdatedAt: t,
		DeletedAt: nil,
		TeamId:    "teamid1",
		Slug:      "slug1",
		Topic:     "topic1",
	}
}
