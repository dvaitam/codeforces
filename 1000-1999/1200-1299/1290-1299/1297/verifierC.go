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

func runProg(bin, input string) (string, error) {
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
		return out.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(a []int) (int, string) {
	n := len(a)
	sumPos := 0
	minAbs := 10001
	idx := -1
	for i, v := range a {
		if v > 0 {
			sumPos += v
		}
		if v != 0 {
			x := v
			if x < 0 {
				x = -x
			}
			if x < minAbs {
				minAbs = x
				idx = i
			}
		}
	}
	ans := make([]byte, n)
	for i, v := range a {
		if v > 0 {
			ans[i] = '1'
		} else {
			ans[i] = '0'
		}
	}
	if a[idx] > 0 {
		ans[idx] = '0'
		sumPos -= a[idx]
	} else {
		ans[idx] = '1'
		sumPos += a[idx]
	}
	return sumPos, string(ans)
}

func genCase(rng *rand.Rand) []int {
	n := rng.Intn(7) + 2
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(21) - 10
	}
	// ensure not all zero
	allZero := true
	for _, v := range arr {
		if v != 0 {
			allZero = false
			break
		}
	}
	if allZero {
		arr[0] = 1
	}
	return arr
}

func runCase(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sum, pattern := solve(arr)
	want := fmt.Sprintf("%d\n%s", sum, pattern)
	out, err := runProg(bin, sb.String())
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(want) {
		return fmt.Errorf("expected:\n%s\ngot:\n%s\ninput:\n%s", want, out, sb.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		arr := genCase(rng)
		if err := runCase(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
