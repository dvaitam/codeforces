package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   var x, y int64
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   fmt.Fscan(reader, &x, &y)
   n := len(s)
   // cost between fixed 0/1 pairs
   var costFixed int64
   var cnt0, cnt1 int64
   // collect fixed counts before each '?'
   zerosB := make([]int64, 0, n)
   onesB := make([]int64, 0, n)
   for i := 0; i < n; i++ {
       switch s[i] {
       case '0':
           costFixed += y * cnt1
           cnt0++
       case '1':
           costFixed += x * cnt0
           cnt1++
       case '?':
           zerosB = append(zerosB, cnt0)
           onesB = append(onesB, cnt1)
       }
   }
   total0 := cnt0
   total1 := cnt1
   K := len(zerosB)
   // compute f0 and f1 for each '?'
   f0 := make([]int64, K)
   f1 := make([]int64, K)
   for i := 0; i < K; i++ {
       zb := zerosB[i]
       ob := onesB[i]
       za := total0 - zb
       oa := total1 - ob
       // assign 0: interacts with fixed 1s before (y) and after (x)
       f0[i] = y*ob + x*oa
       // assign 1: interacts with fixed 0s before (x) and after (y)
       f1[i] = x*zb + y*za
   }
   // prefix sums
   pref0 := make([]int64, K+1)
   pref1 := make([]int64, K+1)
   for i := 0; i < K; i++ {
       pref0[i+1] = pref0[i] + f0[i]
       pref1[i+1] = pref1[i] + f1[i]
   }
   // minimize over t
   var best int64 = -1
   if x <= y {
       for t := 0; t <= K; t++ {
           // first t -> 0, rest -> 1
           var cost int64
           cost = pref0[t] + (pref1[K] - pref1[t]) + x*int64(t)*(int64(K)-int64(t))
           if best < 0 || cost < best {
               best = cost
           }
       }
   } else {
       for t := 0; t <= K; t++ {
           // first t -> 1, rest -> 0
           var cost int64
           cost = pref1[t] + (pref0[K] - pref0[t]) + y*int64(t)*(int64(K)-int64(t))
           if best < 0 || cost < best {
               best = cost
           }
       }
   }
   result := costFixed + best
   fmt.Println(result)
}
