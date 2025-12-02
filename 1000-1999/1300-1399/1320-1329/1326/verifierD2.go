package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func buildOracle() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	src := filepath.Join(dir, "1326D2.go")
	bin := filepath.Join(os.TempDir(), "oracle1326D2.bin")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return bin, nil
}

func runCandidate(bin, input string) (string, error) {
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
	return strings.TrimSpace(out.String()), nil
}

func deterministicCases() []string {
	return []string{
		"1\na\n",
		"1\nabc\n",
		"1\nabcba\n",
	}
}

func genCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		length := rng.Intn(50) + 1
		for j := 0; j < length; j++ {
			sb.WriteByte(byte('a' + rng.Intn(26)))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

// verifyPalindrome checks if the result is valid according to problem statement:
// 1. palindrome
// 2. length <= len(s)
// 3. t = a + b where a is prefix of s, b is suffix of s.
func verifyPalindrome(s, t string) error {
	if len(t) > len(s) {
		return fmt.Errorf("result length %d > original %d", len(t), len(s))
	}
	// Check palindrome
	for i := 0; i < len(t)/2; i++ {
		if t[i] != t[len(t)-1-i] {
			return fmt.Errorf("result is not a palindrome: %s", t)
		}
	}
	// Check prefix+suffix property
	// t = s[:k] + s[n-m:] where k+m = len(t)
	// We need to find if such k, m exist.
	// Since t = a + b, and a is prefix, t starts with a.
	// If t is shorter than s, a could be the whole t, or part of t.
	// Actually, simply check if t can be formed by prefix + suffix of s.
	// A greedy approach works: match longest possible prefix of t with s.
	// Whatever remains must be a suffix of s.
	
	n := len(s)
	tn := len(t)
	
	// Find maximum k such that t[:k] == s[:k]
	k := 0
	for k < tn && k < n && t[k] == s[k] {
		k++
	}
	
	// The remaining part t[k:] must match s[n-(tn-k):]
	rem := tn - k
	if rem > 0 {
		suffix := s[n-rem:]
		if t[k:] != suffix {
			// Maybe we matched too much of the prefix?
			// Actually, if t = a + b, then a is s[:len(a)].
			// So t must start with s[:len(a)].
			// If we matched max prefix, any valid 'a' must be a prefix of that match.
			// Does reducing 'a' help 'b' match?
			// b would become longer.
			// b must be s[n-len(b):].
			// If we take shorter a, we need longer b.
			// t = a' + b' where a' is prefix of a. b' = (rest of a) + b.
			// b' must be suffix of s.
			// b is suffix of s. (rest of a) is suffix of s? Only if s is periodic?
			// This verification is tricky if we don't know the split.
			// BUT, the problem says find longest t.
			// We just need to ensure t is VALID.
			// A robust check:
			// Iterate all possible split points of t into (u, v).
			// Check if u is prefix of s AND v is suffix of s.
			// Also u and v must not overlap in s in a way that uses same characters?
			// "a is prefix of s while b is suffix of s".
			// It implies non-overlapping indices?
			// "s = ... "
			// t = a + b.
			// a = s[0...len(a)-1]
			// b = s[n-len(b)...n-1]
			// Constraint: len(a) + len(b) <= n (implied by len(t) <= len(s)).
			// So we just need to find one split i (0 to len(t)) such that:
			// t[:i] == s[:i] AND t[i:] == s[n-(len(t)-i):]
			found := false
			for i := 0; i <= len(t); i++ {
				partA := t[:i]
				partB := t[i:]
				if len(partA)+len(partB) > n {
					continue
				}
				if s[:len(partA)] == partA && s[n-len(partB):] == partB {
					found = true
					break
				}
			}
			if !found {
				return fmt.Errorf("result %s is not formed by prefix+suffix of %s", t, s)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}

	for i, in := range cases {
		// Run Oracle
		wantOutput, err := runCandidate(oracle, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		
		// Run User Solution
		gotOutput, err := runCandidate(userBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}

		// Split outputs into lines (one per test case)
		wants := strings.Split(wantOutput, "\n")
		gots := strings.Split(gotOutput, "\n")
		inputs := strings.Split(strings.TrimSpace(in), "\n")
		// inputs[0] is t count. inputs[1:] are strings.

		if len(wants) != len(gots) {
			fmt.Printf("case %d failed: line count mismatch. want %d, got %d\n", i+1, len(wants), len(gots))
			os.Exit(1)
		}

		for j, want := range wants {
			got := gots[j]
			inp := inputs[j+1] // Skip t count line
			
			// Length check
			if len(got) < len(want) {
				fmt.Printf("case %d subtest %d failed: shorter result. want len %d (%s), got len %d (%s)\ninput: %s\n", 
					i+1, j+1, len(want), want, len(got), got, inp)
				os.Exit(1)
			}
			
			// Correctness check
			if err := verifyPalindrome(inp, got); err != nil {
				fmt.Printf("case %d subtest %d failed: invalid result. err: %v\ninput: %s\ngot: %s\n", 
					i+1, j+1, err, inp, got)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}