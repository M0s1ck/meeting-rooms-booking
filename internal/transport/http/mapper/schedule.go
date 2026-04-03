package mapper

import (
	"errors"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	createschedule "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/schedule/create"
)

var ErrNilCreateScheduleBody = errors.New("request body is required")

func ToCreateScheduleRequest(
	roomID oapi.RoomIdPath,
	body *oapi.PostRoomsRoomIdScheduleCreateJSONRequestBody,
) (*createschedule.Request, error) {
	if body == nil {
		return nil, ErrNilCreateScheduleBody
	}

	startTime, err := schedule.ParseTimeOfDay(body.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := schedule.ParseTimeOfDay(body.EndTime)
	if err != nil {
		return nil, err
	}

	return &createschedule.Request{
		RoomID:     roomID,
		DaysOfWeek: body.DaysOfWeek,
		StartTime:  startTime,
		EndTime:    endTime,
	}, nil
}

func ToCreateScheduleResponse(resp *createschedule.Response) oapi.PostRoomsRoomIdScheduleCreate201JSONResponse {
	if resp == nil {
		return oapi.PostRoomsRoomIdScheduleCreate201JSONResponse{}
	}

	id := resp.ID

	return oapi.PostRoomsRoomIdScheduleCreate201JSONResponse{
		Schedule: &oapi.Schedule{
			Id:         &id,
			RoomId:     resp.RoomID,
			DaysOfWeek: resp.DaysOfWeek,
			StartTime:  resp.StartTime.String(),
			EndTime:    resp.EndTime.String(),
		},
	}
}
