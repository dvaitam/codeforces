package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   input, err := reader.ReadString('\n')
   if err != nil && err.Error() != "EOF" {
       fmt.Println("No solution")
       return
   }
   s := strings.TrimSpace(input)
   n := len(s)
   // find positions of '@'
   var pos []int
   for i, c := range s {
       if c == '@' {
           pos = append(pos, i)
       }
   }
   m := len(pos)
   if m == 0 {
       fmt.Println("No solution")
       return
   }
   // check first and last
   if pos[0] < 1 || pos[m-1] > n-2 {
       fmt.Println("No solution")
       return
   }
   // check spacing between '@'
   for i := 0; i < m-1; i++ {
       if pos[i+1]-pos[i] < 3 {
           fmt.Println("No solution")
           return
       }
   }
   // build segments
   var parts []string
   start := 0
   for i := 0; i < m-1; i++ {
       // end of this segment
       end := pos[i+1] - 2
       if end < start {
           fmt.Println("No solution")
           return
       }
       parts = append(parts, s[start:end+1])
       start = end + 1
   }
   // last segment
   parts = append(parts, s[start:])
   // verify each part has one '@' and valid
   for _, p := range parts {
       idx := strings.Index(p, "@")
       if idx <= 0 || idx >= len(p)-1 || strings.Count(p, "@") != 1 {
           fmt.Println("No solution")
           return
       }
   }
   // output
   fmt.Println(strings.Join(parts, ","))
}
