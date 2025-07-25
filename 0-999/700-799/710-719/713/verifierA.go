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

// numKey computes the parity pattern key for a decimal number string.
func numKey(s string) int {
	key := 0
	for i, j := 0, len(s)-1; j >= 0; i, j = i+1, j-1 {
		if (s[j]-'0')&1 == 1 {
			key |= 1 << i
		}
	}
	return key
}

// patKey computes the key from the pattern string of '0' and '1'.
func patKey(s string) int {
	key := 0
	for i, j := 0, len(s)-1; j >= 0; i, j = i+1, j-1 {
		if s[j] == '1' {
			key |= 1 << i
		}
	}
	return key
}

type testCase struct {
	in  string
	out string
}

func generateCase(rng *rand.Rand) testCase {
	t := rng.Intn(20) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	counts := make(map[int]int)
	nums := make([]string, 0)
	var out strings.Builder
	for i := 0; i < t; i++ {
		op := rng.Intn(3)
		if op == 2 && len(nums) == 0 {
			op = 0
		}
		switch op {
		case 0: // add
			val := fmt.Sprintf("%d", rng.Int63n(1000000))
			fmt.Fprintf(&sb, "+ %s\n", val)
			counts[numKey(val)]++
			nums = append(nums, val)
		case 1: // query
			l := rng.Intn(18) + 1
			pat := make([]byte, l)
			for j := 0; j < l; j++ {
				if rng.Intn(2) == 1 {
					pat[j] = '1'
				} else {
					pat[j] = '0'
				}
			}
			ps := string(pat)
			fmt.Fprintf(&sb, "? %s\n", ps)
			out.WriteString(fmt.Sprintf("%d\n", counts[patKey(ps)]))
		case 2: // remove
			idx := rng.Intn(len(nums))
			val := nums[idx]
			nums = append(nums[:idx], nums[idx+1:]...)
			fmt.Fprintf(&sb, "- %s\n", val)
			counts[numKey(val)]--
		}
	}
	return testCase{in: sb.String(), out: strings.TrimSpace(out.String())}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.out {
		return fmt.Errorf("expected %q got %q", tc.out, got)
	}
	return nil
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
