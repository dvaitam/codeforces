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

   // read n, m
   line, _ := reader.ReadString('\n')
   parts := splitInts(line)
   n, m := parts[0], parts[1]
   // read costs
   cost := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       line, _ = reader.ReadString('\n')
       v, _ := strconv.ParseInt(trim(line), 10, 64)
       cost[i] = v
   }
   // prefix sums of cost
   prefix := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       prefix[i] = prefix[i-1] + cost[i]
   }
   // races by starting point: for each l, store races starting at l (r, profit)
   type race struct{ r int; p int64 }
   races := make([][]race, n+2)
   for j := 0; j < m; j++ {
       line, _ = reader.ReadString('\n')
       vals := splitInts(line)
       l, r := vals[0], vals[1]
       profit := int64(vals[2])
       if l >= 1 && l <= n {
           races[l] = append(races[l], race{r: r, p: profit})
       }
   }
   // dp[i] = max profit from roads i..n
   dp := make([]int64, n+3)
   // dp[n+1] = 0
   for i := n; i >= 1; i-- {
       // option: skip road i (no repair if not needed)
       best := dp[i+1]
       // option: take any race starting at i
       for _, rc := range races[i] {
           r, p := rc.r, rc.p
           // profit = dp[r+1] + p - cost of roads [i..r]
           costSegment := prefix[r] - prefix[i-1]
           temp := dp[r+1] + p - costSegment
           if temp > best {
               best = temp
           }
       }
       dp[i] = best
   }
   res := dp[1]
   if res < 0 {
       res = 0
   }
   fmt.Fprint(writer, res)
}

func splitInts(s string) []int {
   fields := make([]int, 0, 4)
   num := 0
   neg := false
   started := false
   for i := 0; i < len(s); i++ {
       c := s[i]
       if c == '-' {
           neg = true
       } else if c >= '0' && c <= '9' {
           started = true
           num = num*10 + int(c-'0')
       } else {
           if started {
               if neg {
                   num = -num
               }
               fields = append(fields, num)
               num = 0
               neg = false
               started = false
           }
       }
   }
   if started {
       if neg {
           num = -num
       }
       fields = append(fields, num)
   }
   return fields
}

func trim(s string) string {
   // remove spaces and newline
   i := len(s)
   for i > 0 && (s[i-1] == '\n' || s[i-1] == '\r' || s[i-1] == ' ' || s[i-1] == '\t') {
       i--
   }
   return s[:i]
}
