package postgres

import "time"

// pg tags documented here: https://pg.uptrace.dev/models/

type (
	Image struct {
		Path     string    `pg:"path,pk" json:"path"`
		Uploaded time.Time `pg:"uploaded,notnull,default:now()" json:"-"`
	}

	User struct {
		Email      string    `pg:"email,pk,unique" json:"email"`
		ID         string    `pg:"id,nopk,notnull,unique,default:nextval('user_id_seq')::text" json:"id"`
		Username   string    `pg:"username" json:"username"`
		Bio        string    `pg:"bio" json:"bio"`
		Hashtags   []string  `pg:"hashtags,array,notnull,default:'{}'::text[]" json:"hashtags"`
		Image      string    `pg:"image,pk,default:'placeholder.png'" json:"image"`
		Birth      time.Time `pg:"birth,type:date" json:"birth"`
		Registered time.Time `pg:"registered,pk,default:now()" json:"registered"`
		Updated    time.Time `pg:"updated,notnull,default:now()" json:"-"`
		DeletedAt  time.Time `pg:"deleted_at,soft_delete" json:"-"`
	}
)
