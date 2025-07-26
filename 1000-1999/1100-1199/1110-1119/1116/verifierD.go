package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		n := rand.Intn(3) + 2 // N in [2,4]
		cmd := exec.Command(bin, strconv.Itoa(n))
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "run %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out.String())
		size := 1 << n
		if len(fields) != size*size {
			fmt.Fprintf(os.Stderr, "test %d: expected %d numbers, got %d\n", i+1, size*size, len(fields))
			os.Exit(1)
		}
		// ensure all fields parse as floats
		for _, f := range fields {
			if _, err := strconv.ParseFloat(f, 64); err != nil {
				fmt.Fprintf(os.Stderr, "test %d: invalid float %q\n", i+1, f)
				os.Exit(1)
			}
		}
	}
	fmt.Println("ok")
}
