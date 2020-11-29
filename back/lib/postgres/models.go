package postgres

import "time"

// pg tags documented here: https://pg.uptrace.dev/models/

type (
	User struct {
		tableName  struct{}  `pg:"users"`
		Email      string    `pg:"email,pk" json:"-" form:"-"`
		ID         string    `gp:"id,nopk,notnull,unique" json:"id" form:"id"`
		Username   string    `gp:"username" json:"username" form:"username"`
		Bio        string    `gp:"bio" json:"bio" form:"bio"`
		Hashtags   []string  `gp:"hashtags,array" json:"hashtags" form:"hashtags"`
		Image      string    `gp:"image" json:"image" form:"image"`
		Birth      time.Time `gp:"birth,type:date" json:"birth" form:"birth"`
		Registered time.Time `gp:"registered,notnull" json:"-" form:"-"`
		DeletedAt  time.Time `pg:"deleted_at,soft_delete" json:"-" form:"-"`
	}
)
