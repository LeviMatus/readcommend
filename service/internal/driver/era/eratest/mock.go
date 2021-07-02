package eratest

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/mock"
)

// DriverMock is used to mock the Driver interface.
type DriverMock struct {
	mock.Mock
}

// ListEras is a mock routine that returns items as instructed.
func (d *DriverMock) ListEras(ctx context.Context) ([]entity.Era, error) {
	args := d.Called(ctx)
	return args.Get(0).([]entity.Era), args.Error(1)
}
