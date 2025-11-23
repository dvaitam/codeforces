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

type Query struct {
	c, l, r int
}

func main() {
	rand.Seed(time.Now().UnixNano())

	const numTestCases = 100
	inputBuf := new(bytes.Buffer)
	fmt.Fprintln(inputBuf, numTestCases)

	type TestCase struct {
		N, K, Q int
		Queries []Query
	}

	testCases := make([]TestCase, numTestCases)

	for i := 0; i < numTestCases; i++ {
		var n, k, q int
		var queries []Query
		var validQueries []Query

		for attempt := 0; attempt < 1000; attempt++ {
			n = rand.Intn(20) + 1
			k = rand.Intn(n) + 1 // k usually <= n
			// Generate random array A
			// Values should be somewhat relevant to k.
			// 0..k+5
			a := make([]int, n)
			for x := 0; x < n; x++ {
				// Bias towards values around k and low values for MEX
				if rand.Intn(3) == 0 {
					a[x] = rand.Intn(k + 2)
				} else {
					a[x] = rand.Intn(n + 2)
				}
			}

			// Find valid queries
			validQueries = []Query{}
			for l := 1; l <= n; l++ {
				for r := l; r <= n; r++ {
					sub := a[l-1 : r]
					
					// Check Min
					minVal := sub[0]
					for _, v := range sub {
						if v < minVal {
							minVal = v
						}
					}
					if minVal == k {
						validQueries = append(validQueries, Query{1, l, r})
					}

					// Check MEX
					present := make(map[int]bool)
					for _, v := range sub {
						present[v] = true
					}
					mex := 0
					for present[mex] {
						mex++
					}
					if mex == k {
						validQueries = append(validQueries, Query{2, l, r})
					}
				}
			}

			if len(validQueries) > 0 {
				break
			}
		}

		if len(validQueries) == 0 {
			// Fallback if extremely unlucky, though highly unlikely with 1000 attempts
			// Create a trivial case: [k, k, k], min=k
			n = 5
			k = 1
			queries = []Query{{1, 1, 2}}
		} else {
			q = rand.Intn(10) + 1
			if q > len(validQueries) {
				q = len(validQueries)
			}
			queries = make([]Query, q)
			perm := rand.Perm(len(validQueries))
			for j := 0; j < q; j++ {
				queries[j] = validQueries[perm[j]]
			}
		}

		testCases[i].N = n
		testCases[i].K = k
		testCases[i].Q = len(queries)
		testCases[i].Queries = queries

		fmt.Fprintf(inputBuf, "%d %d %d\n", n, k, len(queries))
		for _, qu := range queries {
			fmt.Fprintf(inputBuf, "%d %d %d\n", qu.c, qu.l, qu.r)
		}
	}

	var solutionBin string
	var cleanup func()

	if len(os.Args) > 1 {
		solutionBin = os.Args[1]
		if !strings.Contains(solutionBin, "/") {
			solutionBin = "./" + solutionBin
		}
		cleanup = func() {}
	} else {
		solutionSrc := "2157C.go"
		solutionBin = "./solution_c_temp"
		if _, err := os.Stat(solutionSrc); os.IsNotExist(err) {
			fmt.Printf("Error: %s not found.\n", solutionSrc)
			os.Exit(1)
		}
		cmdBuild := exec.Command("go", "build", "-o", solutionBin, solutionSrc)
		cmdBuild.Stderr = os.Stderr
		cmdBuild.Stdout = os.Stdout
		if err := cmdBuild.Run(); err != nil {
			fmt.Printf("Failed to build: %v\n", err)
			os.Exit(1)
		}
		cleanup = func() { os.Remove(solutionBin) }
	}
	defer cleanup()

	cmdRun := exec.Command(solutionBin)
	cmdRun.Stdin = inputBuf
	var outBuf bytes.Buffer
	cmdRun.Stdout = &outBuf
	cmdRun.Stderr = os.Stderr

	if err := cmdRun.Run(); err != nil {
		fmt.Printf("Failed to run solution: %v\n", err)
		os.Exit(1)
	}

	outputLines := strings.Split(strings.TrimSpace(outBuf.String()), "\n")
	var validLines []string
	for _, l := range outputLines {
		if strings.TrimSpace(l) != "" {
			validLines = append(validLines, strings.TrimSpace(l))
		}
	}

	if len(validLines) != numTestCases {
		fmt.Printf("Expected %d output lines, got %d\n", numTestCases, len(validLines))
		os.Exit(1)
	}

	passed := 0
	failed := 0
	for i, line := range validLines {
		parts := strings.Fields(line)
		if len(parts) != testCases[i].N {
			fmt.Printf("Test case %d FAILED: Expected length %d, got %d\n", i+1, testCases[i].N, len(parts))
			failed++
			continue
		}

		a := make([]int, testCases[i].N)
		parseErr := false
		for idx, p := range parts {
			val, err := strconv.Atoi(p)
			if err != nil {
				fmt.Printf("Test case %d FAILED: Invalid number '%s'\n", i+1, p)
				parseErr = true
				break
			}
			if val < 0 {
				fmt.Printf("Test case %d FAILED: Negative number %d\n", i+1, val)
				parseErr = true
				break
			}
			a[idx] = val
		}
		if parseErr {
			failed++
			continue
		}

		if err := verify(a, testCases[i].K, testCases[i].Queries); err != nil {
			fmt.Printf("Test case %d FAILED: %v\n", i+1, err)
			fmt.Printf("Input N=%d, K=%d\n", testCases[i].N, testCases[i].K)
			fmt.Printf("Queries: %v\n", testCases[i].Queries)
			fmt.Printf("Output: %v\n", a)
			failed++
			if failed >= 5 {
				break
			}
		} else {
			passed++
		}
	}

	fmt.Printf("Passed %d/%d test cases.\n", passed, numTestCases)
	if failed > 0 {
		os.Exit(1)
	}
}

func verify(a []int, k int, queries []Query) error {
	for _, q := range queries {
		// 0-indexed logic
		sub := a[q.l-1 : q.r] // Go slice is [start, end)
		
		if q.c == 1 {
			// min(sub) == k
			minVal := sub[0]
			for _, x := range sub {
				if x < minVal {
					minVal = x
				}
			}
			if minVal != k {
				return fmt.Errorf("Query %v failed: min is %d, expected %d", q, minVal, k)
			}
		} else {
			// MEX(sub) == k
			present := make(map[int]bool)
			for _, x := range sub {
				present[x] = true
			}
			mex := 0
			for {
				if !present[mex] {
					break
				}
				mex++
			}
			if mex != k {
				return fmt.Errorf("Query %v failed: MEX is %d, expected %d", q, mex, k)
			}
		}
	}
	return nil
}
