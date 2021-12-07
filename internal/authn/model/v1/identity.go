package modelv1

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Identity struct {
	CreatedAt     time.Time    `gorm:"column:created_at;autoCreateTime" json:"created_at"`           //
	DomainID      string       `gorm:"column:domain_id" json:"domain_id"`                            //
	ID            string       `gorm:"column:id;primary_key" json:"id"`                              //
	Stat          string       `gorm:"column:stat;default:actived" json:"stat"`                      //
	StatChangedAt time.Time    `gorm:"column:stat_changed_at;autoUpdateTime" json:"stat_changed_at"` //
	UpdatedAt     time.Time    `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`           //
	Credentials   []Credential `gorm:"foreignKey:IdentityID;AssociationForeignKey:ID" json:"credentials"`
}

// TableName sets the insert table name for this struct type.
func (i *Identity) TableName() string {
	return "tb_identity"
}

func (i *Identity) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.NewString()

	return nil
}
