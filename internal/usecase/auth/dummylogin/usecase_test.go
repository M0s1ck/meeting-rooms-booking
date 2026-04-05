package dummylogin

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/dummylogin/mocks"
)

func TestUsecase_Execute_User(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockuserRepo(ctrl)
	tokenGen := mocks.NewMocktokenGenerator(ctrl)

	uc := NewUsecase(tokenGen, repo)

	req := &Request{
		Role: user.RoleUser,
	}

	repo.EXPECT().
		Ensure(gomock.Any(), userID, user.RoleUser).
		Return(nil)

	tokenGen.EXPECT().
		Generate(userID, user.RoleUser).
		Return("token-user", nil)

	resp, err := uc.Execute(t.Context(), req)

	require.NoError(t, err)
	require.Equal(t, "token-user", resp.Token)
}

func TestUsecase_Execute_Admin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockuserRepo(ctrl)
	tokenGen := mocks.NewMocktokenGenerator(ctrl)

	uc := NewUsecase(tokenGen, repo)

	req := &Request{
		Role: user.RoleAdmin,
	}

	repo.EXPECT().
		Ensure(gomock.Any(), adminID, user.RoleAdmin).
		Return(nil)

	tokenGen.EXPECT().
		Generate(adminID, user.RoleAdmin).
		Return("token-admin", nil)

	resp, err := uc.Execute(t.Context(), req)

	require.NoError(t, err)
	require.Equal(t, "token-admin", resp.Token)
}

func TestUsecase_Execute_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockuserRepo(ctrl)
	tokenGen := mocks.NewMocktokenGenerator(ctrl)

	uc := NewUsecase(tokenGen, repo)

	req := &Request{
		Role: user.RoleUser,
	}

	repo.EXPECT().
		Ensure(gomock.Any(), userID, user.RoleUser).
		Return(errors.New("db error"))

	resp, err := uc.Execute(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
}

func TestUsecase_Execute_TokenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockuserRepo(ctrl)
	tokenGen := mocks.NewMocktokenGenerator(ctrl)

	uc := NewUsecase(tokenGen, repo)

	req := &Request{
		Role: user.RoleUser,
	}

	repo.EXPECT().
		Ensure(gomock.Any(), userID, user.RoleUser).
		Return(nil)

	tokenGen.EXPECT().
		Generate(userID, user.RoleUser).
		Return("", errors.New("token error"))

	resp, err := uc.Execute(t.Context(), req)

	require.Error(t, err)
	require.Nil(t, resp)
}
