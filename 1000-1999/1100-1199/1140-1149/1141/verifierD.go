package main

import (
	"bufio"
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

func compileRef() (string, error) {
	bin := filepath.Join(os.TempDir(), "1141D_ref")
	cmd := exec.Command("go", "build", "-o", bin, "1141D.go")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("compile reference: %v\n%s", err, out)
	}
	return bin, nil
}

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
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func randChar(rng *rand.Rand) byte {
	if rng.Intn(4) == 0 {
		return '?'
	}
	return byte('a' + rng.Intn(26))
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var a, b strings.Builder
	for i := 0; i < n; i++ {
		a.WriteByte(randChar(rng))
	}
	for i := 0; i < n; i++ {
		b.WriteByte(randChar(rng))
	}
	return fmt.Sprintf("%d\n%s\n%s\n", n, a.String(), b.String())
}

func parseInput(in string) (int, string, string, error) {
	var n int
	var a, b string
	_, err := fmt.Fscan(strings.NewReader(in), &n, &a, &b)
	return n, a, b, err
}

func parseOutput(out string) (int, [][2]int, error) {
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out)))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		return 0, nil, fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(sc.Text())
	if err != nil {
		return 0, nil, fmt.Errorf("invalid k")
	}
	pairs := make([][2]int, k)
	for i := 0; i < k; i++ {
		if !sc.Scan() {
			return 0, nil, fmt.Errorf("missing a")
		}
		a, err := strconv.Atoi(sc.Text())
		if err != nil {
			return 0, nil, fmt.Errorf("invalid a")
		}
		if !sc.Scan() {
			return 0, nil, fmt.Errorf("missing b")
		}
		b, err := strconv.Atoi(sc.Text())
		if err != nil {
			return 0, nil, fmt.Errorf("invalid b")
		}
		pairs[i] = [2]int{a, b}
	}
	if sc.Scan() {
		return 0, nil, fmt.Errorf("extra output")
	}
	return k, pairs, nil
}

func checkOutput(n int, l, r string, expectK int, out string) error {
	k, pairs, err := parseOutput(out)
	if err != nil {
		return err
	}
	if k != expectK {
		return fmt.Errorf("expected %d pairs got %d", expectK, k)
	}
	usedL := make([]bool, n)
	usedR := make([]bool, n)
	for _, p := range pairs {
		a, b := p[0], p[1]
		if a < 1 || a > n || b < 1 || b > n {
			return fmt.Errorf("index out of range")
		}
		if usedL[a-1] || usedR[b-1] {
			return fmt.Errorf("index repeated")
		}
		usedL[a-1] = true
		usedR[b-1] = true
		ca := l[a-1]
		cb := r[b-1]
		if ca != cb && ca != '?' && cb != '?' {
			return fmt.Errorf("incompatible pair")
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		wantOut, err := run(ref, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failure: %v\n", err)
			os.Exit(1)
		}
		got, err := run(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		n, a, b, err := parseInput(input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal parse error: %v\n", err)
			os.Exit(1)
		}
		sc := bufio.NewScanner(strings.NewReader(wantOut))
		sc.Split(bufio.ScanWords)
		sc.Scan()
		expectK, _ := strconv.Atoi(sc.Text())
		if err := checkOutput(n, a, b, expectK, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
