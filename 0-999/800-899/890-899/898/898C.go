package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
   "strings"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   // map from name to set of phone numbers
   contacts := make(map[string]map[string]struct{})
   for i := 0; i < n; i++ {
       var name string
       var k int
       fmt.Fscan(reader, &name, &k)
       if _, ok := contacts[name]; !ok {
           contacts[name] = make(map[string]struct{})
       }
       for j := 0; j < k; j++ {
           var num string
           fmt.Fscan(reader, &num)
           contacts[name][num] = struct{}{}
       }
   }
   // collect and sort names
   names := make([]string, 0, len(contacts))
   for name := range contacts {
       names = append(names, name)
   }
   sort.Strings(names)
   // output number of distinct friends
   fmt.Fprintln(writer, len(names))
   // for each friend, process numbers
   for _, name := range names {
       numsMap := contacts[name]
       nums := make([]string, 0, len(numsMap))
       for num := range numsMap {
           nums = append(nums, num)
       }
       // remove numbers that are suffix of another
       keep := make([]string, 0, len(nums))
       for _, x := range nums {
           drop := false
           for _, y := range nums {
               if x != y && len(y) >= len(x) && strings.HasSuffix(y, x) {
                   drop = true
                   break
               }
           }
           if !drop {
               keep = append(keep, x)
           }
       }
       sort.Strings(keep)
       // print record
       fmt.Fprintf(writer, "%s %d", name, len(keep))
       for _, num := range keep {
           fmt.Fprintf(writer, " %s", num)
       }
       fmt.Fprintln(writer)
   }
}
