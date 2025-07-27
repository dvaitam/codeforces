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

func solve(n, l, r int64) []int64 {
	result := make([]int64, 0, r-l+1)
	var pos int64 = 1
	for i := int64(1); i < n && pos <= r; i++ {
		blockLen := 2 * (n - i)
		if l > pos+blockLen-1 {
			pos += blockLen
			continue
		}
		for j := i + 1; j <= n && pos <= r; j++ {
			if pos >= l {
				result = append(result, i)
			}
			pos++
			if pos > r {
				break
			}
			if pos >= l {
				result = append(result, j)
			}
			pos++
		}
	}
	if pos <= r {
		if pos >= l {
			result = append(result, 1)
		}
	}
	return result
}

func expected(n, l, r int64) string {
	res := solve(n, l, r)
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := int64(rng.Intn(5) + 2)
		totalLen := n*(n-1) + 1
		l := int64(rng.Intn(int(totalLen)) + 1)
		maxRange := int(totalLen - l + 1)
		if maxRange > 20 {
			maxRange = 20
		}
		r := l + int64(rng.Intn(maxRange))
		input := fmt.Sprintf("1\n%d %d %d\n", n, l, r)
		exp := expected(n, l, r)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
