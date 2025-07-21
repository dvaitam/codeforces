package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(reader, &s); err != nil {
       return
   }
   n := len(s)
   ans := 1
   // For each starting position p0 as bottom of pile
   for p0 := 0; p0 < n; p0++ {
       // Memoization for state (p_tip, width)
       memo := make(map[int]int)
       // dfs returns max height from state (p_tip, w)
       var dfs func(pTip, w int) int
       dfs = func(pTip, w int) int {
           key := pTip*(n+1) + w
           if v, ok := memo[key]; ok {
               return v
           }
           best := 1
           // f is fold index: between f and f+1, left width = f+1
           // require pTip <= f and new pos pNew < w
           // f <= floor((w + pTip - 2)/2)
           maxF := (w + pTip - 2) / 2
           for f := pTip; f <= maxF; f++ {
               // new position from reflection of pTip over f
               pNew := 2*f - pTip + 1
               if pNew < 0 || pNew >= w {
                   continue
               }
               if s[pNew] != s[p0] {
                   continue
               }
               h := 1 + dfs(pNew, f+1)
               if h > best {
                   best = h
               }
           }
           memo[key] = best
           return best
       }
       h := dfs(p0, n)
       if h > ans {
           ans = h
       }
   }
   fmt.Println(ans)
}
