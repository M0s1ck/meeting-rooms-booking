package conferenceservice

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"

	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/booking/create"
)

type MockProvider struct {
	FailureRate float64
	Random      *rand.Rand
}

func NewMockProvider(failureRate float64) *MockProvider {
	return &MockProvider{
		FailureRate: failureRate,
		Random:      rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *MockProvider) CreateLink(_ context.Context, slotID uuid.UUID, userID uuid.UUID) (string, error) {
	if m.Random.Float64() < m.FailureRate {
		if m.Random.Float64() < 0.5 {
			return "", create.ErrConferenceServiceUnavailable
		}
		return "", create.ErrConferenceInternal
	}

	link := fmt.Sprintf(
		"https://mock-conference.local/%s/%s/%s",
		slotID.String(),
		userID.String(),
		uuid.New().String(),
	)

	return link, nil
}
