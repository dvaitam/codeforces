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

func reverseBits8(b byte) byte {
	var r byte
	for i := 0; i < 8; i++ {
		if (b>>i)&1 == 1 {
			r |= 1 << (7 - i)
		}
	}
	return r
}

func decode(text string) []int {
	res := make([]int, len(text))
	var prevRev byte
	for i := 0; i < len(text); i++ {
		y := reverseBits8(text[i])
		a := int(prevRev) - int(y)
		a = (a%256 + 256) % 256
		res[i] = a
		prevRev = reverseBits8(text[i])
	}
	return res
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randomText(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rng.Intn(95) + 32)
	}
	return string(b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		text := randomText(rng)
		expected := decode(text)
		input := text + "\n"
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		scanner := bufio.NewScanner(strings.NewReader(out))
		var nums []int
		for scanner.Scan() {
			v, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d bad integer: %v\noutput:\n%s", i+1, err, out)
				os.Exit(1)
			}
			nums = append(nums, v)
		}
		if len(nums) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d length mismatch: expected %d got %d\n", i+1, len(expected), len(nums))
			os.Exit(1)
		}
		for j := range expected {
			if nums[j] != expected[j] {
				fmt.Fprintf(os.Stderr, "case %d mismatch at %d: expected %d got %d\n", i+1, j, expected[j], nums[j])
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
