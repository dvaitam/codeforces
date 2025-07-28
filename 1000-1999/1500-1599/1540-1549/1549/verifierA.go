package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runProg(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	primes := []int64{5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67, 71, 73, 79, 83, 89, 97, 101, 103, 107, 109, 113, 127, 131, 137, 139, 149, 151, 157, 163, 167, 173, 179, 181, 191, 193, 197, 199, 211, 223, 227, 229, 233, 239, 241, 251, 257, 263, 269, 271, 277, 281, 283, 293, 307, 311, 313, 317, 331, 337, 347, 349, 353, 359, 367, 373, 379, 383, 389, 397, 401, 409, 419, 421, 431, 433, 439, 443, 449, 457, 461, 463, 467, 479, 487, 491, 499, 503, 509, 521, 523, 541, 547, 557, 999999937}
	for idx, p := range primes {
		input := fmt.Sprintf("1\n%d\n", p)
		out, err := runProg(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		var a, b int64
		if n, _ := fmt.Sscan(out, &a, &b); n != 2 {
			fmt.Fprintf(os.Stderr, "case %d: expected two integers, got %q\n", idx+1, out)
			os.Exit(1)
		}
		if !(2 <= a && a < b && b <= p) {
			fmt.Fprintf(os.Stderr, "case %d: invalid range a=%d b=%d for p=%d\n", idx+1, a, b, p)
			os.Exit(1)
		}
		if p%a != p%b {
			fmt.Fprintf(os.Stderr, "case %d: condition failed for p=%d a=%d b=%d\n", idx+1, p, a, b)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
