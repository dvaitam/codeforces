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
   var n int
   fmt.Fscan(in, &n)
   a := make([]int, n)
   b := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &a[i], &b[i])
   }
   // prev occurrence of same participant
   prevOcc := make([]int, n)
   // last index of participant
   lastIdx := make([]int, n+1)
   for i := range lastIdx {
       lastIdx[i] = -1
   }
   for i := 0; i < n; i++ {
       p := a[i]
       prevOcc[i] = lastIdx[p]
       lastIdx[p] = i
   }
   // prepare list of participants with their last bid positions
   type pair struct{ pos, x int }
   LB := make([]pair, 0, n)
   for x := 1; x <= n; x++ {
       if lastIdx[x] >= 0 {
           LB = append(LB, pair{lastIdx[x], x})
       }
   }
   // sort descending by pos
   sort.Slice(LB, func(i, j int) bool {
       return LB[i].pos > LB[j].pos
   })
   var q int
   fmt.Fscan(in, &q)
   excluded := make([]bool, n+1)
   exList := make([]int, 0, 100)
   for qi := 0; qi < q; qi++ {
       var k int
       fmt.Fscan(in, &k)
       exList = exList[:0]
       for j := 0; j < k; j++ {
           var x int
           fmt.Fscan(in, &x)
           if !excluded[x] {
               excluded[x] = true
               exList = append(exList, x)
           }
       }
       // find primary candidate
       idx0 := 0
       for idx0 < len(LB) && excluded[LB[idx0].x] {
           idx0++
       }
       if idx0 >= len(LB) {
           fmt.Fprintln(out, "0 0")
           // clear exclusions
           for _, x := range exList {
               excluded[x] = false
           }
           continue
       }
       x0 := LB[idx0].x
       pos1 := LB[idx0].pos
       // find best other participant
       idx1 := 0
       for idx1 < len(LB) && (excluded[LB[idx1].x] || LB[idx1].x == x0) {
           idx1++
       }
       pos2Other := -1
       if idx1 < len(LB) {
           pos2Other = LB[idx1].pos
       }
       // x0's previous bids
       pX := prevOcc[pos1]
       var winnerX, winnerPos int
       for {
           if pX <= pos2Other {
               if pos2Other < 0 {
                   winnerX = x0
                   winnerPos = pos1
               } else {
                   winnerX = LB[idx1].x
                   winnerPos = pos2Other
               }
               break
           }
           // pX > pos2Other implies same x0
           pos1 = pX
           pX = prevOcc[pos1]
       }
       fmt.Fprintln(out, winnerX, b[winnerPos])
       // clear exclusions
       for _, x := range exList {
           excluded[x] = false
       }
   }
}
