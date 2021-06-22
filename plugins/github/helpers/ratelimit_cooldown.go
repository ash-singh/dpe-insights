package helpers

import (
	"log"
	"time"

	"github.com/google/go-github/v32/github"
)

// RateLimitCoolDown sleep till the rate limit is reset.
func RateLimitCoolDown(rate *github.Rate) {
	if rate.Remaining < 5 {

		reset := rate.Reset
		now := time.Now()

		d := reset.Sub(now)

		log.Println("Rate limited for ", d)

		time.Sleep(d)
	}
}
