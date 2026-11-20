package TestIssues

import (
	"errors"
	"testing"

	"github.com/sijms/go-ora/v2/network"
)

// TestIssue576 checks that a PL/SQL block ending with an unhandled
// NO_DATA_FOUND returns ORA-01403, while a query without rows still succeeds.
func TestIssue576(t *testing.T) {
	db, err := getDB()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Error(err)
		}
	}()

	blocks := []string{
		"BEGIN RAISE NO_DATA_FOUND; END;",
		`DECLARE
	n NUMBER;
BEGIN
	SELECT 1 INTO n FROM dual WHERE 1 = 0;
END;`,
		`DECLARE
	n NUMBER;
BEGIN
	SELECT 1 INTO n FROM dual WHERE 1 = 0;
EXCEPTION
	WHEN OTHERS THEN
		RAISE;
END;`,
	}
	for _, block := range blocks {
		_, err := db.Exec(block)
		var oraErr *network.OracleError
		if !errors.As(err, &oraErr) || oraErr.ErrCode != 1403 {
			t.Errorf("expected ORA-01403 from:\n%s\ngot: %v", block, err)
		}
	}

	rows, err := db.Query("SELECT 1 FROM dual WHERE 1 = 0")
	if err != nil {
		t.Fatalf("query without rows: %v", err)
	}
	defer rows.Close()
	if rows.Next() {
		t.Error("query without rows returned a row")
	}
	if err := rows.Err(); err != nil {
		t.Errorf("query without rows: %v", err)
	}
}
