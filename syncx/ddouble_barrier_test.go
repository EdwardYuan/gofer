package syncx

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestDSync_DoubleBarrier(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(3)

	for i := 0; i < 3; i++ {
		i := i
		go func() {
			defer wg.Done()

			endpoints := []string{"127.0.0.1:2379"}
			dsync, err := NewDSync(clientv3.Config{Endpoints: endpoints})
			assert.NoError(t, err)
			defer func() {
				err = dsync.Close()
				assert.NoError(t, err)
			}()

			doubleBarrier := dsync.NewDoubleBarrier("/defer/doublebarrier1", 3)
			t.Logf("#%d enter", i)
			doubleBarrier.Enter()

			time.Sleep(time.Second)

			t.Logf("#%d leave", i)
			doubleBarrier.Leave()
		}()
	}

	wg.Wait()
}
