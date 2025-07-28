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

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solve(a []int) int64 {
	freq := make(map[int]int)
	for _, v := range a {
		freq[v]++
	}
	var ans int64
	for _, c := range freq {
		ans += int64(c*(c-1)) / 2
	}
	if c1, ok1 := freq[1]; ok1 {
		if c2, ok2 := freq[2]; ok2 {
			ans += int64(c1 * c2)
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(30) + 1
		a := make([]int, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Intn(6) + 1
			a[j] = val
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := fmt.Sprintf("%d\n", solve(a))
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%sq\n got:%q\n", i+1, input, strings.TrimSpace(expected), got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
