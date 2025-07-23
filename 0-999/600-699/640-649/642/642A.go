package main

import (
	"bufio"
	"fmt"
	"os"
)

// This program implements a very simple scheduler for the interactive
// challenge described in problemA.txt. It processes the input in ticks.
// On every tick it receives the list of new submissions and the
// results from invokers. Then it assigns the next test for every
// pending submission if there are free invokers.

// Problem parameters for each problem: only the number of tests is
// required for this naive solution.
type Problem struct {
	tests int
}

type Submission struct {
	problem int
	next    int  // next test to run
	running bool // if a test is currently running on an invoker
	done    bool // whether all tests finished or rejected
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var invokers, p int
	if _, err := fmt.Fscan(in, &invokers, &p); err != nil {
		return
	}
	problems := make([]Problem, p)
	for i := 0; i < p; i++ {
		var tl, tests int
		if _, err := fmt.Fscan(in, &tl, &tests); err != nil {
			return
		}
		problems[i].tests = tests
	}

	submissions := make([]Submission, 0)
	// main loop: read blocks until EOF
	for {
		// read new submissions (list of problem indices terminated by -1)
		var probIdx int
		if _, err := fmt.Fscan(in, &probIdx); err != nil {
			// EOF
			break
		}
		for probIdx != -1 {
			s := Submission{problem: probIdx}
			submissions = append(submissions, s)
			if _, err := fmt.Fscan(in, &probIdx); err != nil {
				return
			}
		}

		// read results from invokers terminated by line "-1 -1"
		for {
			var sid, test int
			if _, err := fmt.Fscan(in, &sid, &test); err != nil {
				return
			}
			if sid == -1 && test == -1 {
				break
			}
			var verdict string
			fmt.Fscan(in, &verdict)
			if sid >= 0 && sid < len(submissions) {
				invokers++
				sub := &submissions[sid]
				sub.running = false
				if verdict == "RJ" {
					sub.done = true
				} else {
					sub.next = test + 1
					if sub.next >= problems[sub.problem].tests {
						sub.done = true
					}
				}
			}
		}

		// schedule new tests while invokers are free
		for i := range submissions {
			if invokers == 0 {
				break
			}
			s := &submissions[i]
			if s.done || s.running {
				continue
			}
			fmt.Fprintf(out, "%d %d\n", i, s.next)
			out.Flush()
			s.running = true
			invokers--
		}
		fmt.Fprintln(out, "-1 -1")
		out.Flush()
	}
}
