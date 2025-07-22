package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1 // up to 10 tasks
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	type task struct{ t, c int64 }
	tasks := make([]task, n)
	var curT int64
	for i := 0; i < n; i++ {
		curT += int64(rng.Intn(5) + 1)
		c := int64(rng.Intn(5) + 1)
		tasks[i] = task{curT, c}
		fmt.Fprintf(&sb, "%d %d\n", curT, c)
	}
	// compute expected
	var curTime, queue, lastTime, maxQueue int64
	for _, tt := range tasks {
		dt := tt.t - curTime
		processed := queue
		if processed > dt {
			processed = dt
		}
		if processed > 0 {
			lastTime = curTime + processed
			queue -= processed
		}
		curTime = tt.t
		queue += tt.c
		if queue > maxQueue {
			maxQueue = queue
		}
	}
	if queue > 0 {
		lastTime = curTime + queue
	}
	expected := fmt.Sprintf("%d %d", lastTime, maxQueue)
	return sb.String(), expected
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
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
