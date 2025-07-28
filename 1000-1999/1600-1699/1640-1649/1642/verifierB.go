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

func generateCase(rng *rand.Rand) (string, string) {
	t := rng.Intn(3) + 1
	var in strings.Builder
	in.WriteString(fmt.Sprintf("%d\n", t))
	var out strings.Builder
	for tc := 0; tc < t; tc++ {
		n := rng.Intn(10) + 1
		in.WriteString(fmt.Sprintf("%d\n", n))
		arr := make([]int, n)
		types := make(map[int]struct{})
		for i := 0; i < n; i++ {
			if i > 0 {
				in.WriteByte(' ')
			}
			val := rng.Intn(10) + 1
			arr[i] = val
			in.WriteString(fmt.Sprintf("%d", val))
			types[val] = struct{}{}
		}
		in.WriteByte('\n')
		m := len(types)
		for k := 1; k <= n; k++ {
			if k > 1 {
				out.WriteByte(' ')
			}
			if m > k {
				out.WriteString(fmt.Sprintf("%d", m))
			} else {
				out.WriteString(fmt.Sprintf("%d", k))
			}
		}
		if tc+1 < t {
			out.WriteByte('\n')
		}
	}
	return in.String(), out.String()
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected: %s got: %s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
