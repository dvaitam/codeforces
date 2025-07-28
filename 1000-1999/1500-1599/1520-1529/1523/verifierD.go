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
)

const numTestsD = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "candD")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func buildOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleD")
	cmd := exec.Command("go", "build", "-o", tmp, "1523D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	m := rng.Intn(10) + 1
	p := rng.Intn(m) + 1
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		ones := rng.Intn(p + 1)
		idx := rng.Perm(m)[:ones]
		for _, x := range idx {
			b[x] = '1'
		}
		for j := 0; j < m; j++ {
			if b[j] == 0 {
				b[j] = '0'
			}
		}
		lines[i] = string(b)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, p))
	for _, s := range lines {
		sb.WriteString(s)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseInput(input string) (n, m int, grid []string) {
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	sc.Scan()
	n, _ = strconv.Atoi(sc.Text())
	sc.Scan()
	m, _ = strconv.Atoi(sc.Text())
	sc.Scan() // skip p
	grid = make([]string, n)
	for i := 0; i < n; i++ {
		sc.Scan()
		grid[i] = sc.Text()
	}
	return
}

func checkCase(input, want, got string) error {
	n, m, grid := parseInput(input)
	if len(got) != m {
		return fmt.Errorf("output length %d expected %d", len(got), m)
	}
	onesWant := strings.Count(want, "1")
	onesGot := strings.Count(got, "1")
	if onesGot != onesWant {
		return fmt.Errorf("expected %d ones, got %d", onesWant, onesGot)
	}
	half := (n + 1) / 2
	for j := 0; j < m; j++ {
		if got[j] == '1' {
			cnt := 0
			for i := 0; i < n; i++ {
				if grid[i][j] == '1' {
					cnt++
				}
			}
			if cnt < half {
				return fmt.Errorf("bit %d liked by %d < %d", j+1, cnt, half)
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin, clean, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if clean != nil {
		defer clean()
	}
	oracle, c2, err := buildOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c2()

	rng := rand.New(rand.NewSource(3))
	for i := 0; i < numTestsD; i++ {
		input := genCase(rng)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error on case %d: %v\n", i+1, err)
			return
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if err := checkCase(input, want, got); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, err, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
