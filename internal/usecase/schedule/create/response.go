package create

import (
    "github.com/google/uuid"
    "github.com/internships-backend/test-backend-M0s1ck/internal/domain/schedule"
)

type Response struct {
    ID         uuid.UUID
    RoomID     uuid.UUID
    DaysOfWeek []int
    StartTime  schedule.TimeOfDay
    EndTime    schedule.TimeOfDay
}
