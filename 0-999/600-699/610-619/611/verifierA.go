package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runCandidate(bin string, input string) (string, error) {
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

func expected(x int, typ string) int {
	if typ == "week" {
		if x == 5 || x == 6 {
			return 53
		}
		return 52
	}
	days := []int{31, 29, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	cnt := 0
	for _, d := range days {
		if x <= d {
			cnt++
		}
	}
	return cnt
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		week := rng.Intn(2) == 0
		var sb strings.Builder
		var x int
		if week {
			x = rng.Intn(7) + 1
			fmt.Fprintf(&sb, "%d of week\n", x)
			want := expected(x, "week")
			out, err := runCandidate(bin, sb.String())
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
				os.Exit(1)
			}
			got, err := strconv.Atoi(strings.TrimSpace(out))
			if err != nil || got != want {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, want, out)
				os.Exit(1)
			}
		} else {
			x = rng.Intn(31) + 1
			fmt.Fprintf(&sb, "%d of month\n", x)
			want := expected(x, "month")
			out, err := runCandidate(bin, sb.String())
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
				os.Exit(1)
			}
			got, err := strconv.Atoi(strings.TrimSpace(out))
			if err != nil || got != want {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", i+1, want, out)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
