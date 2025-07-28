package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Case struct {
	arr []int64
}

func expected(arr []int64) string {
	var altSum int64
	for i, v := range arr {
		if i%2 == 0 {
			altSum += v
		} else {
			altSum -= v
		}
	}
	if len(arr)%2 == 1 || altSum <= 0 {
		return "YES"
	}
	return "NO"
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(1815))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(19) + 2 // 2..20
		arr := make([]int64, n)
		for j := range arr {
			arr[j] = rng.Int63n(1_000_000_000) + 1
		}
		cases[i] = Case{arr}
	}
	return cases
}

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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		input := fmt.Sprintf("1\n%d\n", len(c.arr))
		for j, v := range c.arr {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		want := expected(c.arr)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, want, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
