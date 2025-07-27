package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // odd number cannot be fully paired
   if n%2 == 1 {
       outNo()
       return
   }
   // read cost matrix
   cost := make([][]int64, n)
   for i := 0; i < n; i++ {
       cost[i] = make([]int64, n)
   }
   for i := 0; i < n; i++ {
       for cnt, j := 0, 0; j < n; j++ {
           if j == i {
               continue
           }
           var c int64
           fmt.Fscan(reader, &c)
           cost[i][j] = c
           cnt++
       }
   }
   // build preference lists
   prefs := make([][]int, n)
   invRank := make([][]int, n)
   for i := 0; i < n; i++ {
       arr := make([]int, 0, n-1)
       for j := 0; j < n; j++ {
           if j != i {
               arr = append(arr, j)
           }
       }
       sort.Slice(arr, func(a, b int) bool {
           return cost[i][arr[a]] < cost[i][arr[b]]
       })
       prefs[i] = arr
       inv := make([]int, n)
       for idx, j := range arr {
           inv[j] = idx
       }
       invRank[i] = inv
   }
   // Phase 1: proposal and list reduction
   queue := make([]int, 0, n)
   inQueue := make([]bool, n)
   for i := 0; i < n; i++ {
       queue = append(queue, i)
       inQueue[i] = true
   }
   for len(queue) > 0 {
       i := queue[0]
       queue = queue[1:]
       inQueue[i] = false
       if len(prefs[i]) == 0 {
           outNo()
           return
       }
       j := prefs[i][0]
       // find position of i in prefs[j]
       r := invRank[j][i]
       // eliminate all k worse than i in j's list
       for idx := r + 1; idx < len(prefs[j]); idx++ {
           k := prefs[j][idx]
           // remove j from prefs[k]; if head removed, requeue k
           removedHead := removePref(&prefs, &invRank, k, j)
           if len(prefs[k]) == 0 {
               outNo()
               return
           }
           if removedHead && !inQueue[k] {
               queue = append(queue, k)
               inQueue[k] = true
           }
       }
       // trim prefs[j]
       prefs[j] = prefs[j][:r+1]
       // no need to update invRank[j] for kept elements
   }
   // Phase 2: eliminate rotations
   for {
       // find person with list length >1
       s := -1
       for i := 0; i < n; i++ {
           if len(prefs[i]) > 1 {
               s = i
               break
           }
       }
       if s == -1 {
           break
       }
       // find rotation
       P := []int{}
       Q := []int{}
       p := s
       for {
           P = append(P, p)
           q := prefs[p][1]
           Q = append(Q, q)
           // find next p
           pos := invRank[q][p]
           p = prefs[q][pos+1]
           if p == s {
               break
           }
       }
       r := len(P)
       // eliminate rotation
       for k := 0; k < r; k++ {
           pk := P[k]
           qk := Q[k]
           pk1 := P[(k+1)%r]
           removePref(&prefs, &invRank, pk, qk)
           if len(prefs[pk]) == 0 {
               outNo()
               return
           }
           removePref(&prefs, &invRank, qk, pk1)
           if len(prefs[qk]) == 0 {
               outNo()
               return
           }
       }
   }
   // output matching
   out := make([]int, n)
   for i := 0; i < n; i++ {
       if len(prefs[i]) == 0 {
           outNo()
           return
       }
       out[i] = prefs[i][0] + 1
   }
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   for i, v := range out {
       if i > 0 {
           fmt.Fprint(w, " ")
       }
       fmt.Fprint(w, v)
   }
   fmt.Fprintln(w)
}

func outNo() {
   fmt.Println(-1)
}

// removePref removes j from prefs[i] and updates invRank; if removing head, person is enqueued by caller if needed
// removePref removes j from prefs[i], updates invRank, and returns true if head was removed
func removePref(prefs *[][]int, invRank *[][]int, i, j int) bool {
   arr := (*prefs)[i]
   pos := (*invRank)[i][j]
   headRemoved := (pos == 0)
   // remove element at pos
   copy(arr[pos:], arr[pos+1:])
   arr = arr[:len(arr)-1]
   (*prefs)[i] = arr
   // update invRank for shifted elements
   for k := pos; k < len(arr); k++ {
       (*invRank)[i][arr[k]] = k
   }
   return headRemoved
}
