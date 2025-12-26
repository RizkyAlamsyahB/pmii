package domain

import (
	"database/sql"
	"time"
)

// DocumentType represents the type of document
type DocumentType string

const (
	DocumentTypeProdukHukum    DocumentType = "produk_hukum"
	DocumentTypeLaguOrganisasi DocumentType = "lagu_organisasi"
	DocumentTypeLogoOrganisasi DocumentType = "logo_organisasi"
)

// ValidDocumentTypes returns all valid document types
func ValidDocumentTypes() []DocumentType {
	return []DocumentType{
		DocumentTypeProdukHukum,
		DocumentTypeLaguOrganisasi,
		DocumentTypeLogoOrganisasi,
	}
}

// IsValid checks if the document type is valid
func (dt DocumentType) IsValid() bool {
	for _, valid := range ValidDocumentTypes() {
		if dt == valid {
			return true
		}
	}
	return false
}

// GetLabel returns human-readable label for document type
func (dt DocumentType) GetLabel() string {
	switch dt {
	case DocumentTypeProdukHukum:
		return "Produk Hukum"
	case DocumentTypeLaguOrganisasi:
		return "Lagu Organisasi"
	case DocumentTypeLogoOrganisasi:
		return "Logo Organisasi"
	default:
		return string(dt)
	}
}

// GetCloudinaryFolder returns the Cloudinary folder path for this document type
func (dt DocumentType) GetCloudinaryFolder() string {
	switch dt {
	case DocumentTypeProdukHukum:
		return "documents/produk_hukum"
	case DocumentTypeLaguOrganisasi:
		return "documents/lagu_organisasi"
	case DocumentTypeLogoOrganisasi:
		return "documents/logo_organisasi"
	default:
		return "documents"
	}
}

// Document represents a downloadable file/document
type Document struct {
	ID        int          `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string       `gorm:"type:varchar(255);not null" json:"name"`
	FileType  DocumentType `gorm:"type:document_type;not null" json:"file_type"`
	FileURI   string       `gorm:"type:varchar(255);not null" json:"file_uri"`
	CreatedAt time.Time    `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time    `gorm:"default:now()" json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName specifies the table name for Document
func (Document) TableName() string {
	return "documents"
}
