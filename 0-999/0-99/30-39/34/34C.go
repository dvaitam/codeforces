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
       fmt.Fprintln(os.Stderr, "Error reading input:", err)
       os.Exit(1)
   }
   line = strings.TrimSpace(line)
   if line == "" {
       return
   }
   parts := strings.Split(line, ",")
   seen := make(map[int]bool)
   var nums []int
   for _, p := range parts {
       v, err := strconv.Atoi(p)
       if err != nil {
           continue
       }
       if !seen[v] {
           seen[v] = true
           nums = append(nums, v)
       }
   }
   if len(nums) == 0 {
       return
   }
   var ranges []string
   start := nums[0]
   prev := nums[0]
   for i := 1; i < len(nums); i++ {
       v := nums[i]
       if v == prev+1 {
           prev = v
           continue
       }
       // end current range
       if start == prev {
           ranges = append(ranges, strconv.Itoa(start))
       } else {
           ranges = append(ranges, fmt.Sprintf("%d-%d", start, prev))
       }
       // start new range
       start = v
       prev = v
   }
   // append last range
   if start == prev {
       ranges = append(ranges, strconv.Itoa(start))
   } else {
       ranges = append(ranges, fmt.Sprintf("%d-%d", start, prev))
   }
   fmt.Println(strings.Join(ranges, ","))
}
