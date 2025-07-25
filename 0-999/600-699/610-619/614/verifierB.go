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

func isBeautiful(s string) bool {
	count1 := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '1' {
			count1++
			if count1 > 1 {
				return false
			}
		} else if s[i] != '0' {
			return false
		}
	}
	return true
}

func expectedOutput(arr []string) string {
	special := "1"
	zeros := 0
	for _, s := range arr {
		if s == "0" {
			return "0\n"
		}
		if isBeautiful(s) {
			if s != "1" {
				zeros += len(s) - 1
			}
		} else {
			special = s
		}
	}
	return special + strings.Repeat("0", zeros) + "\n"
}

func randBeautiful(rng *rand.Rand) string {
	if rng.Intn(2) == 0 {
		return "1"
	}
	zeros := rng.Intn(5) + 1
	return "1" + strings.Repeat("0", zeros)
}

func randNonBeautiful(rng *rand.Rand) string {
	length := rng.Intn(5) + 1
	var sb strings.Builder
	for i := 0; i < length; i++ {
		d := byte(rng.Intn(10) + '0')
		sb.WriteByte(d)
	}
	s := sb.String()
	if isBeautiful(s) {
		// ensure non-beautiful
		s = "2" + s
	}
	return s
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	arr := make([]string, n)
	nonBeautiful := -1
	for i := 0; i < n; i++ {
		if rng.Float64() < 0.9 {
			arr[i] = randBeautiful(rng)
		} else {
			arr[i] = randNonBeautiful(rng)
			nonBeautiful = i
		}
	}
	if nonBeautiful == -1 {
		idx := rng.Intn(n)
		arr[idx] = randNonBeautiful(rng)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, s := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(s)
	}
	sb.WriteByte('\n')
	expect := expectedOutput(arr)
	return sb.String(), expect
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	want := strings.TrimSpace(expected)
	if got != want {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
