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

func toBase(n, base int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		d := n % base
		digits = append(digits, byte('0'+d))
		n /= base
	}
	for l, r := 0, len(digits)-1; l < r; l, r = l+1, r-1 {
		digits[l], digits[r] = digits[r], digits[l]
	}
	return string(digits)
}

func solveH(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var k int
	if _, err := fmt.Fscan(r, &k); err != nil {
		return ""
	}
	var sb strings.Builder
	for i := 1; i < k; i++ {
		for j := 1; j < k; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(toBase(i*j, k))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateCaseH(rng *rand.Rand) string {
	k := rng.Intn(9) + 2
	return fmt.Sprintf("%d\n", k)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseH(rng)
	}
	for i, tc := range cases {
		expect := solveH(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
