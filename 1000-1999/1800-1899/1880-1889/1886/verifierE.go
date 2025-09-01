package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
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

// parseInput1886E parses the generated input back into n, m, values a (1..n), b (1..m).
func parseInput1886E(input string) (int, int, []int, []int, error) {
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) < 3 {
		return 0, 0, nil, nil, fmt.Errorf("bad input")
	}
	hdr := strings.Fields(lines[0])
	if len(hdr) < 2 {
		return 0, 0, nil, nil, fmt.Errorf("bad header")
	}
	n, err1 := strconv.Atoi(hdr[0])
	m, err2 := strconv.Atoi(hdr[1])
	if err1 != nil || err2 != nil {
		return 0, 0, nil, nil, fmt.Errorf("bad n/m")
	}
	fa := strings.Fields(lines[1])
	if len(fa) != n {
		return 0, 0, nil, nil, fmt.Errorf("bad a size")
	}
	a := make([]int, n+1)
	for i := 0; i < n; i++ {
		v, err := strconv.Atoi(fa[i])
		if err != nil {
			return 0, 0, nil, nil, fmt.Errorf("bad a val")
		}
		a[i+1] = v
	}
	fb := strings.Fields(lines[2])
	if len(fb) != m {
		return 0, 0, nil, nil, fmt.Errorf("bad b size")
	}
	b := make([]int, m+1)
	for i := 0; i < m; i++ {
		v, err := strconv.Atoi(fb[i])
		if err != nil {
			return 0, 0, nil, nil, fmt.Errorf("bad b val")
		}
		b[i+1] = v
	}
	return n, m, a, b, nil
}

// validateYes1886E checks candidate YES output semantically.
func validateYes1886E(input, got string) bool {
	n, m, a, b, err := parseInput1886E(input)
	if err != nil {
		return false
	}
	val := make([]int, n+1)
	for i := 1; i <= n; i++ {
		val[i] = a[i]
	}

	lines := strings.Split(strings.TrimSpace(got), "\n")
	// find header line
	idx := 0
	for idx < len(lines) && strings.TrimSpace(lines[idx]) == "" {
		idx++
	}
	if idx >= len(lines) || strings.ToUpper(strings.TrimSpace(lines[idx])) != "YES" {
		return false
	}
	idx++

	groups := make([][]int, 0, m)
	used := make([]bool, n+1)
	for len(groups) < m && idx < len(lines) {
		ln := strings.TrimSpace(lines[idx])
		idx++
		if ln == "" {
			continue
		}
		fields := strings.Fields(ln)
		cnt, err := strconv.Atoi(fields[0])
		if err != nil || cnt < 0 || len(fields)-1 < cnt {
			return false
		}
		ids := make([]int, 0, cnt)
		for k := 0; k < cnt; k++ {
			id, err := strconv.Atoi(fields[1+k])
			if err != nil || id < 1 || id > n || used[id] {
				return false
			}
			used[id] = true
			ids = append(ids, id)
		}
		groups = append(groups, ids)
	}
	if len(groups) != m {
		return false
	}

	for gi := 0; gi < m; gi++ {
		ids := groups[gi]
		cnt := len(ids)
		if cnt == 0 {
			if b[gi+1] != 0 {
				return false
			}
			continue
		}
		minVal := val[ids[0]]
		for _, id := range ids[1:] {
			if val[id] < minVal {
				minVal = val[id]
			}
		}
		if int64(minVal)*int64(cnt) < int64(b[gi+1]) {
			return false
		}
	}
	return true
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
		wHead := strings.ToUpper(strings.TrimSpace(strings.Split(strings.TrimSpace(want), "\n")[0]))
		gHead := strings.ToUpper(strings.TrimSpace(strings.Split(strings.TrimSpace(got), "\n")[0]))
		if wHead == "NO" {
			// Accept candidate NO or any semantically valid YES
			if gHead == "NO" {
				continue
			}
			if gHead == "YES" && validateYes1886E(input, got) {
				continue
			}
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, strings.TrimSpace(want), strings.TrimSpace(got))
			os.Exit(1)
		}
		// YES case: semantically validate candidate
		if gHead != "YES" || !validateYes1886E(input, got) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, strings.TrimSpace(want), strings.TrimSpace(got))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
