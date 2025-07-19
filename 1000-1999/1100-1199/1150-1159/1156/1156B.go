package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       var s string
       fmt.Fscan(reader, &s)
       fmt.Fprintln(writer, solve(s))
   }
}

func solve(s string) string {
   n := len(s)
   // count frequencies
   orig := make([]int, 26)
   for i := 0; i < n; i++ {
       orig[s[i]-'a']++
   }
   // first attempt: evens then odds
   if res := build(orig, n, true); res != "" {
       return res
   }
   // second attempt: odds then evens
   if res := build(orig, n, false); res != "" {
       return res
   }
   return "No answer"
}

// build constructs string with evens first if evensFirst is true, then odds; otherwise odds then evens.
// Returns empty string if invalid (adjacent diff == 1).
func build(orig []int, n int, evensFirst bool) string {
   co := make([]int, 26)
   copy(co, orig)
   ans := make([]byte, 0, n)
   if evensFirst {
       for i := 0; i < 26; i += 2 {
           for co[i] > 0 {
               ans = append(ans, byte('a'+i))
               co[i]--
           }
       }
       for i := 1; i < 26; i += 2 {
           for co[i] > 0 {
               ans = append(ans, byte('a'+i))
               co[i]--
           }
       }
   } else {
       for i := 1; i < 26; i += 2 {
           for co[i] > 0 {
               ans = append(ans, byte('a'+i))
               co[i]--
           }
       }
       for i := 0; i < 26; i += 2 {
           for co[i] > 0 {
               ans = append(ans, byte('a'+i))
               co[i]--
           }
       }
   }
   // validate
   for i := 1; i < n; i++ {
       if absDiff(ans[i], ans[i-1]) == 1 {
           return ""
       }
   }
   return string(ans)
}

func absDiff(a, b byte) int {
   if a > b {
       return int(a - b)
   }
   return int(b - a)
}
