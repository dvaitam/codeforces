package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// applyOps verifies the operations output and checks that the array becomes all zeros.
func applyOps(n int, arr []int64, out string) error {
	reader := bufio.NewReader(strings.NewReader(out))
	for op := 0; op < 3; op++ {
		var l, r int
		if _, err := fmt.Fscan(reader, &l, &r); err != nil {
			return fmt.Errorf("failed to read l r for op %d: %v", op+1, err)
		}
		if l < 1 || r < l || r > n {
			return fmt.Errorf("invalid segment %d: %d %d", op+1, l, r)
		}
		segLen := int64(r - l + 1)
		for i := l - 1; i < r; i++ {
			var bStr string
			if _, err := fmt.Fscan(reader, &bStr); err != nil {
				return fmt.Errorf("failed to read value for op %d: %v", op+1, err)
			}
			b, err := strconv.ParseInt(bStr, 10, 64)
			if err != nil {
				return fmt.Errorf("invalid integer in op %d: %v", op+1, err)
			}
			if b%segLen != 0 {
				return fmt.Errorf("value %d not divisible by %d in op %d", b, segLen, op+1)
			}
			arr[i] += b
		}
	}
	var tmp string
	if _, err := fmt.Fscan(reader, &tmp); err == nil {
		return fmt.Errorf("extra output")
	}
	for i, v := range arr {
		if v != 0 {
			return fmt.Errorf("a[%d]=%d not zero", i+1, v)
		}
	}
	return nil
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
	return out.String(), nil
}

func genTest(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", rng.Intn(21)-10)
	}
	buf.WriteByte('\n')
	return buf.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genTest(rng)
		reader := bufio.NewReader(strings.NewReader(input))
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int64, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &arr[j])
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := applyOps(n, append([]int64(nil), arr...), out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
