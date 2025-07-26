package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveCase(n, m, d int, c []int) (bool, []int) {
	sum := 0
	for _, v := range c {
		sum += v
	}
	rem := n - sum
	if rem > (m+1)*d {
		return false, nil
	}
	gaps := make([]int, m+1)
	for i := 0; i <= m; i++ {
		if rem > d {
			gaps[i] = d
			rem -= d
		} else {
			gaps[i] = rem
			rem = 0
		}
	}
	ans := make([]int, 0, n)
	for i := 0; i <= m; i++ {
		for j := 0; j < gaps[i]; j++ {
			ans = append(ans, 0)
		}
		if i < m {
			for j := 0; j < c[i]; j++ {
				ans = append(ans, i+1)
			}
		}
	}
	return true, ans
}

func generateCase() (string, string) {
	m := rand.Intn(5) + 1
	c := make([]int, m)
	sum := 0
	for i := 0; i < m; i++ {
		c[i] = rand.Intn(3) + 1
		sum += c[i]
	}
	d := rand.Intn(3) + 1
	var n int
	if rand.Intn(2) == 0 {
		n = sum + rand.Intn(d*(m+1)+1)
	} else {
		n = sum + (m+1)*d + rand.Intn(5) + 1
	}
	ok, ans := solveCase(n, m, d, c)
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d %d\n", n, m, d)
	for i := 0; i < m; i++ {
		if i+1 == m {
			fmt.Fprintf(&in, "%d\n", c[i])
		} else {
			fmt.Fprintf(&in, "%d ", c[i])
		}
	}
	if !ok {
		return in.String(), "NO\n"
	}
	var out strings.Builder
	fmt.Fprintln(&out, "YES")
	for i := 0; i < len(ans); i++ {
		if i+1 == len(ans) {
			fmt.Fprintf(&out, "%d\n", ans[i])
		} else {
			fmt.Fprintf(&out, "%d ", ans[i])
		}
	}
	return in.String(), out.String()
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, buf.String())
	}
	return strings.TrimSpace(buf.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(3)
	for i := 0; i < 100; i++ {
		in, exp := generateCase()
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected:\n%s\ngot:\n%s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
