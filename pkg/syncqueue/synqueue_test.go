package syncqueue_test

import (
	"context"
	"strconv"
	"sync"

	"github.com/emilesalem/sync/v2/pkg/syncqueue"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("sync queue", func() {

	var syncQueue syncqueue.Syncqueue[string]
	var ctx context.Context
	var cancel context.CancelFunc

	BeforeEach(func() {
		ctx, cancel = context.WithCancel(context.Background())
	})

	AfterEach(func() {
		cancel()
	})

	Context(`sequential`, func() {

		When(`capacity not exceeded`, func() {
			var capacity int
			var messages []string

			BeforeEach(func() {
				capacity = 10

				messages = make([]string, 10)
				for i := range messages {
					messages[i] = `message_` + strconv.Itoa(i)
				}
				syncQueue = *syncqueue.NewSyncqueue[string](ctx, syncqueue.Options{Capacity: capacity})
			})

			It(`should add to queue and pop head`, func() {
				for _, m := range messages {
					ok := syncQueue.Add(m)
					Expect(ok).To(BeTrue())
				}
				for _, m := range messages {
					msg := <-syncQueue.Read()
					Expect(msg).To(Equal(m))
				}
			})
		})

		When(`capacity exceeded`, func() {
			var capacity int
			var messages []string

			BeforeEach(func() {
				capacity = 10

				messages = make([]string, 11)
				for i := range messages {
					messages[i] = `message_` + strconv.Itoa(i)
				}
				syncQueue = *syncqueue.NewSyncqueue[string](ctx, syncqueue.Options{Capacity: capacity})
			})

			It(`should add to queue and return false`, func() {
				for i, m := range messages {
					if i < capacity {
						ok := syncQueue.Add(m)
						Expect(ok).To(BeTrue())
					}
				}
				ok := syncQueue.Add(messages[capacity])
				Expect(ok).To(BeFalse())

				for i := 1; i <= capacity; i++ {
					msg := <-syncQueue.Read()
					Expect(msg).To(Equal(messages[i]))
				}
				lastMessage := `last_message`

				syncQueue.Add(lastMessage)

				msg := <-syncQueue.Read()

				Expect(msg).To(Equal(lastMessage))
			})
		})
	})

	Context(`concurrent`, func() {
		var capacity int
		var messages []string

		BeforeEach(func() {
			capacity = 10

			messages = make([]string, 10)
			for i := range messages {
				messages[i] = `message_` + strconv.Itoa(i)
			}
			syncQueue = *syncqueue.NewSyncqueue[string](ctx, syncqueue.Options{Capacity: capacity})
		})

		It(`should add and pop`, func() {
			nbWorkers := 100
			var w sync.WaitGroup

			for i := 0; i < nbWorkers; i++ {
				w.Add(1)
				if i < 10 {
					go func(i int) {
						<-syncQueue.Read()
						w.Done()
					}(i)
				} else {
					go func() {
						for _, m := range messages {
							_ = syncQueue.Add(m)
						}
						w.Done()
					}()
				}
			}
			w.Wait()
			ok := syncQueue.Add(`message`)
			Expect(ok).To(BeFalse())
		})
	})

})
