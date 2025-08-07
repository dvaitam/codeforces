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

func compileRef() (string, error) {
	out := "refF.bin"
	cmd := exec.Command("go", "build", "-o", out, "748F.go")
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return "./" + out, nil
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

func runRef(ref, input string) (string, error) {
	cmd := exec.Command(ref)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ref runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

// normalize converts the output into a canonical form so that the
// order of the k triplets and the order of endpoints within each
// triplet do not affect comparison.
func normalize(s string) (string, error) {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	if len(lines) < 2 {
		return "", fmt.Errorf("output too short")
	}
	first := strings.TrimSpace(lines[0])
	second := strings.TrimSpace(lines[1])
	type triplet struct{ a, b, c int }
	trips := make([]triplet, 0, len(lines)-2)
	for _, ln := range lines[2:] {
		fields := strings.Fields(ln)
		if len(fields) != 3 {
			return "", fmt.Errorf("invalid line: %q", ln)
		}
		a, err1 := strconv.Atoi(fields[0])
		b, err2 := strconv.Atoi(fields[1])
		c, err3 := strconv.Atoi(fields[2])
		if err1 != nil || err2 != nil || err3 != nil {
			return "", fmt.Errorf("invalid integers in line: %q", ln)
		}
		if a > b {
			a, b = b, a
		}
		trips = append(trips, triplet{a, b, c})
	}
	sort.Slice(trips, func(i, j int) bool {
		if trips[i].a != trips[j].a {
			return trips[i].a < trips[j].a
		}
		if trips[i].b != trips[j].b {
			return trips[i].b < trips[j].b
		}
		return trips[i].c < trips[j].c
	})
	var sb strings.Builder
	sb.WriteString(first)
	sb.WriteByte('\n')
	sb.WriteString(second)
	for _, t := range trips {
		sb.WriteByte('\n')
		sb.WriteString(fmt.Sprintf("%d %d %d", t.a, t.b, t.c))
	}
	return sb.String(), nil
}

func generateCase(r *rand.Rand) string {
	k := r.Intn(3) + 1
	n := 2*k + r.Intn(5) // ensure n >= 2k
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i := 2; i <= n; i++ {
		p := r.Intn(i-1) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", p, i))
	}
	nodes := r.Perm(n)
	specials := nodes[:2*k]
	for i, v := range specials {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v+1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Println("failed to compile reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(r)
		exp, err := runRef(ref, in)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		nOut, err1 := normalize(out)
		nExp, err2 := normalize(exp)
		if err1 != nil || err2 != nil || nOut != nExp {
			fmt.Printf("test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
