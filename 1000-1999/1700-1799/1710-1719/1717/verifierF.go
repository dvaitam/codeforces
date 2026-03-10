package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"log"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func prepareBinary(path, tag string) (string, error) {
	if strings.HasSuffix(path, ".go") {
		bin := filepath.Join(os.TempDir(), tag+"_"+fmt.Sprint(time.Now().UnixNano()))
		cmd := exec.Command("go", "build", "-o", bin, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build %s: %v\n%s", path, err, out)
		}
		return bin, nil
	}
	return path, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2 // 2..6
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(5) + 1
	if m > maxEdges {
		m = maxEdges
	}
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	s := make([]int, n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if rng.Intn(2) == 0 {
			sb.WriteByte('0')
			s[i] = 0
		} else {
			sb.WriteByte('1')
			s[i] = 1
		}
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		if s[i] == 0 {
			sb.WriteByte('0')
		} else {
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(2*m+1)-m))
		}
	}
	sb.WriteByte('\n')
	used := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		for {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			p1 := [2]int{u, v}
			p2 := [2]int{v, u}
			if used[p1] || used[p2] {
				continue
			}
			used[p1] = true
			sb.WriteString(fmt.Sprintf("%d %d\n", u, v))
			break
		}
	}
	return sb.String()
}

func validateOutput(input, output string) error {
	// Parse input
	inLines := strings.Split(strings.TrimSpace(input), "\n")
	header := strings.Fields(inLines[0])
	n, _ := strconv.Atoi(header[0])
	m, _ := strconv.Atoi(header[1])
	sFields := strings.Fields(inLines[1])
	aFields := strings.Fields(inLines[2])
	s := make([]int, n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		s[i], _ = strconv.Atoi(sFields[i])
		a[i], _ = strconv.Atoi(aFields[i])
	}
	type edge struct{ u, v int }
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		fields := strings.Fields(inLines[3+i])
		u, _ := strconv.Atoi(fields[0])
		v, _ := strconv.Atoi(fields[1])
		edges[i] = edge{u, v}
	}
	// Parse output
	outLines := strings.Split(output, "\n")
	if len(outLines) < 1+m {
		return fmt.Errorf("expected %d lines of output after YES, got %d", m, len(outLines)-1)
	}
	// Build b array from candidate's edge orientations
	b := make([]int, n+1)
	usedEdges := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		fields := strings.Fields(outLines[1+i])
		if len(fields) != 2 {
			return fmt.Errorf("line %d: expected 2 integers, got %q", i+1, outLines[1+i])
		}
		u, err1 := strconv.Atoi(fields[0])
		v, err2 := strconv.Atoi(fields[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("line %d: parse error", i+1)
		}
		// Check this is a valid orientation of some input edge
		e1 := [2]int{u, v}
		e2 := [2]int{v, u}
		foundIdx := -1
		for j, e := range edges {
			if (e.u == u && e.v == v) || (e.u == v && e.v == u) {
				foundIdx = j
				break
			}
		}
		if foundIdx == -1 {
			return fmt.Errorf("line %d: edge (%d,%d) not in input", i+1, u, v)
		}
		if usedEdges[e1] || usedEdges[e2] {
			return fmt.Errorf("line %d: duplicate edge (%d,%d)", i+1, u, v)
		}
		usedEdges[e1] = true
		// b[u] -= 1, b[v] += 1
		b[u]--
		b[v]++
	}
	if len(usedEdges) != m {
		return fmt.Errorf("expected %d edges, got %d", m, len(usedEdges))
	}
	// Check b[i] == a[i] for all i where s[i] == 1
	for i := 0; i < n; i++ {
		if s[i] == 1 && b[i+1] != a[i] {
			return fmt.Errorf("b[%d]=%d but a[%d]=%d (s[%d]=1)", i+1, b[i+1], i+1, a[i], i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := prepareBinary(os.Args[1], "candF")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		log.Fatal("REFERENCE_SOURCE_PATH environment variable is not set")
	}
	refPath, err := prepareBinary(refSrc, "refF")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runBinary(refPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		expTrim := strings.TrimSpace(exp)
		gotTrim := strings.TrimSpace(got)
		expUpper := strings.ToUpper(strings.Split(expTrim, "\n")[0])
		gotUpper := strings.ToUpper(strings.Split(gotTrim, "\n")[0])
		if expUpper == "NO" {
			if gotUpper != "NO" {
				fmt.Printf("case %d failed: expected NO, got:\n%s\ninput:\n%s", i+1, got, input)
				os.Exit(1)
			}
			continue
		}
		// Reference says YES, validate candidate output
		if gotUpper != "YES" {
			fmt.Printf("case %d failed: expected YES, got:\n%s\ninput:\n%s", i+1, got, input)
			os.Exit(1)
		}
		// Parse input to get n, m, s, a, edges
		if err := validateOutput(input, gotTrim); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%sgot:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
