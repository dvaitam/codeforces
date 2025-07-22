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

func expected(n, m int, a []int, events []int) string {
	seq := make([]byte, 0)
	for _, ev := range events {
		if ev >= 0 {
			if ev == 0 {
				seq = append(seq, '0')
			} else {
				seq = append(seq, '1')
			}
		} else {
			L := len(seq)
			k := sort.SearchInts(a, L+1)
			for j := 0; j < k; j++ {
				idx := a[j] - 1 - j
				if idx < 0 || idx >= len(seq) {
					continue
				}
				seq = append(seq[:idx], seq[idx+1:]...)
			}
		}
	}
	if len(seq) == 0 {
		return "Poor stack!"
	}
	return string(seq)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(20) + 1
		m := rng.Intn(4) + 1
		a := make([]int, m)
		for i := 0; i < m; i++ {
			a[i] = i + 1 + rng.Intn(3)
		}
		sort.Ints(a)
		events := make([]int, n)
		for i := 0; i < n; i++ {
			events[i] = []int{-1, 0, 1}[rng.Intn(3)]
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for _, ev := range events {
			sb.WriteString(fmt.Sprintf("%d\n", ev))
		}
		input := sb.String()
		exp := expected(n, m, a, events)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", caseNum+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
