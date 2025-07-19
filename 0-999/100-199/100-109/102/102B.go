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
       // ignore
   }
   s := strings.TrimSpace(line)
   count := 0
   for len(s) > 1 {
       count++
       sum := 0
       for i := 0; i < len(s); i++ {
           sum += int(s[i] - '0')
       }
       s = strconv.Itoa(sum)
   }
   fmt.Print(count)
}
