package genretest

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/mock"
)

// DriverMock is used to mock the Driver interface.
type DriverMock struct {
	mock.Mock
}

// ListGenres is a mock routine that returns items as instructed.
func (d *DriverMock) ListGenres(ctx context.Context) ([]entity.Genre, error) {
	args := d.Called(ctx)
	return args.Get(0).([]entity.Genre), args.Error(1)
}
