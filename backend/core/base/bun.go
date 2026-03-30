package base

import (
	"context"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
)

type QueryOption func(*bun.SelectQuery) *bun.SelectQuery

func WithWhere(cond string, args ...any) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Where(cond, args...)
	}
}

func WithWhereOr(cond string, args ...any) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.WhereOr(cond, args...)
	}
}

func WithEqual(col string, val any) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		if val == nil {
			return q
		}
		return q.Where("? = ?", bun.Ident(col), val)
	}
}

func WithBetween(col string, from, to any) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		if from == nil || to == nil {
			return q
		}
		return q.Where("? BETWEEN ? AND ?", bun.Ident(col), from, to)
	}
}

func WithSearch(col, keyword string) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		if keyword == "" {
			return q
		}
		return q.Where("? LIKE ?", bun.Ident(col), "%"+keyword+"%")
	}
}

func WithSearchOr(cols []string, keyword string) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		if keyword == "" || len(cols) == 0 {
			return q
		}
		parts := make([]string, len(cols))
		args := make([]any, len(cols))
		for i, col := range cols {
			parts[i] = col + " LIKE ?"
			args[i] = "%" + keyword + "%"
		}
		return q.Where("("+strings.Join(parts, " OR ")+")", args...)
	}
}

func WithOrder(order string) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.OrderExpr(order)
	}
}

func WithLimit(n int) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Limit(n)
	}
}

func WithOffset(n int) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Offset(n)
	}
}

func WithRelation(rel string) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Relation(rel)
	}
}

func WithJoin(join string, args ...any) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Join(join, args...)
	}
}

func WithGroup(cols ...string) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.GroupExpr(strings.Join(cols, ", "))
	}
}

func WithHaving(cond string, args ...any) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Having(cond, args...)
	}
}

func WithSelect(cols ...string) QueryOption {
	return func(q *bun.SelectQuery) *bun.SelectQuery {
		return q.Column(cols...)
	}
}

type Repository[T any] struct {
	DB *bun.DB
}

type txKey struct{}

func (r *Repository[T]) IDB(ctx context.Context) bun.IDB {
	if tx, ok := ctx.Value(txKey{}).(bun.Tx); ok {
		return tx
	}
	return r.DB
}

func (r *Repository[T]) applyOpts(q *bun.SelectQuery, opts []QueryOption) *bun.SelectQuery {
	for _, opt := range opts {
		q = opt(q)
	}
	return q
}

func (r *Repository[T]) FindAll(ctx context.Context, out *[]T, opts ...QueryOption) error {
	q := r.IDB(ctx).NewSelect().Model(out)
	q = r.applyOpts(q, opts)
	return q.Scan(ctx)
}

func (r *Repository[T]) FindOne(ctx context.Context, out *T, opts ...QueryOption) error {
	q := r.IDB(ctx).NewSelect().Model(out)
	q = r.applyOpts(q, opts)
	return q.Limit(1).Scan(ctx)
}

