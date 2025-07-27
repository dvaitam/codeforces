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

func buildOracle() (string, error) {
	exe := "oracleD"
	cmd := exec.Command("go", "build", "-o", exe, "1418D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle: %v\n%s", err, out)
	}
	return exe, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	coords := make([]int, 0, n)
	used := make(map[int]bool)
	for len(coords) < n {
		x := rng.Intn(20) + 1
		if !used[x] {
			used[x] = true
			coords = append(coords, x)
		}
	}
	sort.Ints(coords)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i, v := range coords {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	cur := make(map[int]bool)
	for _, v := range coords {
		cur[v] = true
	}
	for i := 0; i < q; i++ {
		if len(cur) == 0 || (len(cur) < 20 && rng.Intn(2) == 0) {
			// add
			t := 1
			var x int
			for {
				x = rng.Intn(20) + 1
				if !cur[x] {
					break
				}
			}
			cur[x] = true
			fmt.Fprintf(&sb, "%d %d\n", t, x)
		} else {
			// remove
			t := 0
			idx := rng.Intn(len(cur))
			var x int
			j := 0
			for k := range cur {
				if j == idx {
					x = k
					break
				}
				j++
			}
			delete(cur, x)
			fmt.Fprintf(&sb, "%d %d\n", t, x)
		}
	}
	return sb.String()
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp, err := runProg("./"+oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failure on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		got, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d mismatch\nexpected:%s\ngot:%s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
