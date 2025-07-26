package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func solveB(n int, w int64, cups []int64) float64 {
	sort.Slice(cups, func(i, j int) bool { return cups[i] < cups[j] })
	a1 := float64(cups[0])
	mid := float64(cups[n])
	res := 0.0
	if a1*2 <= mid {
		res = a1 * float64(n) * 3
	} else {
		res = mid * 1.5 * float64(n)
	}
	mw := float64(w)
	if res > mw {
		res = mw
	}
	return res
}

func genCase() (string, float64) {
	n := rand.Intn(5) + 1
	w := rand.Int63n(1000) + 1
	cups := make([]int64, 2*n)
	for i := range cups {
		cups[i] = rand.Int63n(1000) + 1
	}
	res := solveB(n, w, append([]int64(nil), cups...))
	input := fmt.Sprintf("%d %d\n", n, w)
	for i, v := range cups {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", v)
	}
	input += "\n"
	return input, res
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(1)
	for i := 0; i < 100; i++ {
		input, expected := genCase()
		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		var got float64
		if _, err := fmt.Sscan(gotStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: failed to parse output\n", i+1)
			os.Exit(1)
		}
		if math.Abs(got-expected) > 1e-6*math.Max(1.0, math.Abs(expected)) {
			fmt.Fprintf(os.Stderr, "test %d failed\ninput:\n%sexpected: %.7f\ngot: %s\n", i+1, input, expected, gotStr)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
