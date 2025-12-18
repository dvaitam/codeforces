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

func solve(n int64) string {
	mp := make(map[int64]string)
	// Precompute powers
	// i*i can overflow if i is large int64, but i <= sqrt(n) <= 10^5, so i*i <= 10^10. Safe.
	for i := int64(2); i*i <= n; i++ {
		val := i * i
		for k := 2; val <= n; k++ {
			s := fmt.Sprintf("%d^%d", i, k)
			if best, ok := mp[val]; !ok {
				mp[val] = strconv.FormatInt(val, 10)
				if len(s) < len(mp[val]) {
					mp[val] = s
				}
			} else {
				if len(s) < len(best) {
					mp[val] = s
				}
			}
			
			// Check for overflow before multiplying
			if n/i < val {
				break
			}
			val *= i
		}
	}

	ans := strconv.FormatInt(n, 10)

	getStr := func(v int64) string {
		if s, ok := mp[v]; ok {
			return s
		}
		return strconv.FormatInt(v, 10)
	}

	for val, rep := range mp {
		if val > n { continue } // Should not happen by loop condition
		
		k := n / val
		b := n % val
		
		var curr string
		if k == 0 {
			// Should not happen since val <= n
			continue 
		} else if k == 1 {
			curr = rep
		} else {
			curr = getStr(k) + "*" + rep
		}
		
		if b > 0 {
			curr += "+" + getStr(b)
		}
		
		if len(curr) < len(ans) {
			ans = curr
		}
	}
	
	return ans
}

func generateCaseF(rng *rand.Rand) (string, string) {
	// Generate n up to 10^10
	n := rng.Int63n(10000000000) + 1
	input := fmt.Sprintf("%d\n", n)
	expected := solve(n)
	return input, expected
}

func eval(s string, target int64) (int64, error) {
	for _, c := range s {
		if !strings.ContainsRune("0123456789+*^", c) {
			return 0, fmt.Errorf("invalid char %c", c)
		}
	}
	
	parts := strings.Split(s, "+")
	var sum int64
	for _, p := range parts {
		if p == "" { return 0, fmt.Errorf("empty term") }
		factors := strings.Split(p, "*")
		prod := int64(1)
		for _, f := range factors {
			if f == "" { return 0, fmt.Errorf("empty factor") }
			baseExp := strings.Split(f, "^")
			if len(baseExp) > 2 {
				return 0, fmt.Errorf("multiple carets")
			}
			base, err := strconv.ParseInt(baseExp[0], 10, 64)
			if err != nil { return 0, err }
			val := base
			if len(baseExp) == 2 {
				exp, err := strconv.ParseInt(baseExp[1], 10, 64)
				if err != nil { return 0, err }
				// val = pow(base, exp)
				val = 1
				for i := int64(0); i < exp; i++ {
					if base != 0 && target/base < val { // Overflow check relative to target
						val = target + 1
						break
					}
					val *= base
				}
			}
			// prod *= val
			// Check for overflow
			if prod == 0 { // 0 * anything = 0
				prod = 0
			} else {
				if val > 0 && target/prod < val {
					prod = target + 1
				} else {
					prod *= val
				}
			}
		}
		// sum += prod
		if target-sum < prod {
			sum = target + 1
		} else {
			sum += prod
		}
	}
	return sum, nil
}

func runCaseF(bin, input, expected string) error {
	nStr := strings.TrimSpace(input)
	n, err := strconv.ParseInt(nStr, 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse input: %v", err)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	
	// Check validity and value
	val, err := eval(got, n)
	if err != nil {
		return fmt.Errorf("eval error: %v", err)
	}
	if val != n {
		return fmt.Errorf("wrong value: got %d, expected %d", val, n)
	}
	
	// Check length
	if len(got) > len(expected) {
		return fmt.Errorf("suboptimal length: got %d (%q), expected %d (%q)", len(got), got, len(expected), expected)
	}
	
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCaseF(rng)
		if err := runCaseF(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}