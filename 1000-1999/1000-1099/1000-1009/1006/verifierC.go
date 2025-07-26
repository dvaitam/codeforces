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

func solve(arr []int64) int64 {
	n := len(arr)
	if n == 0 {
		return 0
	}
	suml := arr[0]
	sumr := arr[n-1]
	ans := int64(0)
	l, r := 0, n-1
	for l < r {
		if suml < sumr {
			l++
			suml += arr[l]
		} else if sumr < suml {
			r--
			sumr += arr[r]
		} else {
			ans = suml
			l++
			suml += arr[l]
			r--
			sumr += arr[r]
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(50) + 1
		arr := make([]int64, n)
		input := fmt.Sprintf("%d\n", n)
		for j := 0; j < n; j++ {
			arr[j] = int64(rng.Intn(1000) + 1)
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", arr[j])
		}
		input += "\n"
		expected := fmt.Sprintf("%d", solve(arr))
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
