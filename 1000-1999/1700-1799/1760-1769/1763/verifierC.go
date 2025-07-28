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

func absInt64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func expected(n int, a []int64) string {
	if n == 2 {
		return fmt.Sprintf("%d", 2*absInt64(a[0]-a[1]))
	}
	if n == 3 {
		maxVal, minVal := a[0], a[0]
		for i := 1; i < 3; i++ {
			if a[i] > maxVal {
				maxVal = a[i]
			}
			if a[i] < minVal {
				minVal = a[i]
			}
		}
		boundaryMax := a[0]
		if a[2] > boundaryMax {
			boundaryMax = a[2]
		}
		diff := maxVal - minVal
		v := boundaryMax
		if diff > v {
			v = diff
		}
		ans := int64(3) * v
		return fmt.Sprintf("%d", ans)
	}
	maxVal := a[0]
	for i := 1; i < n; i++ {
		if a[i] > maxVal {
			maxVal = a[i]
		}
	}
	ans := int64(n) * maxVal
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(30) + 1)
	}
	var input strings.Builder
	input.WriteString("1\n")
	input.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			input.WriteByte(' ')
		}
		input.WriteString(fmt.Sprintf("%d", v))
	}
	input.WriteString("\n")
	return input.String(), expected(n, arr)
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
