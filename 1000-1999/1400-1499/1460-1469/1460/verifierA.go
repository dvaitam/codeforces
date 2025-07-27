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

func buildRef() (string, error) {
	ref := "./refA.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1460A.go")
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
	cmd.Stdout = &out
	cmd.Stderr = new(bytes.Buffer)
	if err := cmd.Run(); err != nil {
		return out.String(), err
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 2  // 2..9
	m := rng.Intn(15) + 1 // 1..15
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sc := rng.Intn(401) + 100
		sm := rng.Intn(801) + 200
		sb.WriteString(fmt.Sprintf("%d %d\n", sc, sm))
	}
	for j := 0; j < m; j++ {
		vc := rng.Intn(200) + 1
		vm := rng.Intn(500) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", vc, vm))
	}
	for j := 0; j < m; j++ {
		old := rng.Intn(n)
		newSrv := rng.Intn(n)
		sb.WriteString(fmt.Sprintf("%d %d\n", old, newSrv))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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

	cases := []string{
		// deterministic simple cases
		"2 2\n100 200\n150 300\n1 1\n2 2\n0 1\n1 0\n",
		"3 1\n100 200\n200 300\n150 250\n1 2\n0 2\n",
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, genCase(rng))
	}

	for idx, in := range cases {
		exp, err := runBinary(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := runBinary(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, in, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
