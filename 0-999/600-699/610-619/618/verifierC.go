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

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
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

func genCase(rng *rand.Rand, n int) string {
	points := make([][2]int, n)
	// ensure first three points are non-collinear
	points[0] = [2]int{0, 0}
	points[1] = [2]int{1, 0}
	points[2] = [2]int{0, 1}
	used := map[[2]int]bool{{0, 0}: true, {1, 0}: true, {0, 1}: true}
	for i := 3; i < n; i++ {
		for {
			x := rng.Intn(2001) - 1000
			y := rng.Intn(2001) - 1000
			pt := [2]int{x, y}
			if !used[pt] {
				points[i] = pt
				used[pt] = true
				break
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", points[i][0], points[i][1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 120; t++ {
		n := rng.Intn(20) + 3
		input := genCase(rng, n)
		expected, err := run("618C.go", input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal error on test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
