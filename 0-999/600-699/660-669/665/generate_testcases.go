package main

import (
	"fmt"
	"math/rand"
	"os"
)

func genA() {
	rand.Seed(1)
	f, _ := os.Create("testcasesA.txt")
	defer f.Close()
	const dayStart = 5 * 60
	const dayEnd = 23*60 + 59
	t := 100
	fmt.Fprintln(f, t)
	for i := 0; i < t; i++ {
		a := rand.Intn(120) + 1
		ta := rand.Intn(120) + 1
		b := rand.Intn(120) + 1
		tb := rand.Intn(120) + 1
		maxTrips := (dayEnd-dayStart)/a + 1
		start := dayStart + rand.Intn(maxTrips)*a
		hh := start / 60
		mm := start % 60
		fmt.Fprintf(f, "%d %d %d %d %d %d\n", a, ta, b, tb, hh, mm)
	}
}

func randPerm(n int) []int {
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = i + 1
	}
	rand.Shuffle(n, func(i, j int) { arr[i], arr[j] = arr[j], arr[i] })
	return arr
}

func pickDistinct(k, m int) []int {
	arr := randPerm(k)
	return arr[:m]
}

func genB() {
	rand.Seed(2)
	f, _ := os.Create("testcasesB.txt")
	defer f.Close()
	t := 100
	fmt.Fprintln(f, t)
	for i := 0; i < t; i++ {
		k := rand.Intn(8) + 1
		m := rand.Intn(k) + 1
		n := rand.Intn(5) + 1
		row := randPerm(k)
		fmt.Fprintf(f, "%d %d %d ", n, m, k)
		for j := 0; j < k; j++ {
			fmt.Fprintf(f, "%d ", row[j])
		}
		for c := 0; c < n; c++ {
			order := pickDistinct(k, m)
			for j := 0; j < m; j++ {
				fmt.Fprintf(f, "%d ", order[j])
			}
		}
		fmt.Fprintln(f)
	}
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte('a' + rand.Intn(26))
	}
	return string(b)
}

func genC() {
	rand.Seed(3)
	f, _ := os.Create("testcasesC.txt")
	defer f.Close()
	t := 100
	fmt.Fprintln(f, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(50) + 1
		s := randString(n)
		fmt.Fprintln(f, s)
	}
}

func genD() {
	rand.Seed(4)
	f, _ := os.Create("testcasesD.txt")
	defer f.Close()
	t := 100
	fmt.Fprintln(f, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(30) + 1
		fmt.Fprintf(f, "%d ", n)
		for j := 0; j < n; j++ {
			v := rand.Intn(100) + 1
			fmt.Fprintf(f, "%d ", v)
		}
		fmt.Fprintln(f)
	}
}

func genE() {
	rand.Seed(5)
	f, _ := os.Create("testcasesE.txt")
	defer f.Close()
	t := 100
	fmt.Fprintln(f, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(50) + 1
		k := rand.Intn(1000) + 1
		fmt.Fprintf(f, "%d %d ", n, k)
		for j := 0; j < n; j++ {
			v := rand.Intn(1000)
			fmt.Fprintf(f, "%d ", v)
		}
		fmt.Fprintln(f)
	}
}

func genF() {
	rand.Seed(6)
	f, _ := os.Create("testcasesF.txt")
	defer f.Close()
	t := 100
	fmt.Fprintln(f, t)
	for i := 0; i < t; i++ {
		n := rand.Int63n(1000000) + 1
		fmt.Fprintln(f, n)
	}
}

func main() {
	genA()
	genB()
	genC()
	genD()
	genE()
	genF()
}
