package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

//
// Helper function to count mergesort calls for a given permutation
//
func countMergesortCalls(a []int) int {
    var count int
    
    var mergesort func(l, r int)
    mergesort = func(l, r int) {
        count++
        
        isSorted := true
        for i := l; i < r-1; i++ {
            if a[i] > a[i+1] {
                isSorted = false
                break
            }
        }
        
        if isSorted {
            return
        }
        
        if r - l <= 1 {
             return
        }

        mid := (l + r) / 2
        mergesort(l, mid)
        mergesort(mid, r)
    }
    
    mergesort(0, len(a))
    return count
}


//
// Original verifier functions needed for test generation and running binaries
//
type TestD struct {
	n, k int
}

func (t TestD) Input() string {
	return fmt.Sprintf("%d %d\n", t.n, t.k)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func genTests() []TestD {
	rand.Seed(4)
	tests := make([]TestD, 0, 100)
    // Add the failing test case
    tests = append(tests, TestD{n: 10, k: 15})
	for i := 0; i < 99; i++ {
		n := rand.Intn(20) + 1
		k := rand.Intn(2*n - 1)
		if k%2 == 0 {
			k++
			if k >= 2*n {
				k -= 2
			}
		}
		if k <= 0 {
			k = 1
		}
		tests = append(tests, TestD{n: n, k: k})
	}
	return tests
}

// Keep original oracle just for printing a valid answer on failure
func solveRec(l, r, k, L, R int, ans *[]int) {
	if k == 1 {
		for i := L; i <= R; i++ {
			*ans = append(*ans, i)
		}
		return
	}
	mid := (l + r - 1) >> 1
	Mid := (L + R + 2) >> 1
	leftCnt := mid - l + 1
	if 2*leftCnt-1 >= k-2 {
		solveRec(l, mid, k-2, Mid, R, ans)
		solveRec(mid+1, r, 1, L, Mid-1, ans)
	} else {
		solveRec(l, mid, 2*leftCnt-1, Mid, R, ans)
		solveRec(mid+1, r, k-2*leftCnt, L, Mid-1, ans)
	}
}
func expectedD(t TestD) string {
	n, k := t.n, t.k
	if k%2 == 0 || n*2 <= k {
		return "-1"
	}
	ans := make([]int, 0, n)
	solveRec(1, n, k, 1, n, &ans)
	var sb strings.Builder
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	return sb.String()
}


//
// New main verifier function
//
func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		gotRaw, err := run(bin, tc.Input())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n%s", i+1, err, gotRaw)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(gotRaw)
        
        if gotStr == "-1" {
            if tc.k%2 == 0 || tc.k > 2*tc.n-1 {
                continue
            } else {
                exp := expectedD(tc)
                fmt.Printf("test %d failed: incorrectly output -1\ninput:%sexpected: %s\ngot: %s\n", i+1, tc.Input(), exp, gotStr)
		        os.Exit(1)
            }
        }

        fields := strings.Fields(gotStr)
        if len(fields) != tc.n {
            fmt.Printf("test %d failed: wrong number of elements. expected %d, got %d\n", i+1, tc.n, len(fields))
            os.Exit(1)
        }
        
        perm := make([]int, tc.n)
        used := make(map[int]bool)
        for j, s := range fields {
            p, err := strconv.Atoi(s)
            if err != nil {
                fmt.Printf("test %d failed: could not parse element %s\n", i+1, s)
                os.Exit(1)
            }
            if p < 1 || p > tc.n || used[p] {
                 fmt.Printf("test %d failed: output is not a valid permutation of 1..%d\n", i+1, tc.n)
                 os.Exit(1)
            }
            used[p] = true
            perm[j] = p
        }

        callCount := countMergesortCalls(perm)

        if callCount != tc.k {
            exp := expectedD(tc)
            fmt.Printf("test %d failed: wrong number of mergesort calls. expected %d, got %d\ninput:%spermutation: %s\n(a valid answer would be: %s)\n", i+1, tc.k, callCount, tc.Input(), gotStr, exp)
            os.Exit(1)
        }
	}
	fmt.Println("All tests passed")
}