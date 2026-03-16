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

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// Embedded correct solver for 1914B
func solve1914B(n, k int) string {
	var parts []string
	first := true
	_ = first
	var sb strings.Builder
	for i := n - k; i >= 1; i-- {
		if sb.Len() > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(i))
	}
	for i := n - k + 1; i <= n; i++ {
		if sb.Len() > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(i))
	}
	_ = parts
	return sb.String()
}

func genTest() ([]byte, int, int) {
	n := rand.Intn(49) + 2 // 2..50
	k := rand.Intn(n)
	return []byte(fmt.Sprintf("1\n%d %d\n", n, k)), n, k
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierB /path/to/binary")
		return
	}
	cand := os.Args[1]

	for i := 1; i <= 100; i++ {
		in, n, k := genTest()
		expected := solve1914B(n, k)
		got, err := run(cand, in)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Printf("wrong answer on test %d (n=%d, k=%d):\nexpected: %s\ngot: %s\n", i, n, k, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
