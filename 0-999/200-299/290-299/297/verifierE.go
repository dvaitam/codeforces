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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
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

func checkValid(a []int) bool {
	for rot := 0; rot < 6; rot++ {
		ok1 := true
		for i := 0; i < 3; i++ {
			if a[(rot+2*i)%6] != a[(rot+2*i+1)%6] {
				ok1 = false
				break
			}
		}
		if ok1 {
			return true
		}
		ok2 := true
		for i := 0; i < 3; i++ {
			if a[(rot+i)%6] != a[(rot+i+3)%6] {
				ok2 = false
				break
			}
		}
		if ok2 {
			return true
		}
	}
	return false
}

func countWays(n int, l, r []int) int64 {
	pos := make([]int, 2*n)
	for i := 0; i < n; i++ {
		a := l[i] - 1
		b := r[i] - 1
		pos[a] = i + 1
		pos[b] = -(i + 1)
	}
	var res int64
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				seq := make([]int, 0, 6)
				for t := 0; t < 2*n; t++ {
					x := pos[t]
					if x == i+1 || x == -(i+1) || x == j+1 || x == -(j+1) || x == k+1 || x == -(k+1) {
						if x < 0 {
							x = -x
						}
						if x == i+1 {
							seq = append(seq, 1)
						} else if x == j+1 {
							seq = append(seq, 2)
						} else {
							seq = append(seq, 3)
						}
					}
				}
				if checkValid(seq) {
					res++
				}
			}
		}
	}
	return res
}

func randomMatching(rng *rand.Rand, n int) ([]int, []int) {
	arr := rng.Perm(2 * n)
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		a := arr[2*i] + 1
		b := arr[2*i+1] + 1
		if a < b {
			l[i] = a
			r[i] = b
		} else {
			l[i] = b
			r[i] = a
		}
	}
	return l, r
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		n := rng.Intn(5) + 3
		l, r := randomMatching(rng, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", l[j], r[j]))
		}
		input := sb.String()
		expected := countWays(n, l, r)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		var got int64
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output\ninput:\n%soutput:\n%s", i+1, input, out)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
