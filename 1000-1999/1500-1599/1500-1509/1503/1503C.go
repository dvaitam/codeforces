package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type City struct {
	a int
	c int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	cities := make([]City, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &cities[i].a, &cities[i].c)
	}

	sort.Slice(cities, func(i, j int) bool {
		return cities[i].a < cities[j].a
	})

	var ans int64
	ans += int64(cities[0].c)
	reach := int64(cities[0].a + cities[0].c)

	for i := 1; i < n; i++ {
		if int64(cities[i].a) > reach {
			ans += int64(cities[i].a) - reach
			reach = int64(cities[i].a)
		}
		ans += int64(cities[i].c)
		if int64(cities[i].a+cities[i].c) > reach {
			reach = int64(cities[i].a + cities[i].c)
		}
	}

	fmt.Fprintln(writer, ans)
}
