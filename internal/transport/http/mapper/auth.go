package mapper

import (
	"errors"

	"github.com/oapi-codegen/runtime/types"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	dummylogin "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/dummylogin"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/login"
	"github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/register"
)

var (
	ErrNilDummyLoginBody = errors.New("request body is required")
	ErrNilRegisterBody   = errors.New("request body is required")
	ErrNilLoginBody      = errors.New("request body is required")
)

func ToDummyLoginRequest(body *oapi.PostDummyLoginJSONRequestBody) (*dummylogin.Request, error) {
	if body == nil {
		return nil, ErrNilDummyLoginBody
	}

	role, err := user.ParseRole(string(body.Role))
	if err != nil {
		return nil, err
	}

	return &dummylogin.Request{
		Role: role,
	}, nil
}

func ToDummyLoginResponse(resp *dummylogin.Response) oapi.PostDummyLogin200JSONResponse {
	if resp == nil {
		return oapi.PostDummyLogin200JSONResponse{}
	}

	return oapi.PostDummyLogin200JSONResponse{
		Token: resp.Token,
	}
}

func ToRegisterRequest(body *oapi.PostRegisterJSONRequestBody) (*register.Request, error) {
	if body == nil {
		return nil, ErrNilRegisterBody
	}

	return &register.Request{
		Role:     string(body.Role),
		Email:    string(body.Email),
		Password: body.Password,
	}, nil
}

func ToRegisterResponse(resp *register.Response) oapi.PostRegister201JSONResponse {
	return oapi.PostRegister201JSONResponse{
		User: &oapi.User{
			Id:        resp.ID,
			Email:     types.Email(resp.Email),
			Role:      oapi.UserRole(resp.Role),
			CreatedAt: &resp.CreatedAt,
		},
	}
}

func ToLoginRequest(body *oapi.PostLoginJSONRequestBody) (*login.Request, error) {
	if body == nil {
		return nil, ErrNilLoginBody
	}

	return &login.Request{
		Email:    string(body.Email),
		Password: body.Password,
	}, nil
}

func ToLoginResponse(resp *login.Response) oapi.PostLogin200JSONResponse {
	return oapi.PostLogin200JSONResponse{
		Token: resp.Token,
	}
}
