package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runCmd(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = io.Discard
	err := cmd.Run()
	return out.String(), err
}

func solve(x, n int64) int64 {
	r := n % 4
	if x%2 == 0 {
		switch r {
		case 0:
			return x
		case 1:
			return x - n
		case 2:
			return x + 1
		default:
			return x + n + 1
		}
	} else {
		switch r {
		case 0:
			return x
		case 1:
			return x + n
		case 2:
			return x - 1
		default:
			return x - n - 1
		}
	}
}

func randCase() (int64, int64) {
	x := rand.Int63n(2_000_000_001) - 1_000_000_000
	n := rand.Int63n(1_000_000_000_000_000 + 1)
	return x, n
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	refBin := "./refB.bin"
	if err := exec.Command("go", "build", "-o", refBin, "1607B.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(2)
	type tc struct{ x, n int64 }
	cases := []tc{
		{0, 0},
		{0, 1},
		{-1_000_000_000_000_00, 1_000_000_000_000_00},
	}
	for len(cases) < 100 {
		x, n := randCase()
		cases = append(cases, tc{x, n})
	}

	for i, c := range cases {
		input := fmt.Sprintf("1\n%d %d\n", c.x, c.n)
		exp, err := runCmd(refBin, input)
		if err != nil {
			fmt.Println("reference solution failed:", err)
			os.Exit(1)
		}
		got, err := runCmd(candidate, input)
		if err != nil {
			fmt.Printf("test %d: candidate runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		exp = strings.TrimSpace(exp)
		got = strings.TrimSpace(got)
		if exp != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
