package myerrors

type RoomError struct {
	Code int
	Description string
}

func (err *RoomError) Error() string {
	return err.Description
}
func NewRoomError() *RoomError {
	return &RoomError{
		Code: 0,
		Description: "room name too long",
	}
}
