#include <cstdio>

int k, l, w;

int main() {
	scanf("%d %d", &k, &l);
	while (l % k == 0) l /= k, w++;
	if (!w || l > 1) puts("NO");
	else printf("YES\n%d\n", w - 1);
	return 0;
}