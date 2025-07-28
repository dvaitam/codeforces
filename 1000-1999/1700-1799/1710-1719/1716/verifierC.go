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
	prefTop := make([]int, m+1)
	prefBottom := make([]int, m+1)
	for i := 1; i <= m; i++ {
		prefTop[i] = max(prefTop[i-1]+1, top[i]+1)
		prefBottom[i] = max(prefBottom[i-1]+1, bottom[i]+1)
	}
	sufTop := make([]int, m+2)
	sufBottom := make([]int, m+2)
	for i := m; i >= 1; i-- {
		sufTop[i] = max(sufTop[i+1]+1, top[i]+1)
		sufBottom[i] = max(sufBottom[i+1]+1, bottom[i]+1)
	}
	ans := int(^uint(0) >> 1)
	for i := 1; i <= m; i++ {
		cur1 := max(prefTop[i-1], sufBottom[i+1])
		cur2 := max(prefBottom[i-1], sufTop[i+1])
		cur := cur1
		if cur2 < cur {
			cur = cur2
		}
		if cur < ans {
			ans = cur
		}
	}
	return ans
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
