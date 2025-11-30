package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesG = `3 1 1 1 0 0
1 4 1101 1010
1 3 011 011
2 4 0010 0011 0010
3 4 1100 0011 0110 0000
1 1 1 1
2 4 1101 0011 0100
3 4 1101 1110 0111 0000
2 3 011 111 111
3 4 1101 1111 1011 1110
2 1 1 0 0
1 1 1 1
2 1 0 0 1
3 1 1 1 0 1
1 1 1 1
1 2 10 00
3 1 0 1 1 0
3 1 0 0 0 1
3 4 1000 1111 1101 1110
1 3 111 111
1 1 0 0
3 1 0 1 1 0
1 2 11 00
3 1 1 0 1 1
3 1 0 0 1 0
1 2 01 01
1 2 01 01
1 1 1 0
2 4 1010 0110 0001
1 1 0 0
2 1 0 0 1
3 1 0 0 1 0
3 3 010 001 000 001
3 3 110 001 011 001
2 1 0 0 1
3 3 101 000 101 011
1 3 000 100
3 4 1111 0100 1000 0000
3 4 1110 1010 1111 0111
1 1 1 0
1 4 0001 1000
3 4 1111 0100 0000 0001
2 2 00 00 01
3 4 0010 1011 0100 0111
1 1 1 0
2 3 010 011 100
1 4 1001 1010
2 2 01 11 11
3 3 100 011 001 001
1 2 10 11
3 3 000 010 101 011
3 3 001 001 100 111
3 3 111 010 101 000
1 4 1001 1001
3 3 111 001 111 100
1 4 0010 0100
1 4 1100 1000
1 1 0 1
2 4 1111 0000 0101
2 3 000 101 001
2 4 1111 1101 0111
1 1 1 0
2 4 1101 1100 0110
3 1 1 1 0 0
2 1 1 0 0
2 2 11 11 01
2 3 011 000 000
2 1 1 1 1
2 2 01 00 00
1 3 101 110
2 3 000 010 100
3 1 0 0 0 1
2 2 10 00 10
3 3 000 101 001 101
3 2 00 11 11 01
1 2 10 11
1 3 111 111
2 2 10 10 00
1 2 11 00
3 3 001 101 101 010
2 3 010 001 000
2 1 1 1 0
3 4 1000 0010 1011 1000
1 4 0000 1111
2 2 11 11 10
1 1 1 1
1 4 1000 0001
3 2 01 00 01 01
2 1 0 1 0
1 1 1 1
1 2 11 10
2 2 01 10 01
3 3 010 111 111 111
1 1 1 1
1 2 10 00
1 1 1 0
2 1 0 0 0
1 3 000 000
3 3 100 101 011 011
3 4 0110 0010 0011 0011`

var reader *bufio.Reader
var writer *bufio.Writer

func readToken() string {
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return ""
		}
		if b > ' ' {
			buf := []byte{b}
			for {
				b2, err2 := reader.ReadByte()
				if err2 != nil || b2 <= ' ' {
					if err2 == nil {
						reader.UnreadByte()
					}
					break
				}
				buf = append(buf, b2)
			}
			return string(buf)
		}
	}
}

func readInt() int {
	tok := readToken()
	v, _ := strconv.Atoi(tok)
	return v
}

// readChar skips whitespace and returns next byte
func readChar() byte {
	for {
		b, err := reader.ReadByte()
		if err != nil {
			return 0
		}
		if b > ' ' {
			return b
		}
	}
}

func solve1994G(input string) (string, error) {
	reader = bufio.NewReader(strings.NewReader(input))
	var out bytes.Buffer
	writer = bufio.NewWriter(&out)
	defer writer.Flush()

	t := readInt()
	for ; t > 0; t-- {
		n := readInt()
		K := readInt()
		S := readToken()
		cnt := make([]int, K+1)
		// read n x K matrix
		for i := 0; i < n; i++ {
			for j := 1; j <= K; j++ {
				if readChar() == '1' {
					cnt[K-j+1]++
				}
			}
		}
		size := 2*n + 1
		dp0 := make([]bool, size)
		dp1 := make([]bool, size)
		dp0[0] = true
		opt := make([][]int, K+1)
		for i := 0; i <= K; i++ {
			opt[i] = make([]int, size)
		}
		// DP transitions
		for i := 0; i < K; i++ {
			// clear next
			for j := 0; j < size; j++ {
				dp1[j] = false
			}
			c := cnt[i+1]
			for j := 0; j < size; j++ {
				if dp0[j] {
					v1 := j/2 + n - c
					dp1[v1] = true
					opt[i+1][v1] = j
					v2 := j/2 + c
					dp1[v2] = true
					opt[i+1][v2] = j
				}
			}
			// copy back and prune by parity
			copy(dp0, dp1)
			bit := int(S[K-i-1] - '0')
			for j := 0; j < size; j++ {
				if j&1 != bit {
					dp0[j] = false
				}
			}
		}
		startBit := int(S[0] - '0')
		if startBit >= len(dp0) || !dp0[startBit] {
			fmt.Fprintln(writer, -1)
			continue
		}
		j := startBit
		// backtrack and output
		for i := K; i >= 1; i-- {
			prev := opt[i][j]
			if prev/2+n-cnt[i] == j {
				writer.WriteByte('1')
			} else {
				writer.WriteByte('0')
			}
			j = prev
		}
		writer.WriteByte('\n')
	}
	writer.Flush()
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	idx := 0
	for _, raw := range strings.Split(testcasesG, "\n") {
		line := strings.TrimSpace(raw)
		if line == "" {
			continue
		}
		idx++
		input := "1\n" + strings.Join(strings.Fields(line), " ") + "\n"
		expected, err := solve1994G(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle solve error: %v\n", err)
			os.Exit(1)
		}
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err = cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expected {
			fmt.Printf("test %d failed. Expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
