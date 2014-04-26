package annoying

import (
	"fmt"
	"time"
)

func WaitFor(name string, condition func() (bool, error), conditionCheckInterval time.Duration, conditionTimeoutInterval time.Duration) error {
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
		break
	case err <- failure:
		close(done)
		return err
	case <-timeout:
		close(done)
		return fmt.Errorf("Condition %s timed out after %v", name, conditionTimeoutInterval)
	}
	return nil
}
