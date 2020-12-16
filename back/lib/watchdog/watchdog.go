// This package prints error when defer *.Close() returns it.
package watchdog

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"

	"github.com/FTi130/keep-the-moment-app/back/lib/firebase"
	"github.com/FTi130/keep-the-moment-app/back/lib/minio"
	"github.com/FTi130/keep-the-moment-app/back/lib/postgres"
)

func clearUnusedImages(ctx context.Context, db *postgres.Postgres, mn *minio.Minio) {
	var unused []string
	err := db.ModelContext(ctx, (*postgres.Image)(nil)).
		ColumnExpr("array_agg(path)").
		Where("path != 'placeholder.png'").
		Where("(now() - uploaded) > (INTERVAL '5 minute')").
		Where("path NOT IN (SELECT image FROM users)").
		Where("path NOT IN (SELECT image_1 FROM posts WHERE image_1 IS NOT NULL)").
		Where("path NOT IN (SELECT image_2 FROM posts WHERE image_2 IS NOT NULL)").
		Where("path NOT IN (SELECT image_3 FROM posts WHERE image_3 IS NOT NULL)").
		Where("path NOT IN (SELECT image_4 FROM posts WHERE image_4 IS NOT NULL)").
		Where("path NOT IN (SELECT image_5 FROM posts WHERE image_5 IS NOT NULL)").
		Select(pg.Array(&unused))
	if err != nil {
		return
	}

	for _, img := range unused {
		_ = db.RunInTransaction(ctx, func(tx *pg.Tx) error {
			_, err = tx.ModelContext(ctx, &postgres.Image{Path: img}).
				WherePK().
				Delete()
			if err != nil {
				return err
			}

			return mn.DeleteImage(ctx, img)
		})
	}
}

func hideOldPosts(ctx context.Context, db *postgres.Postgres) {
	_, err := db.ModelContext(ctx, (*postgres.Post)(nil)).
		Where("hidden_at IS NULL").
		Where("(now() - created_at) > (INTERVAL '12 hour')").
		Delete()
	if err != nil {
		return
	}

	type LastIDPair struct {
		ID    uint64 `pg:"last_id"`
		Email string `pg:"email"`
	}
	var list []LastIDPair
	err = db.ModelContext(ctx, (*postgres.Post)(nil)).
		Where("hidden_at IS NULL").
		ColumnExpr("max(id) AS last_id, email").
		GroupExpr("id, email").
		Select(&list)
	if err != nil {
		return
	}

	for _, pair := range list {
		_, _ = db.ModelContext(ctx, (*postgres.Post)(nil)).
			Where("hidden_at IS NULL").
			Where("email = ?", pair.Email).
			Where("id < ?", pair.ID).
			Where("(now() - created_at) > (INTERVAL '30 minute')").
			Delete()
	}
}

func sendPushNotifications(ctx context.Context, db *postgres.Postgres, f *firebase.Firebase) {

}

func Watch(db *postgres.Postgres, mn *minio.Minio, f *firebase.Firebase) {
	ctx := context.Background()

	go func(ctx context.Context, db *postgres.Postgres, mn *minio.Minio) {
		for {
			time.Sleep(5 * time.Second)
			clearUnusedImages(ctx, db, mn)
		}
	}(ctx, db, mn)

	go func(ctx context.Context, db *postgres.Postgres) {
		for {
			time.Sleep(30 * time.Second)
			hideOldPosts(ctx, db)
		}
	}(ctx, db)

	go func(ctx context.Context, db *postgres.Postgres, f *firebase.Firebase) {
		for {
			time.Sleep(2 * time.Minute)
			sendPushNotifications(ctx, db, f)
		}
	}(ctx, db, f)
}
