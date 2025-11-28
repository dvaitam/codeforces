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

func run(bin, input string) (string, error) {
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

// solveNaive implements the product transformation simulation directly.
func solveNaive(input string) string {
	var N, M, a, Q int64
	fmt.Sscan(input, &N, &M, &a, &Q)

	arr := make([]int64, N)
	for i := range arr {
		arr[i] = a % Q
	}

	for step := int64(0); step < M; step++ {
		// Simultaneous update
		nextArr := make([]int64, N)
		for i := 0; i < int(N)-1; i++ {
			nextArr[i] = (arr[i] * arr[i+1]) % Q
		}
		nextArr[N-1] = arr[N-1]
		arr = nextArr
	}

	var sb strings.Builder
	for i, val := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(val))
	}
	return sb.String()
}

func gcd(a, b int64) int64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func getOrder(a, Q int64) int64 {
	if gcd(a, Q) != 1 {
		return -1
	}
	curr := a % Q
	for i := int64(1); i <= Q+1; i++ {
		if curr == 1 {
			return i
		}
		curr = (curr * a) % Q
	}
	return -1
}

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func generateValidCase(rng *rand.Rand) string {
	// Constraints:
	// 7 <= Q (keep small for verification speed, e.g., up to 100)
	// 2 <= a
	// phi(a, Q) is prime P
	// 1 <= N, M < P
	for {
		Q := int64(rng.Intn(100) + 7)
		a := int64(rng.Intn(int(Q)-2) + 2)

		P := getOrder(a, Q)
		if P == -1 {
			continue
		}
		if !isPrime(P) {
			continue
		}

		// Found valid P
		if P <= 2 {
			// N, M must be < P, so if P=2, N,M=1. 
			// If P is larger we have more freedom.
			// Let's try to find slightly larger P if possible, or just accept.
			if rng.Float64() < 0.5 { continue } 
		}
		
		maxVal := P
		if maxVal > 1 {
			N := int64(rng.Intn(int(maxVal-1)) + 1)
			M := int64(rng.Intn(int(maxVal-1)) + 1)
			return fmt.Sprintf("%d %d %d %d\n", N, M, a, Q)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	
	// Run 20 tests
	for i := 1; i <= 20; i++ {
		in := generateValidCase(rng)
		want := solveNaive(in)
		got, err := run(candidate, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, in)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d wrong answer\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i, want, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
