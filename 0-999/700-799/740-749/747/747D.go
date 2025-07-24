package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(in, &n, &k); err != nil {
       return
   }
   t := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &t[i])
   }
   // Identify cold segments
   type seg struct{ l, r int }
   segs := make([]seg, 0)
   inSeg := false
   var l int
   coldDays := 0
   for i := 0; i < n; i++ {
       if t[i] < 0 {
           coldDays++
           if !inSeg {
               inSeg = true
               l = i
           }
       } else {
           if inSeg {
               segs = append(segs, seg{l, i - 1})
               inSeg = false
           }
       }
   }
   if inSeg {
       segs = append(segs, seg{l, n - 1})
   }
   m := len(segs)
   // Impossible if not enough winter tire days
   if coldDays > k {
       fmt.Println(-1)
       return
   }
   // No cold days -> no switches
   if m == 0 {
       fmt.Println(0)
       return
   }
   // Compute interior gaps
   gaps := make([]int, 0, m)
   for i := 0; i < m-1; i++ {
       gap := segs[i+1].l - segs[i].r - 1
       if gap > 0 {
           gaps = append(gaps, gap)
       }
   }
   // Tail gap
   tail := 0
   lastR := segs[m-1].r
   if lastR < n-1 {
       tail = n - 1 - lastR
   }
   // Default transitions: sign changes count
   // = m for warm->cold starts + m or m-1 for cold->warm
   defaultCost := 2*m
   if tail == 0 {
       defaultCost--
   }
   // Greedy cover smallest gaps
   budget := k - coldDays
   sort.Ints(gaps)
   save := 0
   for _, gap := range gaps {
       if gap <= budget {
           budget -= gap
           save += 2
       } else {
           break
       }
   }
   // Cover tail if possible
   if tail > 0 && tail <= budget {
       save += 1
   }
   // Answer
   ans := defaultCost - save
   fmt.Println(ans)
}
