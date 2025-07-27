package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
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

func genCase(rng *rand.Rand) []int64 {
	n := rng.Intn(50) + 1
	arr := make([]int64, n)
	for i := range arr {
		v := rng.Int63n(200) - 100
		if v == 0 {
			v = 1
		}
		arr[i] = v
	}
	return arr
}

func expected(arr []int64) int64 {
	if len(arr) == 0 {
		return 0
	}
	currMax := arr[0]
	ans := int64(0)
	currSign := arr[0] > 0
	for i := 1; i < len(arr); i++ {
		sign := arr[i] > 0
		if sign == currSign {
			if arr[i] > currMax {
				currMax = arr[i]
			}
		} else {
			ans += currMax
			currMax = arr[i]
			currSign = sign
		}
	}
	ans += currMax
	return ans
}

func runCase(bin string, arr []int64) error {
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
	out, err := runProg(bin, sb.String())
	if err != nil {
		return err
	}
	got, err := strconv.ParseInt(strings.TrimSpace(out), 10, 64)
	if err != nil {
		return fmt.Errorf("cannot parse output: %v", err)
	}
	exp := expected(arr)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
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
