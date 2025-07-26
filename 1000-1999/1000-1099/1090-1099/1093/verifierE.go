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
	ref := "refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1093E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

type query struct {
	t  int
	a1 int
	a2 int
	a3 int
	a4 int
}

type caseE struct {
	n, m int
	a    []int
	b    []int
	q    []query
}

func genCase(r *rand.Rand) caseE {
	n := r.Intn(5) + 2 // 2..6
	m := r.Intn(5) + 1 // 1..5
	a := r.Perm(n)
	for i := range a {
		a[i]++
	}
	b := r.Perm(n)
	for i := range b {
		b[i]++
	}
	qs := make([]query, m)
	for i := 0; i < m; i++ {
		if r.Intn(2) == 0 {
			la := r.Intn(n) + 1
			ra := r.Intn(n-la+1) + la
			lb := r.Intn(n) + 1
			rb := r.Intn(n-lb+1) + lb
			qs[i] = query{1, la, ra, lb, rb}
		} else {
			x := r.Intn(n) + 1
			y := r.Intn(n) + 1
			for y == x {
				y = r.Intn(n) + 1
			}
			qs[i] = query{2, x, y, 0, 0}
		}
	}
	return caseE{n: n, m: m, a: a, b: b, q: qs}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
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
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for j, v := range tc.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for j, v := range tc.b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for _, q := range tc.q {
			if q.t == 1 {
				sb.WriteString(fmt.Sprintf("1 %d %d %d %d\n", q.a1, q.a2, q.a3, q.a4))
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
