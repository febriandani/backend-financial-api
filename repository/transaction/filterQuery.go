package transaction

import (
	"strings"

	mt "github.com/febriandani/backend-financial-api/domain/model/transaction"
)

func BuildQueryGetTransactions(baseQuery string, filter mt.Filter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.UserID != 0 {
		conditions = append(conditions, "t.user_id = ?")
		args = append(args, filter.UserID)
	}

	if filter.CategoryType.Valid && filter.CategoryType.String != "" {
		conditions = append(conditions, "t.category_type = ?")
		args = append(args, filter.CategoryType.String)
	}

	if filter.StartDate.Valid && filter.StartDate.String != "" && filter.EndDate.Valid && filter.EndDate.String != "" {
		conditions = append(conditions, "t.created_at >= ? and t.created_at <= ?")
		args = append(args, filter.StartDate.String, filter.EndDate.String)
	} else {
		conditions = append(conditions, `
        EXTRACT(YEAR FROM t.created_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') = EXTRACT(YEAR FROM CURRENT_TIMESTAMP) AND
        EXTRACT(MONTH FROM t.created_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') = EXTRACT(MONTH FROM CURRENT_TIMESTAMP)`)
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
	}

	if filter.Offset.Valid && filter.Limit.Valid && filter.Offset.Int64 != 0 && filter.Limit.Int64 != 0 {
		baseQuery += " ORDER BY t.created_at desc OFFSET ((? - 1) * ?) ROWS FETCH NEXT ? ROWS ONLY"
		args = append(args, filter.Offset.Int64, filter.Limit.Int64, filter.Limit.Int64)
	} else {
		baseQuery += " ORDER BY t.created_at desc"
	}

	return baseQuery, args
}

func BuildQueryGetCurrentBalanceTransactions(baseQuery string, filter mt.Filter) (string, []interface{}) {
	var conditions []string
	var args []interface{}

	if filter.UserID != 0 {
		conditions = append(conditions, "t.user_id = ?")
		args = append(args, filter.UserID)
	}

	if filter.StartDate.Valid && filter.StartDate.String != "" && filter.EndDate.Valid && filter.EndDate.String != "" {
		conditions = append(conditions, " t.created_at >= ? and t.created_at <= ?")
		args = append(args, filter.StartDate.String, filter.EndDate.String)
	} else {
		conditions = append(conditions, `
        EXTRACT(YEAR FROM t.created_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') = EXTRACT(YEAR FROM CURRENT_TIMESTAMP) AND
        EXTRACT(MONTH FROM t.created_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta') = EXTRACT(MONTH FROM CURRENT_TIMESTAMP)`)
	}

	if len(conditions) > 0 {
		whereClause := " WHERE " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
	}

	return baseQuery, args
}
