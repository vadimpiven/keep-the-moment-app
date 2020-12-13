package postgres

import (
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/FTi130/keep-the-moment-app/back/lib/minio"
)

// Query methods documented here: https://pg.uptrace.dev/queries/

func RegisterIfNewUser(c echo.Context, email string) (err error) {
	db, ctx := extract(c)

	return db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		exists, err := tx.ModelContext(ctx, (*User)(nil)).
			Where("email = ?", email).
			Exists()
		if err != nil {
			return err
		}

		if exists == false {
			_, err = tx.ModelContext(ctx, &User{Email: email}).
				Insert()
			if err != nil {
				return err
			}

			_, err = tx.ModelContext(ctx, &Location{Email: email}).
				Insert()
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func GetUser(c echo.Context, email string) (res User, err error) {
	db, ctx := extract(c)

	err = db.ModelContext(ctx, &res).
		Where("email = ?", email).
		Select()
	return
}

func CheckUserValid(c echo.Context, val *User) (res bool, err error) {
	db, ctx := extract(c)

	err = db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		res, err = tx.ModelContext(ctx, val).
			WherePK().
			Exists()
		if err != nil {
			return err
		} else if res == false {
			return nil
		}

		res, err = tx.ModelContext(ctx, &Image{Path: val.Image}).
			WherePK().
			Exists()
		return err
	})
	return
}

func UpdateUser(c echo.Context, val *User) (err error) {
	db, ctx := extract(c)

	val.Updated = time.Now()

	err = db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		var hashtags []string
		err = tx.ModelContext(ctx, val).
			Column("hashtags").
			WherePK().
			Select(pg.Array(&hashtags))
		if err != nil {
			return err
		}

		for _, hashtag := range hashtags {
			tmp := Hashtag{hashtag, 1}
			_, err = tx.ModelContext(ctx, &tmp).
				OnConflict("DO NOTHING").
				Insert()
			if err != nil {
				return err
			}
			_, err = tx.ModelContext(ctx, &tmp).
				WherePK().
				Set("counter = counter - 1").
				Update()
			if err != nil {
				return err
			}
		}

		_, err = tx.ModelContext(ctx, val).
			WherePK().
			Returning("*").
			Update()
		if err != nil {
			return err
		}

		for _, hashtag := range val.Hashtags {
			tmp := Hashtag{hashtag, 0}
			_, err = tx.ModelContext(ctx, &tmp).
				OnConflict("DO NOTHING").
				Insert()
			if err != nil {
				return err
			}
			_, err = tx.ModelContext(ctx, &tmp).
				WherePK().
				Set("counter = counter + 1").
				Update()
			if err != nil {
				return err
			}
		}

		return nil
	})
	return
}

func GetHashtagsBeginningWith(c echo.Context, key string) (val []string, err error) {
	db, ctx := extract(c)

	err = db.ModelContext(ctx, (*Hashtag)(nil)).
		ColumnExpr("array_agg(name)").
		Group("counter").
		Where("name LIKE ? || '%'", key).
		Order("counter DESC").
		Limit(10).
		Select(pg.Array(&val))
	return
}

func UploadNewImage(c echo.Context, img []byte) (name string, err error) {
	db, ctx := extract(c)

	for i := 0; i < 3; i++ {
		if tmp, err := uuid.NewRandom(); err != nil {
			return "", err
		} else {
			name = tmp.String() + ".png"
		}

		err = db.RunInTransaction(ctx, func(tx *pg.Tx) error {
			_, err = db.ModelContext(ctx, &Image{Path: name}).
				Insert()
			if err != nil {
				return err
			}

			return minio.UploadImage(c, img, name)
		})

		if err == nil {
			break
		} else {
			name = ""
		}
	}
	return
}

func CheckImagesExist(c echo.Context, img []string) (exists bool, err error) {
	db, ctx := extract(c)

	exists = true
	err = db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		for _, path := range img {
			exists, err = tx.ModelContext(ctx, &Image{Path: path}).
				WherePK().
				Exists()
			if err != nil {
				return err
			}
			if exists != true {
				return nil
			}
		}
		return nil
	})
	return
}

func UpdateUserLocation(c echo.Context, email string, lat, lon float64) (err error) {
	db, ctx := extract(c)

	_, err = db.ModelContext(ctx, &Location{email, lat, lon, time.Now()}).
		WherePK().
		Update()
	return
}

func GetUserLocation(c echo.Context, email string) (lat, lon float64, err error) {
	db, ctx := extract(c)

	tmp := Location{Email: email}
	err = db.ModelContext(ctx, &tmp).
		WherePK().
		Select()
	if err != nil {
		return 0, 0, err
	}
	return tmp.Latitude, tmp.Longitude, nil
}

