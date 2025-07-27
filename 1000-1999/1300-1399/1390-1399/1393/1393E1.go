package main

import (
   "bufio"
   "fmt"
   "os"
)

const mod = 1000000007

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   s := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &s[i])
   }
   // dp over words, dp[j]: number of ways for previous word with deletion state j
   prev := make([]int, len(s[0])+1)
   // initial: any deletion (including none) is valid
   for j := range prev {
       prev[j] = 1
   }
   // process each adjacent pair
   for i := 1; i < n; i++ {
       a := s[i-1]
       b := s[i]
       la := len(a)
       lb := len(b)
       // states: removal pos 0..la (la means no removal)
       // and 0..lb
       curr := make([]int, lb+1)
       // precompute lengths
       // for each k in [0..la], compute len ta = la - (k < la ? 1 : 0)
       lena := make([]int, la+1)
       for k := 0; k <= la; k++ {
           if k < la {
               lena[k] = la - 1
           } else {
               lena[k] = la
           }
       }
       lenb := make([]int, lb+1)
       for j := 0; j <= lb; j++ {
           if j < lb {
               lenb[j] = lb - 1
           } else {
               lenb[j] = lb
           }
       }
       // check transitions
       for k := 0; k <= la; k++ {
           if prev[k] == 0 {
               continue
           }
           // for each j
           for j := 0; j <= lb; j++ {
               // compare ta <= tb
               // compare first min lengths
               la2 := lena[k]
               lb2 := lenb[j]
               minl := la2
               if lb2 < minl {
                   minl = lb2
               }
               ok := false
               // compare char by char
               p := 0
               for p < minl {
                   // get char from a at T_a[p]
                   ia := p
                   if k < la && p >= k {
                       ia = p + 1
                   }
                   ib := p
                   if j < lb && p >= j {
                       ib = p + 1
                   }
                   ca := a[ia]
                   cb := b[ib]
                   if ca < cb {
                       ok = true
                       break
                   }
                   if ca > cb {
                       ok = false
                       break
                   }
                   p++
               }
               if p == minl {
                   // all equal prefix
                   if la2 <= lb2 {
                       ok = true
                   } else {
                       ok = false
                   }
               }
               if ok {
                   curr[j] = (curr[j] + prev[k]) % mod
               }
           }
       }
       // move curr to prev
       prev = curr
   }
   // sum up final dp
   ans := 0
   for _, v := range prev {
       ans = (ans + v) % mod
   }
   fmt.Println(ans)
}
