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

func buildRef() (string, error) {
	ref := "refG.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1093G.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

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
		exp, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\n got:%s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
