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

type query struct {
	typ int
	l   int
	r   int
	x   int
}

type testCase struct {
	n, m, p int
	arr     []int
	qs      []query
	input   string
}

func ceilMul(a, b int) int {
	return (a*b + 99) / 100
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	p0 := 20 + rng.Intn(81) // 20..100
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(5) + 1
	}
	var qs []query
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, p0)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		typ := rng.Intn(2) + 1
		l := rng.Intn(n) + 1
		r := rng.Intn(n-l+1) + l
		if typ == 1 {
			x := rng.Intn(5) + 1
			qs = append(qs, query{1, l, r, x})
			fmt.Fprintf(&sb, "1 %d %d %d\n", l, r, x)
			for j := l - 1; j < r; j++ {
				arr[j] = x
			}
		} else {
			qs = append(qs, query{2, l, r, 0})
			fmt.Fprintf(&sb, "2 %d %d\n", l, r)
		}
	}
	return testCase{n, m, p0, arr, qs, sb.String()}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func validate(tc testCase, out string) error {
	outputs := strings.Split(strings.TrimSpace(out), "\n")
	idx := 0
	arr := append([]int(nil), tc.arr...)
	limit := 100 / tc.p
	for _, q := range tc.qs {
		if q.typ == 1 {
			for j := q.l - 1; j < q.r; j++ {
				arr[j] = q.x
			}
		} else {
			if idx >= len(outputs) {
				return fmt.Errorf("missing output line")
			}
			line := strings.TrimSpace(outputs[idx])
			idx++
			fields := strings.Fields(line)
			if len(fields) == 0 {
				return fmt.Errorf("empty output line")
			}
			cnt, err := strconv.Atoi(fields[0])
			if err != nil {
				return err
			}
			if cnt != len(fields)-1 {
				return fmt.Errorf("count mismatch")
			}
			if cnt > limit {
				return fmt.Errorf("too many ids")
			}
			ids := make([]int, cnt)
			for i := 0; i < cnt; i++ {
				v, err := strconv.Atoi(fields[1+i])
				if err != nil {
					return err
				}
				ids[i] = v
			}
			// compute heavy advertisers
			freq := make(map[int]int)
			for j := q.l - 1; j < q.r; j++ {
				freq[arr[j]]++
			}
			need := ceilMul(tc.p, q.r-q.l+1)
			for id, f := range freq {
				if f >= need {
					found := false
					for _, v := range ids {
						if v == id {
							found = true
							break
						}
					}
					if !found {
						return fmt.Errorf("id %d missing", id)
					}
				}
			}
		}
	}
	if idx != len(outputs) {
		return fmt.Errorf("extra output lines")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		c := genCase(rng)
		out, err := run(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		if err := validate(c, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, c.input, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
