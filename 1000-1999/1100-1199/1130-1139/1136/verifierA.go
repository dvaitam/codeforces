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

type chapter struct{ l, r int }

func expected(n int, seg []chapter, k int) string {
	cnt := 0
	for _, c := range seg {
		if c.r >= k {
			cnt++
		}
	}
	return fmt.Sprintf("%d", cnt)
}

func runCase(bin, input, want string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	type test struct {
		n   int
		seg []chapter
		k   int
	}
	tests := []test{
		{n: 3, seg: []chapter{{1, 3}, {4, 7}, {8, 11}}, k: 2},
		{n: 3, seg: []chapter{{1, 3}, {4, 7}, {8, 11}}, k: 5},
		{n: 1, seg: []chapter{{1, 5}}, k: 3},
	}
	for i := 0; i < 100; i++ {
		n := rng.Intn(100) + 1
		seg := make([]chapter, n)
		l := 1
		for j := 0; j < n; j++ {
			len := rng.Intn(100) + 1
			seg[j] = chapter{l, l + len - 1}
			l += len
		}
		k := rng.Intn(l-1) + 1
		tests = append(tests, test{n: n, seg: seg, k: k})
	}
	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, c := range tc.seg {
			sb.WriteString(fmt.Sprintf("%d %d", c.l, c.r))
			if i+1 < len(tc.seg) {
				sb.WriteByte('\n')
			}
		}
		sb.WriteString(fmt.Sprintf("\n%d\n", tc.k))
		want := expected(tc.n, tc.seg, tc.k)
		if err := runCase(bin, sb.String(), want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
