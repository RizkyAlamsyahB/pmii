package service

import (
"context"
"errors"
"testing"
"time"

"github.com/garuda-labs-1/pmii-be/internal/domain"
"github.com/garuda-labs-1/pmii-be/internal/dto/requests"
)

// MockAboutRepository adalah mock untuk AboutRepository
type MockAboutRepository struct {
GetFunc    func() (*domain.About, error)
UpsertFunc func(about *domain.About) error
}

func (m *MockAboutRepository) Get() (*domain.About, error) {
if m.GetFunc != nil {
return m.GetFunc()
}
return nil, errors.New("mock not configured")
}

func (m *MockAboutRepository) Upsert(about *domain.About) error {
if m.UpsertFunc != nil {
return m.UpsertFunc(about)
}
return errors.New("mock not configured")
}

// ==================== GET TESTS ====================

// Test: Get berhasil dengan data existing
func TestAboutGet_Success(t *testing.T) {
title := "Tentang PMII"
history := "Sejarah PMII"
vision := "Visi PMII"
mission := "Misi PMII"
videoURL := "https://youtube.com/watch?v=123"

mockRepo := &MockAboutRepository{
GetFunc: func() (*domain.About, error) {
return &domain.About{
ID:        1,
Title:     &title,
History:   &history,
Vision:    &vision,
Mission:   &mission,
VideoURL:  &videoURL,
UpdatedAt: time.Now(),
}, nil
},
}

service := NewAboutService(mockRepo)
result, err := service.Get(context.Background())

if err != nil {
t.Errorf("Expected no error, got: %v", err)
}
if result == nil {
t.Fatal("Expected result, got nil")
}
if result.ID != 1 {
t.Errorf("Expected ID 1, got: %d", result.ID)
}
if *result.Title != title {
t.Errorf("Expected Title '%s', got: '%s'", title, *result.Title)
}
if *result.History != history {
t.Errorf("Expected History '%s', got: '%s'", history, *result.History)
}
if *result.Vision != vision {
t.Errorf("Expected Vision '%s', got: '%s'", vision, *result.Vision)
}
if *result.Mission != mission {
t.Errorf("Expected Mission '%s', got: '%s'", mission, *result.Mission)
}
}

// Test: Get ketika belum ada data (return empty response)
func TestAboutGet_NoData(t *testing.T) {
mockRepo := &MockAboutRepository{
GetFunc: func() (*domain.About, error) {
return nil, errors.New("record not found")
},
}

service := NewAboutService(mockRepo)
result, err := service.Get(context.Background())

if err != nil {
t.Errorf("Expected no error for empty data, got: %v", err)
}
if result == nil {
t.Fatal("Expected empty result, got nil")
}
if result.ID != 0 {
t.Errorf("Expected ID 0 for empty, got: %d", result.ID)
}
}

// ==================== UPDATE TESTS ====================

// Test: Update berhasil dengan Title
func TestAboutUpdate_Success(t *testing.T) {
history := "Sejarah Lama"
mockRepo := &MockAboutRepository{
GetFunc: func() (*domain.About, error) {
return &domain.About{
ID:      1,
History: &history,
}, nil
},
UpsertFunc: func(about *domain.About) error {
return nil
},
}

service := NewAboutService(mockRepo)
req := requests.UpdateAboutRequest{
Title:   "Tentang PMII",
History: "Sejarah Baru",
Vision:  "Visi Baru",
}

result, err := service.Update(context.Background(), req)

if err != nil {
t.Errorf("Expected no error, got: %v", err)
}
if result == nil {
t.Fatal("Expected result, got nil")
}
if *result.Title != "Tentang PMII" {
t.Errorf("Expected Title 'Tentang PMII', got: '%s'", *result.Title)
}
if *result.History != "Sejarah Baru" {
t.Errorf("Expected History 'Sejarah Baru', got: '%s'", *result.History)
}
if *result.Vision != "Visi Baru" {
t.Errorf("Expected Vision 'Visi Baru', got: '%s'", *result.Vision)
}
}

// Test: Update ketika belum ada data (create baru)
func TestAboutUpdate_CreateNew(t *testing.T) {
upsertCalled := false

mockRepo := &MockAboutRepository{
GetFunc: func() (*domain.About, error) {
return nil, errors.New("record not found")
},
UpsertFunc: func(about *domain.About) error {
upsertCalled = true
if about.History == nil || *about.History != "Sejarah Baru" {
t.Error("Expected History to be set")
}
return nil
},
}

service := NewAboutService(mockRepo)
req := requests.UpdateAboutRequest{
History: "Sejarah Baru",
}

_, err := service.Update(context.Background(), req)

if err != nil {
t.Errorf("Expected no error, got: %v", err)
}
if !upsertCalled {
t.Error("Expected Upsert to be called")
}
}

// Test: Database error harus return error
func TestAboutUpdate_ErrorDatabase(t *testing.T) {
mockRepo := &MockAboutRepository{
GetFunc: func() (*domain.About, error) {
return &domain.About{ID: 1}, nil
},
UpsertFunc: func(about *domain.About) error {
return errors.New("database error")
},
}

service := NewAboutService(mockRepo)
req := requests.UpdateAboutRequest{History: "Test"}

_, err := service.Update(context.Background(), req)

if err == nil || err.Error() != "gagal menyimpan about" {
t.Errorf("Expected 'gagal menyimpan about' error, got: %v", err)
}
}

// Test: Update dengan semua fields
func TestAboutUpdate_AllFields(t *testing.T) {
mockRepo := &MockAboutRepository{
GetFunc: func() (*domain.About, error) {
return &domain.About{ID: 1}, nil
},
UpsertFunc: func(about *domain.About) error {
// Verify semua field di-update
if about.Title == nil || *about.Title != "Title" {
t.Error("Expected Title to be updated")
}
if about.History == nil || *about.History != "History" {
t.Error("Expected History to be updated")
}
if about.Vision == nil || *about.Vision != "Vision" {
t.Error("Expected Vision to be updated")
}
if about.Mission == nil || *about.Mission != "Mission" {
t.Error("Expected Mission to be updated")
}
if about.VideoURL == nil || *about.VideoURL != "https://youtube.com/test" {
t.Error("Expected VideoURL to be updated")
}
return nil
},
}

service := NewAboutService(mockRepo)
req := requests.UpdateAboutRequest{
Title:    "Title",
History:  "History",
Vision:   "Vision",
Mission:  "Mission",
VideoURL: "https://youtube.com/test",
}

result, err := service.Update(context.Background(), req)

if err != nil {
t.Errorf("Expected no error, got: %v", err)
}
if result == nil {
t.Fatal("Expected result, got nil")
}
}

// Test: Update partial fields (hanya update yang dikirim)
func TestAboutUpdate_PartialFields(t *testing.T) {
existingHistory := "Existing History"
existingVision := "Existing Vision"

mockRepo := &MockAboutRepository{
GetFunc: func() (*domain.About, error) {
return &domain.About{
ID:      1,
History: &existingHistory,
Vision:  &existingVision,
}, nil
},
UpsertFunc: func(about *domain.About) error {
// Vision harus update karena dikirim
if *about.Vision != "New Vision" {
t.Errorf("Expected Vision to be updated to 'New Vision', got: %s", *about.Vision)
}
return nil
},
}

service := NewAboutService(mockRepo)
req := requests.UpdateAboutRequest{
Vision: "New Vision",
// History tidak dikirim (empty)
}

_, err := service.Update(context.Background(), req)

if err != nil {
t.Errorf("Expected no error, got: %v", err)
}
}
