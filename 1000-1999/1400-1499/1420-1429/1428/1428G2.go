package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   // read k
   line, _ := reader.ReadString('\n')
   k, _ := strconv.Atoi(line[:len(line)-1])
   // read fortunes F0..F5
   F := make([]int64, 6)
   for i := 0; i < 6; i++ {
       var x int64
       fmt.Fscan(reader, &x)
       F[i] = x
   }
   // read q
   var q int
   fmt.Fscan(reader, &q)
   ns := make([]int, q)
   maxN := 0
   for i := 0; i < q; i++ {
       fmt.Fscan(reader, &ns[i])
       if ns[i] > maxN {
           maxN = ns[i]
       }
   }
   // if k == 1, compute directly per number, else dp knapsack
   if k == 1 {
       for _, ni := range ns {
           var res int64
           n := ni
           pos := 0
           for n > 0 {
               d := n % 10
               if d%3 == 0 {
                   res += int64(d/3) * F[pos]
               }
               n /= 10
               pos++
           }
           writer.WriteString(strconv.FormatInt(res, 10))
           writer.WriteByte('\n')
       }
       return
   }
   // dp[j] = max fortune using lucky digits for sum j
   maxNv := maxN
   dp := make([]int64, maxNv+1)
   pow10 := []int{1, 10, 100, 1000, 10000, 100000}
   bound := 3 * int64(k)
   for p := 0; p < 6; p++ {
       w := 3 * pow10[p]
       if w > maxNv {
           continue
       }
       v := F[p]
       if bound*int64(w) >= int64(maxNv) {
           // complete knapsack
           for j := w; j <= maxNv; j++ {
               if dp[j-w]+v > dp[j] {
                   dp[j] = dp[j-w] + v
               }
           }
       } else {
           dpPrev := make([]int64, maxNv+1)
           copy(dpPrev, dp)
           for r := 0; r < w; r++ {
               type pair struct{ t int; val int64 }
               dq := make([]pair, 0)
               head := 0
               for t, j := 0, r; j <= maxNv; t, j = t+1, j+w {
                   val := dpPrev[j] - int64(t)*v
                   for head < len(dq) && dq[head].t < t-int(bound) {
                       head++
                   }
                   if head < len(dq) {
                       cand := dq[head].val + int64(t)*v
                       if cand > dpPrev[j] {
                           dp[j] = cand
                       } else {
                           dp[j] = dpPrev[j]
                       }
                   } else {
                       dp[j] = dpPrev[j]
                   }
                   for tail := len(dq) - 1; tail >= head; tail-- {
                       if dq[tail].val <= val {
                           dq = dq[:tail]
                       } else {
                           break
                       }
                   }
                   dq = append(dq, pair{t: t, val: val})
               }
           }
       }
   }
   // answer queries
   for _, ni := range ns {
       writer.WriteString(strconv.FormatInt(dp[ni], 10))
       writer.WriteByte('\n')
   }
}
