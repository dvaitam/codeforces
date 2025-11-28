package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Case struct{ input string }

func buildRef() (string, error) {
	ref := "refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1658E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	abspath, err := filepath.Abs(ref)
	if err != nil {
		return "", fmt.Errorf("failed to resolve reference path: %v", err)
	}
	return abspath, nil
}

func runExe(bin, input string) (string, error) {
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

func genCases() []Case {
	rng := rand.New(rand.NewSource(6))
	cases := make([]Case, 0, 105)
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 3
		k := rng.Intn(n-2) + 1
		nums := make([]int, n*n)
		for j := range nums {
			nums[j] = j + 1
		}
		rng.Shuffle(len(nums), func(a, b int) { nums[a], nums[b] = nums[b], nums[a] })
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", n, k)
		idx := 0
		for r := 0; r < n; r++ {
			for c := 0; c < n; c++ {
				if c > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(nums[idx]))
				idx++
			}
			sb.WriteByte('\n')
		}
		cases = append(cases, Case{sb.String()})
	}
	cases = append(cases, Case{"3 1\n1 2 3\n4 5 6\n7 8 9\n"})
	cases = append(cases, Case{"4 2\n1 2 3 4\n5 6 7 8\n9 10 11 12\n13 14 15 16\n"})
	cases = append(cases, Case{"3 1\n9 1 2\n3 4 5\n6 7 8\n"})
	cases = append(cases, Case{"5 2\n1 2 3 4 5\n6 7 8 9 10\n11 12 13 14 15\n16 17 18 19 20\n21 22 23 24 25\n"})
	cases = append(cases, Case{"3 1\n1 3 5\n7 9 2\n4 6 8\n"})
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	cases := genCases()
	for i, c := range cases {
		exp, err := runExe(ref, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if exp != got {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, c.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
