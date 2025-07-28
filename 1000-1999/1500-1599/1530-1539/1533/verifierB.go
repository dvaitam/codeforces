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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(arr []int64) string {
	for i := 0; i < len(arr)-1; i++ {
		if (arr[i+1]-arr[i])%2 == 0 {
			return "YES"
		}
	}
	return "NO"
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 2
	arr := make([]int64, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(1000) + 1)
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	expect := solveCase(arr)
	return sb.String(), expect
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" {
		bin = os.Args[2]
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
