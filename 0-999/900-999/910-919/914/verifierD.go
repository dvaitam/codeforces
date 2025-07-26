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
	typ int
	a   int
	b   int
	c   int
}

type testCase struct {
	n   int
	arr []int
	qs  []query
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func isAlmostCorrect(arr []int, l, r, x int) bool {
	cnt := 0
	for i := l; i <= r; i++ {
		if arr[i]%x != 0 {
			cnt++
			if cnt > 1 {
				return false
			}
		}
	}
	return true
}

func expected(tc testCase) []string {
	res := []string{}
	arr := append([]int(nil), tc.arr...)
	for _, q := range tc.qs {
		if q.typ == 1 {
			if isAlmostCorrect(arr, q.a, q.b, q.c) {
				res = append(res, "YES")
			} else {
				res = append(res, "NO")
			}
		} else {
			arr[q.a] = q.b
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, testCase) {
	n := rng.Intn(5) + 1
	arr := make([]int, n+1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		arr[i] = rng.Intn(10) + 1
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	q := rng.Intn(5) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	qs := make([]query, q)
	for i := 0; i < q; i++ {
		typ := rng.Intn(2) + 1
		if typ == 1 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			x := rng.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d %d\n", l, r, x))
			qs[i] = query{typ: 1, a: l, b: r, c: x}
		} else {
			idx := rng.Intn(n) + 1
			val := rng.Intn(10) + 1
			sb.WriteString(fmt.Sprintf("2 %d %d\n", idx, val))
			qs[i] = query{typ: 2, a: idx, b: val}
		}
	}
	return sb.String(), testCase{n: n, arr: arr, qs: qs}
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
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
	for i := 1; i <= 100; i++ {
		input, tc := generateCase(rng)
		exp := strings.Join(expected(tc), "\n")
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
