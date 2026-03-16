package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func expected(n, a, b int) string {
	rounds := 0
	for a != b {
		a = (a + 1) / 2
		b = (b + 1) / 2
		rounds++
	}
	total := 0
	for tmp := n; tmp > 1; tmp /= 2 {
		total++
	}
	if rounds == total {
		return "Final!"
	}
	return fmt.Sprintf("%d", rounds)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	const testcasesRaw = `100
8 2 5
4 4 2
4 4 1
128 111 1
256 137 118
4 3 1
2 1 2
16 14 1
16 15 16
16 12 8
16 15 10
2 2 1
8 5 2
64 55 25
32 19 32
128 9 123
16 13 14
8 6 2
256 56 84
128 95 126
2 2 1
32 26 11
8 4 1
16 8 13
64 46 59
32 1 25
8 4 7
2 2 1
128 125 92
128 89 1
64 59 4
16 6 3
32 3 5
4 1 4
2 2 1
32 8 12
64 38 9
8 3 5
8 5 8
64 64 61
4 1 3
128 88 108
16 9 4
32 14 28
2 1 2
8 1 3
256 219 113
256 115 16
128 83 110
2 2 1
16 2 10
4 1 3
32 11 27
32 9 1
2 1 2
8 1 7
16 12 4
16 14 7
256 54 200
32 32 2
64 52 37
2 1 2
8 6 7
16 9 4
128 89 125
16 3 2
4 2 3
64 33 48
64 44 15
32 16 32
8 2 6
2 2 1
128 38 33
64 15 49
4 2 1
32 24 19
4 4 3
4 1 3
2 1 2
4 1 2
16 14 6
4 4 2
16 6 4
128 97 76
32 31 21
4 2 3
2 1 2
64 58 51
64 52 9
4 3 4
4 3 2
256 183 133
8 4 5
16 8 12
4 3 1
256 47 174
16 13 10
2 2 1
64 39 32
64 13 12`

	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("empty test file")
		os.Exit(1)
	}
	var t int
	fmt.Sscan(scan.Text(), &t)
	for i := 0; i < t; i++ {
		var n, a, b int
		if !scan.Scan() {
			fmt.Printf("case %d missing n\n", i+1)
			os.Exit(1)
		}
		fmt.Sscan(scan.Text(), &n)
		if !scan.Scan() {
			fmt.Printf("case %d missing a\n", i+1)
			os.Exit(1)
		}
		fmt.Sscan(scan.Text(), &a)
		if !scan.Scan() {
			fmt.Printf("case %d missing b\n", i+1)
			os.Exit(1)
		}
		fmt.Sscan(scan.Text(), &b)
		exp := expected(n, a, b)
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d\n", n, a, b)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got := strings.TrimSpace(string(out))
		if got != exp {
			fmt.Printf("case %d failed: expected %s got %s\n", i+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
