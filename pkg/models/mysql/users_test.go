package mysql

import (
	"github.com/gitalek/gistpaste/pkg/models"
	"reflect"
	"testing"
	"time"
)

func TestUserModel_Insert(t *testing.T) {
	if testing.Short() {
		t.Skip("mysql: skipping integration test")
	}

	type args struct {
		name string
		userID int
	}

	type want struct {
		user *models.User
		error error
	}

	tests := []struct {
		name string
		args args
		want want

	} {
		{
			name: "ValidID",
			args: args{userID: 1},
			want: want{user: &models.User{
				ID: 1, Name: "Alice Jones", Email: "alice@example.com",
				Created: time.Date(2018, 12, 23, 17, 25, 22, 0, time.UTC),
				Active: true,
			}},
		},
		{name: "Zero ID", args: args{userID: 0}, want: want{user: nil, error: models.ErrNoRecord}},
		{name: "Non-existent ID", args: args{userID: 2}, want: want{user: nil, error: models.ErrNoRecord}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, teardown := newTestDB(t)
			defer teardown()

			m := UserModel{db}

			user, err := m.Get(tt.args.userID)
			if err != tt.want.error {
				t.Errorf("error: want: %v; got %s", tt.want.error, err)
			}

			if !reflect.DeepEqual(user, tt.want.user) {
				t.Errorf("user: want: %v; got %v", tt.want.user, user)
			}
		})
	}
}
