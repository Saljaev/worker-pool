package workers

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func do(id uuid.UUID, payload interface{}, onFinish func()) {
	time.Sleep(500 * time.Millisecond)
	fmt.Printf("worker %s work on %v\n", id, payload)

	onFinish()
}
