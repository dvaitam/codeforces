package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

func buildRef() (string, error) {
	ref := "refH.bin"
	if err := exec.Command("go", "build", "-o", ref, "1764H.go").Run(); err != nil {
		return "", err
	}
	return ref, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rand.Seed(time.Now().UnixNano())
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 1
		m := rand.Intn(5) + 1
		k := rand.Intn(m) + 1
		fmt.Fprintf(&input, "%d %d %d\n", n, m, k)
		for j := 0; j < m; j++ {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			fmt.Fprintf(&input, "%d %d\n", l, r)
		}
	}

	run := func(bin string) ([]byte, error) {
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		return cmd.CombinedOutput()
	}

	candOut, candErr := run(cand)
	if candErr != nil {
		fmt.Println("candidate run error:", candErr)
		os.Exit(1)
	}
	refOut, refErr := run("./" + ref)
	if refErr != nil {
		fmt.Println("reference run error:", refErr)
		os.Exit(1)
	}

	candScanner := bufio.NewScanner(bytes.NewReader(candOut))
	refScanner := bufio.NewScanner(bytes.NewReader(refOut))
	for i := 0; i < t; i++ {
		if !candScanner.Scan() || !refScanner.Scan() {
			fmt.Println("output missing for test", i+1)
			os.Exit(1)
		}
		if candScanner.Text() != refScanner.Text() {
			fmt.Printf("mismatch on test %d: expected %q got %q\n", i+1, refScanner.Text(), candScanner.Text())
			os.Exit(1)
		}
	}
	if candScanner.Scan() || refScanner.Scan() {
		fmt.Println("extra output")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
