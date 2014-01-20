// This is a derivative work of roger peppe's stackgraph command.
// For more details, see http://code.google.com/p/rog-go/
//

package traceback

type StackStatus string

const (
	StackStatusChanReceive = "chan receive"
	StackStatusSemAcquire  = "semacquire"
	StackStatusRunning     = "running"
	StackStatusRunnable    = "runnable"
)

// Call represents a function call.
type Call struct {
	Func   string
	Source string
	Line   int
	Args   []uint64
}

// Stack represents the call stack of a goroutine.
type Stack struct {
	ID     int
	Status StackStatus
	Calls  []Call
}