func (r *Repository[T]) FindWithPagination(ctx context.Context, out *[]T, page, pageSize int, opts ...QueryOption) (total int, err error) {
	q := r.IDB(ctx).NewSelect().Model(out)
	q = r.applyOpts(q, opts)

	total, err = q.Count(ctx)
	if err != nil {
		return
	}

	if page > 0 && pageSize > 0 {
		q = q.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	err = q.Scan(ctx)
	return
}

func (r *Repository[T]) Count(ctx context.Context, opts ...QueryOption) (int, error) {
	q := r.IDB(ctx).NewSelect().Model(new(T))
	q = r.applyOpts(q, opts)
	return q.Count(ctx)
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	_, err := r.IDB(ctx).NewInsert().Model(entity).Exec(ctx)
	return err
}

func (r *Repository[T]) CreateBulk(ctx context.Context, entities *[]T) error {
	_, err := r.IDB(ctx).NewInsert().Model(entities).Exec(ctx)
	return err
}

func (r *Repository[T]) Upsert(ctx context.Context, entity *T, conflictCols, updateCols []string) error {
	conflict := fmt.Sprintf("CONFLICT (%s) DO UPDATE", strings.Join(conflictCols, ", "))
	q := r.IDB(ctx).NewInsert().Model(entity).On(conflict)
	for _, col := range updateCols {
		q = q.Set(fmt.Sprintf("%s = EXCLUDED.%s", col, col))
	}
	_, err := q.Exec(ctx)
	return err
}

func (r *Repository[T]) UpsertBulk(ctx context.Context, entities *[]T, conflictCols, updateCols []string) error {
	conflict := fmt.Sprintf("CONFLICT (%s) DO UPDATE", strings.Join(conflictCols, ", "))
	q := r.IDB(ctx).NewInsert().Model(entities).On(conflict)
	for _, col := range updateCols {
		q = q.Set(fmt.Sprintf("%s = EXCLUDED.%s", col, col))
	}
	_, err := q.Exec(ctx)
	return err
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	_, err := r.IDB(ctx).NewUpdate().Model(entity).WherePK().Exec(ctx)
	return err
}

func (r *Repository[T]) Delete(ctx context.Context, entity *T) error {
	_, err := r.IDB(ctx).NewDelete().Model(entity).WherePK().Exec(ctx)
	return err
}

type QueryBuilder[T any] struct {
	ctx context.Context
	q   *bun.SelectQuery
}

func (r *Repository[T]) NewQuery(ctx context.Context) *QueryBuilder[T] {
	return &QueryBuilder[T]{
		ctx: ctx,
		q:   r.IDB(ctx).NewSelect().Model(new(T)),
	}
}

// Bisa inject QueryOption yang sama ke builder
func (qb *QueryBuilder[T]) Apply(opts ...QueryOption) *QueryBuilder[T] {
	for _, opt := range opts {
		qb.q = opt(qb.q)
	}
	return qb
}

func (qb *QueryBuilder[T]) Select(cols ...string) *QueryBuilder[T] {
	qb.q = qb.q.Column(cols...)
	return qb
}

func (qb *QueryBuilder[T]) Where(cond string, args ...any) *QueryBuilder[T] {
	qb.q = qb.q.Where(cond, args...)
	return qb
}

func (qb *QueryBuilder[T]) WhereOr(cond string, args ...any) *QueryBuilder[T] {
	qb.q = qb.q.WhereOr(cond, args...)
	return qb
}

func (qb *QueryBuilder[T]) Equal(col string, val any) *QueryBuilder[T] {
	qb.q = qb.q.Where("? = ?", bun.Ident(col), val)
	return qb
}

func (qb *QueryBuilder[T]) Between(col string, from, to any) *QueryBuilder[T] {
	qb.q = qb.q.Where("? BETWEEN ? AND ?", bun.Ident(col), from, to)
	return qb
}

func (qb *QueryBuilder[T]) Search(col, keyword string) *QueryBuilder[T] {
	if keyword != "" {
		qb.q = qb.q.Where("? LIKE ?", bun.Ident(col), "%"+keyword+"%")
	}
	return qb
}

func (qb *QueryBuilder[T]) SearchOr(cols []string, keyword string) *QueryBuilder[T] {
	if keyword == "" || len(cols) == 0 {
		return qb
	}
	parts := make([]string, len(cols))
	args := make([]any, len(cols))
	for i, col := range cols {
		parts[i] = col + " LIKE ?"
		args[i] = "%" + keyword + "%"
	}
	qb.q = qb.q.Where("("+strings.Join(parts, " OR ")+")", args...)
	return qb
}

func (qb *QueryBuilder[T]) Order(order string) *QueryBuilder[T] {
	qb.q = qb.q.OrderExpr(order)
	return qb
}

func (qb *QueryBuilder[T]) Limit(n int) *QueryBuilder[T] {
	qb.q = qb.q.Limit(n)
	return qb
}

func (qb *QueryBuilder[T]) Offset(n int) *QueryBuilder[T] {
	qb.q = qb.q.Offset(n)
	return qb
}

func (qb *QueryBuilder[T]) Relation(rel string) *QueryBuilder[T] {
	qb.q = qb.q.Relation(rel)
	return qb
}

func (qb *QueryBuilder[T]) Join(join string, args ...any) *QueryBuilder[T] {
	qb.q = qb.q.Join(join, args...)
	return qb
}

func (qb *QueryBuilder[T]) Group(cols ...string) *QueryBuilder[T] {
	qb.q = qb.q.GroupExpr(strings.Join(cols, ", "))
	return qb
}

func (qb *QueryBuilder[T]) Having(cond string, args ...any) *QueryBuilder[T] {
	qb.q = qb.q.Having(cond, args...)
	return qb
}

func (qb *QueryBuilder[T]) Find(out *[]T) error {
	return qb.q.Model(out).Scan(qb.ctx)
}

func (qb *QueryBuilder[T]) First(out *T) error {
	return qb.q.Model(out).Limit(1).Scan(qb.ctx)
}

func (qb *QueryBuilder[T]) Count() (int, error) {
	return qb.q.Count(qb.ctx)
}

func (qb *QueryBuilder[T]) Scan(out any) error {
	return qb.q.Scan(qb.ctx, out)
}

func (qb *QueryBuilder[T]) Paginate(out *[]T, page, pageSize int) (total int, err error) {
	total, err = qb.q.Count(qb.ctx)
	if err != nil {
		return
	}
	if page > 0 && pageSize > 0 {
		qb.q = qb.q.Offset((page - 1) * pageSize).Limit(pageSize)
	}
	err = qb.q.Model(out).Scan(qb.ctx)
	return
}

func WithTx(ctx context.Context, tx bun.Tx) context.Context {
	return context.WithValue(ctx, txKey{}, tx)
}

func ExecuteInTransaction(ctx context.Context, db *bun.DB, fn func(ctx context.Context) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	if err := fn(WithTx(ctx, tx)); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}
