package models

type CreateRoom struct {
	//OwnerID     uuid.UUID `json:"owner_id,omitempty" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description,omitempty" db:"description"`
}

type Room struct {
	ID int `json:"id" db:"id"`
	//OwnerID     uuid.UUID `json:"owner_id,omitempty" db:"user_id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

func NewRoom(id int, hub *CreateRoom) *Room {
	return &Room{
		ID:          id,
		Name:        hub.Name,
		Description: hub.Description,
	}
}
