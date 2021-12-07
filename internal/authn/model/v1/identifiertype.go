package modelv1

import "time"

type IdentifierType struct {
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"` //
	Descr        string    `gorm:"column:descr" json:"descr"`                          //
	ID           string    `gorm:"column:id;primary_key" json:"id"`                    //
	IsVerifiable int       `gorm:"column:is_verifiable" json:"is_verifiable"`          //
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"` //
}

// TableName sets the insert table name for this struct type.
func (i *IdentifierType) TableName() string {
	return "tb_identifier_type"
}
