package main

import (
	"bufio"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func makePal(x int) int64 {
	res := int64(x)
	y := x
	for y > 0 {
		res = res*10 + int64(y%10)
		y /= 10
	}
	return res
}

func solve(reader *bufio.Reader) string {
	var k int
	var p int64
	if _, err := fmt.Fscan(reader, &k, &p); err != nil {
		return ""
	}
	var sum int64
	for i := 1; i <= k; i++ {
		sum = (sum + makePal(i)) % p
	}
	return fmt.Sprint(sum % p)
}

func generateCase(rng *rand.Rand) (string, string) {
	k := rng.Intn(100000) + 1
	p := rng.Int63n(1_000_000_000) + 1
	input := fmt.Sprintf("%d %d\n", k, p)
	expect := solve(bufio.NewReader(strings.NewReader(input)))
	return input, expect
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
