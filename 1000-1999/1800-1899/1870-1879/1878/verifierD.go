package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const timeout = 2 * time.Minute

func runExe(path string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else if strings.HasSuffix(path, ".py") {
		cmd = exec.CommandContext(ctx, "python3", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildRef() (string, error) {
	srcPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if srcPath == "" {
		srcPath = "1878D.cpp"
	}
	ref, err := filepath.Abs("./refD.bin")
	if err != nil {
		return "", err
	}
	// Detect language from content, not extension (worker saves all as .go)
	content, _ := os.ReadFile(srcPath)
	lang := detectLang(string(content))
	actualSrc := srcPath
	var cmd *exec.Cmd
	switch lang {
	case "cpp":
		// Copy to .cpp so g++ doesn't try the Go frontend
		actualSrc = filepath.Join(filepath.Dir(ref), "ref.cpp")
		os.WriteFile(actualSrc, content, 0644)
		cmd = exec.Command("g++", "-O2", "-o", ref, actualSrc)
	case "c":
		actualSrc = filepath.Join(filepath.Dir(ref), "ref.c")
		os.WriteFile(actualSrc, content, 0644)
		cmd = exec.Command("gcc", "-O2", "-o", ref, actualSrc)
	default:
		cmd = exec.Command("go", "build", "-o", ref, srcPath)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func detectLang(code string) string {
	if strings.Contains(code, "#include") {
		return "cpp"
	}
	if strings.Contains(code, "int main") && !strings.Contains(code, "package main") {
		return "c"
	}
	return "go"
}

func randString(r *rand.Rand, n int) string {
	b := make([]byte, n)
	for i := 0; i < n; i++ {
		b[i] = byte('a' + r.Intn(26))
	}
	return string(b)
}

func genCase(r *rand.Rand) string {
	n := r.Intn(10) + 1
	k := r.Intn(n) + 1
	s := randString(r, n)
	l := make([]int, k)
	rarr := make([]int, k)
	start := 1
	for i := 0; i < k; i++ {
		remain := k - i
		maxLen := n - start + 1 - (remain - 1)
		length := r.Intn(maxLen) + 1
		l[i] = start
		rarr[i] = start + length - 1
		start = rarr[i] + 1
	}
	q := r.Intn(10) + 1
	xs := make([]int, q)
	for i := range xs {
		xs[i] = r.Intn(n) + 1
	}

	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	sb.WriteString(s)
	sb.WriteByte('\n')
	for i, v := range l {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i, v := range rarr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i, v := range xs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierD /path/to/candidate")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := runExe(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := runExe(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
