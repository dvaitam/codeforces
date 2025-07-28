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

func refSolveD(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var sb strings.Builder
	for ; t > 0; t-- {
		var n int
		var v int64
		fmt.Fscan(in, &n, &v)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		var q int
		fmt.Fscan(in, &q)
		for ; q > 0; q-- {
			var typ int
			fmt.Fscan(in, &typ)
			if typ == 1 {
				var idx int
				var x int64
				fmt.Fscan(in, &idx, &x)
				b[idx-1] = x
			} else if typ == 2 {
				var l, r int
				fmt.Fscan(in, &l, &r)
				l--
				r--
				const INF int64 = 1<<63 - 1
				ans := INF
				for L := l; L <= r; L++ {
					orVal := int64(0)
					maxA := int64(0)
					for R := L; R <= r; R++ {
						orVal |= b[R]
						if a[R] > maxA {
							maxA = a[R]
						}
						if orVal >= v {
							if maxA < ans {
								ans = maxA
							}
							break
						}
					}
				}
				if ans == INF {
					sb.WriteString("-1\n")
				} else {
					sb.WriteString(fmt.Sprintln(ans))
				}
			}
		}
	}
	return strings.TrimSpace(sb.String())
}

func genCaseD(rng *rand.Rand) string {
	t := 1
	n := rng.Intn(3) + 1
	v := rng.Int63n(8) + 1
	a := make([]int64, n)
	b := make([]int64, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Int63n(10)
		b[i] = rng.Int63n(10)
	}
	q := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n%d %d\n", t, n, v))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(b[i]))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		typ := rng.Intn(2) + 1
		if typ == 1 {
			idx := rng.Intn(n) + 1
			x := rng.Int63n(10)
			sb.WriteString(fmt.Sprintf("1 %d %d\n", idx, x))
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("2 %d %d\n", l, r))
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCaseD(rng)
		expect := refSolveD(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
