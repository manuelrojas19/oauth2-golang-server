package store

import "time"

type OauthResource struct {
	ResourceId  string    `gorm:"primaryKey;type:varchar(255);unique;not null"`
	Name        string    `gorm:"type:varchar(255);unique;not null"`
	Description string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	Scopes []Scope `gorm:"many2many:oauth_resource_scopes;foreignKey:ResourceId;joinForeignKey:ResourceId;References:Id;JoinReferences:ScopeId"`
}
