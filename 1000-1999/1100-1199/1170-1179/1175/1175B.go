package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   // Read number of commands
   line, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   line = strings.TrimSpace(line)
   l, err := strconv.Atoi(line)
   if err != nil {
       return
   }

   // Stack of multipliers
   limit := uint64(1) << 32
   stack := make([]uint64, 0, l+1)
   stack = append(stack, 1)
   var res uint64
   overflow := false

   for i := 0; i < l; i++ {
       line, err = reader.ReadString('\n')
       if err != nil {
           break
       }
       line = strings.TrimSpace(line)
       switch {
       case line == "add":
           res += stack[len(stack)-1]
           if res >= limit {
               overflow = true
               // overflow detected, stop processing
               i = l // break outer loop
           }
       case strings.HasPrefix(line, "for"): // for n
           parts := strings.Fields(line)
           n, _ := strconv.ParseUint(parts[1], 10, 64)
           top := stack[len(stack)-1]
           prod := top * n
           if top >= limit || prod >= limit {
               stack = append(stack, limit)
           } else {
               stack = append(stack, prod)
           }
       default: // end
           // pop
           if len(stack) > 1 {
               stack = stack[:len(stack)-1]
           }
       }
   }

   if overflow {
       fmt.Fprintln(writer, "OVERFLOW!!!")
   } else {
       fmt.Fprintln(writer, res)
   }
}
