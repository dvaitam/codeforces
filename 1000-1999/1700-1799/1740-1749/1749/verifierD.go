package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const MOD int64 = 998244353

func buildExecutable(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "bin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func oracle(input string) string {
	r := strings.NewReader(input)
	var n int
	var m int64
	if _, err := fmt.Fscan(r, &n, &m); err != nil {
		return ""
	}
	isPrime := make([]bool, n+1)
	if n >= 2 {
		for i := 2; i <= n; i++ {
			isPrime[i] = true
		}
		for i := 2; i*i <= n; i++ {
			if isPrime[i] {
				for j := i * i; j <= n; j += i {
					isPrime[j] = false
				}
			}
		}
	}
	mm := m % MOD
	total := int64(0)
	pw := int64(1)
	for i := 1; i <= n; i++ {
		pw = pw * mm % MOD
		total = (total + pw) % MOD
	}
	prod := int64(1)
	ways := mm
	unamb := ways % MOD
	for i := 2; i <= n; i++ {
		if isPrime[i] {
			prod *= int64(i)
		}
		if prod > m {
			break
		}
		q := m / prod
		ways = ways * (q % MOD) % MOD
		unamb = (unamb + ways) % MOD
	}
	ans := (total - unamb) % MOD
	if ans < 0 {
		ans += MOD
	}
	return fmt.Sprint(ans)
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 1
	m := int64(rng.Intn(1000) + 1)
	return fmt.Sprintf("%d %d\n", n, m)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	binPath := os.Args[1]
	bin, cleanup, err := buildExecutable(binPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to prepare binary: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(45))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		expect := oracle(tc)
		got, err := run(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
