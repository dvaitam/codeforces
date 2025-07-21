package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   line, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   s := strings.TrimSpace(line)
   line, err = reader.ReadString('\n')
   if err != nil {
       return
   }
   var k int
   fmt.Sscanf(line, "%d", &k)
   const ALPH = 26
   forbidden := make([][ALPH]bool, ALPH)
   used := make([]bool, ALPH)
   for i := 0; i < k; i++ {
       line, err = reader.ReadString('\n')
       if err != nil {
           return
       }
       line = strings.TrimSpace(line)
       if len(line) >= 2 {
           a := line[0] - 'a'
           b := line[1] - 'a'
           if a >= 0 && a < ALPH && b >= 0 && b < ALPH {
               forbidden[a][b] = true
               forbidden[b][a] = true
               used[a] = true
               used[b] = true
           }
       }
   }
   n := len(s)
   // dp[last], last = 0..25 for letters, 26 for initial state
   const INF = 1000000000
   const INIT = 26
   dpPrev := make([]int, ALPH+1)
   dpCurr := make([]int, ALPH+1)
   for i := 0; i <= ALPH; i++ {
       dpPrev[i] = INF
   }
   dpPrev[INIT] = 0
   for i := 0; i < n; i++ {
       c := s[i]
       if c < 'a' || c > 'z' {
           // treat unexpected as neutral
           c = 0
       }
       ci := int(c - 'a')
       for j := 0; j <= ALPH; j++ {
           dpCurr[j] = INF
       }
       // delete s[i]
       for last := 0; last <= ALPH; last++ {
           if dpPrev[last]+1 < dpCurr[last] {
               dpCurr[last] = dpPrev[last] + 1
           }
       }
       // keep s[i]
       for last := 0; last <= ALPH; last++ {
           // check if allowed
           if last == INIT || !forbidden[last][ci] {
               if dpPrev[last] < dpCurr[ci] {
                   dpCurr[ci] = dpPrev[last]
               }
           }
       }
       // swap
       dpPrev, dpCurr = dpCurr, dpPrev
   }
   // result
   res := INF
   for last := 0; last <= ALPH; last++ {
       if dpPrev[last] < res {
           res = dpPrev[last]
       }
   }
   fmt.Println(res)
}
