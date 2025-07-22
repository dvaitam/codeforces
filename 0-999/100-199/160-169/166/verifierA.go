package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type team struct{ problems, time int }

func expected(n, k int, teams []team) int {
	sort.Slice(teams, func(i, j int) bool {
		if teams[i].problems != teams[j].problems {
			return teams[i].problems > teams[j].problems
		}
		return teams[i].time < teams[j].time
	})
	p := teams[k-1].problems
	t := teams[k-1].time
	cnt := 0
	for _, tm := range teams {
		if tm.problems == p && tm.time == t {
			cnt++
		}
	}
	return cnt
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(50) + 1
	k := rng.Intn(n) + 1
	teams := make([]team, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		teams[i] = team{rng.Intn(50) + 1, rng.Intn(50) + 1}
		fmt.Fprintf(&sb, "%d %d\n", teams[i].problems, teams[i].time)
	}
	return sb.String(), expected(n, k, teams)
}

func runCase(exe, input string, expectedAns int) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expectedAns {
		return fmt.Errorf("expected %d got %d", expectedAns, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
