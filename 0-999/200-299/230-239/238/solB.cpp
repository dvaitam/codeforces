#include <cstdio>
#include <cstring>
#include <algorithm>

const int maxn = 100010;
int n, h, a[maxn];

int main() {
	scanf("%d%d", &n, &h);
	for (int i = 0; i < n; i++) scanf("%d", a + i);
	int pos = std::min_element(a, a + n) - a;
	std::sort(a, a + n);
	int val1 = a[n - 1] + a[n - 2] - a[0] - a[1], 
		val2 = std::max(a[n - 1] + a[n - 2], a[0] + a[n - 1] + h)
			- std::min(a[1] + a[2], a[0] + a[1] + h);
	if (val1 < val2) pos = -1;
	printf("%d\n", std::min(val1, val2));
	for (int i = 0; i < n; i++)
		printf("%d ", i == pos ? 2 : 1);
}