package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

type pair struct{ x, y int }

func genA() {
	const T = 100
	f, _ := os.Create("testcasesA.txt")
	defer f.Close()
	fmt.Fprintln(f, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(9) + 1
		m := rand.Intn(9) + 1
		fmt.Fprintf(f, "%d %d\n", n, m)
		permA := rand.Perm(9)
		for j := 0; j < n; j++ {
			if j > 0 {
				fmt.Fprint(f, " ")
			}
			fmt.Fprint(f, permA[j]+1)
		}
		fmt.Fprintln(f)
		permB := rand.Perm(9)
		for j := 0; j < m; j++ {
			if j > 0 {
				fmt.Fprint(f, " ")
			}
			fmt.Fprint(f, permB[j]+1)
		}
		fmt.Fprintln(f)
	}
}

func genB() {
	const T = 100
	f, _ := os.Create("testcasesB.txt")
	defer f.Close()
	fmt.Fprintln(f, T)
	for i := 0; i < T; i++ {
		n := rand.Intn(10) + 1
		k := rand.Intn(5) + 1
		fmt.Fprintf(f, "%d %d\n", n, k)
		for j := 0; j < n; j++ {
			if j > 0 {
				fmt.Fprint(f, " ")
			}
			fmt.Fprint(f, rand.Intn(100)-50)
		}
		fmt.Fprintln(f)
	}
}

func genC() {
	const T = 100
	f, _ := os.Create("testcasesC.txt")
	defer f.Close()
	fmt.Fprintln(f, T)
	for i := 0; i < T; i++ {
		q := rand.Intn(4) + 1
		fmt.Fprintln(f, q)
		for j := 0; j < q; j++ {
			fmt.Fprintln(f, rand.Int63n(40)+1)
		}
	}
}

func genD() {
	const T = 100
	f, _ := os.Create("testcasesD.txt")
	defer f.Close()
	fmt.Fprintln(f, T)
	for t := 0; t < T; t++ {
		n := rand.Intn(4) + 2
		p := rand.Perm(n)
		b := make([]int, n)
		for i, v := range p {
			b[v] = i
		}
		pos0 := b[0]
		fmt.Fprintln(f, n)
		for i := 0; i < n; i++ {
			fmt.Fprintln(f, p[i]^pos0)
		}
		for j := 0; j < n; j++ {
			fmt.Fprintln(f, p[0]^b[j])
		}
	}
}

func genE() {
	const T = 100
	f, _ := os.Create("testcasesE.txt")
	defer f.Close()
	fmt.Fprintln(f, T)
	for t := 0; t < T; t++ {
		n := rand.Intn(8) + 1
		fmt.Fprintln(f, n)
		for i := 0; i < n; i++ {
			fmt.Fprintf(f, "%d %d\n", rand.Intn(5), rand.Intn(5))
		}
	}
}

func genF() {
	const T = 100
	f, _ := os.Create("testcasesF.txt")
	defer f.Close()
	fmt.Fprintln(f, T)
	for i := 0; i < T; i++ {
		fmt.Fprintln(f, rand.Intn(200)+1)
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	genA()
	genB()
	genC()
	genD()
	genE()
	genF()
}
