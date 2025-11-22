package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	n int
	s string
}

func buildRef() (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "refJ.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "2068J.go"))
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
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d\n%s\n", tc.n, tc.s)
	}
	return sb.String()
}

func parseOutputs(out string, t int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	res := make([]string, t)
	for i, f := range fields {
		f = strings.ToUpper(f)
		if f != "YES" && f != "NO" {
			return nil, fmt.Errorf("invalid verdict %q", f)
		}
		res[i] = f
	}
	return res, nil
}

func randomCase(r *rand.Rand) testCase {
	n := r.Intn(100) + 1
	total := 2 * n
	chars := make([]byte, total)
	wLeft := n
	for i := 0; i < total; i++ {
		remain := total - i
		if wLeft == remain {
			chars[i] = 'W'
			wLeft--
			continue
		}
		if wLeft == 0 {
			chars[i] = 'R'
			continue
		}
		if r.Intn(2) == 0 {
			chars[i] = 'W'
			wLeft--
		} else {
			chars[i] = 'R'
		}
	}
	return testCase{n: n, s: string(chars)}
}

func genCases() []testCase {
	cases := []testCase{
		{n: 1, s: "RW"}, // YES
		{n: 1, s: "WR"}, // NO
		{n: 4, s: "WRRWWWRR"},
		{n: 3, s: "WWWRRR"},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN := 0
	for _, tc := range cases {
		totalN += tc.n
	}
	for len(cases) < 80 && totalN < 15000 {
		tc := randomCase(r)
		if totalN+tc.n > 200000 {
			break
		}
		totalN += tc.n
		cases = append(cases, tc)
	}
	return cases
}

func check(bin, ref string, cases []testCase) error {
	input := buildInput(cases)

	refOut, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	expect, err := parseOutputs(refOut, len(cases))
	if err != nil {
		return fmt.Errorf("reference output invalid: %v", err)
	}

	out, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	got, err := parseOutputs(out, len(cases))
	if err != nil {
		return err
	}
	for i := range expect {
		if expect[i] != got[i] {
			return fmt.Errorf("case %d mismatch: expected %s got %s", i+1, expect[i], got[i])
		}
	}
	return nil
}

func main() {
	exitCode := 0
	cleanup := func() {}
	defer func() {
		cleanup()
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierJ.go /path/to/binary")
		exitCode = 1
		return
	}
	bin := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	cleanup = func() { _ = os.Remove(ref) }

	cases := genCases()
	if err := check(bin, ref, cases); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, buildInput(cases))
		exitCode = 1
		return
	}
	fmt.Printf("All %d test cases passed\n", len(cases))
}
