#include <bits/stdc++.h>

using namespace std;

int32_t main() {
	int n;
	scanf("%d", &n);
	int start = n * 10;
	printf("%d\n2 %d 1\n1 %d %d\n", n + 1, n, n, start--);
	for (int i = 0; i < n - 1; ++i)
		printf("2 %d %d\n", i + 1, start--);
}