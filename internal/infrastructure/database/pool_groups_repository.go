package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupsRepository interface {
	UpsertGroup(ctx context.Context, groupName string) (int, error)
}

type GroupsPoolRepository struct {
	Pool *pgxpool.Pool
}

func NewGroupsPoolRepository(pool *pgxpool.Pool) *GroupsPoolRepository {
	return &GroupsPoolRepository{Pool: pool}
}

func (r *GroupsPoolRepository) UpsertGroup(ctx context.Context, groupName string) (int, error) {
	var groupID int

	err := r.Pool.QueryRow(ctx, `
    INSERT INTO groups (name)
    VALUES ($1)
    ON CONFLICT (name) DO UPDATE SET name = $1
    RETURNING id
  `, groupName).Scan(&groupID)
	if err != nil {
		return 0, fmt.Errorf("upserting group: %w", err)
	}

	return groupID, nil
}
