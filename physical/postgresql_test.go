package physical

import (
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestPostgreSQLBackend(t *testing.T) {
	connURL := os.Getenv("PGURL")
	if connURL == "" {
		t.SkipNow()
	}

	table := os.Getenv("PGTABLE")
	if table == "" {
		table = "vault_kv_store"
	}

	// Run vault tests
	logger := log.New(os.Stderr, "", log.LstdFlags)
	b, err := NewBackend("postgresql", logger, map[string]string{
		"connection_url": connURL,
		"table":          table,
	})

	if err != nil {
		t.Fatalf("Failed to create new backend: %v", err)
	}

	defer func() {
		pg := b.(*PostgreSQLBackend)
		_, err := pg.client.Exec("DROP TABLE " + pg.table)
		if err != nil {
			t.Fatalf("Failed to drop table: %v", err)
		}
	}()

	testBackend(t, b)
	testBackend_ListPrefix(t, b)

}
