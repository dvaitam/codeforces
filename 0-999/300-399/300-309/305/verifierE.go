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

func solveE(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var s string
	fmt.Fscan(reader, &s)
	n := len(s)
	rad := make([]int, n)
	for i := 0; i < n; i++ {
		l, r := i, i
		for l-1 >= 0 && r+1 < n && s[l-1] == s[r+1] {
			l--
			r++
		}
		rad[i] = (r - l) / 2
	}
	dp := make([][]uint16, n)
	for i := range dp {
		dp[i] = make([]uint16, n)
	}
	const maxG = 512
	seen := make([]int, maxG)
	iter := 1
	for length := 1; length <= n; length++ {
		for l := 0; l+length-1 < n; l++ {
			r := l + length - 1
			var g uint16 = 0
			if length >= 3 {
				for k := l + 1; k <= r-1; k++ {
					if rad[k] > 0 {
						x := dp[l][k-1] ^ dp[k+1][r]
						if int(x) < maxG {
							seen[x] = iter
						}
					}
				}
				for j := 0; j < maxG; j++ {
					if seen[j] != iter {
						g = uint16(j)
						break
					}
				}
				iter++
			}
			dp[l][r] = g
		}
	}
	full := dp[0][n-1]
	if full == 0 {
		return "Second"
	}
	var buf bytes.Buffer
	buf.WriteString("First\n")
	for i := 1; i+1 < n; i++ {
		if rad[i] > 0 {
			if dp[0][i-1]^dp[i+1][n-1] == 0 {
				fmt.Fprintf(&buf, "%d", i+1)
				break
			}
		}
	}
	return strings.TrimSpace(buf.String())
}

func genTestE() string {
	n := rand.Intn(20) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rand.Intn(26))
	}
	return string(b) + "\n"
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 1; i <= 100; i++ {
		in := genTestE()
		expected := solveE(in)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(in)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\noutput: %s\n", i, err, string(out))
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
