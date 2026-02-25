package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Case struct{ input string }

func buildRef() (string, error) {
	ref := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "613B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
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

func genCases() []Case {
	rng := rand.New(rand.NewSource(62))
	cases := make([]Case, 100)
	for i := range cases {
		n := rng.Intn(6) + 1
		A := rng.Int63n(20) + 1
		cf := rng.Int63n(10)
		cm := rng.Int63n(10)
		m := rng.Int63n(100)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d %d\n", n, A, cf, cm, m)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			v := rng.Int63n(A + 1)
			fmt.Fprintf(&sb, "%d", v)
		}
		sb.WriteByte('\n')
		cases[i] = Case{sb.String()}
	}
	return cases
}

func checkCandidate(input, expected, got string) error {
	inScan := bufio.NewScanner(strings.NewReader(input))
	inScan.Split(bufio.ScanWords)

	readInt := func() int64 {
		inScan.Scan()
		val, _ := strconv.ParseInt(inScan.Text(), 10, 64)
		return val
	}

	n := readInt()
	A := readInt()
	cf := readInt()
	cm := readInt()
	m := readInt()

	orig := make([]int64, n)
	for i := 0; i < int(n); i++ {
		orig[i] = readInt()
	}

	expScan := bufio.NewScanner(strings.NewReader(expected))
	expScan.Split(bufio.ScanWords)
	expScan.Scan()
	expForce, err := strconv.ParseInt(expScan.Text(), 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse expected force")
	}

	gotScan := bufio.NewScanner(strings.NewReader(got))
	gotScan.Split(bufio.ScanWords)
	
	if !gotScan.Scan() {
		return fmt.Errorf("candidate output is empty")
	}
	gotForce, err := strconv.ParseInt(gotScan.Text(), 10, 64)
	if err != nil {
		return fmt.Errorf("could not parse candidate force")
	}

	if gotForce != expForce {
		return fmt.Errorf("expected force %d got %d", expForce, gotForce)
	}

	candArr := make([]int64, n)
	var minVal int64 = -1
	var countA int64 = 0
	var cost int64 = 0

	for i := 0; i < int(n); i++ {
		if !gotScan.Scan() {
			return fmt.Errorf("candidate array too short")
		}
		val, err := strconv.ParseInt(gotScan.Text(), 10, 64)
		if err != nil {
			return fmt.Errorf("could not parse array element")
		}
		candArr[i] = val

		if candArr[i] < orig[i] {
			return fmt.Errorf("skill %d decreased from %d to %d", i, orig[i], candArr[i])
		}
		if candArr[i] > A {
			return fmt.Errorf("skill %d exceeded A (%d > %d)", i, candArr[i], A)
		}

		cost += candArr[i] - orig[i]

		if minVal == -1 || candArr[i] < minVal {
			minVal = candArr[i]
		}
		if candArr[i] == A {
			countA++
		}
	}

	if gotScan.Scan() {
		return fmt.Errorf("extra tokens in candidate output")
	}

	if cost > m {
		return fmt.Errorf("used %d coins, but m is %d", cost, m)
	}

	actualForce := countA*cf + minVal*cm
	if actualForce != expForce {
		return fmt.Errorf("array has force %d but printed force is %d", actualForce, expForce)
	}

	return nil
}

func runCase(bin, ref string, c Case) error {
	expected, err := runBinary(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	return checkCandidate(c.input, expected, got)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, ref, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}