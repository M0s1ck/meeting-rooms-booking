package login_test

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/login"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/login/mocks"
)

func TestUsecase_Execute(t *testing.T) {
	validEmail := "test@example.com"
	normalizedEmail, _ := user.NewEmail(validEmail)

	userID := uuid.New()

	usr := &user.User{
		ID:       userID,
		Email:    normalizedEmail,
		PassHash: "hash",
		Role:     user.RoleUser,
	}

	errDB := errors.New("db error")
	errInternal := errors.New("internal")

	tests := []struct {
		name    string
		req     *login.Request
		setup   func(repo *mocks.MockuserRepo, token *mocks.MocktokenGenerator, pass *mocks.MockpassComparor)
		wantErr error
		wantRes *login.Response
	}{
		{
			name: "success",
			req: &login.Request{
				Email:    validEmail,
				Password: "password",
			},
			setup: func(repo *mocks.MockuserRepo, token *mocks.MocktokenGenerator, pass *mocks.MockpassComparor) {
				repo.EXPECT().
					GetByEmail(gomock.Any(), normalizedEmail).
					Return(usr, nil)

				pass.EXPECT().
					Compare("hash", "password").
					Return(nil)

				token.EXPECT().
					Generate(userID, user.RoleUser).
					Return("token", nil)
			},
			wantRes: &login.Response{Token: "token"},
		},
		{
			name: "invalid email format",
			req: &login.Request{
				Email:    "bad-email",
				Password: "password",
			},
			wantErr: user.ErrInvalidCredentials,
		},
		{
			name: "user not found",
			req: &login.Request{
				Email:    validEmail,
				Password: "password",
			},
			setup: func(repo *mocks.MockuserRepo, _ *mocks.MocktokenGenerator, _ *mocks.MockpassComparor) {
				repo.EXPECT().
					GetByEmail(gomock.Any(), normalizedEmail).
					Return(nil, user.ErrNotFound)
			},
			wantErr: user.ErrInvalidCredentials,
		},
		{
			name: "repo error",
			req: &login.Request{
				Email:    validEmail,
				Password: "password",
			},
			setup: func(repo *mocks.MockuserRepo, _ *mocks.MocktokenGenerator, _ *mocks.MockpassComparor) {
				repo.EXPECT().
					GetByEmail(gomock.Any(), normalizedEmail).
					Return(nil, errDB)
			},
			wantErr: errDB,
		},
		{
			name: "wrong password",
			req: &login.Request{
				Email:    validEmail,
				Password: "wrong",
			},
			setup: func(repo *mocks.MockuserRepo, _ *mocks.MocktokenGenerator, pass *mocks.MockpassComparor) {
				repo.EXPECT().
					GetByEmail(gomock.Any(), normalizedEmail).
					Return(usr, nil)

				pass.EXPECT().
					Compare("hash", "wrong").
					Return(user.ErrWrongPassword)
			},
			wantErr: user.ErrInvalidCredentials,
		},
		{
			name: "password compare internal error",
			req: &login.Request{
				Email:    validEmail,
				Password: "password",
			},
			setup: func(repo *mocks.MockuserRepo, _ *mocks.MocktokenGenerator, pass *mocks.MockpassComparor) {
				repo.EXPECT().
					GetByEmail(gomock.Any(), normalizedEmail).
					Return(usr, nil)

				pass.EXPECT().
					Compare("hash", "password").
					Return(errInternal)
			},
			wantErr: errInternal,
		},
		{
			name: "token error",
			req: &login.Request{
				Email:    validEmail,
				Password: "password",
			},
			setup: func(repo *mocks.MockuserRepo, token *mocks.MocktokenGenerator, pass *mocks.MockpassComparor) {
				repo.EXPECT().
					GetByEmail(gomock.Any(), normalizedEmail).
					Return(usr, nil)

				pass.EXPECT().
					Compare("hash", "password").
					Return(nil)

				token.EXPECT().
					Generate(userID, user.RoleUser).
					Return("", errInternal)
			},
			wantErr: errInternal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mocks.NewMockuserRepo(ctrl)
			token := mocks.NewMocktokenGenerator(ctrl)
			pass := mocks.NewMockpassComparor(ctrl)

			if tt.setup != nil {
				tt.setup(repo, token, pass)
			}

			uc := login.NewUsecase(repo, token, pass)

			resp, err := uc.Execute(context.Background(), tt.req)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
				require.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.wantRes, resp)
		})
	}
}
