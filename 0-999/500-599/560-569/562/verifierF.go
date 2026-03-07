package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "562F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

const letters = "abcdefghijklmnopqrstuvwxyz"

func randWord(rng *rand.Rand) string {
	l := rng.Intn(3) + 1
	b := make([]byte, l)
	for i := range b {
		b[i] = letters[rng.Intn(len(letters))]
	}
	return string(b)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(randWord(rng))
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		sb.WriteString(randWord(rng))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}

func lcpLen(a, b string) int {
	l := 0
	for l < len(a) && l < len(b) && a[l] == b[l] {
		l++
	}
	return l
}

func parseInput(in string) (int, []string, []string) {
	sc := bufio.NewScanner(strings.NewReader(in))
	sc.Scan()
	var n int
	fmt.Sscan(sc.Text(), &n)
	names := make([]string, n)
	pseudos := make([]string, n)
	for i := 0; i < n; i++ {
		sc.Scan()
		names[i] = strings.TrimSpace(sc.Text())
	}
	for i := 0; i < n; i++ {
		sc.Scan()
		pseudos[i] = strings.TrimSpace(sc.Text())
	}
	return n, names, pseudos
}

// checkOutput validates the output and returns the actual quality, or an error.
func checkOutput(n int, names, pseudos []string, output string) (int, error) {
	sc := bufio.NewScanner(strings.NewReader(output))
	if !sc.Scan() {
		return 0, fmt.Errorf("output is empty")
	}
	var claimedQuality int
	if _, err := fmt.Sscan(sc.Text(), &claimedQuality); err != nil {
		return 0, fmt.Errorf("failed to parse quality: %v", err)
	}
	usedStudent := make([]bool, n+1)
	usedPseudo := make([]bool, n+1)
	actualQuality := 0
	for i := 0; i < n; i++ {
		if !sc.Scan() {
			return 0, fmt.Errorf("expected %d pairs, got %d", n, i)
		}
		var a, b int
		if _, err := fmt.Sscan(sc.Text(), &a, &b); err != nil {
			return 0, fmt.Errorf("failed to parse pair %d: %v", i+1, err)
		}
		if a < 1 || a > n {
			return 0, fmt.Errorf("student index %d out of range [1,%d]", a, n)
		}
		if b < 1 || b > n {
			return 0, fmt.Errorf("pseudo index %d out of range [1,%d]", b, n)
		}
		if usedStudent[a] {
			return 0, fmt.Errorf("student %d used more than once", a)
		}
		if usedPseudo[b] {
			return 0, fmt.Errorf("pseudo %d used more than once", b)
		}
		usedStudent[a] = true
		usedPseudo[b] = true
		actualQuality += lcpLen(names[a-1], pseudos[b-1])
	}
	if actualQuality != claimedQuality {
		return 0, fmt.Errorf("claimed quality %d but actual LCP sum is %d", claimedQuality, actualQuality)
	}
	return actualQuality, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		n, names, pseudos := parseInput(in)

		refOut, err := runBinary(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		refQuality, err := checkOutput(n, names, pseudos, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d: %v\n", i+1, err)
			os.Exit(1)
		}

		got, err := runBinary(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		candQuality, err := checkOutput(n, names, pseudos, got)
		if err != nil {
			fmt.Printf("test %d failed: invalid output: %v\ninput:\n%sgot:\n%s\n", i+1, err, in, got)
			os.Exit(1)
		}
		if candQuality != refQuality {
			fmt.Printf("test %d failed\ninput:\n%sexpected quality: %d\ngot quality: %d\ngot:\n%s\n", i+1, in, refQuality, candQuality, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
