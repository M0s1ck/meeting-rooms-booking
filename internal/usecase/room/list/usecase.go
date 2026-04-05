package list

import (
	"context"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/room"
)

type Usecase struct {
	repo roomRepo
}

func NewUsecase(repo roomRepo) *Usecase {
	return &Usecase{repo: repo}
}

func (u *Usecase) Execute(ctx context.Context) (*Response, error) {
	rooms, err := u.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	resp := buildResp(rooms)
	return resp, nil
}

func buildResp(rooms []*room.Room) *Response {
	resp := &Response{
		Rooms: make([]Room, 0, len(rooms)),
	}

	for _, item := range rooms {
		if item == nil {
			continue
		}

		resp.Rooms = append(resp.Rooms, Room{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			Capacity:    item.Capacity,
			CreatedAt:   item.CreatedAt,
		})
	}

	return resp
}
