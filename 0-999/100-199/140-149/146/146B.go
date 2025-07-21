package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var a int
   var b string
   if _, err := fmt.Fscan(reader, &a, &b); err != nil {
       return
   }
   aStr := fmt.Sprintf("%d", a)
   L := len(aStr)
   bStr := b
   m := len(bStr)
   // Try same length as a
   if s, ok := solveSameLength(aStr, bStr); ok {
       fmt.Println(s)
       return
   }
   // Otherwise, construct minimal for length > len(a)
   // minimal length is max(L+1, m)
   N := L + 1
   if m > N {
       N = m
   }
   // build result of length N: non-lucky positions as minimal digits, lucky digits bStr at end
   res := make([]byte, N)
   prefix := N - m
   for i := 0; i < N; i++ {
       if i < prefix {
           if i == 0 {
               res[i] = '1'
           } else {
               res[i] = '0'
           }
       } else {
           res[i] = bStr[i-prefix]
       }
   }
   fmt.Println(string(res))
}

// solveSameLength attempts to find minimal string s of same length as aStr, s > aStr, mask(s)==bStr
func solveSameLength(aStr, bStr string) (string, bool) {
   N := len(aStr)
   m := len(bStr)
   // memo[pos][k][gt] marks visited impossible states
   memo := make([][][]bool, N+1)
   for i := 0; i <= N; i++ {
       memo[i] = make([][]bool, m+1)
       for k2 := 0; k2 <= m; k2++ {
           memo[i][k2] = make([]bool, 2)
       }
   }
   var dfs func(pos, k int, gt bool) (string, bool)
   dfs = func(pos, k int, gt bool) (string, bool) {
       if pos == N {
           if k == m && gt {
               return "", true
           }
           return "", false
       }
       gtIdx := 0
       if gt {
           gtIdx = 1
       }
       if memo[pos][k][gtIdx] {
           return "", false
       }
       // try digits
       for d := byte('0'); d <= '9'; d++ {
           // no leading zero
           if pos == 0 && d == '0' {
               continue
           }
           // respect lower bound
           if !gt && d < aStr[pos] {
               continue
           }
           isLucky := (d == '4' || d == '7')
           nk := k
           if isLucky {
               if k < m && d == bStr[k] {
                   nk = k + 1
               } else {
                   continue
               }
           }
           // ensure no extra lucky digits after mask exhausted
           if !isLucky && k == m {
               // ok, non-lucky
           }
           // check enough space for remaining lucky digits
           rem := N - pos - 1
           need := m - nk
           if need < 0 || need > rem {
               continue
           }
           ngt := gt || (d > aStr[pos])
           if rest, ok := dfs(pos+1, nk, ngt); ok {
               return string(d) + rest, true
           }
       }
       memo[pos][k][gtIdx] = true
       return "", false
   }
   return dfs(0, 0, false)
}
