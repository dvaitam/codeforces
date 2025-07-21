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
   line, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Fprintln(os.Stderr, err)
       return
   }
   fields := strings.Fields(line)
   if len(fields) < 3 {
       return
   }
   n, _ := strconv.Atoi(fields[0])
   p, _ := strconv.Atoi(fields[1])
   k, _ := strconv.Atoi(fields[2])

   start := p - k
   if start < 1 {
       start = 1
   }
   end := p + k
   if end > n {
       end = n
   }

   var tokens []string
   if start > 1 {
       tokens = append(tokens, "<<")
   }
   for i := start; i <= end; i++ {
       if i == p {
           tokens = append(tokens, fmt.Sprintf("(%d)", i))
       } else {
           tokens = append(tokens, strconv.Itoa(i))
       }
   }
   if end < n {
       tokens = append(tokens, ">>")
   }

   fmt.Println(strings.Join(tokens, " "))
}
