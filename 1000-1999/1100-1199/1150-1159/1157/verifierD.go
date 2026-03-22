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

// hasSolution uses the same greedy algorithm as the accepted CF solution
// to determine if a valid sequence exists for given n and k.
func hasSolution(n, k int64) bool {
	if k <= 0 {
		return false
	}
	if n < k*(k+1)/2 {
		return false
	}

	var prev int64 = 0
	currentN := n

	for i := int64(1); i <= k; i++ {
		remDays := k - i

		sumMinRem := remDays * (remDays + 1) / 2
		numerator := currentN - sumMinRem
		if numerator < 0 {
			return false
		}
		limit := numerator / (remDays + 1)

		var minNeeded int64
		if remDays >= 31 {
			minNeeded = 1
		} else {
			denom := (int64(1) << (remDays + 1)) - 1
			minNeeded = (currentN + denom - 1) / denom
		}

		low := prev + 1
		if i == 1 {
			low = 1
		}
		if minNeeded > low {
			low = minNeeded
		}

		high := prev * 2
		if i == 1 {
			high = limit + 1
		}
		if limit < high {
			high = limit
		}

		if low > high {
			return false
		}

		val := low
		currentN -= val
		prev = val
	}

	return currentN == 0
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
