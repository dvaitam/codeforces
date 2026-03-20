package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func GetOperationsToBuildT(T [][]int, n int) [][2]int {
	opsReduce := [][2]int{}
	M := make([][]int, n)
	for i := 0; i < n; i++ {
		M[i] = make([]int, n)
		copy(M[i], T[i])
	}

	r := 0
	pivots := []int{}
	for c := 0; c < n; c++ {
		pivot := -1
		for i := r; i < n; i++ {
			if M[i][c] == 1 {
				pivot = i
				break
			}
		}
		if pivot != -1 {
			if pivot != r {
				opsReduce = append(opsReduce, [2]int{r, pivot})
				opsReduce = append(opsReduce, [2]int{pivot, r})
				opsReduce = append(opsReduce, [2]int{r, pivot})
				M[r], M[pivot] = M[pivot], M[r]
			}
			for k := 0; k < n; k++ {
				if k != r && M[k][c] == 1 {
					opsReduce = append(opsReduce, [2]int{k, r})
					for j := 0; j < n; j++ {
						M[k][j] ^= M[r][j]
					}
				}
			}
			pivots = append(pivots, c)
			r++
		}
	}

	opsBuild := [][2]int{}
	rowContains := make([]int, n)
	for i := 0; i < n; i++ {
		rowContains[i] = i
	}

	for i := 0; i < r; i++ {
		targetCol := pivots[i]
		k := i
		for ; k < n; k++ {
			if rowContains[k] == targetCol {
				break
			}
		}
		if k != i {
			opsBuild = append(opsBuild, [2]int{i, k})
			opsBuild = append(opsBuild, [2]int{k, i})
			opsBuild = append(opsBuild, [2]int{i, k})
			rowContains[i], rowContains[k] = rowContains[k], rowContains[i]
		}
	}

	pos := make([]int, n)
	for i := 0; i < n; i++ {
		pos[rowContains[i]] = i
	}

	isPivot := make([]bool, n)
	for _, p := range pivots {
		isPivot[p] = true
	}
	S := []int{}
	for i := 0; i < n; i++ {
		if !isPivot[i] {
			S = append(S, i)
		}
	}

	for i := 0; i < r; i++ {
		for _, j := range S {
			if M[i][j] == 1 {
				opsBuild = append(opsBuild, [2]int{i, pos[j]})
			}
		}
	}

	for k := r; k < n; k++ {
		opsBuild = append(opsBuild, [2]int{k, k})
	}

	for i := len(opsReduce) - 1; i >= 0; i-- {
		opsBuild = append(opsBuild, opsReduce[i])
	}

	return opsBuild
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return
	}
	n, _ := strconv.Atoi(scanner.Text())

	X := make([]int, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		X[i], _ = strconv.Atoi(scanner.Text())
	}

	Y := make([]int, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		Y[i], _ = strconv.Atoi(scanner.Text())
	}

	ops := [][2]int{}
	basisIdx := make([]int, 30)
	for i := 0; i < 30; i++ {
		basisIdx[i] = -1
	}
	isBasis := make([]bool, n)

	for b := 29; b >= 0; b-- {
		pivot := -1
		for i := 0; i < n; i++ {
			if !isBasis[i] && ((X[i]>>b)&1) == 1 {
				pivot = i
				break
			}
		}
		if pivot != -1 {
			basisIdx[b] = pivot
			isBasis[pivot] = true
			for i := 0; i < n; i++ {
				if i != pivot && ((X[i]>>b)&1) == 1 {
					ops = append(ops, [2]int{i, pivot})
					X[i] ^= X[pivot]
				}
			}
		}
	}

	for i := 0; i < n; i++ {
		temp := Y[i]
		for b := 29; b >= 0; b-- {
			if ((temp>>b)&1) == 1 {
				if basisIdx[b] != -1 {
					temp ^= X[basisIdx[b]]
				}
			}
		}
		if temp != 0 {
			fmt.Println("-1")
			return
		}
	}

	for i := 0; i < n; i++ {
		if !isBasis[i] {
			temp := Y[i]
			for b := 29; b >= 0; b-- {
				if ((temp>>b)&1) == 1 {
					ops = append(ops, [2]int{i, basisIdx[b]})
					temp ^= X[basisIdx[b]]
				}
			}
		}
	}

	B := []int{}
	BBit := []int{}
	for i := 0; i < n; i++ {
		if isBasis[i] {
			B = append(B, i)
			for b := 29; b >= 0; b-- {
				if basisIdx[b] == i {
					BBit = append(BBit, b)
					break
				}
			}
		}
	}

	k := len(B)
	for i := 0; i < k; i++ {
		for j := i + 1; j < k; j++ {
			if BBit[i] < BBit[j] {
				BBit[i], BBit[j] = BBit[j], BBit[i]
				B[i], B[j] = B[j], B[i]
			}
		}
	}

	TB := make([][]int, k)
	for i := 0; i < k; i++ {
		TB[i] = make([]int, k)
		temp := Y[B[i]]
		for j := 0; j < k; j++ {
			bit := BBit[j]
			if ((temp>>bit)&1) == 1 {
				TB[i][j] = 1
				temp ^= X[B[j]]
			}
		}
	}

	opsMatrix := GetOperationsToBuildT(TB, k)
	for _, op := range opsMatrix {
		ops = append(ops, [2]int{B[op[0]], B[op[1]]})
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	fmt.Fprintln(out, len(ops))
	for _, op := range ops {
		fmt.Fprintf(out, "%d %d\n", op[0]+1, op[1]+1)
	}
}