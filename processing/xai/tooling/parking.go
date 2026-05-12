package tooling

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	_ "github.com/duckdb/duckdb-go/v2"
	_ "modernc.org/sqlite"
	"machin8/processing/xai"
)

const (
	toolingDBPath     = "processing/tooling.db"
	parkingDBPath     = "storage/1/ak.db"
	defaultQueryLimit = 100
	maxQueryLimit     = 500
)

var simpleIdentifierPattern = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)

var (
	parkingDBOnce sync.Once
	parkingDB     *sql.DB
	parkingDBErr  error
)

var parkingHandlers = map[string]xai.ToolHandler{
	"parking_list_tables":      parkingListTables,
	"parking_get_table_schema": parkingGetTableSchema,
	"parking_query":            parkingQuery,
}

func ParkingTools() ([]xai.Tool, error) {
	dbPath, err := filepath.Abs(toolingDBPath)
	if err != nil {
		return nil, fmt.Errorf("resolve tooling db path: %w", err)
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open tooling db: %w", err)
	}
	defer db.Close()

	rows, err := db.Query(`
		SELECT t.type, t.name, t.description, tp.type, tp.properties, tp.required
		FROM tool t
		JOIN tool_param tp ON tp.tool_id = t.id
		ORDER BY t.id
	`)
	if err != nil {
		return nil, fmt.Errorf("query tooling db: %w", err)
	}
	defer rows.Close()

	var tools []xai.Tool
	for rows.Next() {
		var (
			toolType, name, description string
			paramType, propertiesJSON, requiredJSON string
		)
		if err := rows.Scan(&toolType, &name, &description, &paramType, &propertiesJSON, &requiredJSON); err != nil {
			return nil, fmt.Errorf("scan tooling row: %w", err)
		}

		var properties map[string]any
		if err := json.Unmarshal([]byte(propertiesJSON), &properties); err != nil {
			return nil, fmt.Errorf("parse properties for %s: %w", name, err)
		}

		var required []string
		if err := json.Unmarshal([]byte(requiredJSON), &required); err != nil {
			return nil, fmt.Errorf("parse required for %s: %w", name, err)
		}

		params := map[string]any{
			"type":       paramType,
			"properties": properties,
		}
		if len(required) > 0 {
			params["required"] = required
		}

		tools = append(tools, xai.Tool{
			Type:        toolType,
			Name:        name,
			Description: description,
			Parameters:  params,
			Handler:     parkingHandlers[name],
		})
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tooling rows: %w", err)
	}

	return tools, nil
}

func parkingListTables(ctx context.Context, _ json.RawMessage) (any, error) {
	const listTablesSQL = `
SELECT
	table_schema,
	table_name,
	table_type
FROM information_schema.tables
WHERE table_schema NOT IN ('information_schema', 'pg_catalog')
ORDER BY table_schema, table_name
`

	return runDuckDBJSON(ctx, listTablesSQL)
}

func parkingGetTableSchema(ctx context.Context, arguments json.RawMessage) (any, error) {
	var input struct {
		TableName string `json:"table_name"`
		Schema    string `json:"schema"`
	}
	if err := json.Unmarshal(arguments, &input); err != nil {
		return nil, fmt.Errorf("invalid parking_get_table_schema arguments: %w", err)
	}

	if input.Schema == "" {
		input.Schema = "main"
	}

	schema, err := quoteIdentifier(input.Schema)
	if err != nil {
		return nil, fmt.Errorf("invalid schema: %w", err)
	}

	whereClause := fmt.Sprintf("table_schema = %s", schema)
	if input.TableName != "" {
		tableName, err := quoteIdentifier(input.TableName)
		if err != nil {
			return nil, fmt.Errorf("invalid table_name: %w", err)
		}
		whereClause += fmt.Sprintf(" AND table_name = %s", tableName)
	}

	query := fmt.Sprintf(`
SELECT
	table_schema,
	table_name,
	column_name,
	data_type,
	is_nullable,
	column_default,
	ordinal_position
FROM information_schema.columns
WHERE %s
ORDER BY table_name, ordinal_position
`, whereClause)

	return runDuckDBJSON(ctx, query)
}

