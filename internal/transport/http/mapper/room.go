package mapper

import (
	"errors"

	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	createroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/create"
	listroom "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/room/list"
)

var ErrNilCreateRoomBody = errors.New("request body is required")

func ToCreateRoomRequest(body *oapi.PostRoomsCreateJSONRequestBody) (*createroom.Request, error) {
	if body == nil {
		return nil, ErrNilCreateRoomBody
	}

	return &createroom.Request{
		Name:        body.Name,
		Description: body.Description,
		Capacity:    body.Capacity,
	}, nil
}

func ToCreateRoomResponse(resp *createroom.Response) oapi.PostRoomsCreate201JSONResponse {
	if resp == nil {
		return oapi.PostRoomsCreate201JSONResponse{}
	}

	createdAt := resp.CreatedAt

	return oapi.PostRoomsCreate201JSONResponse{
		Room: &oapi.Room{
			Id:          resp.ID,
			Name:        resp.Name,
			Description: resp.Description,
			Capacity:    resp.Capacity,
			CreatedAt:   &createdAt,
		},
	}
}

func ToListRoomsResponse(resp *listroom.Response) oapi.GetRoomsList200JSONResponse {
	if resp == nil {
		return oapi.GetRoomsList200JSONResponse{}
	}

	rooms := make([]oapi.Room, 0, len(resp.Rooms))
	for _, item := range resp.Rooms {
		createdAt := item.CreatedAt

		rooms = append(rooms, oapi.Room{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Capacity:    item.Capacity,
			CreatedAt:   &createdAt,
		})
	}

	return oapi.GetRoomsList200JSONResponse{
		Rooms: &rooms,
	}
}
