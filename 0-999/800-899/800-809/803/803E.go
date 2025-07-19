package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   n, k, t int
   s        []byte
)

func max(a, b int) int {
   if a > b {
       return a
   }
   return b
}

func min(a, b int) int {
   if a < b {
       return a
   }
   return b
}

// can determines if it's possible to assign outcomes up to position pos
// so that cumulative score t stays within [lo, hi] bounds, modifying s and t on success.
func can(pos, lo, hi int) bool {
   if pos < n {
       lo = max(lo, -k+1)
       hi = min(hi, k-1)
   }
   if lo > hi {
       return false
   }
   if pos == 0 {
       if lo <= 0 && 0 <= hi {
           t = 0
           return true
       }
       return false
   }
   pos--
   switch s[pos] {
   case 'L':
       if can(pos, lo+1, hi+1) {
           t--
           return true
       }
   case 'W':
       if can(pos, lo-1, hi-1) {
           t++
           return true
       }
   case 'D':
       if can(pos, lo, hi) {
           return true
       }
   case '?':
       if can(pos, lo-1, hi+1) {
           if lo <= t && t <= hi {
               s[pos] = 'D'
           } else if t <= lo {
               s[pos] = 'W'
               t++
           } else if t >= hi {
               s[pos] = 'L'
               t--
           }
           return true
       }
   }
   return false
}

func main() {
   in := bufio.NewReader(os.Stdin)
   fmt.Fscan(in, &n, &k)
   var str string
   fmt.Fscan(in, &str)
   s = []byte(str)
   if can(n, k, k) {
       fmt.Println(string(s))
       return
   }
   if can(n, -k, -k) {
       fmt.Println(string(s))
       return
   }
   fmt.Println("NO")
}
