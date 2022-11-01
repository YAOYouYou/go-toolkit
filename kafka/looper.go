package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/YAOYouYou/go-toolkit/logging"
)

type Looper struct {
	consumer   kafka.Consumer
	topic      string
	sigC       chan os.Signal
	eventCb    func(e *kafka.Message) error
	workers    int
	count      int
	gen        int32
	offsetsMap map[string]kafka.TopicPartition
	wg         sync.WaitGroup

	filter  func(m *kafka.Message) bool
	errorCb func(m *kafka.Message, err error) error
	closeCb func()
}

func (l *Looper) Run() {
	rebalanceCb := func(consumer *kafka.Consumer, ev kafka.Event) error {
		switch e := ev.(type) {
		case kafka.AssignedPartitions:
			logging.Infof("Rebalance - Assigned:", e.Partitions)

			// reset
			atomic.AddInt32(&l.gen, 1)
			l.offsetsMap = make(map[string]kafka.TopicPartition)

			consumer.Assign(e.Partitions)

		case kafka.RevokedPartitions:
			logging.Infof("Rebalance - Revoked:", e.Partitions)
			consumer.Unassign()
		}
		return nil
	}
	l.consumer.Subscribe(l.topic, rebalanceCb)

	signal.Notify(l.sigC, syscall.SIGINT, syscall.SIGTERM)

Loop:
	for {
		select {
		case sig := <-l.sigC:
			logging.Infof("System interupt dected:", sig)
			l.closeCb()
			break Loop
		default:
			ev := l.consumer.Poll(100)
			gen := atomic.LoadInt32(&l.gen)
			if ev == nil {
				logging.Debugf("Heartbeat")
				continue
			}
			switch e := ev.(type) {
			case *kafka.Message:
				logging.Debugf("receive messages@", e.TopicPartition)
				l.wg.Add(1)
				go func(e *kafka.Message, g int32) {
					defer func() {
						l.wg.Done()

						// check rebalance
						if g == atomic.LoadInt32(&l.gen) {
							// Store current offset based on current message topic partition
							key := fmt.Sprintf("%s[%d]", *e.TopicPartition.Topic, e.TopicPartition.Partition)
							l.offsetsMap[key] = e.TopicPartition
						}
					}()

					if l.filter != nil {
						if !l.filter(e) {
							return
						}
					}
					// call eventCb process message
					// eventCb err are not handled by default
					err := l.eventCb(e)
					if l.errorCb != nil {
						l.errorCb(e, err)
					}
				}(e, gen)

				if l.count%l.workers == 0 {
					l.wg.Wait()
					logging.Infof("receive mesaage count: %d", l.count)
					l.commit(gen)
				}
			default:
				logging.Debugf("other events:", e)
			}
		}
	}
}

func (l *Looper) commit(gen int32) {
	if len(l.offsetsMap) == 0 {
		return
	}
	// check rebalanced
	if gen != atomic.LoadInt32(&l.gen) {
		return
	}
	tps := make([]kafka.TopicPartition, len(l.offsetsMap))
	index := 0
	for _, tp := range l.offsetsMap {
		tp.Offset = tp.Offset + 1
		tps[index] = tp
		index++
	}
	if _, err := l.consumer.CommitOffsets(tps); err != nil {
		logging.Errorf("CommitOffsets Error: %v", err)
		l.Close()
		return
	}
}

func (l *Looper) Close() {
	l.consumer.Close()
	if l.closeCb != nil {
		l.closeCb()
	}
}