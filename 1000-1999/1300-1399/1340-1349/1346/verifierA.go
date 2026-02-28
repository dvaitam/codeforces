package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) (string, []int64, []int64) {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	ns := make([]int64, t)
	ks := make([]int64, t)
	for i := 0; i < t; i++ {
		k := int64(rng.Intn(500) + 1)
		denom := int64(1) + k + k*k + k*k*k
		maxN1 := int64(1e9) / denom
		if maxN1 == 0 {
			maxN1 = 1
		}
		n1 := rng.Int63n(maxN1) + 1
		n := n1 * denom
		ns[i] = n
		ks[i] = k
		fmt.Fprintf(&sb, "%d %d\n", n, k)
	}
	return sb.String(), ns, ks
}

func check(output string, ns, ks []int64) error {
	r := strings.NewReader(output)
	for i, n := range ns {
		k := ks[i]
		var n1, n2, n3, n4 int64
		if _, err := fmt.Fscan(r, &n1, &n2, &n3, &n4); err != nil {
			return fmt.Errorf("test %d: failed to parse output: %v", i+1, err)
		}
		if n2 != k*n1 {
			return fmt.Errorf("test %d: n2=%d != k*n1=%d*%d=%d", i+1, n2, k, n1, k*n1)
		}
		if n3 != k*n2 {
			return fmt.Errorf("test %d: n3=%d != k*n2=%d*%d=%d", i+1, n3, k, n2, k*n2)
		}
		if n4 != k*n3 {
			return fmt.Errorf("test %d: n4=%d != k*n3=%d*%d=%d", i+1, n4, k, n3, k*n3)
		}
		if n1+n2+n3+n4 != n {
			return fmt.Errorf("test %d: sum=%d != n=%d", i+1, n1+n2+n3+n4, n)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		input, ns, ks := genCase(rng)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := check(got, ns, ks); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%sgot:\n%s\n", i+1, err, input, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
