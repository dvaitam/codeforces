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

func solveCase(x int, doors [4]int) string {
	first := doors[x]
	if first == 0 {
		return "NO"
	}
	second := doors[first]
	if second == 0 {
		return "NO"
	}
	return "YES"
}

func genCase(rng *rand.Rand) (string, string) {
	perm := []int{1, 2, 3}
	rng.Shuffle(3, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
	x := perm[0]
	vals := []int{perm[1], perm[2], 0}
	rng.Shuffle(3, func(i, j int) { vals[i], vals[j] = vals[j], vals[i] })
	doors := [4]int{0, vals[0], vals[1], vals[2]}
	input := fmt.Sprintf("1\n%d\n%d %d %d\n", x, doors[1], doors[2], doors[3])
	expect := solveCase(x, doors)
	return input, expect
}

func main() {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
