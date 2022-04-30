package eventLogger

import (
	"os"
	"testing"

	"time"
)

func init() {

	os.Setenv("LOGGER_HOST", "localhost")
	os.Setenv("LOGGER_PORT", "27017")
	os.Setenv("LOGGER_DB_NAME", "logs")
	os.Setenv("LICENCE_KEY", "62345678uhbvcfgxdtfhujnbvcdtyuijknbvesrtyuhjbvdfgiuhbjvchxrdtyu")

}

func Test_EventLogger(t *testing.T) {

	GetInstance().Listen()
	ticker := time.NewTicker(2 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				b, err := time.Now().MarshalJSON()
				if err != nil {
					t.Fatal(err)
				}
				GetInstance().sender(b)
			case <-quit:
				t.Log("stop")
				ticker.Stop()
				return
			}
		}
	}()
	<-quit
}
