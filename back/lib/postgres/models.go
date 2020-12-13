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
		Latitude  float64   `pg:"latitude,notnull,use_zero,default:\"'0'::double precision\"" json:"latitude"`
		Longitude float64   `pg:"longitude,notnull,use_zero,default:\"'0'::double precision\"" json:"longitude"`
		Updated   time.Time `pg:"updated,notnull,default:now()" json:"-"`
	}

	Post struct {
		ID           uint64    `pg:"id,pk,unique,type:bigserial,default:\"nextval('posts_id_seq')::bigserial\"" json:"id"`
		Email        string    `pg:"email,notnull" json:"-"`
		UserID       string    `pg:"-" json:"user_id"`
		UserImage    string    `pg:"-" json:"user_image"`
		Background   []int32   `pg:"background,array,notnull,default:\"'{}'::integer[]\"" json:"background"`
		Content      string    `pg:"content" json:"content"`
		Hashtags     []string  `pg:"hashtags,array,notnull,default:\"'{}'::text[]\"" json:"hashtags"`
		UserHashtags []string  `pg:"-" json:"user_hashtags"`
		Image1       string    `pg:"image_1" json:"image_1"`
		Image2       string    `pg:"image_2" json:"image_2"`
		Image3       string    `pg:"image_3" json:"image_3"`
		Image4       string    `pg:"image_4" json:"image_4"`
		Image5       string    `pg:"image_5" json:"image_5"`
		Latitude     float64   `pg:"latitude,use_zero,notnull" json:"latitude"`
		Longitude    float64   `pg:"longitude,use_zero,notnull" json:"longitude"`
		CreatedAt    time.Time `pg:"created_at,default:now()" json:"created"`
		HiddenAt     time.Time `pg:"hidden_at,soft_delete" json:"-"`
		Likes        uint64    `pg:"likes,notnull,default:\"'0'::bigint\""`
	}

	PostLike struct {
		PostId  uint64    `pg:"post_id" json:"-"`
		Email   string    `pg:"email" json:"-"`
		LikedAt time.Time `pg:"liked_at,notnull,default:now()" json:"-"`
	}

	PostComment struct {
		ID          uint64    `pg:"id,pk,unique,type:bigserial,default:\"nextval('post_comments_id_seq')::bigserial\"" json:"id"`
		PostId      uint64    `pg:"post_id" json:"-"`
		Email       string    `pg:"email" json:"-"`
		UserID      string    `pg:"-" json:"user_id"`
		UserImage   string    `pg:"-" json:"user_image"`
		Content     string    `pg:"content" json:"content"`
		CommentedAt time.Time `pg:"commented_at,default:now()" json:"commented_at"`
		DeletedAt   time.Time `pg:"deleted_at,soft_delete" json:"-"`
	}

	PostAssembled struct {
		Post     Post          `pg:"-" json:"post"`
		Comments []PostComment `pg:"-" json:"comments"`
		IsLiked  bool          `pg:"-" json:"is_liked"`
	}

	PostBrief struct {
		ID        uint64  `pg:"id" json:"id"`
		Latitude  float64 `pg:"latitude,use_zero" json:"latitude"`
		Longitude float64 `pg:"longitude,use_zero" json:"longitude"`
		Mine      bool    `pg:"mine" json:"mine"`
	}
)
