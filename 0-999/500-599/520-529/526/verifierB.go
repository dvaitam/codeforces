package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solveB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	size := 1 << (n + 1)
	a := make([]int, size)
	for i := 2; i < size; i++ {
		fmt.Fscan(in, &a[i])
	}
	dp := make([]int, size)
	var res int64
	for i := size/2 - 1; i >= 1; i-- {
		left, right := 2*i, 2*i+1
		sL := a[left] + dp[left]
		sR := a[right] + dp[right]
		if sL > sR {
			res += int64(sL - sR)
			sR = sL
		} else {
			res += int64(sR - sL)
			sL = sR
		}
		dp[i] = sL
	}
	return fmt.Sprint(res)
}

func genTestB(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	size := 1 << (n + 1)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d\n", n)
	for i := 2; i < size; i++ {
		if i > 2 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", rng.Intn(100)+1)
	}
	buf.WriteByte('\n')
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestB(rng)
		expect := solveB(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
