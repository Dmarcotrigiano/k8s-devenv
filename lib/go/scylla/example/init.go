package example

import (
	"fmt"
	"k8s-devenv/scylla"

	"github.com/scylladb/gocqlx/table"
)

func CreateExampleKeyspace(c *scylla.Cluster, keyspace string) {
	c.Session.ExecStmt(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};`, keyspace))
}

func CreateExampleTable(s *scylla.Cluster) error {
	err := s.Session.ExecStmt(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s.example (
		uuid uuid,
		name text,
		age int,
		PRIMARY KEY (uuid));`, s.Keyspace))
	if err != nil {
		return err
	}
	return nil
}

func GetExampleTable(keyspace string) *table.Table {
	MetadataTableMetadata := table.Metadata{
		Name: fmt.Sprintf("%s.worker", keyspace),
		Columns: []string{
			"uuid",
			"name",
			"age",
		},
		PartKey: []string{"uuid"},
	}
	return table.New(MetadataTableMetadata)
}
