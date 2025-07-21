package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	n       int
	queries []string
	answers []string
}

func generateTest() (string, string) {
	n := rand.Intn(20) + 1
	q := rand.Intn(20) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d %d\n", n, q)
	occupied := make([]bool, n+1)
	empPos := make(map[int]int)
	answers := []string{}
	existingIDs := []int{}
	nextID := 1
	ensureQuery := false
	for i := 0; i < q; i++ {
		typ := rand.Intn(3)
		if !ensureQuery && i == q-1 {
			typ = 0 // ensure at least one query
		}
		switch typ {
		case 0: // query
			l := rand.Intn(n) + 1
			r := rand.Intn(n) + 1
			if l > r {
				l, r = r, l
			}
			count := 0
			for p := l; p <= r; p++ {
				if occupied[p] {
					count++
				}
			}
			fmt.Fprintf(&buf, "0 %d %d\n", l, r)
			answers = append(answers, fmt.Sprintf("%d", count))
			ensureQuery = true
		case 1: // arrival
			id := nextID
			nextID++
			// find longest free segment
			bestLen := -1
			bestStart := -1
			bestEnd := -1
			start := -1
			for p := 1; p <= n; p++ {
				if !occupied[p] {
					if start == -1 {
						start = p
					}
				} else {
					if start != -1 {
						length := p - start
						if length > bestLen || (length == bestLen && start > bestStart) {
							bestLen = length
							bestStart = start
							bestEnd = p - 1
						}
						start = -1
					}
				}
			}
			if start != -1 {
				length := n - start + 1
				if length > bestLen || (length == bestLen && start > bestStart) {
					bestLen = length
					bestStart = start
					bestEnd = n
				}
			}
			length := bestEnd - bestStart + 1
			pos := bestStart + (length-1)/2
			if length%2 == 0 {
				pos = bestStart + length/2
			}
			occupied[pos] = true
			empPos[id] = pos
			existingIDs = append(existingIDs, id)
			fmt.Fprintf(&buf, "%d\n", id)
		default: // departure
			if len(existingIDs) == 0 {
				id := nextID
				nextID++
				fmt.Fprintf(&buf, "%d\n", id)
				continue
			}
			idx := rand.Intn(len(existingIDs))
			id := existingIDs[idx]
			existingIDs = append(existingIDs[:idx], existingIDs[idx+1:]...)
			pos := empPos[id]
			occupied[pos] = false
			delete(empPos, id)
			fmt.Fprintf(&buf, "%d\n", id)
		}
	}
	return buf.String(), strings.Join(answers, "\n")
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	const tests = 100
	for t := 1; t <= tests; t++ {
		inp, want := generateTest()
		cmd := exec.Command(binary)
		cmd.Stdin = bytes.NewBufferString(inp)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		want = strings.TrimSpace(want)
		if got != want {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected:\n%s\nGot:\n%s\n", t, inp, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
