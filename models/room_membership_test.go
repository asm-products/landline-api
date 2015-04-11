package models

import (
	"fmt"
	"testing"
	"time"
)

func TestDeleteRoomMembership(t *testing.T) {
	membership := makeFakeRoomMembership()
	membership.Id = "TestDeleteRoomMembership-1"
	membership.RoomId = "TestDeleteRoomMembership-room1"
	membership.UserId = "TestDeleteRoomMembership-user1"
	_ = insertFakeRoomMembership(membership, t)

	err := DeleteRoomMembership(membership.RoomId, membership.UserId)
	if err != nil {
		t.Error("TestDeleteRoomMembership error:", err)
	}

	result := RoomMembership{}
	err = Db.SelectOne(&result, "select * from room_memberships where id=$1", membership.Id)
	if err != nil {
		t.Error("TestDeleteRoomMembership error:", err)
	}
	if result.DeletedAt == nil {
		t.Error("TestDeleteRoomMembership: DeletedAt should be updated (%+v)", result)
	}
	result.setTime(membership.CreatedAt)
	if *membership != result {
		t.Errorf("TestDeleteRoomMembership: got (%+v), wanted (%+v)", result, *membership)
	}
}

func TestFindOrCreateRoomMembership(t *testing.T) {
	membership := makeFakeRoomMembership()
	membership.Id = "TestFindOrCreateRoomMembership-1"
	membership.RoomId = "TestFindOrCreateRoomMembership-room1"
	membership.UserId = "TestFindOrCreateRoomMembership-user1"
	_ = insertFakeRoomMembership(membership, t)
	err := DeleteRoomMembership(membership.RoomId, membership.UserId)
	if err != nil {
		t.Error("TestFindOrCreateRoomMembership error:", err)
	}

	result, err := FindOrCreateRoomMembership(membership)
	if err != nil {
		t.Error("TestFindOrCreateRoomMembership error:", err)
	}
	if result.DeletedAt != nil {
		t.Error("TestFindOrCreateRoomMembership: DeletedAt should be nil (%+v)", result)
	}
	result.setTime(membership.CreatedAt)
	if *membership != *result {
		t.Errorf("TestFindOrCreateRoomMembership: got (%+v), wanted (%+v)", result, membership)
	}
}

func TestFindRoomMemberships(t *testing.T) {
	userId := "TestFindRoomMemberships-user"
	memberships := []*RoomMembership{makeFakeRoomMembership(), makeFakeRoomMembership(), makeFakeRoomMembership()}
	for i, membership := range memberships {
		membership.Id = fmt.Sprintf("TestFindRoomMemberships-%d", i)
		membership.RoomId = fmt.Sprintf("TestFindRoomMemberships-room%d", i)
		membership.UserId = userId
		_ = insertFakeRoomMembership(membership, t)
	}
	membership := makeFakeRoomMembership()
	membership.Id = "TestFindRoomMemberships-100"
	_ = insertFakeRoomMembership(membership, t)

	result, err := FindRoomMemberships(userId)
	if err != nil {
		t.Error("TestFindRoomMemberships error:", err)
	}
	if len(result) != len(memberships) {
		t.Fatalf("TestFindRoomMemberships result length: got %d, want %d", len(result), len(memberships))
	}
	for i, roomId := range result {
		if memberships[i].RoomId != roomId {
			t.Errorf("TestFindRoomMemberships: got (%s), want (%s)", roomId, memberships[i].RoomId)
		}
	}
}

func (o *RoomMembership) setTime(t time.Time) {
	o.CreatedAt = t
	o.UpdatedAt = t
	o.DeletedAt = nil
}

func insertFakeRoomMembership(roomMembership *RoomMembership, t *testing.T) *RoomMembership {
	err := Db.Insert(roomMembership)
	if err != nil {
		t.Fatal("Insert fake room membership error:", err)
	}
	return roomMembership
}

func makeFakeRoomMembership() *RoomMembership {
	t := time.Now()
	return &RoomMembership{
		Id:        "1",
		CreatedAt: t,
		DeletedAt: nil,
		UpdatedAt: t,
		RoomId:    "roomid",
		UserId:    "userid",
	}
}
