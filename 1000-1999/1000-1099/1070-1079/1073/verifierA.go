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

func runProg(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return out.String() + stderr.String(), fmt.Errorf("%v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + rng.Intn(26))
	}
	return fmt.Sprintf("%d\n%s\n", n, string(b))
}

// isDiverse returns true if no letter appears strictly more than len(t)/2 times.
func isDiverse(t string) bool {
	n := len(t)
	freq := [26]int{}
	for i := 0; i < n; i++ {
		freq[t[i]-'a']++
	}
	for _, c := range freq {
		if 2*c > n {
			return false
		}
	}
	return true
}

// hasDiverseSubstring returns true if s has any diverse substring.
// Since any 2-char substring with two different chars is diverse, this is
// equivalent to checking whether s has at least two distinct characters.
func hasDiverseSubstring(s string) bool {
	for i := 1; i < len(s); i++ {
		if s[i] != s[0] {
			return true
		}
	}
	return false
}

func validate(input, output string) error {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	s := strings.TrimSpace(lines[1])

	outLines := strings.Split(strings.TrimSpace(output), "\n")

	verdict := strings.TrimSpace(outLines[0])
	if verdict != "YES" && verdict != "NO" {
		return fmt.Errorf("first line must be YES or NO, got %q", verdict)
	}

	if verdict == "NO" {
		if hasDiverseSubstring(s) {
			return fmt.Errorf("output NO but a diverse substring exists in %q", s)
		}
		return nil
	}

	// verdict == "YES"
	if !hasDiverseSubstring(s) {
		return fmt.Errorf("output YES but no diverse substring exists in %q", s)
	}
	if len(outLines) < 2 {
		return fmt.Errorf("YES verdict missing substring on second line")
	}
	sub := strings.TrimSpace(outLines[1])
	if len(sub) == 0 {
		return fmt.Errorf("returned empty substring")
	}
	if !strings.Contains(s, sub) {
		return fmt.Errorf("returned %q is not a substring of %q", sub, s)
	}
	if !isDiverse(sub) {
		return fmt.Errorf("returned %q is not diverse", sub)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 1000; i++ {
		input := genCase(rng)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		if err := validate(input, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput: %s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
