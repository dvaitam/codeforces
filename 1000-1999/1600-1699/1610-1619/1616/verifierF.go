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
const refSource = "1616F.go"

func buildRef() (string, error) {
	tmp := filepath.Join(os.TempDir(), "refF_bin")
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
	n := rand.Intn(5) + 1
	m := rand.Intn(n*(n-1)/2 + 1)
	edges := make([][3]int, m)
	used := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		for {
			a := rand.Intn(n) + 1
			b := rand.Intn(n) + 1
			if a == b {
				continue
			}
			if a > b {
				a, b = b, a
			}
			if !used[[2]int{a, b}] {
				used[[2]int{a, b}] = true
				edges[i] = [3]int{a, b, rand.Intn(4) - 1}
				break
			}
		}
	}
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "1\n%d %d\n", n, m)
	for i := 0; i < m; i++ {
		fmt.Fprintf(&buf, "%d %d %d\n", edges[i][0], edges[i][1], edges[i][2])
	}
	return buf.Bytes()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
		if !bytes.Equal(bytes.TrimSpace(want), bytes.TrimSpace(got)) {
			fmt.Printf("WA on test %d\ninput:\n%sexpected:\n%s\nreceived:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
