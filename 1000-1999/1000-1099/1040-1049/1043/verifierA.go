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

func expected(a []int) string {
	n := len(a)
	sum := 0
	maxa := 0
	for _, v := range a {
		sum += v
		if v > maxa {
			maxa = v
		}
	}
	k := (2*sum)/n + 1
	if k < maxa {
		k = maxa
	}
	return fmt.Sprintf("%d", k)
}

func runCase(bin, input, want string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if strings.TrimSpace(want) != got {
		return fmt.Errorf("expected %s got %s", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := [][]int{{1}, {5, 1, 1, 1, 5}}
	for i := 0; i < 100; i++ {
		n := rng.Intn(100) + 1
		arr := make([]int, n)
		for j := range arr {
			arr[j] = rng.Intn(100) + 1
		}
		tests = append(tests, arr)
	}

	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", len(tc)))
		for j, v := range tc {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		want := expected(tc)
		if err := runCase(bin, sb.String(), want); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
