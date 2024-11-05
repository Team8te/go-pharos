package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/Team8te/go-pharos/app/ds"
	perror "github.com/Team8te/go-pharos/pkg/platform/error"
)

func (r *Repository) AddDevice(ctx context.Context, d *ds.Device) error {
	query := r.qb.Insert(deviceTable).
		SetMap(newDeviceMap(d)).
		Suffix("RETURNING id, create_at, update_at")

	sql, args, err := query.ToSql()
	if err != nil {
		return perror.MakeError(err, perror.ErrorInternal)
	}

	err = r.db.QueryRowxContext(ctx, sql, args...).StructScan(d)
	if err != nil {
		return perror.MakeError(err, perror.ErrorCodeDB)
	}

	return nil
}

func newDeviceMap(d *ds.Device) map[string]interface{} {
	return map[string]interface{}{
		"uuid": d.UUID,
		"ip":   d.IP,
	}
}

func (r *Repository) ListDevice(ctx context.Context, limit, offset int) ([]*ds.Device, error) {
	sql, args, err := r.qb.Select("*").
		From(deviceTable).
		Limit(uint64(limit)).
		Offset(uint64(offset)).ToSql()
	if err != nil {
		return nil, perror.MakeError(err, perror.ErrorCodeDB)
	}

	result := []*ds.Device{}
	err = r.db.SelectContext(ctx, &result, sql, args...)
	if err != nil {
		return nil, perror.MakeError(err, perror.ErrorCodeDB)
	}

	return result, nil
}

func (r *Repository) FindDevice(ctx context.Context, ip string) (*ds.Device, error) {
	sql, args, err := r.qb.Select("*").
		From(deviceTable).
		Where(sq.And{
			sq.Eq{"ip": ip},
		}).ToSql()
	if err != nil {
		return nil, perror.MakeError(err, perror.ErrorCodeDB)
	}

	result := ds.Device{}
	err = r.db.SelectContext(ctx, &result, sql, args...)
	if err != nil {
		return nil, perror.MakeError(err, perror.ErrorCodeDB)
	}

	return &result, nil
}
