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

func days(a, b, h int64) int64 {
	if h <= a {
		return 1
	}
	return 1 + (h-a+(a-b-1))/(a-b)
}

type query struct {
	typ int
	a   int64
	b   int64
	n   int64
}

type testCase struct {
	queries []query
}

func expected(tc testCase) string {
	low, high := int64(1), int64(1e18)
	var res []string
	for _, q := range tc.queries {
		if q.typ == 1 {
			var L, R int64
			if q.n == 1 {
				L, R = 1, q.a
			} else {
				L = (q.n-1)*q.a - (q.n-2)*q.b + 1
				R = q.n*q.a - (q.n-1)*q.b
			}
			if L > high || R < low {
				res = append(res, "0")
			} else {
				if L > low {
					low = L
				}
				if R < high {
					high = R
				}
				res = append(res, "1")
			}
		} else {
			d1 := days(q.a, q.b, low)
			d2 := days(q.a, q.b, high)
			if d1 == d2 {
				res = append(res, fmt.Sprint(d1))
			} else {
				res = append(res, "-1")
			}
		}
	}
	return strings.Join(res, " ")
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprint(len(tc.queries)))
	sb.WriteByte('\n')
	for i, q := range tc.queries {
		sb.WriteString(fmt.Sprint(q.typ))
		if q.typ == 1 {
			sb.WriteString(fmt.Sprintf(" %d %d %d", q.a, q.b, q.n))
		} else {
			sb.WriteString(fmt.Sprintf(" %d %d", q.a, q.b))
		}
		if i+1 == len(tc.queries) {
			sb.WriteByte('\n')
		} else {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func runCase(exe string, tc testCase) error {
	input := buildInput(tc)
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected(tc))
	if got != exp {
		return fmt.Errorf("expected %q got %q\ninput:\n%s", exp, got, input)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		q := rng.Intn(10) + 1
		queries := make([]query, q)
		for j := 0; j < q; j++ {
			typ := rng.Intn(2) + 1
			if typ == 1 {
				a := int64(rng.Intn(5) + 2)
				b := int64(rng.Intn(int(a-1)) + 1)
				n := int64(rng.Intn(5) + 1)
				queries[j] = query{typ: 1, a: a, b: b, n: n}
			} else {
				a := int64(rng.Intn(5) + 2)
				b := int64(rng.Intn(int(a-1)) + 1)
				queries[j] = query{typ: 2, a: a, b: b}
			}
		}
		tc := testCase{queries: queries}
		if err := runCase(exe, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