func parkingQuery(ctx context.Context, arguments json.RawMessage) (any, error) {
	var input struct {
		Query string `json:"query"`
		Limit int    `json:"limit"`
	}
	if err := json.Unmarshal(arguments, &input); err != nil {
		return nil, fmt.Errorf("invalid parking_query arguments: %w", err)
	}

	query, err := normalizeReadOnlyQuery(input.Query)
	if err != nil {
		return nil, err
	}

	limit := input.Limit
	if limit <= 0 {
		limit = defaultQueryLimit
	}
	if limit > maxQueryLimit {
		limit = maxQueryLimit
	}

	wrappedQuery := fmt.Sprintf("SELECT * FROM (%s) AS parking_query LIMIT %d", query, limit)
	return runDuckDBJSON(ctx, wrappedQuery)
}

func runDuckDBJSON(ctx context.Context, query string) (any, error) {
	db, err := parkingReadOnlyDB()
	if err != nil {
		return nil, err
	}

	rows, err := db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("duckdb query failed: %w", err)
	}
	defer rows.Close()

	resultRows, err := scanRows(rows)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"database": parkingDBPath,
		"rows":     resultRows,
		"count":    len(resultRows),
	}, nil
}

func normalizeReadOnlyQuery(query string) (string, error) {
	query = strings.TrimSpace(query)
	query = strings.TrimRight(query, " \t\r\n;")
	if query == "" {
		return "", fmt.Errorf("query is required")
	}
	if strings.Contains(query, ";") {
		return "", fmt.Errorf("multiple statements are not allowed")
	}

	lower := strings.ToLower(query)
	if !strings.HasPrefix(lower, "select") && !strings.HasPrefix(lower, "with") {
		return "", fmt.Errorf("only SELECT or WITH queries are allowed")
	}

	return query, nil
}

func quoteIdentifier(identifier string) (string, error) {
	if !simpleIdentifierPattern.MatchString(identifier) {
		return "", fmt.Errorf("identifier %q contains unsupported characters", identifier)
	}

	return "'" + identifier + "'", nil
}

func parkingReadOnlyDB() (*sql.DB, error) {
	parkingDBOnce.Do(func() {
		dbPath, err := filepath.Abs(parkingDBPath)
		if err != nil {
			parkingDBErr = fmt.Errorf("resolve duckdb path: %w", err)
			return
		}

		db, err := sql.Open("duckdb", dbPath+"?access_mode=read_only")
		if err != nil {
			parkingDBErr = fmt.Errorf("open duckdb: %w", err)
			return
		}

		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		db.SetConnMaxLifetime(0)
		db.SetConnMaxIdleTime(0)

		pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := db.PingContext(pingCtx); err != nil {
			_ = db.Close()
			parkingDBErr = fmt.Errorf("ping duckdb: %w", err)
			return
		}

		parkingDB = db
	})

	if parkingDBErr != nil {
		return nil, parkingDBErr
	}

	return parkingDB, nil
}

func scanRows(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("read duckdb columns: %w", err)
	}

	result := make([]map[string]any, 0, 32)
	for rows.Next() {
		values := make([]any, len(columns))
		dest := make([]any, len(columns))
		for i := range values {
			dest[i] = &values[i]
		}

		if err := rows.Scan(dest...); err != nil {
			return nil, fmt.Errorf("scan duckdb row: %w", err)
		}

		row := make(map[string]any, len(columns))
		for i, column := range columns {
			row[column] = jsonValue(values[i])
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate duckdb rows: %w", err)
	}

	return result, nil
}

func jsonValue(value any) any {
	switch v := value.(type) {
	case nil:
		return nil
	case []byte:
		return string(v)
	case time.Time:
		return v
	default:
		return v
	}
}
