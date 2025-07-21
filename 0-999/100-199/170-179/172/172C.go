package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Student holds arrival time, destination, and original index
type Student struct {
	t   int64
	x   int64
	idx int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	fmt.Fscan(reader, &n, &m)
	students := make([]Student, n)
	for i := 0; i < n; i++ {
		var ti, xi int64
		fmt.Fscan(reader, &ti, &xi)
		students[i] = Student{t: ti, x: xi, idx: i}
	}
	ans := make([]int64, n)
	var queue []Student
	nextI := 0
	var curTime int64
	processed := 0

	for processed < n {
		// If no one waiting, wait for next arrival
		if len(queue) == 0 && nextI < n && curTime < students[nextI].t {
			curTime = students[nextI].t
		}
		// Load all who have arrived by now, up to capacity
		for nextI < n && students[nextI].t <= curTime && len(queue) < m {
			queue = append(queue, students[nextI])
			nextI++
		}
		// If not full and students remain, wait for next arrivals until full or none left
		for len(queue) < m && nextI < n {
			curTime = students[nextI].t
			queue = append(queue, students[nextI])
			nextI++
		}
		// Depart with current batch
		departTime := curTime
		k := len(queue)
		batch := make([]Student, k)
		copy(batch, queue)
		sort.Slice(batch, func(i, j int) bool { return batch[i].x < batch[j].x })
		// Compute total unload time
		var unloadTime int64
		var lastX int64 = -1
		var cntAtX int64
		for _, st := range batch {
			if st.x != lastX {
				if lastX != -1 {
					unloadTime += 1 + cntAtX/2
				}
				lastX = st.x
				cntAtX = 1
			} else {
				cntAtX++
			}
		}
		if lastX != -1 {
			unloadTime += 1 + cntAtX/2
		}
		maxX := batch[k-1].x
		// Record answer for each student
		for _, st := range batch {
			ans[st.idx] = departTime + st.x
		}
		// Return to stop
		curTime = departTime + maxX + unloadTime + maxX
		processed += k
		queue = queue[:0]
	}
	// Output results
	for i := 0; i < n; i++ {
		if i > 0 {
			fmt.Fprint(writer, " ")
		}
		fmt.Fprint(writer, ans[i])
	}
	fmt.Fprintln(writer)
}
