package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const numTestsD = 100

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: verifierD.go <binary>")
		os.Exit(1)
	}
	binPath, cleanup, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		os.Exit(1)
	}
	if cleanup != nil {
		defer cleanup()
	}
	r := rand.New(rand.NewSource(1))
	for t := 1; t <= numTestsD; t++ {
		k := int64(r.Intn(10) + 1)
		min := k * (k + 1) / 2
		n := min + int64(r.Intn(100))
		input := fmt.Sprintf("%d %d\n", n, k)
		out, err := run(binPath, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if err := validateOutput(n, k, out); err != nil {
			fmt.Printf("test %d failed\ninput:%sreason:%v\ngot:%s\n", t, input, err, strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "verify_binD")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, string(out))
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func parseInts(tokens []string) ([]int64, error) {
	out := make([]int64, len(tokens))
	for i, tk := range tokens {
		v, err := strconv.ParseInt(tk, 10, 64)
		if err != nil {
			return nil, err
		}
		out[i] = v
	}
	return out, nil
}

func hasSolution(n, k int64) bool {
	if k <= 0 {
		return false
	}
	minSum := k * (k + 1) / 2
	if n < minSum {
		return false
	}
	if k == 1 {
		return n >= 1
	}

	var dfs func(pos, prev, rem int64) bool
	dfs = func(pos, prev, rem int64) bool {
		if pos == k {
			return rem == 0
		}
		remainDays := k - pos
		low := prev + 1
		if low < 1 {
			low = 1
		}
		high := prev * 2
		if pos == 0 {
			high = rem
		}
		if high > rem {
			high = rem
		}
		if low > high {
			return false
		}

		for x := low; x <= high; x++ {
			remAfter := rem - x
			// Prune by achievable sum bounds for the remaining days.
			nextLow := x + 1
			minRem := remainDays*nextLow + (remainDays*(remainDays-1))/2

			// For tests here k<=10, this is safe and simple.
			maxRem := int64(0)
			cur := x
			for i := int64(0); i < remainDays; i++ {
				cur *= 2
				maxRem += cur
			}
			if remAfter < minRem || remAfter > maxRem {
				continue
			}
			if dfs(pos+1, x, remAfter) {
				return true
			}
		}
		return false
	}

	return dfs(0, 0, n)
}

func validateOutput(n, k int64, outRaw string) error {
	out := strings.TrimSpace(outRaw)
	if out == "" {
		return fmt.Errorf("empty output")
	}
	toks := strings.Fields(out)
	head := strings.ToUpper(toks[0])
	switch head {
	case "NO":
		if len(toks) != 1 {
			return fmt.Errorf("NO output should not contain extra tokens")
		}
		if hasSolution(n, k) {
			return fmt.Errorf("reported NO but a valid sequence exists")
		}
		return nil
	case "YES":
		if len(toks) != int(k)+1 {
			return fmt.Errorf("expected %d numbers after YES, got %d", k, len(toks)-1)
		}
		a, err := parseInts(toks[1:])
		if err != nil {
			return fmt.Errorf("invalid number in sequence")
		}
		var sum int64
		for i := int64(0); i < k; i++ {
			if a[i] <= 0 {
				return fmt.Errorf("a[%d] must be positive", i+1)
			}
			sum += a[i]
			if i+1 < k {
				if !(a[i] < a[i+1] && a[i+1] <= 2*a[i]) {
					return fmt.Errorf("sequence rule violated at day %d", i+1)
				}
			}
		}
		if sum != n {
			return fmt.Errorf("sum mismatch: got %d need %d", sum, n)
		}
		return nil
	default:
		return fmt.Errorf("first token must be YES or NO")
	}
}
