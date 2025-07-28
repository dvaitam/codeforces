package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func solve(reader *bufio.Reader) string {
	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return ""
	}
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	var sb strings.Builder
	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var l, r, x int
			fmt.Fscan(reader, &l, &r, &x)
			l--
			r--
			for i := l; i <= r; i++ {
				a[i] = x
			}
		} else {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			l--
			r--
			res := int64(1<<63 - 1)
			for i := l; i <= r; i++ {
				g := gcd(a[i], b[i])
				val := int64(a[i]) * int64(b[i]) / int64(g*g)
				if val < res {
					res = val
				}
			}
			fmt.Fprintf(&sb, "%d\n", res)
		}
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	q := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, q)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(50)+1)
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", rng.Intn(50)+1)
	}
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		t := rng.Intn(2) + 1
		if t == 1 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			x := rng.Intn(50) + 1
			fmt.Fprintf(&sb, "1 %d %d %d\n", l, r, x)
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			fmt.Fprintf(&sb, "2 %d %d\n", l, r)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solve(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
