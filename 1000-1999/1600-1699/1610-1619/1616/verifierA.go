package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const numTests = 100
const refSource = "1616A.go"

func buildRef() (string, error) {
	tmp := filepath.Join(os.TempDir(), "refA_bin")
	cmd := exec.Command("go", "build", "-o", tmp, refSource)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference: %v\n%s", err, out)
	}
	return tmp, nil
}

func runBinary(path string, in []byte) ([]byte, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(in)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.Bytes(), fmt.Errorf("%v: %s", err, stderr.String())
	}
	return out.Bytes(), nil
}

func genTest() []byte {
	n := rand.Intn(100) + 1
	arr := make([]int, n)
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "1\n%d\n", n)
	for i := range arr {
		arr[i] = rand.Intn(201) - 100
		if i > 0 {
			buf.WriteByte(' ')
		}
		fmt.Fprintf(&buf, "%d", arr[i])
	}
	buf.WriteByte('\n')
	return buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	comp := os.Args[1]
	refBin, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	for i := 0; i < numTests; i++ {
		input := genTest()
		want, err := runBinary(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference run error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(comp, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "competitor run error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if bytes.TrimSpace(got) == nil || !bytes.Equal(bytes.TrimSpace(want), bytes.TrimSpace(got)) {
			fmt.Printf("WA on test %d\ninput:\n%sexpected:\n%s\nreceived:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
