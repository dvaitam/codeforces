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

func solve(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n, k int
	if _, err := fmt.Fscan(r, &n, &k); err != nil {
		return ""
	}
	sums := make([]int64, k)
	for i := 0; i < n; i++ {
		var a int64
		fmt.Fscan(r, &a)
		sums[i%k] += a
	}
	minSum := sums[0]
	minR := 0
	for r := 1; r < k; r++ {
		if sums[r] < minSum {
			minSum = sums[r]
			minR = r
		}
	}
	return fmt.Sprintf("%d\n", minR+1)
}

func genTest() (string, string) {
	k := rand.Intn(9) + 1
	q := rand.Intn(9) + 1
	n := k * q
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", rand.Intn(1000))
	}
	sb.WriteByte('\n')
	inp := sb.String()
	out := solve(inp)
	return inp, out
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genTest()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s\n", i+1, err, in, got)
			return
		}
		if got != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
