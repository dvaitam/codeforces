package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // consume end of line
   // read each short IPv6 and restore full form
   for i := 0; i < n; i++ {
       var s string
       fmt.Fscan(reader, &s)
       full := expandIPv6(s)
       fmt.Println(full)
   }
}

// expandIPv6 restores full IPv6 notation of 8 blocks of 4 hex digits
func expandIPv6(s string) string {
   var blocks []string
   if strings.Contains(s, "::") {
       parts := strings.SplitN(s, "::", 2)
       var left, right []string
       if parts[0] != "" {
           left = strings.Split(parts[0], ":")
       }
       if parts[1] != "" {
           right = strings.Split(parts[1], ":")
       }
       // count missing zero blocks
       missing := 8 - (len(left) + len(right))
       blocks = make([]string, 0, 8)
       blocks = append(blocks, left...)
       for j := 0; j < missing; j++ {
           blocks = append(blocks, "0")
       }
       blocks = append(blocks, right...)
   } else {
       blocks = strings.Split(s, ":")
   }
   // now blocks should be exactly 8
   // pad each block to 4 hex digits
   for idx, b := range blocks {
       if len(b) < 4 {
           blocks[idx] = strings.Repeat("0", 4-len(b)) + b
       }
   }
   return strings.Join(blocks, ":")
}
