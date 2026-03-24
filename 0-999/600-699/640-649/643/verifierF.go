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

const solverSrc = `package main

import (
	"fmt"
	"io"
	"os"
)

func gcd(a, b uint32) uint32 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func readUint32(b []byte, idx int) (uint32, int) {
	var res uint32
	for idx < len(b) && (b[idx] < '0' || b[idx] > '9') {
		idx++
	}
	for idx < len(b) && b[idx] >= '0' && b[idx] <= '9' {
		res = res*10 + uint32(b[idx]-'0')
		idx++
	}
	return res, idx
}

func main() {
	b, _ := io.ReadAll(os.Stdin)
	var n, p, q uint32
	idx := 0
	n, idx = readUint32(b, idx)
	p, idx = readUint32(b, idx)
	q, idx = readUint32(b, idx)

	k := p
	if n-1 < k {
		k = n - 1
	}

	C := make([]uint32, k+1)
	C[0] = 1
	for j := uint32(1); j <= k; j++ {
		num := make([]uint32, j)
		for m := uint32(0); m < j; m++ {
			num[m] = n - m
		}
		for m := uint32(1); m <= j; m++ {
			div := m
			for x := uint32(0); x < j && div > 1; x++ {
				g := gcd(num[x], div)
				num[x] /= g
				div /= g
			}
		}
		var res uint32 = 1
		for x := uint32(0); x < j; x++ {
			res *= num[x]
		}
		C[j] = res
	}

	var ans uint32
	for i := uint32(1); i <= q; i++ {
		Ri := C[k]
		for j := int(k) - 1; j >= 0; j-- {
			Ri = Ri*i + C[j]
		}
		ans ^= Ri * i
	}

	fmt.Println(ans)
}
`

type testCase struct {
	input  string
	expect uint32
}

func buildSolver() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verF643")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(dir) }
	src := filepath.Join(dir, "solver.go")
	if err := os.WriteFile(src, []byte(solverSrc), 0644); err != nil {
		cleanup()
		return "", nil, err
	}
	bin := filepath.Join(dir, "solver")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("build solver: %v\n%s", err, out)
	}
	return bin, cleanup, nil
}

func genCase(rng *rand.Rand, ref string) testCase {
	n := uint64(rng.Intn(100) + 1)
	p := rng.Intn(20) + 1
	q := rng.Intn(50) + 1
	input := fmt.Sprintf("%d %d %d\n", n, p, q)

	cmd := exec.Command(ref)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Run()
	val, _ := strconv.ParseUint(strings.TrimSpace(out.String()), 10, 32)
	return testCase{input, uint32(val)}
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, cleanup, err := buildSolver()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		c := genCase(rng, ref)
		out, err := run(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
		val, err := strconv.ParseUint(strings.TrimSpace(out), 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: invalid integer output\n", i+1)
			os.Exit(1)
		}
		if uint32(val) != c.expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", i+1, c.expect, val, c.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
