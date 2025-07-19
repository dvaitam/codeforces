package main

import (
   "bufio"
   "fmt"
   "os"
)

// Pi represents a point or line parameters
type Pi struct { x, y int64 }

var op []int
var QA, QB []int64

func getans(A Pi, k, b int64) int64 {
   return A.x*k + A.y + b
}

func calc(A, B Pi) int64 {
   dety := A.y - B.y
   detx := B.x - A.x
   if dety%detx == 0 {
       return dety / detx
   }
   return dety/detx + 1
}

func check(A, B, C Pi) bool {
   return calc(A, B) > calc(B, C)
}

// Solve processes operations in [l, r] with initial total initot
func Solve(l, r int, initot int64, w *bufio.Writer) {
   V := make([]Pi, 0, r-l+3)
   V = append(V, Pi{0, 0})
   tot := initot
   var adds, addb int64
   for i := l; i <= r; i++ {
       if op[i] == 2 {
           k := QA[i]
           // add line
           x := Pi{tot, -adds*tot - addb}
           if x.y < V[len(V)-1].y {
               for len(V) > 1 && !check(V[len(V)-2], V[len(V)-1], x) {
                   V = V[:len(V)-1]
               }
               V = append(V, x)
           }
           tot += k
       } else {
           // type 3 operation
           adds += QB[i]
           addb += QA[i]
       }
       // query best line
       for len(V) > 1 && getans(V[len(V)-2], adds, addb) <= getans(V[len(V)-1], adds, addb) {
           V = V[:len(V)-1]
       }
       last := V[len(V)-1]
       fmt.Fprintf(w, "%d %d\n", last.x+1, getans(last, adds, addb))
   }
}

func main() {
   rdr := bufio.NewReader(os.Stdin)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()

   var n0, m0 int
   fmt.Fscan(rdr, &n0, &m0)
   origM := m0
   // one extra delimiter
   m := origM + 1
   op = make([]int, m+1)
   QA = make([]int64, m+1)
   QB = make([]int64, m+1)
   for i := 1; i <= origM; i++ {
       fmt.Fscan(rdr, &op[i], &QA[i])
       if op[i] == 3 {
           fmt.Fscan(rdr, &QB[i])
       }
   }
   op[m] = 1
   // process operations
   var tot int64 = int64(n0)
   for i := 1; i <= m; {
       tmptot := tot
       nxt := i
       for ; nxt <= m && op[nxt] != 1; nxt++ {
           if op[nxt] == 2 {
               tot += QA[nxt]
           }
       }
       Solve(i, nxt-1, tmptot, w)
       if nxt < m {
           fmt.Fprintln(w, "1 0")
       }
       tot += QA[nxt]
       i = nxt + 1
   }
}
