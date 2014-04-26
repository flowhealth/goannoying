package annoying

import (
	"fmt"
	"time"
)

func WaitUntil(name string, condition func() (bool, error), conditionCheckInterval time.Duration, conditionTimeoutInterval time.Duration) (ok bool, err error) {
	failure := make(chan error)
	done := make(chan bool)
	timeout := time.After(conditionTimeoutInterval)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				ok, err := condition()
				if err != nil {
					failure <- err
					return
				}
				if ok {
					done <- true
					return
				}
				time.Sleep(conditionCheckInterval)
			}
		}
	}()
	select {
	case <-done:
		ok = true
	case err = <-failure:
		close(done)
		ok = false
	case <-timeout:
		close(done)
		err = fmt.Errorf("Condition '%s' check timed out after %v", name, conditionTimeoutInterval)
		ok = false
	}
	return
}
