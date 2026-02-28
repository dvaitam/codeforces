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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func solveC(r *bufio.Reader) string {
	var n int
	fmt.Fscan(r, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(r, &a[i])
	}
	if n < 3 {
		return "0\n"
	}

	var ans int64
	stack := make([]int, 0, n)

	for i := 0; i < n; i++ {
		stack = append(stack, a[i])
		for len(stack) >= 3 {
			m := len(stack)
			if stack[m-3] >= stack[m-2] && stack[m-2] <= stack[m-1] {
				mn := stack[m-3]
				if stack[m-1] < mn {
					mn = stack[m-1]
				}
				ans += int64(mn)
				stack[m-2] = stack[m-1]
				stack = stack[:m-1]
			} else {
				break
			}
		}
	}

	sum := int64(0)
	max1, max2 := 0, 0
	for _, val := range stack {
		sum += int64(val)
		if val > max1 {
			max2 = max1
			max1 = val
		} else if val > max2 {
			max2 = val
		}
	}
	ans += sum - int64(max1) - int64(max2)

	return fmt.Sprintf("%d\n", ans)
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 3
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", rng.Intn(1000)+1)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect := solveC(bufio.NewReader(strings.NewReader(tc)))
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %sinput:\n%s", i+1, expect, out, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
