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
)

func runCandidate(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(n int, strengths [3]int, criminals []int) int {
	s := []int{strengths[0], strengths[1], strengths[2]}
	sort.Ints(s)
	a, b, c := s[0], s[1], s[2]
	ms := append([]int(nil), criminals...)
	sort.Ints(ms)

	removeLE := func(limit int) bool {
		idx := sort.Search(len(ms), func(i int) bool { return ms[i] > limit }) - 1
		if idx < 0 {
			return false
		}
		ms = append(ms[:idx], ms[idx+1:]...)
		return true
	}

	hours := 0
	for len(ms) > 0 {
		mx := ms[len(ms)-1]
		ms = ms[:len(ms)-1]
		hours++
		if mx > a+b+c {
			return -1
		}
		if mx <= c {
			if removeLE(b) {
				removeLE(a)
			} else {
				removeLE(a + b)
			}
		} else if mx <= a+b {
			removeLE(c)
		} else if mx <= a+c {
			removeLE(b)
		} else if mx <= b+c {
			removeLE(a)
		}
	}
	return hours
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(5))
	for t := 0; t < 100; t++ {
		n := rng.Intn(20) + 1
		strengths := [3]int{rng.Intn(50) + 1, rng.Intn(50) + 1, rng.Intn(50) + 1}
		criminals := make([]int, n)
		for i := 0; i < n; i++ {
			criminals[i] = rng.Intn(50) + 1
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		fmt.Fprintf(&sb, "%d %d %d\n", strengths[0], strengths[1], strengths[2])
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", criminals[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		want := expected(n, strengths, append([]int(nil), criminals...))
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(out))
		if err != nil || got != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", t+1, want, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
