package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n, k int
   if _, err := fmt.Fscan(reader, &n, &k); err != nil {
       return
   }
   t := make([]int, n)
   d := make([]int, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &t[i], &d[i])
   }
   // If no calls or ignore all, free whole day
   if n == 0 || k >= n {
       fmt.Println(86400)
       return
   }
   // Precompute cost_pre and cost_post
   costPre := make([]int, n)
   for i := 0; i < n; i++ {
       cntGood := 0
       for j := 0; j < i; j++ {
           if t[j]+d[j] <= t[i] {
               cntGood++
           }
       }
       costPre[i] = i - cntGood
   }
   costPost := make([]int, n)
   for i := 0; i < n; i++ {
       cntBad := 0
       endI := t[i] + d[i]
       for j := i + 1; j < n; j++ {
           if t[j] < endI {
               cntBad++
           }
       }
       costPost[i] = cntBad
   }
   best := 0
   // start segments: from time 1 to first accepted call
   for e := 0; e < n; e++ {
       // ignore e calls before e-th
       if e <= k {
           gap := t[e] - 1
           if gap > best {
               best = gap
           }
       }
   }
   // end segments: from last accepted to end of day
   for s := 0; s < n; s++ {
       // cost = cost to prepare before s + cost to ignore calls causing after arrival in gap
       if costPre[s]+costPost[s] <= k {
           gap := 86400 - (t[s] + d[s]) + 1
           if gap > best {
               best = gap
           }
       }
   }
   // middle segments: between s and e
   for s := 0; s < n; s++ {
       // must ignore bad before s
       pre := costPre[s]
       if pre > k {
           continue
       }
       baseEnd := t[s] + d[s]
       maxMid := k - pre
       // e > s
       // we need ignore all between s and e: (e-s-1) <= maxMid => e <= s+1+maxMid
       limit := s + 1 + maxMid
       if limit > n-1 {
           limit = n - 1
       }
       for e := s + 1; e <= limit; e++ {
           gap := t[e] - baseEnd
           if gap > best {
               best = gap
           }
       }
   }
   if best < 0 {
       best = 0
   }
   if best > 86400 {
       best = 86400
   }
   // Output
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   out.WriteString(strconv.Itoa(best))
   out.WriteByte('\n')
}
