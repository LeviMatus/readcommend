package drivertest

import (
	"context"

	"github.com/LeviMatus/readcommend/service/internal/driver"
	"github.com/LeviMatus/readcommend/service/internal/driver/book"
	"github.com/LeviMatus/readcommend/service/internal/entity"
	"github.com/stretchr/testify/mock"
)

type DriverMock struct {
	mock.Mock
}

var _ driver.Driver = (*DriverMock)(nil)

func (d *DriverMock) ListAuthors(ctx context.Context) ([]entity.Author, error) {
	args := d.Called(ctx)
	return args.Get(0).([]entity.Author), args.Error(1)
}

func (d *DriverMock) ListGenres(ctx context.Context) ([]entity.Genre, error) {
	args := d.Called(ctx)
	return args.Get(0).([]entity.Genre), args.Error(1)
}

func (d *DriverMock) ListSizes(ctx context.Context) ([]entity.Size, error) {
	args := d.Called(ctx)
	return args.Get(0).([]entity.Size), args.Error(1)
}

func (d *DriverMock) ListEras(ctx context.Context) ([]entity.Era, error) {
	args := d.Called(ctx)
	return args.Get(0).([]entity.Era), args.Error(1)
}

func (d *DriverMock) SearchBooks(ctx context.Context, params book.SearchInput) ([]entity.Book, error) {
	args := d.Called(ctx, params)
	return args.Get(0).([]entity.Book), args.Error(1)
}
