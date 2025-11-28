package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

// Envelope represents the dimensions and original index of an envelope
type Envelope struct {
	w, h, id int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, w0, h0 int
	// If input reading fails (empty file), just exit
	if _, err := fmt.Fscan(reader, &n, &w0, &h0); err != nil {
		return
	}

	envelopes := make([]Envelope, 0, n)
	for i := 0; i < n; i++ {
		var w, h int
		fmt.Fscan(reader, &w, &h)
		// Filter out envelopes that cannot even hold the card
		if w > w0 && h > h0 {
			envelopes = append(envelopes, Envelope{w: w, h: h, id: i + 1})
		}
	}

	if len(envelopes) == 0 {
		fmt.Fprintln(writer, 0)
		return
	}

	// Sort envelopes:
	// Primary key: width ascending
	// Secondary key: height ascending
	// This helps in processing potential chains in order.
	sort.Slice(envelopes, func(i, j int) bool {
		if envelopes[i].w != envelopes[j].w {
			return envelopes[i].w < envelopes[j].w
		}
		return envelopes[i].h < envelopes[j].h
	})

	// DP array: dp[i] stores the max length of a chain ending at envelope i
	dp := make([]int, len(envelopes))
	// prev array: stores the index of the previous envelope in the chain for reconstruction
	prev := make([]int, len(envelopes))

	maxLen := 0
	endIdx := -1

	for i := 0; i < len(envelopes); i++ {
		dp[i] = 1
		prev[i] = -1
		for j := 0; j < i; j++ {
			// Check if envelope j fits strictly inside envelope i
			if envelopes[i].w > envelopes[j].w && envelopes[i].h > envelopes[j].h {
				if dp[j]+1 > dp[i] {
					dp[i] = dp[j] + 1
					prev[i] = j
				}
			}
		}
		if dp[i] > maxLen {
			maxLen = dp[i]
			endIdx = i
		}
	}

	if maxLen == 0 {
		fmt.Fprintln(writer, 0)
	} else {
		fmt.Fprintln(writer, maxLen)
		// Reconstruct path
		path := make([]int, 0, maxLen)
		curr := endIdx
		for curr != -1 {
			path = append(path, envelopes[curr].id)
			curr = prev[curr]
		}
		// The path is constructed backwards (largest to smallest), so reverse it
		// Wait, the problem asks to "start with the number of the smallest envelope".
		// Our path slice has [largest_id, ..., smallest_id].
		// Printing in reverse order of slice (i.e., from end to start) gives smallest to largest.
		for i := len(path) - 1; i >= 0; i-- {
			if i < len(path)-1 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, path[i])
		}
		fmt.Fprintln(writer)
	}
}
