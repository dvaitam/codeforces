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
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveB(r *bufio.Reader) string {
	var t int
	if _, err := fmt.Fscan(r, &t); err != nil {
		return ""
	}
	var out strings.Builder
	for ; t > 0; t-- {
		var n int
		var s string
		fmt.Fscan(r, &n, &s)
		var total0, total1 int64
		var cur0, cur1 int64
		var max0, max1 int64
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				total0++
				cur0++
				if cur0 > max0 {
					max0 = cur0
				}
				cur1 = 0
			} else {
				total1++
				cur1++
				if cur1 > max1 {
					max1 = cur1
				}
				cur0 = 0
			}
		}
		prod := total0 * total1
		sq0 := max0 * max0
		sq1 := max1 * max1
		ans := prod
		if sq0 > ans {
			ans = sq0
		}
		if sq1 > ans {
			ans = sq1
		}
		out.WriteString(fmt.Sprintf("%d\n", ans))
	}
	return strings.TrimSpace(out.String())
}

func generateCaseB(rng *rand.Rand) string {
	n := rng.Intn(40) + 1
	var s strings.Builder
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			s.WriteByte('0')
		} else {
			s.WriteByte('1')
		}
	}
	return fmt.Sprintf("1\n%d\n%s\n", n, s.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		expect := solveB(bufio.NewReader(strings.NewReader(tc)))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
