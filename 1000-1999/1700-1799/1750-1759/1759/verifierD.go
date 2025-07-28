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

func solveCase(n, m int64) int64 {
	k := int64(1)
	cnt2, cnt5 := 0, 0
	tmp := n
	for tmp%2 == 0 {
		cnt2++
		tmp /= 2
	}
	for tmp%5 == 0 {
		cnt5++
		tmp /= 5
	}
	for k*2 <= m && cnt2 < cnt5 {
		k *= 2
		cnt2++
	}
	for k*5 <= m && cnt5 < cnt2 {
		k *= 5
		cnt5++
	}
	for k*10 <= m {
		k *= 10
	}
	k *= m / k
	return n * k
}

func generateCase(rng *rand.Rand) (string, string) {
	n := int64(rng.Intn(1000) + 1)
	m := int64(rng.Intn(1000) + 1)
	ans := solveCase(n, m)
	input := fmt.Sprintf("1\n%d %d\n", n, m)
	return input, fmt.Sprint(ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
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
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
