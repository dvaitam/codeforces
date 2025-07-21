package main

import (
   "bufio"
   "fmt"
   "os"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   rows := make([]string, 8)
   for i := 0; i < 8; i++ {
       line, err := reader.ReadString('\n')
       if err != nil && len(line) == 0 {
           fmt.Println("NO")
           return
       }
       line = strings.TrimSpace(line)
       if len(line) != 8 {
           fmt.Println("NO")
           return
       }
       rows[i] = line
   }

   for i := 0; i < 8; i++ {
       // build target pattern for row i
       var targetBuilder strings.Builder
       for j := 0; j < 8; j++ {
           if (i+j)%2 == 0 {
               targetBuilder.WriteByte('W')
           } else {
               targetBuilder.WriteByte('B')
           }
       }
       target := targetBuilder.String()
       // check if target is a rotation of rows[i]
       doubled := rows[i] + rows[i]
       if !strings.Contains(doubled, target) {
           fmt.Println("NO")
           return
       }
   }
   fmt.Println("YES")
}
