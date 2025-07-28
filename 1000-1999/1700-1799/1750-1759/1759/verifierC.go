package main

import (
	"bytes"
	"fmt"
	"math"
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func solveCase(l, r, x, a, b int) int {
	if a == b {
		return 0
	} else if abs(a-b) >= x {
		return 1
	} else if abs(b-l) < x && abs(b-r) < x {
		return -1
	} else if abs(a-l) >= x && abs(b-l) >= x {
		return 2
	} else if abs(a-r) >= x && abs(b-r) >= x {
		return 2
	} else if abs(r-l) >= x && ((abs(a-l) >= x && abs(r-b) >= x) || (abs(a-r) >= x && abs(l-b) >= x)) {
		return 3
	}
	return -1
}

func generateCase(rng *rand.Rand) (string, string) {
	l := rng.Intn(21) - 10
	r := l + rng.Intn(21)
	x := rng.Intn(int(math.Max(float64(r-l), 1))) + 1
	a := l + rng.Intn(r-l+1)
	b := l + rng.Intn(r-l+1)
	ans := solveCase(l, r, x, a, b)
	input := fmt.Sprintf("1\n%d %d %d\n%d %d\n", l, r, x, a, b)
	return input, fmt.Sprint(ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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
