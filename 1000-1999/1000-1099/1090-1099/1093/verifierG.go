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

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type query struct {
	t    int
	a1   int
	a2   int
	vals []int
}

type caseG struct {
	n, k   int
	points [][]int
	q      int
	qs     []query
}

func genCase(r *rand.Rand) caseG {
	n := r.Intn(4) + 1
	k := r.Intn(3) + 1
	points := make([][]int, n)
	for i := 0; i < n; i++ {
		p := make([]int, k)
		for j := 0; j < k; j++ {
			p[j] = r.Intn(21) - 10
		}
		points[i] = p
	}
	q := r.Intn(5) + 1
	qs := make([]query, q)
	for i := 0; i < q; i++ {
		if r.Intn(2) == 0 {
			l := r.Intn(n) + 1
			r2 := r.Intn(n-l+1) + l
			qs[i] = query{2, l, r2, nil}
		} else {
			idx := r.Intn(n) + 1
			vals := make([]int, k)
			for j := 0; j < k; j++ {
				vals[j] = r.Intn(21) - 10
			}
			qs[i] = query{1, idx, 0, vals}
		}
	}
	return caseG{n: n, k: k, points: points, q: q, qs: qs}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattan(a, b []int) int {
	d := 0
	for i := range a {
		d += abs(a[i] - b[i])
	}
	return d
}

// bruteExpected simulates the queries and returns expected answers for type-2 queries.
func bruteExpected(tc caseG) []int {
	// deep copy points
	pts := make([][]int, tc.n)
	for i, p := range tc.points {
		cp := make([]int, len(p))
		copy(cp, p)
		pts[i] = cp
	}
	var answers []int
	for _, q := range tc.qs {
		if q.t == 1 {
			cp := make([]int, len(q.vals))
			copy(cp, q.vals)
			pts[q.a1-1] = cp
		} else {
			l, r := q.a1-1, q.a2-1
			best := 0
			for i := l; i <= r; i++ {
				for j := i + 1; j <= r; j++ {
					if d := manhattan(pts[i], pts[j]); d > best {
						best = d
					}
				}
			}
			answers = append(answers, best)
		}
	}
	return answers
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		for _, p := range tc.points {
			for j, v := range p {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(fmt.Sprintf("%d", v))
			}
			sb.WriteByte('\n')
		}
		sb.WriteString(fmt.Sprintf("%d\n", tc.q))
		for _, q := range tc.qs {
			if q.t == 1 {
				sb.WriteString("1 ")
				sb.WriteString(fmt.Sprintf("%d", q.a1))
				for _, v := range q.vals {
					sb.WriteByte(' ')
					sb.WriteString(fmt.Sprintf("%d", v))
				}
				sb.WriteByte('\n')
			} else {
				sb.WriteString(fmt.Sprintf("2 %d %d\n", q.a1, q.a2))
			}
		}
		input := sb.String()

		expected := bruteExpected(tc)

		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}

		gotLines := strings.Fields(strings.TrimSpace(got))
		if len(gotLines) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d answers, got %d\ninput:\n%s",
				i+1, len(expected), len(gotLines), input)
			os.Exit(1)
		}
		for j, line := range gotLines {
			var val int
			if _, err := fmt.Sscan(line, &val); err != nil || val != expected[j] {
				expStrs := make([]string, len(expected))
				for ei, ev := range expected {
					expStrs[ei] = fmt.Sprintf("%d", ev)
				}
				fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected: %s\ngot: %s\n",
					i+1, input, strings.Join(expStrs, " "), got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
