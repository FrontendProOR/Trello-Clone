package config

import (
	"time"

	"github.com/gocql/gocql"
)

func ConnectToCassandra() (*gocql.Session, error) {
	cluster := gocql.NewCluster("cassandra")
	cluster.Consistency = gocql.Quorum
	cluster.ConnectTimeout = 15 * time.Second

	sysSession, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}
	defer sysSession.Close()

	// napravi keyspace ako ne postoji
	err = sysSession.Query(`
		CREATE KEYSPACE IF NOT EXISTS test_keyspace 
		WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};
	`).Exec()
	if err != nil {
		return nil, err
	}

	// sada se konektujemo NA keyspace
	cluster.Keyspace = "test_keyspace"
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
}
