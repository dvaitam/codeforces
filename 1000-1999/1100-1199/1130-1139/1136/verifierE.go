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

type query struct {
	typ    byte
	x1, x2 int
}

type test struct {
	n  int
	a  []int
	k  []int
	qs []query
}

func expected(t test) string {
	arr := make([]int, t.n)
	copy(arr, t.a)
	kk := make([]int, t.n-1)
	copy(kk, t.k)
	var out []string
	for _, q := range t.qs {
		if q.typ == '+' {
			idx := q.x1 - 1
			add := q.x2
			arr[idx] += add
			for i := idx; i < t.n-1; i++ {
				need := arr[i] + kk[i]
				if arr[i+1] < need {
					arr[i+1] = need
				} else {
					break
				}
			}
		} else {
			l := q.x1 - 1
			r := q.x2 - 1
			sum := 0
			for i := l; i <= r; i++ {
				sum += arr[i]
			}
			out = append(out, fmt.Sprintf("%d", sum))
		}
	}
	return strings.Join(out, "\n")
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
		return fmt.Errorf("expected %q got %q", want, got)
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

	tests := []test{
		{n: 3, a: []int{1, 2, 3}, k: []int{1, 1}, qs: []query{{typ: '+', x1: 1, x2: 2}, {typ: 's', x1: 1, x2: 3}}},
	}

	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2
		a := make([]int, n)
		k := make([]int, n-1)
		a[0] = rng.Intn(21) - 10
		for j := 0; j < n-1; j++ {
			k[j] = rng.Intn(7) - 3
			base := a[j] + k[j]
			a[j+1] = base + rng.Intn(5)
		}
		qn := rng.Intn(15) + 1
		qs := make([]query, qn)
		for j := 0; j < qn; j++ {
			if rng.Intn(2) == 0 {
				idx := rng.Intn(n) + 1
				x := rng.Intn(4)
				qs[j] = query{typ: '+', x1: idx, x2: x}
			} else {
				l := rng.Intn(n) + 1
				r := rng.Intn(n-l+1) + l
				qs[j] = query{typ: 's', x1: l, x2: r}
			}
		}
		tests = append(tests, test{n: n, a: a, k: k, qs: qs})
	}

	for idx, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for i, v := range tc.k {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d\n", len(tc.qs)))
		for _, q := range tc.qs {
			if q.typ == '+' {
				sb.WriteString(fmt.Sprintf("+ %d %d\n", q.x1, q.x2))
			} else {
				sb.WriteString(fmt.Sprintf("s %d %d\n", q.x1, q.x2))
			}
		}
		want := expected(tc)
		if err := runCase(bin, sb.String(), want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
