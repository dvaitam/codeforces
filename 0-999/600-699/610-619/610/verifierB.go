package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func solveB(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	a := make([]int64, n)
	var minV int64 = 1<<63 - 1
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
		if a[i] < minV {
			minV = a[i]
		}
	}
	longest := 0
	cur := 0
	for i := 0; i < 2*n; i++ {
		if a[i%n] > minV {
			cur++
			if cur > longest {
				longest = cur
			}
		} else {
			if cur > longest {
				longest = cur
			}
			cur = 0
		}
	}
	if cur > longest {
		longest = cur
	}
	result := int64(n)*minV + int64(longest)
	return fmt.Sprint(result)
}

func genTestB(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", rng.Int63n(1000)+1)
	}
	buf.WriteByte('\n')
	return buf.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestB(rng)
		expect := solveB(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
