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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveCase(a, b, c int) string {
	if a == b {
		return fmt.Sprintf("%d", c)
	}
	if a == c {
		return fmt.Sprintf("%d", b)
	}
	return fmt.Sprintf("%d", a)
}

func generateCase(rng *rand.Rand) (string, string) {
	x := rng.Intn(10)
	y := rng.Intn(10)
	for y == x {
		y = rng.Intn(10)
	}
	pos := rng.Intn(3)
	nums := [3]int{}
	for i := 0; i < 3; i++ {
		if i == pos {
			nums[i] = y
		} else {
			nums[i] = x
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d %d\n", nums[0], nums[1], nums[2])
	input := sb.String()
	expected := solveCase(nums[0], nums[1], nums[2])
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := generateCase(rng)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
