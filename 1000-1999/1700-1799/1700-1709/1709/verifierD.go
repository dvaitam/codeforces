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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(n int, a []int, queries [][5]int) []string {
	res := make([]string, len(queries))
	for idx, q := range queries {
		xs, ys, xf, yf, k := q[0], q[1], q[2], q[3], q[4]
		ys--
		yf--
		if (xs-xf)%k != 0 || (ys-yf)%k != 0 {
			res[idx] = "NO"
			continue
		}
		high := xs + (n-xs)/k*k
		l := ys
		r := yf
		if l > r {
			l, r = r, l
		}
		maxA := 0
		for i := l; i <= r; i++ {
			if a[i] > maxA {
				maxA = a[i]
			}
		}
		if maxA < high {
			res[idx] = "YES"
		} else {
			res[idx] = "NO"
		}
	}
	return res
}

func genCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(50) + 10 // 10..59
	m := rng.Intn(10) + 1  // 1..10 columns
	a := make([]int, m)
	for i := range a {
		a[i] = rng.Intn(n)
	}
	q := rng.Intn(5) + 1 // 1..5 queries
	queries := make([][5]int, q)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		ys := rng.Intn(m) + 1
		yf := rng.Intn(m) + 1
		xs := rng.Intn(n-a[ys-1]) + a[ys-1] + 1
		xf := rng.Intn(n-a[yf-1]) + a[yf-1] + 1
		k := rng.Intn(10) + 1
		queries[i] = [5]int{xs, ys, xf, yf, k}
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", xs, ys, xf, yf, k))
	}
	expect := solveCase(n, a, queries)
	return sb.String(), expect
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		outFields := strings.Fields(out)
		if len(outFields) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, len(expect), len(outFields), in)
			os.Exit(1)
		}
		for j, val := range expect {
			if outFields[j] != val {
				fmt.Fprintf(os.Stderr, "case %d failed on query %d: expected %s got %s\ninput:\n%s", i+1, j+1, val, outFields[j], in)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
