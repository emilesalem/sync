package syncqueue_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSyncQueue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sync Queue Suite")
}
