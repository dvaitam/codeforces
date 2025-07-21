package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type domino struct {
	x, h int
	idx  int
}

func expectedAnswer(dom []domino) []int {
	n := len(dom)
	sort.Slice(dom, func(i, j int) bool { return dom[i].x < dom[j].x })
	far := make([]int, n)
	R := make([]int, n)
	stack := make([]int, 0, n)
	for i := n - 1; i >= 0; i-- {
		far[i] = dom[i].x + dom[i].h - 1
		R[i] = i
		for len(stack) > 0 {
			j := stack[len(stack)-1]
			if dom[j].x > far[i] {
				break
			}
			stack = stack[:len(stack)-1]
			if far[j] > far[i] {
				far[i] = far[j]
			}
			if R[j] > R[i] {
				R[i] = R[j]
			}
		}
		stack = append(stack, i)
	}
	ans := make([]int, n)
	for i := 0; i < n; i++ {
		ans[dom[i].idx] = R[i] - i + 1
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(20) + 1
	dom := make([]domino, n)
	x := 0
	for i := 0; i < n; i++ {
		x += rng.Intn(5) + 1
		dom[i] = domino{x: x, h: rng.Intn(10) + 2, idx: i}
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&b, "%d %d\n", dom[i].x, dom[i].h)
	}
	ans := expectedAnswer(append([]domino(nil), dom...))
	return b.String(), ans
}

func runCase(bin, input string, expected []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	parts := strings.Fields(strings.TrimSpace(out.String()))
	if len(parts) != len(expected) {
		return fmt.Errorf("expected %d numbers got %d", len(expected), len(parts))
	}
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil || v != expected[i] {
			return fmt.Errorf("mismatch at %d: expected %d got %s", i+1, expected[i], p)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := generateCase(rng)
		if err := runCase(bin, input, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
