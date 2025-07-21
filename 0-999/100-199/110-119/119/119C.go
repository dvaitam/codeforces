package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type Subject struct {
   a, b int64
   c    int
   id   int
}

type State struct {
   sum   int64
   prevI int
   prevX int64
}

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var n, m, k int
   if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
       return
   }
   subs := make([]Subject, m)
   for i := 0; i < m; i++ {
       fmt.Fscan(in, &subs[i].a, &subs[i].b, &subs[i].c)
       subs[i].id = i + 1
   }
   sort.Slice(subs, func(i, j int) bool {
       if subs[i].c != subs[j].c {
           return subs[i].c < subs[j].c
       }
       return subs[i].id < subs[j].id
   })
   // dp[t][i] = map[x]State
   dp := make([][]map[int64]State, n+1)
   for t := 1; t <= n; t++ {
       dp[t] = make([]map[int64]State, m)
       for i := 0; i < m; i++ {
           dp[t][i] = make(map[int64]State)
       }
   }
   // t=1
   for i := 0; i < m; i++ {
       for x := subs[i].a; x <= subs[i].b; x++ {
           dp[1][i][x] = State{sum: x, prevI: -1, prevX: -1}
       }
   }
   // transitions
   for t := 2; t <= n; t++ {
       for j := 0; j < m; j++ {
           if len(dp[t-1][j]) == 0 {
               continue
           }
           for u, st := range dp[t-1][j] {
               // two options
               vs := []int64{u + int64(k), u * int64(k)}
               for _, v := range vs {
                   for i := j + 1; i < m; i++ {
                       if subs[i].c <= subs[j].c {
                           continue
                       }
                       if v < subs[i].a || v > subs[i].b {
                           continue
                       }
                       newSum := st.sum + v
                       old, ok := dp[t][i][v]
                       if !ok || newSum > old.sum {
                           dp[t][i][v] = State{sum: newSum, prevI: j, prevX: u}
                       }
                   }
               }
           }
       }
   }
   // find best at t=n
   var bestSum int64 = -1
   endI, endX := -1, int64(0)
   for i := 0; i < m; i++ {
       for x, st := range dp[n][i] {
           if st.sum > bestSum {
               bestSum = st.sum
               endI = i
               endX = x
           }
       }
   }
   if bestSum < 0 {
       fmt.Fprintln(out, "NO")
       return
   }
   // reconstruct
   resSub := make([]int, n)
   resX := make([]int64, n)
   ti, tx := endI, endX
   for t := n; t >= 1; t-- {
       resSub[t-1] = subs[ti].id
       resX[t-1] = tx
       st := dp[t][ti][tx]
       ti, tx = st.prevI, st.prevX
   }
   fmt.Fprintln(out, "YES")
   for i := 0; i < n; i++ {
       fmt.Fprintf(out, "%d %d\n", resSub[i], resX[i])
   }
}
