package repository

import (
	"context"
	"shorten-url/internal/model"

	"github.com/stretchr/testify/mock"
)

type (
	MockURLRepository struct {
		mock.Mock
	}
)

func (mr *MockURLRepository) Create(ctx context.Context, url *model.URL) (*model.URLInterpeter, error) {

	args := mr.Called(ctx, url)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.URLInterpeter), args.Error(1)
}
func (mr *MockURLRepository) GetByShortCode(ctx context.Context, shortCode string) (*model.URL, error) {

	args := mr.Called(ctx, shortCode)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.URL), args.Error(1)
}
func (mr *MockURLRepository) UpdateShortUrl(pctx context.Context, shortCode string, updatedUrl string) (*model.URL, error) {

	args := mr.Called(pctx, shortCode, updatedUrl)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*model.URL), args.Error(1)
}
func (mr *MockURLRepository) DeleteByShortCode(pctx context.Context, shortCode string) error {

	args := mr.Called(pctx, shortCode)

	return args.Error(0)
}
func (mr *MockURLRepository) UpdateShortUrlCount(pctx context.Context, shortCode string) error {

	args := mr.Called(pctx, shortCode)

	return args.Error(0)
}
func (mr *MockURLRepository) IsShortCodeExists(pctx context.Context, shortCode string) bool {

	args := mr.Called(pctx, shortCode)

	return args.Bool(0)
}
