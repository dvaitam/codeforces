package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, m int
   var y0, y1 int // unused y0,y1
   if _, err := fmt.Fscan(reader, &n, &m, &y0, &y1); err != nil {
       return
   }
   mice := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &mice[i])
   }
   cheese := make([]int, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(reader, &cheese[i])
   }
   // No cheese: all mice remain hungry
   if m == 0 {
       fmt.Println(n)
       return
   }
   // t[j]: count of mice tied between cheese j and j+1, j=0..m-2
   t := make([]int, m)
   // d2[j], c[j]: minimal doubled distance and count of forced mice for cheese j
   inf := int64(1) << 60
   d2 := make([]int64, m)
   c := make([]int, m)
   for j := 0; j < m; j++ {
       d2[j] = inf
       c[j] = 0
   }
   // Assign mice to forced or tie
   for _, x := range mice {
       // find first cheese >= x
       k := sort.SearchInts(cheese, x)
       if k > 0 && k < m {
           dl := x - cheese[k-1]
           if dl < 0 {
               dl = -dl
           }
           dr := cheese[k] - x
           if dl == dr {
               // tie between k-1 and k
               t[k-1]++
               continue
           }
       }
       var j int
       if k == 0 {
           j = 0
       } else if k == m {
           j = m - 1
       } else {
           dl := x - cheese[k-1]
           if dl < 0 {
               dl = -dl
           }
           dr := cheese[k] - x
           if dl < dr {
               j = k - 1
           } else {
               j = k
           }
       }
       dist2 := int64(x - cheese[j])
       if dist2 < 0 {
           dist2 = -dist2
       }
       dist2 *= 2
       if c[j] == 0 || dist2 < d2[j] {
           d2[j] = dist2
           c[j] = 1
       } else if dist2 == d2[j] {
           c[j]++
       }
   }
   // DP over cheeses
   // dpPrev[s]: max fed up to previous cheese, s=1 if prev segment assigned here
   dpPrev := [2]int64{0, -inf}
   dpCur := [2]int64{0, 0}
   // Iterate cheeses 0..m-2
   for j := 0; j < m-1; j++ {
       // reset dpCur
       dpCur[0], dpCur[1] = -inf, -inf
       for s := 0; s < 2; s++ {
           base := dpPrev[s]
           if base < 0 {
               continue
           }
           // consider assignment of segment j: n_s=1 -> to cheese j; 0 -> to cheese j+1
           for n_s := 0; n_s < 2; n_s++ {
               // compute fed for cheese j with L and R ties
               var Lcnt, Rcnt int
               if s == 1 && j > 0 {
                   Lcnt = t[j-1]
               }
               if n_s == 1 {
                   Rcnt = t[j]
               }
               // distances
               minDist := inf
               if c[j] > 0 {
                   minDist = d2[j]
               }
               var DL, DR int64
               if Lcnt > 0 {
                   DL = int64(cheese[j] - cheese[j-1])
                   if DL < minDist {
                       minDist = DL
                   }
               }
               if Rcnt > 0 {
                   DR = int64(cheese[j+1] - cheese[j])
                   if DR < minDist {
                       minDist = DR
                   }
               }
               // count fed
               fed := int64(0)
               if c[j] > 0 && d2[j] == minDist {
                   fed += int64(c[j])
               }
               if Lcnt > 0 && int64(cheese[j]-cheese[j-1]) == minDist {
                   fed += int64(Lcnt)
               }
               if Rcnt > 0 && int64(cheese[j+1]-cheese[j]) == minDist {
                   fed += int64(Rcnt)
               }
               // update
               val := base + fed
               if val > dpCur[n_s] {
                   dpCur[n_s] = val
               }
           }
       }
       dpPrev = dpCur
   }
   // final cheese m-1
   best := int64(0)
   for s := 0; s < 2; s++ {
       base := dpPrev[s]
       if base < 0 {
           continue
       }
       j := m - 1
       var Lcnt int
       if s == 1 && j > 0 {
           Lcnt = t[j-1]
       }
       // compute fed for cheese j
       minDist := inf
       if c[j] > 0 {
           minDist = d2[j]
       }
       if Lcnt > 0 {
           DL := int64(cheese[j] - cheese[j-1])
           if DL < minDist {
               minDist = DL
           }
       }
       fed := int64(0)
       if c[j] > 0 && d2[j] == minDist {
           fed += int64(c[j])
       }
       if Lcnt > 0 && int64(cheese[j]-cheese[j-1]) == minDist {
           fed += int64(Lcnt)
       }
       total := base + fed
       if total > best {
           best = total
       }
   }
   // hungry mice = n - best fed
   hungry := int64(n) - best
   if hungry < 0 {
       hungry = 0
   }
   fmt.Println(hungry)
}
