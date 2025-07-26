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

func isCrisp(s []byte, A [][]int) bool {
	for i := 0; i < len(s)-1; i++ {
		if A[int(s[i]-'a')][int(s[i+1]-'a')] == 0 {
			return false
		}
	}
	return true
}

func removeAll(s string, ch byte) string {
	var b []byte
	for i := 0; i < len(s); i++ {
		if s[i] != ch {
			b = append(b, s[i])
		}
	}
	return string(b)
}

func dfs(s string, A [][]int, memo map[string]int) int {
	if v, ok := memo[s]; ok {
		return v
	}
	best := len(s)
	used := make(map[byte]bool)
	for i := 0; i < len(s); i++ {
		used[s[i]] = true
	}
	for ch := range used {
		t := removeAll(s, ch)
		if isCrisp([]byte(t), A) {
			l := dfs(t, A, memo)
			if l < best {
				best = l
			}
		}
	}
	memo[s] = best
	return best
}

func expectedF(n, p int, str string, A [][]int) int {
	memo := make(map[string]int)
	memo[str] = len(str)
	return dfs(str, A, memo)
}

func generateCase(rng *rand.Rand) (string, int) {
	p := rng.Intn(3) + 2
	n := rng.Intn(6) + 1
	A := make([][]int, p)
	for i := 0; i < p; i++ {
		A[i] = make([]int, p)
		for j := 0; j <= i; j++ {
			val := rng.Intn(2)
			A[i][j] = val
			A[j][i] = val
		}
		A[i][i] = 1
	}
	letters := make([]byte, p)
	for i := 0; i < p; i++ {
		letters[i] = byte('a' + i)
	}
	var str string
	for {
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = letters[rng.Intn(p)]
		}
		if isCrisp(b, A) {
			str = string(b)
			break
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n%s\n", n, p, str)
	for i := 0; i < p; i++ {
		for j := 0; j < p; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", A[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), expectedF(n, p, str, A)
}

func runCase(bin, input string, exp int) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("failed to parse output: %v", err)
	}
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
