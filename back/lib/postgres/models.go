package postgres

import "time"

// pg tags documented here: https://pg.uptrace.dev/models/

type (
	Image struct {
		Path     string    `pg:"path,pk" json:"path"`
		Uploaded time.Time `pg:"uploaded,notnull,default:\"now()\"" json:"-"`
	}

	User struct {
		Email      string    `pg:"email,pk,unique" json:"email"`
		ID         string    `pg:"id,nopk,notnull,unique,default:\"'id'||nextval('user_id_seq')::text\"" json:"id"`
		Username   string    `pg:"username" json:"username"`
		Bio        string    `pg:"bio" json:"bio"`
		Hashtags   []string  `pg:"hashtags,array,notnull,default:\"'{}'::text[]\"" json:"hashtags"`
		Image      string    `pg:"image,notnull,default:'placeholder.png'" json:"image"`
		Birth      time.Time `pg:"birth,type:date" json:"birth"`
		Registered time.Time `pg:"registered,pk,default:now()" json:"registered"`
		Updated    time.Time `pg:"updated,notnull,default:now()" json:"-"`
		DeletedAt  time.Time `pg:"deleted_at,soft_delete" json:"-"`
	}

	Hashtag struct {
		Name    string `pg:"name,pk,unique" json:"-"`
		Counter uint64 `pg:"counter,notnull,default:\"'0'::bigint\"" json:"-"`
	}

	Location struct {
		Email     string    `pg:"email,pk,unique" json:"email"`
		Latitude  float64   `pg:"latitude,notnull,default:\"'0'::double precision\"" json:"latitude"`
		Longitude float64   `pg:"longitude,notnull,default:\"'0'::double precision\"" json:"longitude"`
		Updated   time.Time `pg:"updated,notnull,default:now()" json:"-"`
	}
)
