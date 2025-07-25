package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Process struct{ start, end int }

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	cells := make([]int, n)
	allZero := true
	for i := range cells {
		if rng.Intn(2) == 0 { // empty
			cells[i] = 0
		} else {
			cells[i] = rng.Intn(4) + 1
			allZero = false
		}
	}
	if allZero {
		cells[0] = 1
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range cells {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

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
	prefix := 0
	for _, p := range processes {
		length := p.end - p.start + 1
		finalStart := prefix + 1
		finalEnd := finalStart + length - 1
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
	exp := fmt.Sprintf("%d\n", totalMoves)
	return sb.String(), exp
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
