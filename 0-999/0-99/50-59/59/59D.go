package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, k int
   fmt.Fscan(in, &n)
   total := 3 * n
   rank := make([]int, total+1)
   // read results: i-th number is student in i-th place (1-based best)
   for i := 1; i <= total; i++ {
       var id int
       fmt.Fscan(in, &id)
       rank[id] = i
   }
   teams := make([][3]int, n)
   for i := 0; i < n; i++ {
       for j := 0; j < 3; j++ {
           fmt.Fscan(in, &teams[i][j])
       }
   }
   fmt.Fscan(in, &k)
   // find if k is ever a captain
   isCap := false
   capIdx := -1
   var mateA, mateB int
   for i := 0; i < n; i++ {
       team := teams[i]
       // find captain: minimal rank
       cap := team[0]
       for j := 1; j < 3; j++ {
           if rank[team[j]] < rank[cap] {
               cap = team[j]
           }
       }
       if cap == k {
           isCap = true
           capIdx = i
           // identify mates
           a, b := -1, -1
           for j := 0; j < 3; j++ {
               if teams[i][j] == k {
                   continue
               }
               if a == -1 {
                   a = teams[i][j]
               } else {
                   b = teams[i][j]
               }
           }
           // order by lexicographically minimal pick
           if a > b {
               a, b = b, a
           }
           mateA, mateB = a, b
           break
       }
   }
   var ans []int
   if !isCap {
       // k never captain: any order, lexicographically minimal
       ans = make([]int, 0, total-1)
       for i := 1; i <= total; i++ {
           if i == k {
               continue
           }
           ans = append(ans, i)
       }
       sort.Ints(ans)
   } else {
       // build setBefore and R
       takenBefore := make([]bool, total+1)
       for i := 0; i < capIdx; i++ {
           for j := 0; j < 3; j++ {
               takenBefore[teams[i][j]] = true
           }
       }
       sBefore := make([]int, 0, total)
       for i := 1; i <= total; i++ {
           if i == k {
               continue
           }
           if takenBefore[i] {
               sBefore = append(sBefore, i)
           }
       }
       sort.Ints(sBefore)
       // build R: unteamed before cap, excluding mates
       sR := make([]int, 0, total)
       for i := 1; i <= total; i++ {
           if i == k || i == mateA || i == mateB {
               continue
           }
           if !takenBefore[i] {
               sR = append(sR, i)
           }
       }
       sort.Ints(sR)
       // greedy merge
       idxB := 0 // index in sBefore
       haveA := false
       haveB := false
       // place mateA
       for !haveA {
           if idxB < len(sBefore) && sBefore[idxB] < mateA {
               ans = append(ans, sBefore[idxB])
               idxB++
           } else {
               ans = append(ans, mateA)
               haveA = true
           }
       }
       // place mateB
       for !haveB {
           if idxB < len(sBefore) && sBefore[idxB] < mateB {
               ans = append(ans, sBefore[idxB])
               idxB++
           } else {
               ans = append(ans, mateB)
               haveB = true
           }
       }
       // merge remaining sBefore and sR
       iB, iR := idxB, 0
       for iB < len(sBefore) && iR < len(sR) {
           if sBefore[iB] < sR[iR] {
               ans = append(ans, sBefore[iB]); iB++
           } else {
               ans = append(ans, sR[iR]); iR++
           }
       }
       for iB < len(sBefore) {
           ans = append(ans, sBefore[iB]); iB++
       }
       for iR < len(sR) {
           ans = append(ans, sR[iR]); iR++
       }
   }
   // output
   for i, v := range ans {
       if i > 0 {
           out.WriteByte(' ')
       }
       fmt.Fprint(out, v)
   }
   out.WriteByte('\n')
}
