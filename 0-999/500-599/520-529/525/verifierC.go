package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
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

func solveCase(sticks []int) string {
	sort.Slice(sticks, func(i, j int) bool { return sticks[i] > sticks[j] })
	sides := make([]int, 0, len(sticks)/2)
	for i := 0; i < len(sticks)-1; {
		if sticks[i]-sticks[i+1] <= 1 {
			sides = append(sides, sticks[i+1])
			i += 2
		} else {
			i++
		}
	}
	var area int64
	for i := 0; i+1 < len(sides); i += 2 {
		area += int64(sides[i]) * int64(sides[i+1])
	}
	return fmt.Sprintf("%d", area)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(30) + 1
	sticks := make([]int, n)
	for i := range sticks {
		sticks[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i, v := range sticks {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	expect := solveCase(append([]int(nil), sticks...))
	return sb.String(), expect
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
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
