package room

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/mrbttf/bridge-server/pkg/core"
)

type RoomService struct {
	rooms core.RoomRepository
	users core.UserRepository
}

func New(rooms core.RoomRepository, users core.UserRepository) *RoomService {
	return &RoomService{
		rooms: rooms,
		users: users,
	}
}

func (rs *RoomService) Create(host_id string) (string, error) {
	room_id := uuid.New().String()

	room := &core.Room{
		Id:    room_id,
		Host:  host_id,
		Users: []string{host_id},
		Open:  true,
	}
	err := rs.rooms.Store(room)
	if err != nil {
		return "", fmt.Errorf("Unable to create room: %w", err)
	}
	return room_id, nil
}

func (rs *RoomService) Get(room_id string) (core.Room, error) {
	return rs.rooms.Get(room_id)
}

func (rs *RoomService) GetUsers(room_id string) ([]core.User, error) {
	return rs.users.GetForRoom(room_id)
}

func (rs *RoomService) Join(room_id string, user_id string) error {
	room, err := rs.rooms.Get(room_id)
	if err != nil {
		return fmt.Errorf("Unable to join room, room_id %s, user_id %s: %w", room_id, user_id, err)
	}
	room.Users = append(room.Users, user_id)
	err = rs.rooms.Store(&room)
	if err != nil {
		return fmt.Errorf("Unable to join room, room_id %s, user_id %s: %w", room_id, user_id, err)
	}
	return nil
}

func (rs *RoomService) List(open bool) ([]core.Room, error) {
	return rs.rooms.List(open)
}

func (rs *RoomService) Close(room_id string) error {
	room, err := rs.rooms.Get(room_id)
	if err != nil {
		return fmt.Errorf("Unable to close room, room_id %s: %w", room_id, err)
	}
	room.Open = false
	err = rs.rooms.Store(&room)
	if err != nil {
		return fmt.Errorf("Unable to close room, room_id %s: %w", room_id, err)
	}
	return nil
}

func (rs *RoomService) Delete(room_id string) error {
	return rs.rooms.Delete(room_id)
}
