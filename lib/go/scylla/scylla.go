package scylla

import (
	"log/slog"
	"os"

	"github.com/gocql/gocql"
	"github.com/scylladb/gocqlx/v2"
)

type Cluster struct {
	cluster  *gocql.ClusterConfig
	Session  *gocqlx.Session
	Keyspace string
}

func NewScyllaCluster(hosts []string, keyspace string) *Cluster {
	slog.Debug("Creating Scylla cluster", "hosts", hosts, "keyspace", keyspace)
	localDC := os.Getenv("SCYLLA_LOCAL_DC")
	c := gocql.NewCluster(hosts...)
	c.Consistency = gocql.Quorum
	fallback := gocql.RoundRobinHostPolicy()
	if localDC != "" {
		fallback = gocql.DCAwareRoundRobinPolicy(localDC)
	}
	c.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(fallback)

	// If using multi-dc cluster use the "local" consistency levels.
	if localDC != "" {
		c.Consistency = gocql.LocalQuorum
	}
	return &Cluster{
		c,
		nil,
		keyspace,
	}
}

func (c *Cluster) CreateSession() error {
	session, err := gocqlx.WrapSession(c.cluster.CreateSession())
	if err != nil {
		return err
	}
	c.Session = &session
	return nil
}
