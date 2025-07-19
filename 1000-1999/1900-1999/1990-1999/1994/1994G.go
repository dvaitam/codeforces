package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

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

func main() {
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
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
		if startBit < len(dp0) && !dp0[startBit] {
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
}
