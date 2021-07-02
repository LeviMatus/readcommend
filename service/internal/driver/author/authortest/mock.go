package authortest

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/mock"
)

// DriverMock is used to mock the Driver interface.
type DriverMock struct {
	mock.Mock
}

// ListAuthors is a mock routine that returns items as instructed.
func (d *DriverMock) ListAuthors(ctx context.Context) ([]entity.Author, error) {
	args := d.Called(ctx)
	return args.Get(0).([]entity.Author), args.Error(1)
}
