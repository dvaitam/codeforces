package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func isSubsequence(target, source string) bool {
	j := 0
	for i := 0; i < len(source) && j < len(target); i++ {
		if source[i] == target[j] {
			j++
		}
	}
	return j == len(target)
}

func solveCase(nStr string) int {
	n, _ := strconv.Atoi(nStr)
	limit := int(math.Sqrt(float64(n)))
	best := len(nStr) + 1
	for i := 1; i <= limit; i++ {
		sq := i * i
		sqStr := strconv.Itoa(sq)
		if isSubsequence(sqStr, nStr) {
			ops := len(nStr) - len(sqStr)
			if ops < best {
				best = ops
			}
		}
	}
	if best == len(nStr)+1 {
		return -1
	}
	return best
}

func genCase(rng *rand.Rand) (string, string) {
	length := rng.Intn(9) + 1
	bs := make([]byte, length)
	bs[0] = byte(rng.Intn(9)+1) + '0'
	for i := 1; i < length; i++ {
		bs[i] = byte(rng.Intn(10)) + '0'
	}
	nStr := string(bs)
	input := fmt.Sprintf("%s\n", nStr)
	out := solveCase(nStr)
	expected := fmt.Sprintf("%d\n", out)
	return input, expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
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
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
