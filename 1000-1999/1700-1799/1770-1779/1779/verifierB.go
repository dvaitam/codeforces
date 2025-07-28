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

func solveB(n int) (bool, []int) {
	if n%2 == 1 && n == 3 {
		return false, nil
	}
	arr := make([]int, n)
	if n%2 == 1 {
		p := n/2 - 1
		q := n / 2
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				arr[i] = p
			} else {
				arr[i] = -q
			}
		}
	} else {
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				arr[i] = 1
			} else {
				arr[i] = -1
			}
		}
	}
	return true, arr
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

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 2
	ok, arr := solveB(n)
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d\n", n))
	var expect strings.Builder
	if !ok {
		expect.WriteString("NO")
	} else {
		expect.WriteString("YES\n")
		for i, v := range arr {
			if i > 0 {
				expect.WriteByte(' ')
			}
			expect.WriteString(fmt.Sprintf("%d", v))
		}
	}
	inputStr := input.String()
	expectedStr := expect.String()
	if ok {
		expectedStr = strings.TrimSpace(expectedStr)
	}
	return inputStr, expectedStr
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := genCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
