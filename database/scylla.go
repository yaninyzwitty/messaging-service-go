package database

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v3"
)

const (
	maxRetries = 60              // Max retry attempts
	retryDelay = 1 * time.Second // 1 second delay between retries
)

func NewDatabaseConnection(hosts string) (*gocqlx.Session, error) {
	cluster := gocql.NewCluster(hosts)
	cluster.Keyspace = "messaging_keyspace"
	cluster.Consistency = gocql.Quorum
	cluster.Timeout = 10 * time.Second

	var session gocqlx.Session
	var err error

	for retries := 0; retries < maxRetries; retries++ {
		session, err = gocqlx.WrapSession(cluster.CreateSession())
		if err == nil {
			slog.Info("Successfully connected to db ")
			break
		}
		slog.Info(fmt.Sprintf("Failed to connect to database (attempt %d/%d): %v\n", retries+1, maxRetries, err))
	}
	if err != nil {
		return nil, err
	}
	err = session.ExecStmt(`CREATE KEYSPACE IF NOT EXISTS messaging_keyspace WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`)
	if err != nil {
		return nil, fmt.Errorf("failed to create messaging_keyspace keyspace: %w", err)
	}

	// Create persons table if it doesn't exist
	err = session.ExecStmt(`CREATE TABLE IF NOT EXISTS messages (
		id UUID PRIMARY KEY,
		conversation_id UUID,
		sender_id UUID,
		created_at TIMESTAMP,
		updated_at TIMESTAMP,
		body TEXT,
		is_soft_deleted BOOLEAN
	)`)

	if err != nil {
		return nil, fmt.Errorf("failed to create messages table: %w", err)
	}

	return &session, nil

}
