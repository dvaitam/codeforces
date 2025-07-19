package main

import (
   "bufio"
   "fmt"
   "os"
)

// Pair holds a grid coordinate
type Pair struct{ r, c int }

var N, M int
var Mp [2][][]byte
var now [][]byte
var Ans [2][]Pair
var nowK int

// changeOnce ensures that now[n][m] becomes k by recursive flips
func changeOnce(n, m int, k byte) {
   if now[n][m] == k {
       return
   }
   if k == 'L' && now[n][m] == 'U' {
       changeOnce(n, m+1, 'U')
       // apply flip to make 2x2 block horizontal
       now[n][m], now[n+1][m] = 'L', 'L'
       now[n][m+1], now[n+1][m+1] = 'R', 'R'
       Ans[nowK] = append(Ans[nowK], Pair{n, m})
       return
   }
   if k == 'U' && now[n][m] == 'L' {
       changeOnce(n+1, m, 'L')
       // apply flip to make 2x2 block vertical
       now[n][m], now[n][m+1] = 'U', 'U'
       now[n+1][m], now[n+1][m+1] = 'D', 'D'
       Ans[nowK] = append(Ans[nowK], Pair{n, m})
       return
   }
   // unexpected case
   panic(fmt.Sprintf("changeOnce invalid at (%d,%d): have %c, want %c", n, m, now[n][m], k))
}

// change traverses grid to convert all to vertical 'U' at (n,m)
func change(n, m int) {
   if n > N || (N%2 == 1 && n == N) {
       return
   }
   if m > M {
       change(n+2, 1)
       return
   }
   if now[n][m] == 'U' {
       change(n, m+1)
       return
   }
   // now[n][m] == 'L'
   changeOnce(n+1, m, 'L')
   // apply flip to make vertical at (n,m)
   now[n][m], now[n][m+1] = 'U', 'U'
   now[n+1][m], now[n+1][m+1] = 'D', 'D'
   Ans[nowK] = append(Ans[nowK], Pair{n, m})
   change(n, m+2)
}

// calc applies changes to Mp[k] grid
func calc(k int) {
   // copy Mp[k] to now
   for i := 1; i <= N; i++ {
       copy(now[i][1:], Mp[k][i][1:])
   }
   nowK = k
   change(1, 1)
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   fmt.Fscan(in, &N, &M)
   // allocate grids
   for k := 0; k < 2; k++ {
       Mp[k] = make([][]byte, N+2)
       for i := 0; i < N+2; i++ {
           Mp[k][i] = make([]byte, M+2)
       }
   }
   now = make([][]byte, N+2)
   for i := 0; i < N+2; i++ {
       now[i] = make([]byte, M+2)
   }
   // read two grids
   for k := 0; k < 2; k++ {
       for i := 1; i <= N; i++ {
           var s string
           fmt.Fscan(in, &s)
           for j := 1; j <= M; j++ {
               Mp[k][i][j] = s[j-1]
           }
       }
   }
   // perform transformations
   calc(0)
   calc(1)
   // output
   total := len(Ans[0]) + len(Ans[1])
   fmt.Fprintln(out, total)
   // first sequence
   for _, p := range Ans[0] {
       fmt.Fprintln(out, p.r, p.c)
   }
   // reverse second sequence
   for i := len(Ans[1]) - 1; i >= 0; i-- {
       p := Ans[1][i]
       fmt.Fprintln(out, p.r, p.c)
   }
}
