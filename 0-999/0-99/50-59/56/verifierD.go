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
	"time"
)

func editDistance(s, t string) int {
	n, m := len(s), len(t)
	dp := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]int, m+1)
	}
	for i := 0; i <= n; i++ {
		dp[i][0] = i
	}
	for j := 0; j <= m; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			cost := 0
			if s[i-1] != t[j-1] {
				cost = 1
			}
			dp[i][j] = dp[i-1][j-1] + cost
			if dp[i-1][j]+1 < dp[i][j] {
				dp[i][j] = dp[i-1][j] + 1
			}
			if dp[i][j-1]+1 < dp[i][j] {
				dp[i][j] = dp[i][j-1] + 1
			}
		}
	}
	return dp[n][m]
}

type op struct {
	typ string
	pos int
	ch  byte
}

func applyOps(s string, ops []op) (string, error) {
	b := []byte(s)
	for _, o := range ops {
		switch o.typ {
		case "INSERT":
			if o.pos < 1 || o.pos > len(b)+1 {
				return "", fmt.Errorf("bad insert pos")
			}
			b = append(b[:o.pos-1], append([]byte{o.ch}, b[o.pos-1:]...)...)
		case "DELETE":
			if o.pos < 1 || o.pos > len(b) {
				return "", fmt.Errorf("bad delete pos")
			}
			b = append(b[:o.pos-1], b[o.pos:]...)
		case "REPLACE":
			if o.pos < 1 || o.pos > len(b) {
				return "", fmt.Errorf("bad replace pos")
			}
			b[o.pos-1] = o.ch
		default:
			return "", fmt.Errorf("unknown op %s", o.typ)
		}
	}
	return string(b), nil
}

func generateCase(rng *rand.Rand) (string, string) {
	l1 := rng.Intn(6) + 1
	l2 := rng.Intn(6) + 1
	letters := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	sb := make([]byte, l1)
	tb := make([]byte, l2)
	for i := 0; i < l1; i++ {
		sb[i] = letters[rng.Intn(len(letters))]
	}
	for i := 0; i < l2; i++ {
		tb[i] = letters[rng.Intn(len(letters))]
	}
	return string(sb), string(tb)
}

func runCase(bin, s, t string) error {
	input := fmt.Sprintf("%s\n%s\n", s, t)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(scanner.Text())
	if err != nil {
		return fmt.Errorf("bad k: %s", scanner.Text())
	}
	ops := make([]op, k)
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("missing operation type")
		}
		typ := scanner.Text()
		if !scanner.Scan() {
			return fmt.Errorf("missing op pos")
		}
		pos, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("bad pos")
		}
		var ch byte
		if typ != "DELETE" {
			if !scanner.Scan() {
				return fmt.Errorf("missing char")
			}
			token := scanner.Text()
			if len(token) != 1 {
				return fmt.Errorf("bad char")
			}
			ch = token[0]
		}
		ops[i] = op{typ, pos, ch}
	}
	// ignore any extra whitespace
	if scanner.Scan() {
		return fmt.Errorf("extra output")
	}
	if k != editDistance(s, t) {
		return fmt.Errorf("expected %d operations, got %d", editDistance(s, t), k)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		s, t := generateCase(rng)
		if err := runCase(bin, s, t); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\n%s\n", i+1, err, s, t)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
