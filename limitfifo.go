// Modified by ZRY From go.fifo (github.com/foize/go.fifo) on 2017.08.15
// This is a size limited FIFO modified from github.com/foize/go.fifo
/*
	ModifierZRY := AuthorInfo{
		nickname: "ZRY",
		alias: "swzry",
		website: "http://www.swzry.com/",
		email: "admin@z-touhou.org",
		personal_git: "https://git.swzry.com/",
		github: "github.com/swzry/",
	}
 */
// ------------------------------------------------------------------------
// 'github.com/foize/go.fifo' was Created by Yaz Saito on 06/15/12.
// 'github.com/foize/go.fifo' was Modified by Geert-Johan Riemer, Foize B.V.

package limitfifo

import (
	"sync"
)

const chunkSize = 64

// chunks are used to make a queue auto resizeable.
type chunk struct {
	items       [chunkSize]interface{} // list of queue'ed items
	first, last int                    // positions for the first and list item in this chunk
	next        *chunk                 // pointer to the next chunk (if any)
}

// fifo queue
type Queue struct {
	head, tail 		*chunk     // chunk head and tail
	maxsize			int
	count			int        // total amount of items in the queue
	lock			sync.Mutex // synchronisation lock
}

// NewQueue creates a new and empty *fifo.Queue
func NewLimitQueue(maxChunkNum int) (q *Queue) {
	initChunk := new(chunk)
	q = &Queue{
		maxsize: maxChunkNum * chunkSize,
		head: initChunk,
		tail: initChunk,
	}
	return q
}

// Return the number of items in the queue
func (q *Queue) Len() (length int) {
	// locking to make Queue thread-safe
	q.lock.Lock()
	defer q.lock.Unlock()

	// copy q.count and return length
	length = q.count
	return length
}

// Add an item to the end of the queue
func (q *Queue) Add(item interface{}) {
	// locking to make Queue thread-safe
	q.lock.Lock()
	defer q.lock.Unlock()

	// check if item is valid
	if item == nil {
		panic("can not add nil item to fifo queue")
	}

	// if the tail chunk is full, create a new one and add it to the queue.
	if q.tail.last >= chunkSize {
		if q.count < q.maxsize{
			q.tail.next = new(chunk)
			q.tail = q.tail.next
		}else{
			panic("out of max chunk number")
		}
		
	}

	// add item to the tail chunk at the last position
	q.tail.items[q.tail.last] = item
	q.tail.last++
	q.count++
}

const (
	ERR_ADD_NIL = iota
	ERR_OUT_OF_SIZE_LIMIT
)

type ErrLimitFIFO struct {
	eid int
}

func (this ErrLimitFIFO) Error() string {
	switch(this.eid){
		case ERR_ADD_NIL:
			return "can not add nil item to fifo queue"
		case ERR_OUT_OF_SIZE_LIMIT:
			return "out of max chunk number"
		default:
			return "unknown error"

	}
}

func (q *Queue) SafeAdd(item interface{}) error {
	// locking to make Queue thread-safe
	q.lock.Lock()
	defer q.lock.Unlock()

	// check if item is valid
	if item == nil {
		return ErrLimitFIFO{eid: ERR_ADD_NIL}
	}

	// if the tail chunk is full, create a new one and add it to the queue.
	if q.tail.last >= chunkSize {
		if q.count < q.maxsize{
			q.tail.next = new(chunk)
			q.tail = q.tail.next
		}else{
			return ErrLimitFIFO{eid: ERR_OUT_OF_SIZE_LIMIT}
		}
	}

	// add item to the tail chunk at the last position
	q.tail.items[q.tail.last] = item
	q.tail.last++
	q.count++

	return nil
}

// Remove the item at the head of the queue and return it.
// Returns nil when there are no items left in queue.
func (q *Queue) Next() (item interface{}) {
	// locking to make Queue thread-safe
	q.lock.Lock()
	defer q.lock.Unlock()

	// Return nil if there are no items to return
	if q.count == 0 {
		return nil
	}
	// FIXME: why would this check be required?
	if q.head.first >= q.head.last {
		return nil
	}

	// Get item from queue
	item = q.head.items[q.head.first]

	// increment first position and decrement queue item count
	q.head.first++
	q.count--

	if q.head.first >= q.head.last {
		// we're at the end of this chunk and we should do some maintainance
		// if there are no follow up chunks then reset the current one so it can be used again.
		if q.count == 0 {
			q.head.first = 0
			q.head.last = 0
			q.head.next = nil
		} else {
			// set queue's head chunk to the next chunk
			// old head will fall out of scope and be GC-ed
			q.head = q.head.next
		}
	}

	// return the retrieved item
	return item
}
