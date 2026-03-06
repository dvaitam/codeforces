package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Case struct{ input string }

func buildRef() (string, error) {
	ref := "./refC.bin"
	cmd := exec.Command("go", "build", "-o", ref, "613C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin, input string) (string, error) {
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

func genCases() []Case {
	rng := rand.New(rand.NewSource(63))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(5) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Intn(10) + 1
			fmt.Fprintf(&sb, "%d", val)
		}
		sb.WriteByte('\n')
		cases[i] = Case{sb.String()}
	}
	return cases
}

func isPalindrome(s string) bool {
	n := len(s)
	for i := 0; i < n/2; i++ {
		if s[i] != s[n-1-i] {
			return false
		}
	}
	return true
}

func countCuts(s string) int {
	n := len(s)
	count := 0
	doubled := s + s
	for k := 0; k < n; k++ {
		if isPalindrome(doubled[k : k+n]) {
			count++
		}
	}
	return count
}

func parseInput(input string) (n int, counts []int) {
	fields := strings.Fields(input)
	fmt.Sscan(fields[0], &n)
	counts = make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Sscan(fields[1+i], &counts[i])
	}
	return
}

func runCase(bin, ref string, c Case) error {
	expOut, err := runBinary(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}

	// Parse expected cut count (first line of reference)
	expLines := strings.SplitN(strings.TrimSpace(expOut), "\n", 2)
	var expCuts int
	fmt.Sscan(expLines[0], &expCuts)

	// Parse candidate output
	gotLines := strings.SplitN(strings.TrimSpace(got), "\n", 2)
	if len(gotLines) == 0 {
		return fmt.Errorf("empty output")
	}
	var gotCuts int
	if _, err := fmt.Sscan(gotLines[0], &gotCuts); err != nil {
		return fmt.Errorf("cannot parse cut count: %v", err)
	}
	if gotCuts != expCuts {
		return fmt.Errorf("expected %d cuts got %d", expCuts, gotCuts)
	}
	if gotCuts == 0 {
		// Just need any necklace with correct counts
		if len(gotLines) < 2 {
			return fmt.Errorf("missing necklace line")
		}
		// verify counts
		n, counts := parseInput(c.input)
		necklace := strings.TrimSpace(gotLines[1])
		got2 := make([]int, n)
		for _, ch := range necklace {
			idx := int(ch - 'a')
			if idx < 0 || idx >= n {
				return fmt.Errorf("invalid character %c", ch)
			}
			got2[idx]++
		}
		for i := 0; i < n; i++ {
			if got2[i] != counts[i] {
				return fmt.Errorf("wrong bead count for color %c: expected %d got %d", 'a'+i, counts[i], got2[i])
			}
		}
		return nil
	}

	// Verify candidate necklace
	if len(gotLines) < 2 {
		return fmt.Errorf("missing necklace line")
	}
	necklace := strings.TrimSpace(gotLines[1])
	n, counts := parseInput(c.input)

	// Check bead counts
	got2 := make([]int, n)
	for _, ch := range necklace {
		idx := int(ch - 'a')
		if idx < 0 || idx >= n {
			return fmt.Errorf("invalid character %c", ch)
		}
		got2[idx]++
	}
	for i := 0; i < n; i++ {
		if got2[i] != counts[i] {
			return fmt.Errorf("wrong bead count for color %c: expected %d got %d", 'a'+i, counts[i], got2[i])
		}
	}

	// Check actual cut count
	actual := countCuts(necklace)
	if actual != gotCuts {
		return fmt.Errorf("necklace has %d beautiful cuts but claimed %d", actual, gotCuts)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, ref, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
