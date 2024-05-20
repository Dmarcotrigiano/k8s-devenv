package example

import (
	"k8s-devenv/scylla"
	"k8s-devenv/scylla/models"
)

func InsertExample(sc *scylla.Cluster, e *models.Example) error {
	q := sc.Session.Query(GetExampleTable(sc.Keyspace).Insert()).BindStruct(e)
	return q.ExecRelease()
}
