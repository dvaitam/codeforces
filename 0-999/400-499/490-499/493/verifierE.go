package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveRef(t, a, b int64) string {
	if t == 1 {
		if a == 1 {
			if b == 1 {
				return "inf"
			}
			return "0"
		}
		cnt := 0
		var search func(deg int, remSum int64, remVal int64)
		search = func(deg int, remSum int64, remVal int64) {
			powA := int64(1)
			for i := 0; i < deg; i++ {
				if math.MaxInt64/a < powA {
					powA = math.MaxInt64
				} else {
					powA *= a
				}
			}
			if powA > remVal {
				if remSum == 0 && remVal == 0 {
					cnt++
				}
				return
			}
			
			limit := remSum
			if powA > 0 {
				l2 := remVal / powA
				if l2 < limit {
					limit = l2
				}
			}
			
			for c := int64(0); c <= limit; c++ {
				search(deg+1, remSum-c, remVal-c*powA)
			}
		}
		search(0, a, b)
		return fmt.Sprintf("%d", cnt%1000000007)
	} else {
		cnt := 0
		var search func(deg int, remA int64, currentB int64)
		search = func(deg int, remA int64, currentB int64) {
			powA := int64(1)
			for i := 0; i < deg; i++ {
				if math.MaxInt64/a < powA {
					powA = math.MaxInt64
				} else {
					powA *= a
				}
			}
			
			if remA == 0 {
				if currentB == b {
					cnt++
				}
				return
			}
			
			base := remA % t
			for c := base; c <= remA; c += t {
				termB := c
				bad := false
				// termB = c * a^deg
				// Check overflow for a^deg
				if powA == math.MaxInt64 {
					bad = true
				} else {
					// check c * powA
					if powA > 0 && c > math.MaxInt64/powA {
						bad = true
					} else {
						termB = c * powA
					}
				}
				
				if bad || currentB > math.MaxInt64 - termB {
					if currentB + termB > b {
						break
					}
				}
				if currentB + termB > b {
					break
				}
				
				search(deg+1, (remA-c)/t, currentB+termB)
			}
		}
		search(0, a, 0)
		return fmt.Sprintf("%d", cnt)
	}
}

func generateCaseE(rng *rand.Rand) string {
	t := rng.Int63n(5) + 1
	a := rng.Int63n(100) + 1 // Reduced for safety of recursive BF
	b := rng.Int63n(100) + 1
	return fmt.Sprintf("%d %d %d\n", t, a, b)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseE(rng)
		var t, a, b int64
		fmt.Sscan(tc, &t, &a, &b)
		expect := solveRef(t, a, b)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
