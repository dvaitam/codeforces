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

   // read n and m
   line, _ := reader.ReadString('\n')
   parts := splitInts(line)
   if len(parts) < 2 {
       return
   }
   n, m := parts[0], parts[1]
   // read a_i
   line, _ = reader.ReadString('\n')
   a := splitInts(line)
   if len(a) != n {
       // invalid input
       return
   }
   // a is sorted increasing
   if a[n-1] != m {
       fmt.Fprintln(writer, "NO")
       return
   }
   // dp bitset reachable sums
   size := (m>>6) + 1
   dp := make([]uint64, size)
   dp[0] = 1 // sum 0 reachable
   var P []int
   for _, v := range a {
       idx := v >> 6
       off := uint(v & 63)
       if (dp[idx]>>off)&1 == 1 {
           // reachable, skip
           continue
       }
       // new generator
       P = append(P, v)
       // dp |= dp << v
       shiftWords := v >> 6
       shiftBits := uint(v & 63)
       // iterate backwards
       for i := len(dp) - 1; i >= int(shiftWords); i-- {
           word := dp[i-int(shiftWords)] << shiftBits
           if shiftBits != 0 && i-int(shiftWords)-1 >= 0 {
               word |= dp[i-int(shiftWords)-1] >> (64 - shiftBits)
           }
           dp[i] |= word
       }
   }
   // output
   fmt.Fprintln(writer, "YES")
   fmt.Fprintln(writer, len(P))
   for i, v := range P {
       if i > 0 {
           writer.WriteByte(' ')
       }
       writer.WriteString(strconv.Itoa(v))
   }
   writer.WriteByte('\n')
}

// splitInts splits a space-separated line into ints
func splitInts(s string) []int {
   var res []int
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
               res = append(res, num)
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
       res = append(res, num)
   }
   return res
}
