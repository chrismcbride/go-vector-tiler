package metrics

import (
	"log"
	"time"
)

func LogElapsedTime(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("%s took %s", name, elapsed)
}