func CreateNewPost(c echo.Context, post *Post) (err error) {
	db, ctx := extract(c)

	err = db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		tmp := Location{Email: post.Email}
		err = tx.ModelContext(ctx, &tmp).
			WherePK().
			Select()
		if err != nil {
			return err
		} else {
			post.Latitude = tmp.Latitude
			post.Longitude = tmp.Longitude
		}

		_, err = tx.ModelContext(ctx, post).
			Returning("*").
			Insert()
		if err != nil {
			return err
		}

		for _, hashtag := range post.Hashtags {
			tmp := Hashtag{hashtag, 0}
			_, err = tx.ModelContext(ctx, &tmp).
				OnConflict("DO NOTHING").
				Insert()
			if err != nil {
				return err
			}
			_, err = tx.ModelContext(ctx, &tmp).
				WherePK().
				Set("counter = counter + 1").
				Update()
			if err != nil {
				return err
			}
		}

		err = tx.ModelContext(ctx, (*User)(nil)).
			Where("email = ?", post.Email).
			ColumnExpr("id, hashtags, image").
			Select(&post.UserID, pg.Array(&post.UserHashtags), &post.UserImage)
		if err != nil {
			return err
		}

		return nil
	})
	return
}

func CheckPostExists(c echo.Context, id uint64) (exists, visible bool, err error) {
	db, ctx := extract(c)

	exists, err = db.ModelContext(ctx, &Post{ID: id}).
		WherePK().
		Exists()
	if err != nil || exists == false {
		return false, false, err
	}

	num, err := db.ModelContext(ctx, (*Post)(nil)).
		Where("id = ?", id).
		Count()
	if err != nil {
		return false, false, err
	}
	return true, num > 0, nil
}

func GetPostByID(c echo.Context, id uint64, email string) (post PostAssembled, err error) {
	db, ctx := extract(c)

	post.Post.ID = id
	post.Comments = []PostComment{}
	err = db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		err = tx.ModelContext(ctx, &post.Post).
			WherePK().
			Select()
		if err != nil {
			return err
		}

		err = tx.ModelContext(ctx, (*User)(nil)).
			Where("email = ?", post.Post.Email).
			ColumnExpr("id, hashtags, image").
			Select(&post.Post.UserID, pg.Array(&post.Post.UserHashtags), &post.Post.UserImage)
		if err != nil {
			return err
		}

		err = tx.ModelContext(ctx, &post.Comments).
			Where("post_id = ?", id).
			Select()
		if err != nil {
			return err
		}

		for i := range post.Comments {
			err = tx.ModelContext(ctx, (*User)(nil)).
				Where("email = ?", post.Comments[i].Email).
				ColumnExpr("id, image").
				Select(&post.Comments[i].UserID, &post.Comments[i].UserImage)
			if err != nil {
				return err
			}
		}

		if email != "" {
			tmp, err := tx.ModelContext(ctx, (*PostLike)(nil)).
				Where("post_id = ?", id).
				Where("email = ?", email).
				Exists()
			if err != nil {
				return err
			}
			post.IsLiked = tmp
		}

		return nil
	})
	return
}

func LikePostByID(c echo.Context, id uint64, email string) (err error) {
	db, ctx := extract(c)

	return db.RunInTransaction(ctx, func(tx *pg.Tx) error {
		liked, err := tx.ModelContext(ctx, (*PostLike)(nil)).
			Where("post_id = ?", id).
			Where("email = ?", email).
			Exists()
		if err != nil {
			return err
		}

		if liked {
			_, err = tx.ModelContext(ctx, (*PostLike)(nil)).
				Where("post_id = ?", id).
				Where("email = ?", email).
				Delete()
			if err != nil {
				return err
			}

			_, err = tx.ModelContext(ctx, &Post{ID: id}).
				WherePK().
				Set("likes = likes - 1").
				Update()
			if err != nil {
				return err
			}
		} else {
			_, err = tx.ModelContext(ctx, &PostLike{PostId: id, Email: email}).
				Insert()
			if err != nil {
				return err
			}

			_, err = tx.ModelContext(ctx, &Post{ID: id}).
				WherePK().
				Set("likes = likes + 1").
				Update()
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func CommentPostByID(c echo.Context, id uint64, email, comment string) (err error) {
	db, ctx := extract(c)

	_, err = db.ModelContext(ctx, &PostComment{PostId: id, Email: email, Content: comment}).
		Insert()
	return
}

func GetVisiblePostIDs(c echo.Context, email string) (posts []PostBrief, err error) {
	db, ctx := extract(c)

	posts = []PostBrief{}
	err = db.ModelContext(ctx, (*Post)(nil)).
		Where("hidden_at IS NULL").
		ColumnExpr("id, latitude, longitude, CASE WHEN email = ? THEN TRUE ELSE FALSE END mine", email).
		Select(&posts)
	return
}
