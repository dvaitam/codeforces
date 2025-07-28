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

func runBinary(cmdPath string, input string, args ...string) (string, error) {
	cmdArgs := append([]string{cmdPath}, args...)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
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

func runCandidate(bin, input string) (string, error) {
	if strings.HasSuffix(bin, ".go") {
		return runBinary("go", input, "run", bin)
	}
	return runBinary(bin, input)
}

func runReference(input string) (string, error) {
	return runBinary("go", input, "run", "1712F.go")
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(7) + 3
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 2; i <= n; i++ {
		if i > 2 {
			sb.WriteByte(' ')
		}
		parent := rng.Intn(i-1) + 1
		sb.WriteString(strconv.Itoa(parent))
	}
	sb.WriteString("\n")
	q := rng.Intn(3) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	used := make(map[int]struct{})
	for i := 0; i < q; i++ {
		var x int
		for {
			x = rng.Intn(n) + 1
			if _, ok := used[x]; !ok {
				used[x] = struct{}{}
				break
			}
		}
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(x))
	}
	sb.WriteString("\n")
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp, err := runReference(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed: %v\ninput:\n%s", err, in)
			os.Exit(1)
		}
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
