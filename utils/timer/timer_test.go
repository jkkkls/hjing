// This package from github.com/garallz/gotools/timer
// garalinluzhi@gmail.com
package timer

import (
	"fmt"
	"testing"
	"time"
)

func TestTimer(t *testing.T) {
	InitTicker(1)

	SetTimeZone("+00:01")
	time.Sleep(time.Second)

	fmt.Println(time.Now(), time.Now().UTC().Unix())
	if err := NewTimer("00", 4, false, "time [:00] one", display); err != nil {
		t.Error(err)
	}

	// if err := NewTimer("15", 4, false, "time [:15] one", display); err != nil {
	// 	t.Error(err)
	// }

	// if err := NewTimer("30", 4, false, "time [:30] two", display); err != nil {
	// 	t.Error(err)
	// }

	// if err := NewTimer("45", 4, false, "time [:45] two", display); err != nil {
	// 	t.Error(err)
	// }

	if err := NewTimer("06:00", 4, false, "time [06:00] one", display); err != nil {
		t.Error(err)
	}

	if err := NewTimer("07:00", 4, false, "time [07:00] one", display); err != nil {
		t.Error(err)
	}

	if err := NewTimer("08:00", 4, false, "time [08:00] one", display); err != nil {
		t.Error(err)
	}

	//	NewRunDuration(time.Second*3, -1, nil, oneprint)

	//	NewRunTime(time.Now().Add(time.Second*2), nil, twoprint)

	// for i := 1; i < 10; i++ {
	// 	NewRunTime(time.Now().Add(time.Second), i*111, display)
	// 	time.Sleep(time.Millisecond * 200)
	// }

	time.Sleep(time.Minute * 3)
}

func display(data interface{}) {
	fmt.Println(time.Now(), data)
}

func eachprint(data interface{}) {
	fmt.Println(time.Now(), "each second to display")
}

func oneprint(data interface{}) {
	fmt.Println(time.Now(), "one second last display")
}

func twoprint(data interface{}) {
	fmt.Println(time.Now(), "two second last display")
}

func TestTimeZone(t *testing.T) {
	a := time.Now().Local().Unix()
	b := time.Now().UTC().Unix()
	fmt.Println(a, b)
}
