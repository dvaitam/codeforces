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

func runCandidate(bin, input string) (string, error) {
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

func solveCase(arr []int) []int {
	for {
		n := len(arr)
		found := false
		var bestX, bestB int
		for x := 1; x*2 <= n && !found; x++ {
			for b := 0; b+2*x <= n; b++ {
				ok := true
				for i := 0; i < x; i++ {
					if arr[b+i] != arr[b+x+i] {
						ok = false
						break
					}
				}
				if ok {
					bestX = x
					bestB = b
					found = true
					break
				}
			}
		}
		if !found {
			break
		}
		arr = arr[bestB+bestX:]
	}
	return arr
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	counts := make(map[int]int)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		for {
			v := rng.Intn(5)
			if counts[v] < 10 {
				arr[i] = v
				counts[v]++
				break
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	res := solveCase(append([]int(nil), arr...))
	var out strings.Builder
	fmt.Fprintf(&out, "%d\n", len(res))
	for i := 0; i < len(res); i++ {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(&out, "%d", res[i])
	}
	out.WriteByte('\n')
	return sb.String(), out.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
