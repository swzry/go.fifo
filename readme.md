## go.limitfifo

### Description

This package is modified from go.fifo (github.com/foize/go.fifo).
It's a buffer size limited version in order to prevent memory overflow in some application scenes.

e.g.
If you need a fifo buffer for a serial port communicating multiplexer,
when the data are coming from many input serial ports,
we put the data into a fifo buffer. Then we fowarding them into another serial port.
Unfortunately, the speed of output port is not always enough to send those data.
If the buffer is unlimited, the memory overflow may occur. To avoid this failure,
we should drop some input data when the buffer is too large.

### Installation
`go get github.com/swzry/go.limitfifo`

### Usage
```go
package main

import (
	"github.com/swzry/go.limitfifo"
	"fmt"
)

func main() {
	// create a new queue
	numbers := fifo.NewLimitQueue(4)  // 4 is the size of chunk, the size of each chunk is 64.

	// add items to the queue
	numbers.Add(5)
	numbers.Add(6)
	numbers.Add(9)

	// add items to the queue, use 'SafeAdd' to catch if there's an error occur.
    err := numbers.SafeAdd(32)

	// retrieve items from the queue
	fmt.Println(numbers.Next()) // 5
	fmt.Println(numbers.Next()) // 6
	fmt.Println(numbers.Next()) // 9
	fmt.Println(numbers.Next()) // 32
}
```

### Documentations and More

You can refer to the [README.md of 'go.fifo'](https://github.com/foize/go.fifo/blob/master/readme.md).


### The README.md of go.fifo

#### Description
go.limitfifo provides a simple FIFO thread-safe queue.
*fifo.Queue supports pushing an item at the end with Add(), and popping an item from the front with Next().
There is no intermediate type for the stored data. Data is directly added and retrieved as type interface{}
The queue itself is implemented as a single-linked list of chunks containing max 64 items each.

#### Installation
`go get github.com/foize/go.fifo`

#### Usage
```go
package main

import (
	"github.com/foize/go.fifo"
	"fmt"
)

func main() {
	// create a new queue
	numbers := fifo.NewQueue()

	// add items to the queue
	numbers.Add(42)
	numbers.Add(123)
	numbers.Add(456)

	// retrieve items from the queue
	fmt.Println(numbers.Next()) // 42
	fmt.Println(numbers.Next()) // 123
	fmt.Println(numbers.Next()) // 456
}
```

```go
package main

import (
	"github.com/foize/go.fifo"
	"fmt"
)

type thing struct {
	Text string
	Number int
}

func main() {
	// create a new queue
	things := fifo.NewQueue()

	// add items to the queue
	things.Add(&thing{
		Text: "one thing",
		Number: 1,
	})
	things.Add(&thing {
		Text: "another thing",
		Number: 2,
	})

	// retrieve items from the queue
	for  {
		// get a new item from the things queue
		item := things.Next();

		// check if there was an item
		if item == nil {
			fmt.Println("queue is empty")
			return
		}

		// assert the type for the item
		someThing := item.(*thing)

		// print the fields
		fmt.Println(someThing.Text)
		fmt.Printf("with number: %d\n", someThing.Number)
	}
}

/* output: */
// one thing
// with number: 1
// another thing
// with number: 2
// queue is empty
```

#### Documentation
Documentation can be found at [godoc.org/github.com/foize/go.fifo](http://godoc.org/github.com/foize/go.fifo).
For more detailed documentation, read the source.

#### History
This package is based on github.com/yasushi-saito/fifo_queue
There are several differences:
- renamed package to `fifo` to make usage simpler
- removed intermediate type `Item` and now directly using interface{} instead.
- renamed (*Queue).PushBack() to (*Queue).Add()
- renamed (*Queue).PopFront() to (*Queue).Next()
- Next() will not panic on empty queue, will just return nil interface{}
- Add() does not accept nil interface{} and will panic when trying to add nil interface{}.
- Made fifo.Queue thread/goroutine-safe (sync.Mutex)
- Added a lot of comments
- renamed internal variable/field names