package mapper

import (
	openapi_types "github.com/oapi-codegen/runtime/types"

	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	listslot "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/slot/list"
)

func ToListSlotsRequest(roomID oapi.RoomIdPath, date openapi_types.Date) *listslot.Request {
	return &listslot.Request{
		RoomID: roomID,
		Date:   date.Time.UTC(),
	}
}

func ToListSlotsResponse(resp *listslot.Response) oapi.GetRoomsRoomIdSlotsList200JSONResponse {
	slots := make([]oapi.Slot, 0)
	if resp != nil {
		slots = make([]oapi.Slot, 0, len(resp.Slots))
		for _, item := range resp.Slots {
			slots = append(slots, oapi.Slot{
				Id:     item.ID,
				RoomId: item.RoomID,
				Start:  item.StartAt.UTC(),
				End:    item.EndAt.UTC(),
			})
		}
	}

	return oapi.GetRoomsRoomIdSlotsList200JSONResponse{
		Slots: &slots,
	}
}
