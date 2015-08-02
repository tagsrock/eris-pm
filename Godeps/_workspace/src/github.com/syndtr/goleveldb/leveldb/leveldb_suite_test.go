package leveldb

import (
	"testing"

	"github.com/eris-ltd/eris-pm/Godeps/_workspace/src/github.com/syndtr/goleveldb/leveldb/testutil"
)

func TestLevelDB(t *testing.T) {
	testutil.RunSuite(t, "LevelDB Suite")
}
