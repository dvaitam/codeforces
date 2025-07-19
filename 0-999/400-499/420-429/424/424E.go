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

// dfs computes expected moves for state of size n
func dfs(n int) float64 {
   s := myhash(n)
   if f, ok := memo[s]; ok {
       return f
   }
   const Inf = math.Inf(1)
   // c[1..3] stores minimal expected for removing each color
   c := [4]float64{0, Inf, Inf, Inf}
   an := aData[n-1]
   // iterate over layers except top
   for i := 0; i < n-1; i++ {
       ai := aData[i]
       x0 := setStateArr(ai)
       // count blocks in layer i
       cntx := 0
       for _, v := range x0 {
           if v > 0 {
               cntx++
           }
       }
       if cntx <= 1 {
           continue
       }
       // top layer state
       y0 := setStateArr(an)
       nn := n
       if y0[0] > 0 && y0[1] > 0 && y0[2] > 0 {
           nn = n + 1
       }
       // try taking each block from layer i
       for j, vj := range x0 {
           if vj == 0 {
               continue
           }
           if cntx == 2 && (j == 1 || x0[1] == 0) {
               continue
           }
           // copy x and remove block j
           x := x0
           x[j] = 0
           // compute new aData[i]
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
           // prepare y
           y0 := setStateArr(an)
           if nn > n {
               y0 = [3]int{0, 0, 0}
           }
           // place removed block on top
           for k := 0; k < 3; k++ {
               if y0[k] != 0 {
                   continue
               }
               y := y0
               y[k] = vj
               newAn := getState(y[0], y[1], y[2])
               if nn == n {
                   aData[n-1] = newAn
                   c[vj] = math.Min(c[vj], dfs(nn))
                   aData[n-1] = an
               } else {
                   aData = append(aData, newAn)
                   c[vj] = math.Min(c[vj], dfs(nn))
                   aData = aData[:len(aData)-1]
               }
           }
           aData[i] = ai
       }
   }
   // no moves available
   if c[1] == Inf && c[2] == Inf && c[3] == Inf {
       memo[s] = 0
       return 0
   }
   p := 1.0/6.0
   if c[1] == Inf {
       p += 1.0/3.0
       c[1] = 0
   }
   if c[2] == Inf {
       p += 1.0/3.0
       c[2] = 0
   }
   if c[3] == Inf {
       p += 1.0/6.0
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
