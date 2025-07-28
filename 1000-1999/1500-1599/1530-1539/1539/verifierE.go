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
	k  int
	xl int
	xr int
	yl int
	yr int
}

type Test struct {
	n  int
	m  int
	qs []query
}

func (tc Test) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, q := range tc.qs {
		sb.WriteString(fmt.Sprintf("%d\n", q.k))
		sb.WriteString(fmt.Sprintf("%d %d\n", q.xl, q.xr))
		sb.WriteString(fmt.Sprintf("%d %d\n", q.yl, q.yr))
	}
	return sb.String()
}

type state struct {
	l int
	r int
}

func expected(tc Test) string {
	dp := map[state][]int{{0, 0}: {}}
	for _, q := range tc.qs {
		next := make(map[state][]int)
		for st, path := range dp {
			// replace left
			nl := q.k
			nr := st.r
			if nl >= q.xl && nl <= q.xr && nr >= q.yl && nr <= q.yr {
				if _, ok := next[state{nl, nr}]; !ok {
					np := append(append([]int(nil), path...), 0)
					next[state{nl, nr}] = np
				}
			}
			// replace right
			nl = st.l
			nr = q.k
			if nl >= q.xl && nl <= q.xr && nr >= q.yl && nr <= q.yr {
				if _, ok := next[state{nl, nr}]; !ok {
					np := append(append([]int(nil), path...), 1)
					next[state{nl, nr}] = np
				}
			}
		}
		if len(next) == 0 {
			return "No"
		}
		dp = next
	}
	for _, path := range dp {
		var sb strings.Builder
		sb.WriteString("Yes\n")
		for i, v := range path {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		return sb.String()
	}
	return "No"
}

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func genQuery(rng *rand.Rand, m int) query {
	k := rng.Intn(m + 1)
	xl := rng.Intn(m + 1)
	xr := xl + rng.Intn(m-xl+1)
	yl := rng.Intn(m + 1)
	yr := yl + rng.Intn(m-yl+1)
	return query{k, xl, xr, yl, yr}
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(5) + 1
	m := rng.Intn(20) + 1
	qs := make([]query, n)
	for i := 0; i < n; i++ {
		qs[i] = genQuery(rng, m)
	}
	return Test{n: n, m: m, qs: qs}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		tc := genTest(rng)
		expect := expected(tc)
		got, err := runProg(bin, tc.Input())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, tc.Input(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
