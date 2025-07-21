package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read sequence
   s, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   s = s[:len(s)-1]
   // TODO: implement solution for problem E
   // Placeholder: output 0
   fmt.Println(0)
}
