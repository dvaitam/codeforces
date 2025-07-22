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

func solveD(n, k int) string {
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = 1
	}
	a[n] = 0
	mx := 1
	var out strings.Builder
	for step := 0; step < k; step++ {
		for i := 1; i <= n; i++ {
			if a[i] == n-i {
				fmt.Fprintf(&out, "%d ", n)
			} else if a[i]+mx >= n-i {
				fmt.Fprintf(&out, "%d ", i+a[i])
				a[i] = n - i
			} else {
				fmt.Fprintf(&out, "1 ")
				a[i] += mx
			}
		}
		out.WriteByte('\n')
		mx *= 2
	}
	return out.String()
}

func genCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(8) + 1
	k := rng.Intn(5) + 1
	input := fmt.Sprintf("%d %d\n", n, k)
	return input, solveD(n, k)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCaseD(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
