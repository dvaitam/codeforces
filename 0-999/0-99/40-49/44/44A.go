package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   scanner.Split(bufio.ScanWords)
   if !scanner.Scan() {
       return
   }
   n, err := strconv.Atoi(scanner.Text())
   if err != nil {
       return
   }
   seen := make(map[string]struct{})
   for i := 0; i < n; i++ {
       if !scanner.Scan() {
           break
       }
       species := scanner.Text()
       if !scanner.Scan() {
           break
       }
       color := scanner.Text()
       key := species + " " + color
       seen[key] = struct{}{}
   }
   fmt.Println(len(seen))
}
