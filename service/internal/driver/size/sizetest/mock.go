package sizetest

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/mock"
)

// DriverMock is used to mock the Driver interface.
type DriverMock struct {
	mock.Mock
}

// ListSizes is a mock routine that returns items as instructed.
func (d *DriverMock) ListSizes(ctx context.Context) ([]entity.Size, error) {
	args := d.Called(ctx)
	return args.Get(0).([]entity.Size), args.Error(1)
}
