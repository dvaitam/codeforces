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
	"time"
)

func buildRef() (string, error) {
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	ref := filepath.Join(os.TempDir(), "ref730I.bin")
	cmd := exec.Command("go", "build", "-o", ref, refSrc)
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
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

type Case struct{ input string }

// validate checks that the candidate output is a valid solution for the 730I problem.
// The problem: choose p students for programming team and s students for sports team (disjoint).
// Maximize sum of a[i] for programming team + sum of b[i] for sports team.
// Output: first line = total strength, second line = programming team members, third line = sports team members.
func validate(input, output string) error {
	r := strings.NewReader(input)
	var n, p, s int
	fmt.Fscan(r, &n, &p, &s)
	a := make([]int, n+1)
	b := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(r, &b[i])
	}

	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) < 3 {
		return fmt.Errorf("expected 3 lines, got %d", len(lines))
	}

	claimedStr := strings.TrimSpace(lines[0])
	claimed, err := strconv.Atoi(claimedStr)
	if err != nil {
		return fmt.Errorf("cannot parse total strength: %v", err)
	}

	progFields := strings.Fields(strings.TrimSpace(lines[1]))
	if len(progFields) != p {
		return fmt.Errorf("expected %d programming team members, got %d", p, len(progFields))
	}
	sportFields := strings.Fields(strings.TrimSpace(lines[2]))
	if len(sportFields) != s {
		return fmt.Errorf("expected %d sports team members, got %d", s, len(sportFields))
	}

	used := make(map[int]bool)
	totalStrength := 0
	for _, f := range progFields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("cannot parse team member: %v", err)
		}
		if v < 1 || v > n {
			return fmt.Errorf("member %d out of range", v)
		}
		if used[v] {
			return fmt.Errorf("member %d used twice", v)
		}
		used[v] = true
		totalStrength += a[v]
	}
	for _, f := range sportFields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("cannot parse team member: %v", err)
		}
		if v < 1 || v > n {
			return fmt.Errorf("member %d out of range", v)
		}
		if used[v] {
			return fmt.Errorf("member %d used twice", v)
		}
		used[v] = true
		totalStrength += b[v]
	}

	if totalStrength != claimed {
		return fmt.Errorf("claimed strength %d but computed %d", claimed, totalStrength)
	}

	return nil
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(5) + 2
		p := rng.Intn(n-1) + 1
		maxS := n - p
		s := rng.Intn(maxS) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, p, s)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(20)+1)
		}
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(20)+1)
		}
		sb.WriteByte('\n')
		cases[i] = Case{sb.String()}
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
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
		// Get reference answer for optimal value
		expOut, err := runBinary(ref, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		expLines := strings.Split(expOut, "\n")
		expVal := strings.TrimSpace(expLines[0])

		got, err := runBinary(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}

		// Validate structure of output
		if verr := validate(c.input, got); verr != nil {
			fmt.Fprintf(os.Stderr, "case %d validation failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, verr, c.input, got)
			os.Exit(1)
		}

		// Check that the optimal value matches
		gotLines := strings.Split(got, "\n")
		gotVal := strings.TrimSpace(gotLines[0])
		if gotVal != expVal {
			fmt.Fprintf(os.Stderr, "case %d: expected optimal value %s, got %s\ninput:\n%s", i+1, expVal, gotVal, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
