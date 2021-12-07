package modelv1

import "time"

type IdentifierCredential struct {
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`    //
	CredentialID string    `gorm:"column:credential_id" json:"credential_id"`             //
	IdentifierID string    `gorm:"column:identifier_id;primary_key" json:"identifier_id"` //
}

// TableName sets the insert table name for this struct type.
func (i *IdentifierCredential) TableName() string {
	return "tb_identifier_credential"
}
