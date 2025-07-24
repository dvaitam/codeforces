package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   data, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, "read error:", err)
       return
   }
   // remove possible trailing newline
   if len(data) > 0 && data[len(data)-1] == '\n' {
       data = data[:len(data)-1]
   }
   s := data
   n := len(s)
   // result per depth
   var res [][]string
   maxDepth := 0
   // stack for remaining children counts
   var stk []int
   i := 0
   for i < n {
       // parse text until comma
       j := i
       for j < n && s[j] != ',' {
           j++
       }
       text := s[i:j]
       // skip comma
       j++
       // parse number until comma or end
       k := 0
       start := j
       for j < n && s[j] != ',' {
           j++
       }
       numStr := s[start:j]
       k, err = strconv.Atoi(numStr)
       if err != nil {
           // should not happen
           fmt.Fprintln(os.Stderr, "parse int error:", err)
           return
       }
       // skip comma if present
       if j < n && s[j] == ',' {
           j++
       }
       i = j
       // current depth is stack size + 1
       depth := len(stk) + 1
       if depth > maxDepth {
           maxDepth = depth
       }
       // ensure res has enough levels
       if len(res) < depth {
           for d := len(res); d < depth; d++ {
               res = append(res, []string{})
           }
       }
       res[depth-1] = append(res[depth-1], text)
       // decrement parent remaining count
       if len(stk) > 0 {
           stk[len(stk)-1]--
           for len(stk) > 0 && stk[len(stk)-1] == 0 {
               stk = stk[:len(stk)-1]
           }
       }
       // push current's children count
       if k > 0 {
           stk = append(stk, k)
       }
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, maxDepth)
   for d := 0; d < maxDepth; d++ {
       for j, txt := range res[d] {
           if j > 0 {
               writer.WriteByte(' ')
           }
           writer.WriteString(txt)
       }
       writer.WriteByte('\n')
   }
}
