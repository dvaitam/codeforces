package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func factorialMap(limit int64) map[int64]int64 {
	facts := make(map[int64]int64)
	facts[0] = 1
	facts[1] = 1
	v := int64(1)
	for i := int64(2); ; i++ {
		v *= i
		if v > limit {
			break
		}
		facts[i] = v
		if i >= 20 {
			break
		}
	}
	return facts
}

func dfs(a []int64, pos int, sum int64, used, k int, S int64, facts map[int64]int64) int64 {
	if sum > S || used > k {
		return 0
	}
	if pos == len(a) {
		if sum == S {
			return 1
		}
		return 0
	}
	res := dfs(a, pos+1, sum, used, k, S, facts)
	res += dfs(a, pos+1, sum+a[pos], used, k, S, facts)
	if used < k {
		if f, ok := facts[a[pos]]; ok {
			res += dfs(a, pos+1, sum+f, used+1, k, S, facts)
		}
	}
	return res
}

func solveCase(n, k int, S int64, a []int64) string {
	facts := factorialMap(S)
	ans := dfs(a, 0, 0, 0, k, S, facts)
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	k := rng.Intn(n + 1)
	a := make([]int64, n)
	sum := int64(0)
	for i := range a {
		a[i] = int64(rng.Intn(10) + 1)
		sum += a[i]
	}
	S := int64(rng.Intn(int(sum)+50) + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, S))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	expect := solveCase(n, k, S, a)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
