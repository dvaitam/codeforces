package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type pair struct{ first, second int }

var operations []pair

func permuteRec(n, offset int) {
	if n == 1 {
		return
	}
	solveRec(n-1, offset)
	if n > 2 {
		operations = append(operations, pair{offset + 1, offset + n - 2})
	}
	permuteRec(n-1, offset+1)
}

func solveRec(n, offset int) {
	permuteRec(n, offset)
	operations = append(operations, pair{offset, offset + n - 1})
}

func solveCase(arr []int) string {
	n := len(arr)
	mask := make([]int, n)
	var maxMask []int
	maxSum := -1
	for mask[0] < 2 {
		sum := 0
		streak := 0
		for i := 0; i < n; i++ {
			if mask[i] == 0 {
				streak++
			} else {
				sum += streak * streak
				sum += arr[i]
				streak = 0
			}
		}
		sum += streak * streak
		if sum > maxSum {
			maxSum = sum
			maxMask = append([]int(nil), mask...)
		}
		idx := n - 1
		mask[idx]++
		for idx > 0 && mask[idx] == 2 {
			mask[idx] = 0
			idx--
			mask[idx]++
		}
	}
	operations = operations[:0]
	var out bytes.Buffer
	fmt.Fprintf(&out, "%d ", maxSum)
	prev := -1
	for i := 0; i < n; i++ {
		if maxMask[i] == 0 && arr[i] != 0 {
			operations = append(operations, pair{i, i})
		}
		if maxMask[i] == 0 && prev == -1 {
			prev = i
		} else if maxMask[i] == 1 {
			if prev != -1 {
				solveRec(i-prev, prev)
			}
			prev = -1
		}
	}
	if prev != -1 {
		solveRec(n-prev, prev)
	}
	fmt.Fprintf(&out, "%d\n", len(operations))
	for _, op := range operations {
		fmt.Fprintf(&out, "%d %d\n", op.first+1, op.second+1)
	}
	return strings.TrimSpace(out.String())
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 1
	arr := make([]int, n)
	var in bytes.Buffer
	fmt.Fprintf(&in, "%d\n", n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(6)
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", arr[i])
	}
	in.WriteByte('\n')
	exp := solveCase(arr)
	return in.String(), exp
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
