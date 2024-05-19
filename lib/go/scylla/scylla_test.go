package scylla

import (
	"k8s-devenv/kubernetes"
	"os"
	"path/filepath"
	"testing"

	"github.com/gocql/gocql"

	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func TestNewScyllaCluster(t *testing.T) {
	hosts := []string{"host1", "host2", "host3"}
	keyspace := "test_keyspace"

	cluster := NewScyllaCluster(hosts, keyspace)

	if cluster == nil {
		t.Error("Expected non-nil cluster object, got nil")
	} else {
		if cluster.Keyspace != keyspace {
			t.Errorf("Expected keyspace %s, got %s", keyspace, cluster.Keyspace)
		}

		c := cluster.cluster
		// Test default consistency level when SCYLLA_LOCAL_DC is not set
		if c.Consistency != gocql.Quorum {
			t.Errorf("Expected consistency level %v, got %v", gocql.Quorum, c.Consistency)
		}
		if c.PoolConfig.HostSelectionPolicy == nil {
			t.Error("Expected non-nil host selection policy, got nil")
		}

		// Since SCYLLA_LOCAL_DC is now set, the consistency level should be LocalQuorum
		t.Setenv("SCYLLA_LOCAL_DC", "dc1")
		cluster = NewScyllaCluster(hosts, keyspace)
		c = cluster.cluster
		if c.Consistency != gocql.LocalQuorum {
			t.Errorf("Expected consistency level %v, got %v", gocql.LocalQuorum, c.Consistency)
		}
	}
}

func TestCreateSession(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Errorf("Error getting user home directory: %v", err)
	}
	kubeconfig := filepath.Join(homeDir, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		t.Errorf("Error building kubeconfig: %v", err)
	}
	clientset, err := k8s.NewForConfig(config)
	if err != nil {
		t.Errorf("Error creating clientset: %v", err)
	}
	ips, err := kubernetes.GetScyllaLoadBalancerAddrs(clientset, "scylla")
	if err != nil {
		t.Errorf("Error getting scylla hosts: %v", err)
	}
	keyspace := "test_keyspace"
	cluster := NewScyllaCluster(ips, keyspace)
	err = cluster.CreateSession()
	if err != nil {
		t.Errorf("Error creating session: %v", err)
	}
	if cluster.Session == nil {
		t.Error("Expected non-nil session, got nil")
	}
}
