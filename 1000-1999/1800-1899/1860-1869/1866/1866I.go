package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	buffer, _ := io.ReadAll(os.Stdin)
	var bpos int

	readInt := func() int {
		for bpos < len(buffer) && buffer[bpos] <= ' ' {
			bpos++
		}
		if bpos >= len(buffer) {
			return 0
		}
		res := 0
		for bpos < len(buffer) && buffer[bpos] > ' ' {
			res = res*10 + int(buffer[bpos]-'0')
			bpos++
		}
		return res
	}

	N := readInt()
	if N == 0 {
		return
	}
	M := readInt()
	K := readInt()

	specials := make([][]int, N+1)
	for i := 0; i < K; i++ {
		r := readInt()
		c := readInt()
		specials[r] = append(specials[r], c)
	}

	in_S := make([]bool, M+1)
	for i := 1; i <= M; i++ {
		in_S[i] = true
	}

	max_S := M

	for r := N; r >= 1; r-- {
		s_1 := 0
		for _, c := range specials[r] {
			if c > s_1 {
				s_1 = c
			}
		}

		for max_S > 0 && !in_S[max_S] {
			max_S--
		}

		c := max_S
		if c > s_1 && c > 0 {
			in_S[c] = false
			if r == 1 && c == 1 {
				fmt.Println("Bhinneka")
				return
			}
		}

		for _, sp := range specials[r] {
			in_S[sp] = false
		}
	}

	fmt.Println("Chaneka")
}
