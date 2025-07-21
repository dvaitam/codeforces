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

func computeMinimal(s string) string {
	if s == "0" {
		return "0"
	}
	digits := []byte(s)
	for i := 0; i < len(digits); i++ {
		for j := i + 1; j < len(digits); j++ {
			if digits[j] < digits[i] {
				digits[i], digits[j] = digits[j], digits[i]
			}
		}
	}
	// count leading zeros
	z := 0
	for z < len(digits) && digits[z] == '0' {
		z++
	}
	if z == 0 {
		return string(digits)
	}
	first := digits[z]
	res := []byte{first}
	for i := 0; i < z; i++ {
		res = append(res, '0')
	}
	for i := z + 1; i < len(digits); i++ {
		res = append(res, digits[i])
	}
	return string(res)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(1_000_000_000 + 1)
		nStr := strconv.Itoa(n)
		minimal := computeMinimal(nStr)
		var mStr string
		if rng.Intn(2) == 0 {
			mStr = minimal
		} else {
			for {
				m := rng.Intn(1_000_000_000 + 1)
				mStr = strconv.Itoa(m)
				if mStr != minimal {
					break
				}
			}
		}
		input := fmt.Sprintf("%s\n%s\n", nStr, mStr)
		expected := "WRONG_ANSWER"
		if mStr == minimal {
			expected = "OK"
		}
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
