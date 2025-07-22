package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ori struct{ top, bottom, front, back, right, left int }

func orientations() []ori {
	var res []ori
	for t := 1; t <= 6; t++ {
		bt := 7 - t
		for f := 1; f <= 6; f++ {
			if f == t || f == bt {
				continue
			}
			bk := 7 - f
			rem := make([]int, 0, 2)
			for v := 1; v <= 6; v++ {
				if v != t && v != bt && v != f && v != bk {
					rem = append(rem, v)
				}
			}
			res = append(res, ori{t, bt, f, bk, rem[0], rem[1]})
			res = append(res, ori{t, bt, f, bk, rem[1], rem[0]})
		}
	}
	return res
}

func solve(n int, x int, a, b []int) string {
	oris := orientations()
	m := len(oris)
	dp := make([]int, m)
	for j, o := range oris {
		if o.top == x && o.front == a[0] && o.right == b[0] {
			dp[j] = 1
		}
	}
	for i := 1; i < n; i++ {
		next := make([]int, m)
		for j, o := range oris {
			if o.front != a[i] || o.right != b[i] {
				continue
			}
			sum := 0
			for k, p := range oris {
				if dp[k] == 0 {
					continue
				}
				if p.bottom != o.top {
					sum += dp[k]
					if sum >= 2 {
						sum = 2
						break
					}
				}
			}
			next[j] = sum
		}
		dp = next
	}
	total := 0
	for _, v := range dp {
		total += v
		if total >= 2 {
			break
		}
	}
	if total == 1 {
		return "YES"
	}
	return "NO"
}

func runBinary(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(1))
	oris := orientations()
	for t := 1; t <= 100; t++ {
		n := r.Intn(10) + 1
		seq := make([]ori, n)
		seq[0] = oris[r.Intn(len(oris))]
		for i := 1; i < n; i++ {
			options := make([]ori, 0, len(oris))
			for _, o := range oris {
				if o.top != seq[i-1].bottom {
					options = append(options, o)
				}
			}
			seq[i] = options[r.Intn(len(options))]
		}
		x := seq[0].top
		a := make([]int, n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = seq[i].front
			b[i] = seq[i].right
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n%d\n", n, x)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", a[i], b[i])
		}
		input := sb.String()
		expected := solve(n, x, a, b)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\nInput:\n%s", t, err, input)
			return
		}
		if strings.TrimSpace(out) != expected {
			fmt.Printf("Test %d FAILED\nInput:\n%sExpected: %s\nGot: %s\n", t, input, expected, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
