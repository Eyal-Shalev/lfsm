# Lock-Free State Machine
LFSM is a light-weight State Machine implementation that doesn't use any locks.

## Basic Example
```go
package main

import (
    "log"

    "github.com/Eyal-Shalev/lfsm"
)

func main() {
	const (
		closed uint64 = iota
		opened
	)
	s := lfsm.NewState(closed, lfsm.Constraints{
		closed: {opened},
		opened: {closed},
	})

	log.Println(s.Current()) // 2019/12/06 17:34:52 0

	if err := s.Transition(opened); err != nil {
		log.Fatal(err)
	}

	log.Println(s.Current()) // 2019/12/06 17:34:52 1

	if err := s.Transition(closed); err != nil {
		log.Fatal(err)
	}

	log.Println(s.Current()) // 2019/12/06 17:34:52 0
}
```

## Concurrent access - Elevator
```go
package main

import (
    "log"

    "github.com/Eyal-Shalev/lfsm"
)

type Elevator

func main() {
    
}
```