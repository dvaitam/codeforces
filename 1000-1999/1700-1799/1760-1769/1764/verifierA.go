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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	refBin := "refA.bin"
	if err := exec.Command("go", "build", "-o", refBin, "1764A.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rand.Seed(time.Now().UnixNano())
	const t = 100
	var input bytes.Buffer
	fmt.Fprintln(&input, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(20) + 1
		fmt.Fprintln(&input, n)
		for j := 0; j < n; j++ {
			fmt.Fprint(&input, rand.Intn(n)+1)
			if j+1 < n {
				fmt.Fprint(&input, " ")
			}
		}
		fmt.Fprintln(&input)
	}

	candCmd := exec.Command(candidate)
	candCmd.Stdin = bytes.NewReader(input.Bytes())
	candOut, candErr := candCmd.CombinedOutput()
	if candErr != nil {
		fmt.Println("candidate run error:", candErr)
		os.Exit(1)
	}

	refCmd := exec.Command("./" + refBin)
	refCmd.Stdin = bytes.NewReader(input.Bytes())
	refOut, refErr := refCmd.CombinedOutput()
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
