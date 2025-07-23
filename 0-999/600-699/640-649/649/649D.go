package main

import (
	"bufio"
	"fmt"
	"os"
)

// Process represents a contiguous block of memory belonging to a process.
type Process struct {
	start int // 1-based index of the first cell
	end   int // 1-based index of the last cell
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	cells := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &cells[i])
	}

	// Extract processes as ordered intervals of equal non-zero values.
	processes := make([]Process, 0)
	for i := 0; i < n; {
		if cells[i] == 0 {
			i++
			continue
		}
		id := cells[i]
		j := i
		for j+1 < n && cells[j+1] == id {
			j++
		}
		processes = append(processes, Process{start: i + 1, end: j + 1})
		i = j + 1
	}

	totalMoves := 0
	prefix := 0 // length of all previous processes

	for _, p := range processes {
		length := p.end - p.start + 1
		finalStart := prefix + 1
		finalEnd := finalStart + length - 1

		// Calculate overlap of the initial interval with its final position.
		overlapStart := p.start
		if finalStart > overlapStart {
			overlapStart = finalStart
		}
		overlapEnd := p.end
		if finalEnd < overlapEnd {
			overlapEnd = finalEnd
		}
		overlap := 0
		if overlapStart <= overlapEnd {
			overlap = overlapEnd - overlapStart + 1
		}

		totalMoves += length - overlap
		prefix += length
	}

	fmt.Fprintln(out, totalMoves)
}
