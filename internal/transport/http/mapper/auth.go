package mapper

import (
	"errors"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
	"github.com/internships-backend/test-backend-M0s1ck/internal/transport/http/oapi"
	dummylogin "github.com/internships-backend/test-backend-M0s1ck/internal/usecase/auth/dummylogin"
)

var ErrNilDummyLoginBody = errors.New("request body is required")

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
