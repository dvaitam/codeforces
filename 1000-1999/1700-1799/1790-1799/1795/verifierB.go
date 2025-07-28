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

func expected(n, k int, segs [][2]int) string {
	left, right := false, false
	for _, s := range segs {
		if s[0] == k {
			left = true
		}
		if s[1] == k {
			right = true
		}
	}
	if left && right {
		return "YES"
	}
	return "NO"
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
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
	for tc := 0; tc < 100; tc++ {
		n := rng.Intn(10) + 1
		k := rng.Intn(50) + 1
		segs := make([][2]int, n)
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			l := rng.Intn(50) + 1
			r := rng.Intn(50) + 1
			if l > r {
				l, r = r, l
			}
			segs[i] = [2]int{l, r}
			sb.WriteString(fmt.Sprintf("%d %d\n", l, r))
		}
		input := sb.String()
		expectedOut := expected(n, k, segs)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", tc+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expectedOut {
			fmt.Printf("case %d failed: expected %s got %s\ninput:\n%s", tc+1, expectedOut, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
