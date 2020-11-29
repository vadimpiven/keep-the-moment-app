package postgres

import "time"

// pg tags documented here: https://pg.uptrace.dev/models/
type (
	User struct {
		tableName  struct{}  `pg:"users,discard_unknown_columns"`
		Email      string    `pg:"email,pk,type:text" json:"-" form:"-"`
		ID         string    `gp:"id,nopk,notnull,unique,type:text" json:"id" form:"id"`
		Username   string    `gp:"username,notnull,use_zero,type:text" json:"username" form:"username"`
		Bio        string    `gp:"bio,notnull,use_zero,type:text" json:"bio" form:"bio"`
		Hashtags   []string  `gp:"hashtags,notnull,array,type:jsonb" json:"hashtags" form:"hashtags"`
		Image      string    `gp:"image,type:text" json:"image" form:"image"`
		Birth      time.Time `gp:"birth,notnull,type:date" json:"birth" form:"birth"`
		Registered time.Time `gp:"registered,notnull,type:date" json:"-" form:"-"`
		DeletedAt  time.Time `pg:"deleted_at,soft_delete,type:timestamptz" json:"-" form:"-"`
	}
)
