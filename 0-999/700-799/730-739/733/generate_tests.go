package main

import (
	"fmt"
	"math/rand"
	"os"
)

func randString(n int) string {
	letters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func genA() {
	f, _ := os.Create("testcasesA.txt")
	defer f.Close()
	fmt.Fprintln(f, 100)
	for i := 0; i < 100; i++ {
		l := rand.Intn(20) + 1
		fmt.Fprintln(f, randString(l))
	}
}

func genB() {
	f, _ := os.Create("testcasesB.txt")
	defer f.Close()
	fmt.Fprintln(f, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 1
		fmt.Fprintln(f, n)
		for j := 0; j < n; j++ {
			l := rand.Intn(5) + 1
			r := rand.Intn(5) + 1
			fmt.Fprintf(f, "%d %d\n", l, r)
		}
	}
}

func genC() {
	f, _ := os.Create("testcasesC.txt")
	defer f.Close()
	fmt.Fprintln(f, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(6) + 1
		fmt.Fprintln(f, n)
		a := make([]int, n)
		sum := 0
		for j := 0; j < n; j++ {
			a[j] = rand.Intn(5) + 1
			sum += a[j]
		}
		for j := 0; j < n; j++ {
			if j > 0 {
				fmt.Fprint(f, " ")
			}
			fmt.Fprint(f, a[j])
		}
		fmt.Fprintln(f)
		k := rand.Intn(n) + 1
		fmt.Fprintln(f, k)
		remain := sum
		for j := 0; j < k; j++ {
			var val int
			if j == k-1 {
				val = remain
			} else {
				val = rand.Intn(remain-(k-j-1)) + 1
			}
			remain -= val
			if j > 0 {
				fmt.Fprint(f, " ")
			}
			fmt.Fprint(f, val)
		}
		fmt.Fprintln(f)
	}
}

func genD() {
	f, _ := os.Create("testcasesD.txt")
	defer f.Close()
	fmt.Fprintln(f, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 1
		fmt.Fprintln(f, n)
		for j := 0; j < n; j++ {
			a := rand.Intn(20) + 1
			b := rand.Intn(20) + 1
			c := rand.Intn(20) + 1
			fmt.Fprintf(f, "%d %d %d\n", a, b, c)
		}
	}
}

func genE() {
	f, _ := os.Create("testcasesE.txt")
	defer f.Close()
	fmt.Fprintln(f, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 1
		fmt.Fprintln(f, n)
		s := make([]byte, n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				s[j] = 'U'
			} else {
				s[j] = 'D'
			}
		}
		fmt.Fprintln(f, string(s))
	}
}

func genF() {
	f, _ := os.Create("testcasesF.txt")
	defer f.Close()
	fmt.Fprintln(f, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(4) + 2
		maxEdges := n * (n - 1) / 2
		m := (n - 1) + rand.Intn(maxEdges-(n-1)+1)
		if m > 8 {
			m = 8
		}
		fmt.Fprintf(f, "%d %d\n", n, m)
		w := make([]int, m)
		c := make([]int, m)
		for j := 0; j < m; j++ {
			w[j] = rand.Intn(10) + 1
		}
		for j := 0; j < m; j++ {
			c[j] = rand.Intn(5) + 1
		}
		for j := 0; j < m; j++ {
			if j > 0 {
				fmt.Fprint(f, " ")
			}
			fmt.Fprint(f, w[j])
		}
		fmt.Fprintln(f)
		for j := 0; j < m; j++ {
			if j > 0 {
				fmt.Fprint(f, " ")
			}
			fmt.Fprint(f, c[j])
		}
		fmt.Fprintln(f)
		edges := make([][2]int, 0, m)
		for len(edges) < m {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			if u == v {
				continue
			}
			ok := true
			for _, e := range edges {
				if (e[0] == u && e[1] == v) || (e[0] == v && e[1] == u) {
					ok = false
					break
				}
			}
			if ok {
				edges = append(edges, [2]int{u, v})
			}
		}
		for j := 0; j < m; j++ {
			fmt.Fprintf(f, "%d %d\n", edges[j][0], edges[j][1])
		}
		s := rand.Intn(10)
		fmt.Fprintln(f, s)
	}
}

func main() {
	rand.Seed(42)
	genA()
	genB()
	genC()
	genD()
	genE()
	genF()
}
