package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // read all input lines, concatenate to one string
   var s string
   for {
       var line string
       if _, err := fmt.Fscan(reader, &line); err != nil {
           break
       }
       s += line
   }
   // stack of cell counts for open tables
   stack := make([]int, 0)
   // result list of cell counts per table
   res := make([]int, 0)
   // parse tags
   for i := 0; i < len(s); {
       if s[i] != '<' {
           i++
           continue
       }
       // at '<'
       j := i + 1
       closing := false
       if j < len(s) && s[j] == '/' {
           closing = true
           j++
       }
       start := j
       // find '>'
       for j < len(s) && s[j] != '>' {
           j++
       }
       tag := s[start:j]
       if !closing {
           switch tag {
           case "table":
               stack = append(stack, 0)
           case "td":
               if len(stack) > 0 {
                   stack[len(stack)-1]++
               }
           }
       } else {
           if tag == "table" && len(stack) > 0 {
               // close current table
               cnt := stack[len(stack)-1]
               stack = stack[:len(stack)-1]
               res = append(res, cnt)
           }
       }
       // move past '>'
       if j < len(s) {
           i = j + 1
       } else {
           break
       }
   }
   // sort results in non-decreasing order
   sort.Ints(res)
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for i, v := range res {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprintf(writer, "%d", v)
   }
   writer.WriteByte('\n')
}
