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

func solveB(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(reader, &t)
	var sb strings.Builder
	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(reader, &n)
		sum := 0
		maxv := 0
		for j := 0; j < n; j++ {
			var x int
			fmt.Fscan(reader, &x)
			sum += x
			if x > maxv {
				maxv = x
			}
		}
		if maxv > sum-maxv || sum%2 == 1 {
			sb.WriteString("T")
		} else {
			sb.WriteString("HL")
		}
		if i+1 < t {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func genTest(rng *rand.Rand) string {
	t := rng.Intn(5) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d\n", t)
	for i := 0; i < t; i++ {
		n := rng.Intn(5) + 1
		fmt.Fprintf(&buf, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(&buf, "%d", rng.Intn(10)+1)
		}
		buf.WriteByte('\n')
	}
	return buf.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genTest(rng)
		expect := solveB(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
