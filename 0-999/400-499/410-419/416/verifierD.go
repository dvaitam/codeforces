package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscan(in, &n)
	a := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	ans := 0
	i := 0
	for i < n {
		ans++
		firstIdx := -1
		var firstVal, d int64
		haveD := false
		j := i
		for ; j < n; j++ {
			v := a[j]
			if v != -1 {
				if firstIdx < 0 {
					firstIdx = j
					firstVal = v
				} else if !haveD {
					dist := int64(j - firstIdx)
					delta := v - firstVal
					if delta%dist != 0 {
						break
					}
					d = delta / dist
					haveD = true
				} else {
					exp := firstVal + int64(j-firstIdx)*d
					if exp != v {
						break
					}
				}
			}
		}
		if j == i {
			j++
		}
		i = j
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		val := int64(-1)
		if rng.Intn(4) != 0 {
			val = int64(rng.Intn(20) + 1)
		}
		sb.WriteString(fmt.Sprintf("%d", val))
	}
	sb.WriteByte('\n')
	input := sb.String()
	exp := solve(input)
	return input, strings.TrimSpace(exp)
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
