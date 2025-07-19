package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	var A, B string
	fmt.Fscan(reader, &A, &B)

	type1 := make([]int, 0)
	type2 := make([]int, 0)
	for i := 0; i < n; i++ {
		if A[i] == 'a' && B[i] == 'b' {
			type1 = append(type1, i+1)
		} else if A[i] == 'b' && B[i] == 'a' {
			type2 = append(type2, i+1)
		}
	}
	s1 := len(type1)
	s2 := len(type2)
	// Pair within type1 and type2
	ops := s1/2 + s2/2
	// If one odd and other even, impossible
	if (s1%2)^(s2%2) == 1 {
		fmt.Fprintln(writer, -1)
		return
	}
	total := ops
	if s1%2 == 1 {
		total += 2
	}
	fmt.Fprintln(writer, total)

	// Output pairs for type1
	for i := 0; i+1 < s1; i += 2 {
		fmt.Fprintln(writer, type1[i], type1[i+1])
	}
	// Output pairs for type2
	for i := 0; i+1 < s2; i += 2 {
		fmt.Fprintln(writer, type2[i], type2[i+1])
	}
	// If both have one left, do two more ops
	if s1%2 == 1 {
		l1 := type1[s1-1]
		l2 := type2[s2-1]
		fmt.Fprintln(writer, l1, l1)
		fmt.Fprintln(writer, l1, l2)
	}
}
