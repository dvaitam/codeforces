package main

import (
	"bufio"
	"fmt"
	"os"
)

const (
	offset64 = 14695981039346656037
	prime64  = 1099511628211
)

type state struct {
	king int
	wins int
	hash uint64
}

func hashQueue(q []int, head, size, cap int) uint64 {
	h := uint64(offset64)
	for i := 0; i < size; i++ {
		v := uint64(q[(head+i)%cap] + 1)
		h ^= v
		h *= prime64
	}
	h ^= uint64(size)
	h *= prime64
	return h
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	A := make([]int64, n)
	B := make([]int64, n)
	C := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &A[i], &B[i], &C[i])
	}

	if n == 1 {
		fmt.Fprintln(writer, 0, 0)
		return
	}

	q := make([]int, n-1)
	for i := 1; i < n; i++ {
		q[i-1] = i
	}
	head := 0
	size := n - 1
	capQ := n - 1

	king := 0
	wins := 0
	fights := 0

	visited := make(map[state]bool)

	for {
		h := hashQueue(q, head, size, capQ)
		st := state{king, wins, h}
		if visited[st] {
			fmt.Fprintln(writer, -1, -1)
			return
		}
		visited[st] = true

		var ks int64
		if fights == 0 && king == 0 && wins == 0 {
			ks = A[king]
		} else {
			if wins == 0 {
				ks = A[king]
			} else if wins == 1 {
				ks = B[king]
			} else {
				ks = C[king]
			}
		}

		challenger := q[head]
		head = (head + 1) % capQ
		size--
		cs := A[challenger]

		if ks > cs {
			wins++
			q[(head+size)%capQ] = challenger
			size++
			fights++
			if wins == 3 {
				fmt.Fprintln(writer, king, fights)
				return
			}
		} else {
			q[(head+size)%capQ] = king
			size++
			king = challenger
			wins = 1
			fights++
		}
	}
}
