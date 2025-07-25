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

func solveA(a, p []int64) string {
	minPrice := int64(1 << 60)
	var total int64
	for i := range a {
		if p[i] < minPrice {
			minPrice = p[i]
		}
		total += minPrice * a[i]
	}
	return fmt.Sprintf("%d\n", total)
}

func genCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	a := make([]int64, n)
	pr := make([]int64, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		a[i] = int64(rng.Intn(100) + 1)
		pr[i] = int64(rng.Intn(100) + 1)
		fmt.Fprintf(&sb, "%d %d\n", a[i], pr[i])
	}
	return sb.String(), solveA(a, pr)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCaseA(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
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
