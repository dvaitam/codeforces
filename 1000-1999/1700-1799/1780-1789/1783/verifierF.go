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

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genPerm(n int) []int {
	arr := rand.Perm(n)
	for i := range arr {
		arr[i]++
	}
	return arr
}

func genTest() []byte {
	n := rand.Intn(10) + 2 // 2..11 for broader testing
	a := genPerm(n)
	b := genPerm(n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	for i, v := range b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func parseInput(input string) (int, []int, []int) {
	r := strings.NewReader(input)
	var n int
	fmt.Fscan(r, &n)
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(r, &a[i])
	}
	b := make([]int, n)
	for i := range b {
		fmt.Fscan(r, &b[i])
	}
	return n, a, b
}

func validate(input string, output string, expectedK int) error {
	n, a, b := parseInput(input)
	
	// Parse output
	r := strings.NewReader(output)
	var k int
	if _, err := fmt.Fscan(r, &k); err != nil {
		return fmt.Errorf("failed to read k: %v", err)
	}
	
	if k != expectedK {
		return fmt.Errorf("expected %d operations, got %d", expectedK, k)
	}
	
	ops := make([]int, k)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(r, &ops[i]); err != nil {
			return fmt.Errorf("failed to read operation %d: %v", i+1, err)
		}
	}
	
	// Apply operations
	// Operation i: swap a[i] with a[pos] where a[pos] == i
	// Note: Problem uses 1-based values.
	// If we use 0-based arrays, value v corresponds to index v-1? No.
	// The problem says: "choose integer i from 1 to n".
	// "let x be integer such that a_x = i. Swap a_i with a_x."
	// So we are swapping the element at index i with the element at index x.
	// (Using 1-based indexing for description).
	
	// Helper to find index of value val
	indexOf := func(arr []int, val int) int {
		for idx, v := range arr {
			if v == val {
				return idx
			}
		}
		return -1
	}
	
	for _, op := range ops {
		if op < 1 || op > n {
			return fmt.Errorf("operation %d out of bounds", op)
		}
		
		// In 0-based array 'a':
		// We want to swap a[op-1] with a[x] where a[x] == op.
		// Wait, "let x be integer such that a_x = i". i is the operation value 'op'.
		// So x is the position where value 'op' is found.
		// We swap a[op] (index op-1? No. a_i means element at index i) with a[x].
		
		// Yes. Swap element at index (op-1) with element at index (indexOf(op)).
		
		idx1 := op - 1
		idx2 := indexOf(a, op)
		if idx2 == -1 {
			return fmt.Errorf("value %d not found in a", op) // Should not happen
		}
		a[idx1], a[idx2] = a[idx2], a[idx1]
		
		idx3 := indexOf(b, op)
		if idx3 == -1 {
			return fmt.Errorf("value %d not found in b", op)
		}
		b[idx1], b[idx3] = b[idx3], b[idx1]
	}
	
	// Check if sorted
	isSorted := func(arr []int) bool {
		for i := 0; i < n; i++ {
			if arr[i] != i+1 {
				return false
			}
		}
		return true
	}
	
	if !isSorted(a) {
		return fmt.Errorf("permutation a is not sorted after operations: %v", a)
	}
	if !isSorted(b) {
		return fmt.Errorf("permutation b is not sorted after operations: %v", b)
	}
	
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go <binary>")
		os.Exit(1)
	}
	cand := os.Args[1]
	ref := "./refF.bin"
	if err := exec.Command("go", "build", "-o", ref, "1783F.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rand.Seed(time.Now().UnixNano())
	
	const numTests = 100
	for i := 0; i < numTests; i++ {
		inputBytes := genTest()
		inputStr := string(inputBytes)
		
		// Run Reference to get K
		refOut, err := run(ref, inputBytes)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		
		var expectedK int
		if _, err := fmt.Fscan(strings.NewReader(refOut), &expectedK); err != nil {
			fmt.Printf("failed to parse reference output k: %v\nOutput:\n%s", err, refOut)
			os.Exit(1)
		}
		
		// Run Candidate
		candOut, err := run(cand, inputBytes)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", inputStr)
			os.Exit(1)
		}
		
		if err := validate(inputStr, candOut, expectedK); err != nil {
			fmt.Printf("wrong answer on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", inputStr)
			fmt.Println("expected K:", expectedK)
			fmt.Println("got:\n", candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}