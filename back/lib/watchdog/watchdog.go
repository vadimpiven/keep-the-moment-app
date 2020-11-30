// This package prints error when defer *.Close() returns it.
package watchdog

import (
	"context"
	"time"

	"github.com/go-pg/pg/v10"

	"github.com/FTi130/keep-the-moment-app/back/lib/minio"
	"github.com/FTi130/keep-the-moment-app/back/lib/postgres"
)

func clearUnusedImages(db *postgres.Postgres, mn *minio.Minio) {
	ctx := context.Background()

	var unused []string
	err := db.ModelContext(ctx, (*postgres.Image)(nil)).
		ColumnExpr("array_agg(path)").
		Where("path NOT IN (SELECT image FROM users)").
		Where("(now() - uploaded) > (INTERVAL '5 minute')").
		Select(pg.Array(&unused))
	if err != nil {
		return
	}

	for _, img := range unused {
		if img == "placeholder.png" {
			continue
		}

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

func Watch(db *postgres.Postgres, mn *minio.Minio) {
	go func(db *postgres.Postgres, mn *minio.Minio) {
		for {
			time.Sleep(5 * time.Second)
			clearUnusedImages(db, mn)
		}
	}(db, mn)
}
