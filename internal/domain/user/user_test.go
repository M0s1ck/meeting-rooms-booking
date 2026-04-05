package user_test

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/internships-backend/test-backend-M0s1ck/internal/domain/user"
)

func TestUser_New_Success(t *testing.T) {
	id := uuid.New()
	email, err := user.NewEmail("test@example.com")
	require.NoError(t, err)

	now := time.Now()

	u, err := user.New(id, email, "hash123", user.RoleUser, now)

	require.NoError(t, err)
	require.Equal(t, id, u.ID)
	require.Equal(t, email, u.Email)
	require.Equal(t, "hash123", u.PassHash)
	require.Equal(t, user.RoleUser, u.Role)
	require.WithinDuration(t, now.UTC(), u.CreatedAt, time.Second)
}

func TestUser_New_EmptyPassHash(t *testing.T) {
	id := uuid.New()
	email, _ := user.NewEmail("test@example.com")

	u, err := user.New(id, email, "", user.RoleUser, time.Now())

	require.Error(t, err)
	require.ErrorIs(t, err, user.ErrEmptyPassHash)
	require.Nil(t, u)
}

func TestParseRole_Success(t *testing.T) {
	r, err := user.ParseRole("admin")
	require.NoError(t, err)
	require.Equal(t, user.RoleAdmin, r)

	r, err = user.ParseRole("user")
	require.NoError(t, err)
	require.Equal(t, user.RoleUser, r)
}

func TestParseRole_Invalid(t *testing.T) {
	r, err := user.ParseRole("superuser")

	require.Error(t, err)
	require.ErrorIs(t, err, user.ErrInvalidRole)
	require.Equal(t, user.Role(""), r)
}

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    user.Email
		wantErr error
	}{
		{
			name:  "valid email",
			input: "test@example.com",
			want:  user.Email("test@example.com"),
		},
		{
			name:  "normalize",
			input: "  TEST@Example.COM  ",
			want:  user.Email("test@example.com"),
		},
		{
			name:    "blank",
			input:   "   ",
			wantErr: user.ErrBlankEmail,
		},
		{
			name:    "no @",
			input:   "testexample.com",
			wantErr: user.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := user.NewEmail(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
				require.Equal(t, user.Email(""), got)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestValidatePassStrength(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:  "valid",
			input: "abc12345",
		},
		{
			name:    "too short",
			input:   "a1b2",
			wantErr: user.ErrPassTooShort,
		},
		{
			name:    "too long",
			input:   strings.Repeat("a1", 30), // > 50
			wantErr: user.ErrPassTooLong,
		},
		{
			name:    "no letter",
			input:   "12345678",
			wantErr: user.ErrPassNoLetter,
		},
		{
			name:    "no digit",
			input:   "abcdefgh",
			wantErr: user.ErrPassNoDigit,
		},
		{
			name:  "unicode letters",
			input: "пароль123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := user.ValidatePassStrength(tt.input)

			if tt.wantErr != nil {
				require.Error(t, err)
				require.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
		})
	}
}
