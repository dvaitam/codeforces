package main

import (
	"bytes"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

var inputs = []string{
	"4167792762229302596005813\n",
	"5023266485352224530541666346579718188045725802556947855902702063768930526665237598287495728" +
		"2186920039740245544313021979167491414627648054421626450903732301970386214502290904360792618" +
		"5591029614599889902115472391135622402044979347133959392884686037208893694733655782993294168" +
		"1679738555852317096830120847236770822731988661111203691013036774095229665675217827154840019" +
		"9277276899311984129170278649605877582438144407974816241674549565633361834348720814779487433" +
		"7933873576016717726298883519261055062303842274145012056670644839715140659887936321934474824" +
		"687778512706909988484451300384818197143498259061041\n",
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1812F-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", path, "1812F.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseOutput(out string) (*big.Int, error) {
	val := new(big.Int)
	_, ok := val.SetString(strings.TrimSpace(out), 10)
	if !ok {
		return nil, fmt.Errorf("invalid integer output: %q", out)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, input := range inputs {
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expVal, err := parseOutput(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		if expVal.Cmp(gotVal) != 0 {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %s got %s\n", idx+1, expVal.String(), gotVal.String())
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(inputs))
}
