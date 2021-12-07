package modelv1

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Identifier struct {
	CreatedAt      time.Time    `gorm:"column:created_at;autoCreateTime" json:"created_at"` //
	DomainID       string       `gorm:"column:domain_id" json:"domain_id"`                  //
	ID             string       `gorm:"column:id;primary_key" json:"id"`                    //
	Identifier     string       `gorm:"column:identifier" json:"identifier"`                //
	IdentifierType string       `gorm:"column:identifier_type" json:"identifier_type"`      //
	IdentitiyID    string       `gorm:"column:identitiy_id" json:"identitiy_id"`            //
	UpdatedAt      time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"` //
	Credentials    []Credential `gorm:"many2many:tb_identifier_credential" json:"credentials"`
}

// TableName sets the insert table name for this struct type.
func (i *Identifier) TableName() string {
	return "tb_identifier"
}

func (i *Identifier) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.NewString()

	return nil
}
