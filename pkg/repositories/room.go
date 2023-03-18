package repositories

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
	"github.com/mrbttf/bridge-server/pkg/core"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

const SelectRoomById = `
SELECT room_id, host_id, user_ids, open
FROM rooms
WHERE room_id = $1
`

func (rr *RoomRepository) Get(room_id string) (core.Room, error) {
	var room core.Room

	err := rr.db.QueryRow(SelectRoomById, room_id).Scan(
		&room.Id,
		&room.Host,
		pq.Array(&room.Users),
		&room.Open,
	)
	if err != nil {
		return core.Room{}, fmt.Errorf("Unable to get room for id %s: %w", room_id, err)
	}

	return room, nil
}

const SelectRooms = `
SELECT room_id, host_id, user_ids, open
FROM rooms
WHERE open = $1
`

func (rr *RoomRepository) List(open bool) ([]core.Room, error) {
	rows, err := rr.db.Query(SelectRooms, open)
	if err != nil {
		return nil, fmt.Errorf("Unable to list rooms for open %t: %w", open, err)
	}
	defer rows.Close()

	var rooms []core.Room
	for rows.Next() {
		var room core.Room
		if err := rows.Scan(
			&room.Id,
			&room.Host,
			pq.Array(&room.Users),
			&room.Open,
		); err != nil {
			return nil, fmt.Errorf("Unable to list rooms for open %t: %w", open, err)
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("Unable to list rooms for open %t: %w", open, err)
	}
	return rooms, nil
}

const UpsertRoom = `
INSERT INTO rooms (room_id, host_id, user_ids, open)
VALUES($1, $2, $3, $4) 
ON CONFLICT (room_id) 
WHERE room_id = $1 
DO UPDATE
SET 
	host_id = EXCLUDED.host_id, 
	user_ids = EXCLUDED.user_ids, 
	open = EXCLUDED.open 
`

func (rr *RoomRepository) Store(room *core.Room) error {
	_, err := rr.db.Exec(UpsertRoom,
		room.Id,
		room.Host,
		pq.Array(room.Users),
		room.Open,
	)
	if err != nil {
		return err
	}

	return nil
}

const DeleteRoom = `
DELETE FROM rooms
WHERE room_id = $1
`

func (rr *RoomRepository) Delete(room_id string) error {
	_, err := rr.db.Exec(DeleteRoom, room_id)
	if err != nil {
		return fmt.Errorf("Unable to delete room for id %s: %w", room_id, err)
	}

	return nil
}
