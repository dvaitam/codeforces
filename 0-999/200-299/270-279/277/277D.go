package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// problem holds information for each problem
type problem struct {
	scoresmall, scorelarge int64
	timesmall, timelarge   int64
	probpass, probfail     int64
}

// results stores accumulated score and time (in micro units)
type results struct {
	score, time float64
}

// getmax updates a to max(a, b) by score, tie-breaking by smaller time
func getmax(a *results, b results) {
	if a.score < b.score || (a.score == b.score && a.time > b.time) {
		*a = b
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n, t int
	if _, err := fmt.Fscan(reader, &n, &t); err != nil {
		return
	}
	p := make([]problem, n)
	for i := 0; i < n; i++ {
		var d float64
		fmt.Fscan(reader,
			&p[i].scoresmall, &p[i].scorelarge,
			&p[i].timesmall, &p[i].timelarge,
			&d)
		// scale probabilities by 1e6
		pf := int64(d*1e6 + 0.5)
		p[i].probfail = pf
		p[i].probpass = 1000000 - pf
	}
	// sort by custom comparator
	sort.Slice(p, func(i, j int) bool {
		// p[i] < p[j] if p[i].probfail * p[j].probpass * p[i].timelarge < p[j].probfail * p[i].probpass * p[j].timelarge
		return p[i].probfail*p[j].probpass*p[i].timelarge < p[j].probfail*p[i].probpass*p[j].timelarge
	})
	// dp array of results for time 0..t
	dp := make([]results, t+1)
	// iterate problems
	for i := 0; i < n; i++ {
		pr := p[i]
		for j := t - int(pr.timesmall); j >= 0; j-- {
			if j != 0 && dp[j].score == 0 {
				continue
			}
			// take small only
			k := j + int(pr.timesmall)
			var tmp results
			tmp.score = dp[j].score + float64(pr.scoresmall)*1e6
			tmp.time = dp[j].time + float64(pr.timesmall)*1e6
			getmax(&dp[k], tmp)
			// take both small and large
			k2 := k + int(pr.timelarge)
			if k2 > t {
				continue
			}
			// expected additions
			tmp.score += float64(pr.probpass) * float64(pr.scorelarge)
			// additional expected time: probpass * (k2 - current_time_in_seconds)
			tmp.time += float64(pr.probpass) * (float64(k2) - tmp.time/1e6)
			getmax(&dp[k2], tmp)
		}
	}
	// find best over all times
	var ans results
	for i := 0; i <= t; i++ {
		getmax(&ans, dp[i])
	}
	// output scaled back to original units
	fmt.Printf("%.9f %.9f", ans.score/1e6, ans.time/1e6)
}
