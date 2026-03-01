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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func solveCase(m int, top, bottom []int) int {
	a := [2][]int{make([]int, m), make([]int, m)}
	for i := 0; i < m; i++ {
		a[0][i] = top[i+1]
		a[1][i] = bottom[i+1]
	}

	// Match the accepted solution exactly.
	a[0][0] = -1
	p := [2][]int{make([]int, m+1), make([]int, m+1)}
	for i := m - 1; i >= 0; i-- {
		for row := 0; row < 2; row++ {
			p[row][i] = max(
				max(a[1-row][i]+1, a[row][i]+(m-i)*2),
				p[row][i+1]+1,
			)
		}
	}

	ans := int(^uint(0) >> 1)
	n := 0
	for i := 0; i < m; i++ {
		c := i & 1
		ans = min(ans, max(n, p[c][i]))
		n = max(
			n,
			max(
				a[c][i]+(m-i)*2,
				a[1-c][i]+(m-i-1)*2+1,
			),
		)
	}
	return ans
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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

func verifyCase(bin string, m int, top, bottom []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 1; i <= m; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(top[i]))
	}
	sb.WriteByte('\n')
	for i := 1; i <= m; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(bottom[i]))
	}
	sb.WriteByte('\n')
	expected := fmt.Sprint(solveCase(m, top, bottom))
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != expected {
		return fmt.Errorf("expected %s got %s", expected, out)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		m := rng.Intn(50) + 2
		top := make([]int, m+1)
		bottom := make([]int, m+1)
		for j := 1; j <= m; j++ {
			top[j] = rng.Intn(1000000000)
		}
		for j := 1; j <= m; j++ {
			bottom[j] = rng.Intn(1000000000)
		}
		if err := verifyCase(bin, m, top, bottom); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
