package metrics

import (
	"log"
	"time"
)

// LogElapsedTime logs time since start. Typically used with a defer statement:
//  defer LogElapsedTime(time.Now(), "myfunction")
// time.Now() will be evaluated instantly, but the logging won't occur until
// until the end of the function
func LogElapsedTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
