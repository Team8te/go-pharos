package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/go-pharos/app/ds"
	perror "github.com/go-pharos/pkg/platform/error"
)

func (r *Repository) AddDevice(ctx context.Context, d *ds.Device) error {
	query := r.qb.Insert(deviceTable).
		SetMap(newDeviceMap(d)).
		Suffix("RETURNING id")

	sql, args, err := query.ToSql()
	if err != nil {
		return perror.MakeError(err, perror.ErrorInternal)
	}

	err = r.db.QueryRowContext(ctx, sql, args...).Scan(&d.ID)
	if err != nil {
		return perror.MakeError(err, perror.ErrorCodeDB)
	}

	return nil
}

func newDeviceMap(d *ds.Device) map[string]interface{} {
	var id interface{}

	if d.ID == 0 {
		id = sq.Expr("DEFAULT")
	}
	return map[string]interface{}{
		"id":   id,
		"uuid": d.UUID,
		"ip":   d.IP,
	}
}
