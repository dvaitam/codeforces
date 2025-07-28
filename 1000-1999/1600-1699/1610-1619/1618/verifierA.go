package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func generateCase(rng *rand.Rand) (string, string) {
	a := []int{rng.Intn(1000) + 1, rng.Intn(1000) + 1, rng.Intn(1000) + 1}
	sums := []int{a[0], a[1], a[2], a[0] + a[1], a[0] + a[2], a[1] + a[2], a[0] + a[1] + a[2]}
	sort.Ints(sums)
	var sb strings.Builder
	sb.WriteString("1\n")
	for i, v := range sums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := fmt.Sprintf("%d %d %d\n", a[0], a[1], a[2])
	return input, expected
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, strings.TrimSpace(expect), strings.TrimSpace(out), input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
