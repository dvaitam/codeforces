package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
   "strings"
)

func main() {
   scanner := bufio.NewScanner(os.Stdin)
   if !scanner.Scan() {
       return
   }
   n, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
   if err != nil {
       fmt.Fprintln(os.Stderr, "Invalid input")
       return
   }
   seen := make(map[string]struct{})
   for i := 0; i < n; i++ {
       // read two lines of the amulet
       if !scanner.Scan() {
           break
       }
       parts := strings.Fields(scanner.Text())
       if len(parts) < 2 {
           continue
       }
       a, _ := strconv.Atoi(parts[0])
       b, _ := strconv.Atoi(parts[1])
       if !scanner.Scan() {
           break
       }
       parts2 := strings.Fields(scanner.Text())
       if len(parts2) < 2 {
           continue
       }
       c, _ := strconv.Atoi(parts2[0])
       d, _ := strconv.Atoi(parts2[1])
       // skip separator
       if i < n-1 {
           scanner.Scan()
       }
       // generate all rotations and pick minimal representation
       r1 := fmt.Sprintf("%d%d%d%d", a, b, c, d)
       r2 := fmt.Sprintf("%d%d%d%d", c, a, d, b)
       r3 := fmt.Sprintf("%d%d%d%d", d, c, b, a)
       r4 := fmt.Sprintf("%d%d%d%d", b, d, a, c)
       min := r1
       if r2 < min {
           min = r2
       }
       if r3 < min {
           min = r3
       }
       if r4 < min {
           min = r4
       }
       seen[min] = struct{}{}
   }
   fmt.Println(len(seen))
}
