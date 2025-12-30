package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const (
	INF = 1e9 + 7
	N   = 22
)

// Task structure replacing pair<pair<int,int>,pair<int,int>>
type Task struct {
	Start    int // Busy interval start (T[i].st.st)
	Len      int // Busy interval length (T[i].st.nd)
	Duration int // Task execution time (T[i].nd.st)
	ID       int // Original index (T[i].nd.nd)
}

var (
	n     int
	tasks []Task
	// Global DP arrays to handle large size (approx 16MB each)
	dp    [1 << N]int
	jaki  [1 << N]int // 'which' task was added
	gdzie [1 << N]int // 'where' (start time)
	ans   [N][2]int   // Result storage: [processor_id, start_time]
)

// nxt calculates the earliest finish time for a task of duration 't'
// starting at or after 'x', given the mask 'm' of already scheduled tasks.
func nxt(m, x, t int) int {
	l1, l2 := 0, 0

	// 1. Skip fixed busy intervals that end before x
	for l1 < n && x >= tasks[l1].Start {
		l1++
	}

	// 2. Check if x falls inside the previous busy interval
	if l1 > 0 && x <= tasks[l1-1].Start+tasks[l1-1].Len-1 {
		l1-- // Move back to handle the jump in the main loop
	} else {
		// 3. If x is valid wrt fixed intervals, check against task deadlines (from mask m)
		// Find the first task in mask m that starts after x
		for l2 < n && (x >= tasks[l2].Start || (m&(1<<l2)) == 0) {
			l2++
		}
		// If no constraint found or we fit before the next constraint
		if l2 == n || x+t < tasks[l2].Start {
			return x + t
		}
	}

	// 4. Iterate forward, jumping over busy intervals to find a valid slot
	for l1 < n {
		x = tasks[l1].Start + tasks[l1].Len

		// Optimization: If the next interval starts exactly where this one ends, skip effectively merging them
		if l1 < n-1 && tasks[l1+1].Start == x {
			l1++
			continue
		}

		// Check constraints again for the new x
		for l2 < n && (x >= tasks[l2].Start || (m&(1<<l2)) == 0) {
			l2++
		}

		if l2 == n || x+t < tasks[l2].Start {
			return x + t
		}
		l1++
	}
	return x + t
}

func main() {
	// Fast I/O
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var P int
	if _, err := fmt.Fscan(reader, &n, &P); err != nil {
		return
	}

	tasks = make([]Task, n)
	for i := 0; i < n; i++ {
		var s, l, d int
		fmt.Fscan(reader, &s, &l, &d)
		tasks[i] = Task{Start: s, Len: l, Duration: d, ID: i}
		
		// Specific constraint from original code
		if s == 1 {
			fmt.Fprintln(writer, "NO")
			return
		}
	}

	// Sort tasks by Start time (critical for the nxt logic)
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Start != tasks[j].Start {
			return tasks[i].Start < tasks[j].Start
		}
		return tasks[i].ID < tasks[j].ID
	})

	// Initialize DP
	// dp[mask] = earliest finish time for the set of tasks in mask
	limit := 1 << n
	for i := 0; i < limit; i++ {
		dp[i] = INF
	}
	dp[0] = 1

	// Bitmask DP Loop
	for i := 0; i < limit; i++ {
		if dp[i] == INF {
			continue
		}
		// Try adding task j to the set i
		for j := 0; j < n; j++ {
			// If task j is NOT in set i yet
			if (i & (1 << j)) == 0 {
				nextMask := i | (1 << j)
				// Calculate potential finish time
				// We pass 'nextMask' because nxt uses the mask to find deadlines imposed by other tasks
				x := nxt(nextMask, dp[i], tasks[j].Duration)
				
				// Constraint: Task must finish before its own fixed start time (deadline)
				if x >= tasks[j].Start {
					x = INF
				}

				if dp[nextMask] > x {
					dp[nextMask] = x
					jaki[nextMask] = j
					gdzie[nextMask] = x - tasks[j].Duration // Store start time
				}
			}
		}
	}

	// Check results
	fullMask := (1 << n) - 1
	for i := 0; i < limit; i++ {
		// Valid conditions:
		// 1. P=1: Check if full set (i == fullMask) is valid
		// 2. P=2: Check if set i is valid AND complementary set (fullMask ^ i) is valid
		
		complement := fullMask ^ i
		isValid := false
		
		if P == 2 {
			if dp[i] < INF && dp[complement] < INF {
				isValid = true
			}
		} else {
			if i == fullMask && dp[i] < INF {
				isValid = true
			}
		}

		if isValid {
			// Reconstruct Processor 1 schedule (from set i)
			curr := i
			for curr > 0 {
				idx := jaki[curr]
				originalID := tasks[idx].ID
				ans[originalID][0] = 1
				ans[originalID][1] = gdzie[curr]
				curr ^= (1 << idx)
			}

			// Reconstruct Processor 2 schedule (from complement set)
			curr = complement
			for curr > 0 {
				idx := jaki[curr]
				originalID := tasks[idx].ID
				ans[originalID][0] = 2
				ans[originalID][1] = gdzie[curr]
				curr ^= (1 << idx)
			}

			fmt.Fprintln(writer, "YES")
			for k := 0; k < n; k++ {
				fmt.Fprintf(writer, "%d %d\n", ans[k][0], ans[k][1])
			}
			return
		}
	}

	fmt.Fprintln(writer, "NO")
}
