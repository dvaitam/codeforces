package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m int
   fmt.Fscan(in, &n, &m)
   var atkJiro, defJiro []int
   for i := 0; i < n; i++ {
       var pos string
       var s int
       fmt.Fscan(in, &pos, &s)
       if pos == "ATK" {
           atkJiro = append(atkJiro, s)
       } else {
           defJiro = append(defJiro, s)
       }
   }
   var C []int
   for i := 0; i < m; i++ {
       var s int
       fmt.Fscan(in, &s)
       C = append(C, s)
   }
   sort.Ints(atkJiro)
   sort.Ints(defJiro)
   sort.Ints(C)
   // Strategy 1: clear all cards and direct attack remaining
   damageFull := -1
   {
       used := make([]bool, len(C))
       ok := true
       // destroy DEF cards
       for _, d := range defJiro {
           idx := -1
           for j, v := range C {
               if !used[j] && v > d {
                   idx = j
                   break
               }
           }
           if idx < 0 {
               ok = false
               break
           }
           used[idx] = true
       }
       if ok {
           damage := 0
           // destroy ATK cards
           for _, a := range atkJiro {
               idx := -1
               for j, v := range C {
                   if !used[j] && v >= a {
                       idx = j
                       break
                   }
               }
               if idx < 0 {
                   ok = false
                   break
               }
               used[idx] = true
               damage += C[idx] - a
           }
           if ok {
               // direct attack with remaining cards
               for j, v := range C {
                   if !used[j] {
                       damage += v
                   }
               }
               damageFull = damage
           }
       }
   }
   // Strategy 2: match only ATK cards to maximize sum(C - a)
   // use Hungarian algorithm for maximum weight matching
   damageAtk := 0
   k := len(atkJiro)
   if k > 0 {
       mC := len(C)
       N := mC
       if k > N {
           N = k
       }
       // cost matrix for minimization: cost[i][j] = -(weight)
       const INF = 1000000000
       size := N + 1
       u := make([]int, size)
       v := make([]int, size)
       p := make([]int, size)
       way := make([]int, size)
       cost := make([][]int, size)
       for i := range cost {
           cost[i] = make([]int, size)
       }
       for i := 1; i <= N; i++ {
           for j := 1; j <= N; j++ {
               w := 0
               if i <= mC && j <= k && C[i-1] >= atkJiro[j-1] {
                   w = C[i-1] - atkJiro[j-1]
               }
               cost[i][j] = -w
           }
       }
       for i := 1; i <= N; i++ {
           p[0] = i
           j0 := 0
           minv := make([]int, size)
           used := make([]bool, size)
           for j := 0; j <= N; j++ {
               minv[j] = INF
           }
           for {
               used[j0] = true
               i0 := p[j0]
               delta := INF
               j1 := 0
               for j := 1; j <= N; j++ {
                   if !used[j] {
                       cur := cost[i0][j] - u[i0] - v[j]
                       if cur < minv[j] {
                           minv[j] = cur
                           way[j] = j0
                       }
                       if minv[j] < delta {
                           delta = minv[j]
                           j1 = j
                       }
                   }
               }
               for j := 0; j <= N; j++ {
                   if used[j] {
                       u[p[j]] += delta
                       v[j] -= delta
                   } else {
                       minv[j] -= delta
                   }
               }
               j0 = j1
               if p[j0] == 0 {
                   break
               }
           }
           for {
               j1 := way[j0]
               p[j0] = p[j1]
               j0 = j1
               if j0 == 0 {
                   break
               }
           }
       }
       // minimal cost is -v[0]
       damageAtk = -(-v[0])
   }
   ans := damageAtk
   if damageFull > ans {
       ans = damageFull
   }
   fmt.Println(ans)
