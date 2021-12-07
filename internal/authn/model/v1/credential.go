package modelv1

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Credential struct {
	Config         string       `gorm:"column:config" json:"config"`                        //
	CreatedAt      time.Time    `gorm:"column:created_at;autoCreateTime" json:"created_at"` //
	CredentialType string       `gorm:"column:credential_type" json:"credential_type"`      //
	ID             string       `gorm:"column:id;primary_key" json:"id"`                    //
	IdentityID     string       `gorm:"column:identity_id" json:"identity_id"`              //
	UpdatedAt      time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"` //
	Identifiers    []Identifier `gorm:"many2many:tb_identifier_credential" json:"identifiers"`
}

// TableName sets the insert table name for this struct type.
func (c *Credential) TableName() string {
	return "tb_credential"
}

func (c *Credential) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.NewString()

	return nil
}
