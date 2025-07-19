package main

import (
   "bufio"
   "fmt"
   "math"
   "os"
   "sort"
)

const P = 71

var (
   aData []int
   memo  map[uint64]float64
)

func getState(a, b, c int) int {
   if a > c {
       return ((c << 2) + b) << 2 + a
   }
   return ((a << 2) + b) << 2 + c
}

func setStateArr(x int) [3]int {
   return [3]int{(x >> 4) & 3, (x >> 2) & 3, x & 3}
}

func myhash(n int) uint64 {
   b := make([]int, n)
   copy(b, aData[:n])
   sort.Ints(b)
   var rtn uint64
   for i := 0; i < n; i++ {
       rtn = rtn*P + uint64(b[i])
   }
   return rtn
}

func dfs(n int) float64 {
   s := myhash(n)
   if f, ok := memo[s]; ok {
       return f
   }
   const Inf = math.Inf(1)
   c := [4]float64{0, Inf, Inf, Inf}
   an := aData[n-1]
   for i := 0; i < n-1; i++ {
       ai := aData[i]
       x := setStateArr(ai)
       y := setStateArr(an)
       // skip if only one block
       cntx := 0
       for _, v := range x {
           if v > 0 {
               cntx++
           }
       }
       if cntx == 1 {
           continue
       }
       nn := n
       // top layer complete
       if y[0] > 0 && y[1] > 0 && y[2] > 0 {
           nn = n + 1
           y = [3]int{0, 0, 0}
       }
       for j, vj := range x {
           if vj == 0 {
               continue
           }
           // skip invalid two-block middle or j==1 cases
           if cntx == 2 && (j == 1 || x[1] == 0) {
               continue
           }
           t := vj
           xjOld := x[j]
           x[j] = 0
           // compute new ai
           cntx2 := 0
           for _, v := range x {
               if v > 0 {
                   cntx2++
               }
           }
           var newAi int
           if cntx2 == 1 || x[1] == 0 {
               newAi = 0
           } else {
               newAi = getState(x[0], x[1], x[2])
           }
           aData[i] = newAi
           // try placing
           for k := 0; k < 3; k++ {
               if y[k] == 0 {
                   y[k] = t
                   newAn := getState(y[0], y[1], y[2])
                   // set top layer
                   if nn == n {
                       aData[n-1] = newAn
                   } else {
                       aData = append(aData, newAn)
                   }
                   c[t] = math.Min(c[t], dfs(nn))
                   // restore y
                   y[k] = 0
                   // restore aData length if appended
                   if nn != n {
                       aData = aData[:len(aData)-1]
                   } else {
                       // restore top layer
                       aData[n-1] = an
                   }
               }
           }
       }
       // restore aData[i] and x (x is local copy)
       aData[i] = ai
   }
       }
   }
   // no moves
   if c[1] == math.Inf(1) && c[2] == math.Inf(1) && c[3] == math.Inf(1) {
       memo[s] = 0
       return 0
   }
   p := 1.0 / 6.0
   if c[1] == math.Inf(1) {
       p += 1.0 / 3.0
       c[1] = 0
   }
   if c[2] == math.Inf(1) {
       p += 1.0 / 3.0
       c[2] = 0
   }
   if c[3] == math.Inf(1) {
       p += 1.0 / 6.0
       c[3] = 0
   }
   f := (c[1]/3.0 + c[2]/3.0 + c[3]/6.0 + 1.0) / (1.0 - p)
   memo[s] = f
   return f
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   fmt.Fscan(reader, &n)
   // preallocate with extra capacity for potential new layers
   aData = make([]int, n, 2*n+5)
   memo = make(map[uint64]float64)
   tag := map[byte]int{'G': 1, 'B': 2, 'R': 3}
   for i := 0; i < n; i++ {
       var str string
       fmt.Fscan(reader, &str)
       aData[i] = getState(tag[str[0]], tag[str[1]], tag[str[2]])
   }
   res := dfs(n)
   fmt.Printf("%.15f\n", res)
}
