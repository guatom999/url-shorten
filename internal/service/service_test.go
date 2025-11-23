package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"shorten-url/internal/entities"
	appErrors "shorten-url/internal/errors"
	"shorten-url/internal/model"
	"shorten-url/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenURL_Success(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	originalUrl := "http://example.com"
	now := time.Now()

	mockRepo.On("Create", ctx, mock.MatchedBy(func(url *model.URL) bool {
		return url.OriginalURL == originalUrl && len(url.ShortCode) == 6
	})).Return(&model.URLInterpeter{
		ID:        1,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	result, err := service.ShortenURL(ctx, originalUrl)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.Id)
	assert.Equal(t, originalUrl, result.OriginalURL)
	assert.NotEmpty(t, result.ShortUrl)
	assert.Equal(t, now, result.CreatedAt)
	assert.Equal(t, now, result.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestShortenURL_Error(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	originalUrl := "http://example.com"

	mockRepo.On("Create", ctx, mock.AnythingOfType("*model.URL")).
		Return(nil, errors.New("database error"))

	result, err := service.ShortenURL(ctx, originalUrl)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &appErrors.AppError{}, err)

	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_Success(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "abc123"
	expectedURL := &model.URL{
		ID:          1,
		ShortCode:   shortCode,
		OriginalURL: "http://example.com",
		ClickCount:  5,
	}

	mockRepo.On("GetByShortCode", ctx, shortCode).
		Return(expectedURL, nil)
	mockRepo.On("UpdateShortUrlCount", ctx, shortCode).
		Return(nil)

	result, err := service.GetOriginalURL(ctx, shortCode)

	assert.NoError(t, err)
	assert.Equal(t, expectedURL.OriginalURL, result)

	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_NotFound(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "notfound"

	mockRepo.On("GetByShortCode", ctx, shortCode).
		Return(nil, errors.New("not found"))

	result, err := service.GetOriginalURL(ctx, shortCode)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.IsType(t, &appErrors.AppError{}, err)

	appErr := err.(*appErrors.AppError)
	assert.Equal(t, appErrors.NotFound, appErr.Type)

	mockRepo.AssertExpectations(t)
}

func TestGetOriginalURL_UpdateCountError(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "abc123"
	expectedURL := &model.URL{
		ID:          1,
		ShortCode:   shortCode,
		OriginalURL: "http://example.com",
	}

	mockRepo.On("GetByShortCode", ctx, shortCode).
		Return(expectedURL, nil)
	mockRepo.On("UpdateShortUrlCount", ctx, shortCode).
		Return(errors.New("update error"))

	result, err := service.GetOriginalURL(ctx, shortCode)

	assert.Error(t, err)
	assert.Empty(t, result)
	assert.IsType(t, &appErrors.AppError{}, err)

	mockRepo.AssertExpectations(t)
}

func TestRetrieveOriginalURL_Success(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "abc123"
	now := time.Now()
	expectedURL := &model.URL{
		ID:          1,
		ShortCode:   shortCode,
		OriginalURL: "http://example.com",
		ClickCount:  10,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("GetByShortCode", ctx, shortCode).
		Return(expectedURL, nil)

	result, err := service.RetrieveOriginalURL(ctx, shortCode)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedURL.ID, result.Id)
	assert.Equal(t, expectedURL.OriginalURL, result.OriginalUrl)
	assert.Equal(t, expectedURL.ShortCode, result.ShortUrl)
	assert.Equal(t, expectedURL.CreatedAt, result.CreatedAt)
	assert.Equal(t, expectedURL.UpdatedAt, result.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestRetrieveOriginalURL_NotFound(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "notfound"

	mockRepo.On("GetByShortCode", ctx, shortCode).
		Return(nil, errors.New("not found"))

	result, err := service.RetrieveOriginalURL(ctx, shortCode)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &appErrors.AppError{}, err)

	mockRepo.AssertExpectations(t)
}

func TestUpdateShortUrl_Success(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "abc123"
	updatedUrl := "http://newexample.com"
	now := time.Now()
	expectedURL := &model.URL{
		ID:          1,
		ShortCode:   shortCode,
		OriginalURL: updatedUrl,
		UpdatedAt:   now,
	}

	mockRepo.On("UpdateShortUrl", ctx, shortCode, updatedUrl).
		Return(expectedURL, nil)

	result, err := service.UpdateShortUrl(ctx, shortCode, updatedUrl)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedURL.OriginalURL, result.OriginalURL)
	assert.Equal(t, expectedURL.ShortCode, result.ShortCode)

	mockRepo.AssertExpectations(t)
}

func TestUpdateShortUrl_NotFound(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "notfound"
	updatedUrl := "http://newexample.com"

	mockRepo.On("UpdateShortUrl", ctx, shortCode, updatedUrl).
		Return(nil, errors.New("not found"))

	result, err := service.UpdateShortUrl(ctx, shortCode, updatedUrl)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &appErrors.AppError{}, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteShortUrl_Success(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "abc123"

	mockRepo.On("IsShortCodeExists", ctx, shortCode).
		Return(true)
	mockRepo.On("DeleteByShortCode", ctx, shortCode).
		Return(nil)

	err := service.DeleteShortUrl(ctx, shortCode)

	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteShortUrl_NotFound(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "notfound"

	mockRepo.On("IsShortCodeExists", ctx, shortCode).
		Return(false)

	err := service.DeleteShortUrl(ctx, shortCode)

	assert.Error(t, err)
	assert.IsType(t, &appErrors.AppError{}, err)

	appErr := err.(*appErrors.AppError)
	assert.Equal(t, appErrors.NotFound, appErr.Type)

	mockRepo.AssertExpectations(t)
}

func TestDeleteShortUrl_DeleteError(t *testing.T) {
	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "abc123"

	mockRepo.On("IsShortCodeExists", ctx, shortCode).
		Return(true)
	mockRepo.On("DeleteByShortCode", ctx, shortCode).
		Return(errors.New("delete failed"))

	err := service.DeleteShortUrl(ctx, shortCode)

	assert.Error(t, err)
	assert.IsType(t, &appErrors.AppError{}, err)

	mockRepo.AssertExpectations(t)
}

func TestGetUrlStatic_Success(t *testing.T) {
	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "abc123"
	now := time.Now()
	expectedURL := &model.URL{
		ID:          1,
		ShortCode:   shortCode,
		OriginalURL: "http://example.com",
		ClickCount:  15,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	mockRepo.On("GetByShortCode", ctx, shortCode).
		Return(expectedURL, nil)

	result, err := service.GetUrlStatic(ctx, shortCode)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "1", result.Id)
	assert.Equal(t, expectedURL.OriginalURL, result.Url)
	assert.Equal(t, expectedURL.ShortCode, result.ShortCode)
	assert.Equal(t, expectedURL.ClickCount, result.AccessCount)
	assert.Equal(t, expectedURL.CreatedAt, result.CreatedAt)
	assert.Equal(t, expectedURL.UpdatedAt, result.UpdatedAt)

	mockRepo.AssertExpectations(t)
}

func TestGetUrlStatic_NotFound(t *testing.T) {

	mockRepo := new(repository.MockURLRepository)
	service := NewURLService(mockRepo)
	ctx := context.Background()

	shortCode := "notfound"

	mockRepo.On("GetByShortCode", ctx, shortCode).
		Return(nil, errors.New("not found"))

	result, err := service.GetUrlStatic(ctx, shortCode)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.IsType(t, &appErrors.AppError{}, err)

	appErr := err.(*appErrors.AppError)
	assert.Equal(t, appErrors.NotFound, appErr.Type)

	mockRepo.AssertExpectations(t)
}

func TestShortenURL_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		originalURL string
		setupMock   func(*repository.MockURLRepository)
		wantErr     bool
		checkResult func(*testing.T, *entities.CreateShortenUrlRes)
	}{
		{
			name:        "Success - Valid URL",
			originalURL: "http://example.com",
			setupMock: func(m *repository.MockURLRepository) {
				now := time.Now()
				m.On("Create", mock.Anything, mock.MatchedBy(func(url *model.URL) bool {
					return url.OriginalURL == "http://example.com" && len(url.ShortCode) == 6
				})).Return(&model.URLInterpeter{
					ID:        1,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
			},
			wantErr: false,
			checkResult: func(t *testing.T, result *entities.CreateShortenUrlRes) {
				assert.NotNil(t, result)
				assert.Equal(t, "1", result.Id)
				assert.Equal(t, "http://example.com", result.OriginalURL)
				assert.NotEmpty(t, result.ShortUrl)
			},
		},
		{
			name:        "Error - Repository failure",
			originalURL: "http://example.com",
			setupMock: func(m *repository.MockURLRepository) {
				m.On("Create", mock.Anything, mock.AnythingOfType("*model.URL")).
					Return(nil, errors.New("database error"))
			},
			wantErr: true,
			checkResult: func(t *testing.T, result *entities.CreateShortenUrlRes) {
				assert.Nil(t, result)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := new(repository.MockURLRepository)
			tt.setupMock(mockRepo)
			service := NewURLService(mockRepo)
			ctx := context.Background()

			result, err := service.ShortenURL(ctx, tt.originalURL)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteShortUrl_TableDriven(t *testing.T) {
	tests := []struct {
		name         string
		shortCode    string
		setupMock    func(*repository.MockURLRepository)
		wantErr      bool
		expectedType appErrors.ErrorType
	}{
		{
			name:      "Success - Delete existing URL",
			shortCode: "abc123",
			setupMock: func(m *repository.MockURLRepository) {
				m.On("IsShortCodeExists", mock.Anything, "abc123").Return(true)
				m.On("DeleteByShortCode", mock.Anything, "abc123").Return(nil)
			},
			wantErr: false,
		},
		{
			name:      "Error - URL not found",
			shortCode: "notfound",
			setupMock: func(m *repository.MockURLRepository) {
				m.On("IsShortCodeExists", mock.Anything, "notfound").Return(false)
			},
			wantErr:      true,
			expectedType: appErrors.NotFound,
		},
		{
			name:      "Error - Delete failed",
			shortCode: "abc123",
			setupMock: func(m *repository.MockURLRepository) {
				m.On("IsShortCodeExists", mock.Anything, "abc123").Return(true)
				m.On("DeleteByShortCode", mock.Anything, "abc123").
					Return(errors.New("delete failed"))
			},
			wantErr:      true,
			expectedType: appErrors.Internal,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mockRepo := new(repository.MockURLRepository)
			tt.setupMock(mockRepo)
			service := NewURLService(mockRepo)
			ctx := context.Background()

			err := service.DeleteShortUrl(ctx, tt.shortCode)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.expectedType != "" {
					appErr, ok := err.(*appErrors.AppError)
					assert.True(t, ok)
					assert.Equal(t, tt.expectedType, appErr.Type)
				}
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
