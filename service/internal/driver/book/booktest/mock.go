package booktest

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/mock"
)

// DriverMock is used to mock the Driver interface.
type DriverMock struct {
	mock.Mock
}

// SearchBooks is a mock routine that returns items as instructed.
func (d *DriverMock) SearchBooks(ctx context.Context, params book.SearchInput) ([]entity.Book, error) {
	args := d.Called(ctx, params)
	return args.Get(0).([]entity.Book), args.Error(1)
}
