package register_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/register"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/register/mocks"
)

func TestUsecase_Execute_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockuserRepo(ctrl)
	hasher := mocks.NewMockpassHasher(ctrl)

	uc := register.NewUsecase(repo, hasher)

	req := &register.Request{
		Email:    "test@example.com",
		Password: "abc12345",
		Role:     "user",
	}

	hasher.EXPECT().
		Hash("abc12345").
		Return("hash", nil)

	repo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		DoAndReturn(func(_ context.Context, u *user.User) error {
			require.Equal(t, user.Email("test@example.com"), u.Email)
			require.Equal(t, "hash", u.PassHash)
			require.Equal(t, user.RoleUser, u.Role)
			require.NotEqual(t, uuid.Nil, u.ID)
			return nil
		})

	resp, err := uc.Execute(context.Background(), req)

	require.NoError(t, err)
	require.NotNil(t, resp)
	require.Equal(t, user.Email("test@example.com"), resp.Email)
	require.Equal(t, user.RoleUser, resp.Role)
	require.NotEqual(t, uuid.Nil, resp.ID)
}

func TestUsecase_Execute_EmailAlreadyTaken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockuserRepo(ctrl)
	hasher := mocks.NewMockpassHasher(ctrl)

	uc := register.NewUsecase(repo, hasher)

	req := &register.Request{
		Email:    "test@example.com",
		Password: "abc12345",
		Role:     "user",
	}

	hasher.EXPECT().
		Hash("abc12345").
		Return("hash", nil)

	repo.EXPECT().
		Create(gomock.Any(), gomock.Any()).
		Return(user.ErrEmailAlreadyTaken)

	resp, err := uc.Execute(context.Background(), req)

	require.Error(t, err)
	require.ErrorIs(t, err, user.ErrEmailAlreadyTaken)
	require.Nil(t, resp)
}

func TestUsecase_Execute_ValidationErrors(t *testing.T) {
	tests := []struct {
		name    string
		req     *register.Request
		setup   func(hasher *mocks.MockpassHasher)
		wantErr error
	}{
		{
			name: "invalid email",
			req: &register.Request{
				Email:    "bad-email",
				Password: "abc12345",
				Role:     "user",
			},
			wantErr: user.ErrInvalidEmail,
		},
		{
			name: "weak password",
			req: &register.Request{
				Email:    "test@example.com",
				Password: "123",
				Role:     "user",
			},
			wantErr: user.ErrPassTooShort,
		},
		{
			name: "hash error",
			req: &register.Request{
				Email:    "test@example.com",
				Password: "abc12345",
				Role:     "user",
			},
			setup: func(hasher *mocks.MockpassHasher) {
				hasher.EXPECT().
					Hash("abc12345").
					Return("", errors.New("hash error"))
			},
			wantErr: errors.New("hash error"),
		},
		{
			name: "invalid role",
			req: &register.Request{
				Email:    "test@example.com",
				Password: "abc12345",
				Role:     "weird",
			},
			setup: func(hasher *mocks.MockpassHasher) {
				hasher.EXPECT().
					Hash("abc12345").
					Return("hash", nil)
			},
			wantErr: user.ErrInvalidRole,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockuserRepo(ctrl)
			hasher := mocks.NewMockpassHasher(ctrl)

			if tt.setup != nil {
				tt.setup(hasher)
			}

			uc := register.NewUsecase(repo, hasher)

			resp, err := uc.Execute(context.Background(), tt.req)

			require.Error(t, err)

			if !errors.Is(err, tt.wantErr) {
				require.Contains(t, err.Error(), tt.wantErr.Error())
			}

			require.Nil(t, resp)
		})
	}
}
