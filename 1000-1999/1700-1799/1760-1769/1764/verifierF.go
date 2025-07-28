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
	ref := "refF.bin"
	if err := exec.Command("go", "build", "-o", ref, "1764F.go").Run(); err != nil {
		return "", err
	}
	return ref, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
	lines := 0
	sizes := make([]int, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(3) + 2
		sizes[i] = n
		fmt.Fprintln(&input, n)
		mat := make([][]int64, n)
		for r := 0; r < n; r++ {
			mat[r] = make([]int64, n)
		}
		for r := 0; r < n; r++ {
			for c := 0; c <= r; c++ {
				val := rand.Int63n(10)
				mat[r][c] = val
				mat[c][r] = val
			}
		}
		for r := 0; r < n; r++ {
			for c := 0; c <= r; c++ {
				if c > 0 {
					fmt.Fprint(&input, " ")
				}
				fmt.Fprint(&input, mat[r][c])
			}
			fmt.Fprintln(&input)
		}
		lines += n - 1
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
	for i := 0; i < lines; i++ {
		if !candScanner.Scan() || !refScanner.Scan() {
			fmt.Println("output missing line", i+1)
			os.Exit(1)
		}
		if candScanner.Text() != refScanner.Text() {
			fmt.Printf("mismatch on line %d\n", i+1)
			os.Exit(1)
		}
	}
	if candScanner.Scan() || refScanner.Scan() {
		fmt.Println("extra output")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
