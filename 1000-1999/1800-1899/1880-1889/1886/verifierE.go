package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func normalizeOutput1886E(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	for i, ln := range lines {
		lines[i] = strings.Join(strings.Fields(ln), " ")
	}
	return strings.Join(lines, "\n")
}

// canonicalize1886E returns the header (YES/NO) and a sorted list of normalized lines after the header.
func canonicalize1886E(s string) (string, []string) {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	norm := make([]string, 0, len(lines))
	for _, ln := range lines {
		norm = append(norm, strings.Join(strings.Fields(ln), " "))
	}
	if len(norm) == 0 {
		return "", nil
	}
	head := strings.ToUpper(norm[0])
	if head != "YES" {
		return "NO", nil
	}
	rest := append([]string(nil), norm[1:]...)
	sort.Strings(rest)
	return "YES", rest
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1 // 1..6
	m := rng.Intn(3) + 1 // 1..3
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(20)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(20)+1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]

	ref := "refE_bin"
	if err := exec.Command("go", "build", "-o", ref, "1886E.go").Run(); err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		want, err := run("./"+ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := run(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		wHead, wLines := canonicalize1886E(want)
		gHead, gLines := canonicalize1886E(got)
		if wHead == "NO" {
			if gHead != "NO" {
				fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, strings.TrimSpace(want), strings.TrimSpace(got))
				os.Exit(1)
			}
			continue
		}
		// YES case: compare multisets of lines
		if gHead != "YES" || len(wLines) != len(gLines) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, strings.TrimSpace(want), strings.TrimSpace(got))
			os.Exit(1)
		}
		for j := range wLines {
			if wLines[j] != gLines[j] {
				fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, strings.TrimSpace(want), strings.TrimSpace(got))
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
