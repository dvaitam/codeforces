#include <cstdio>
#include <vector>
#include <algorithm>
using namespace std;

typedef long long LL;

LL p10[50]={1}, f[50][50];
char o[50];

int main() {
	int n;
	char s[20];

	for (int i=1; i<50; ++i) p10[i] = p10[i-1] * 10;
	scanf("%d%s", &n, s);
	for (int i=0; s[i]; ++i) s[i] -= '0';
	for (int i=0; i<=n; ++i) for (int j=0; j<=n; ++j) {
		f[i+1][j] = max(f[i+1][j], f[i][j] + s[i+j] * p10[n-i-1]);
		f[i][j+1] = max(f[i][j+1], f[i][j] + s[i+j] * p10[n-j-1]);
	}
//	printf("%lld\n", f[n][n]);
	int x = n, y = n;
	while (x || y) {
//		printf("%d %d\n", x, y);
		if (x >= 1 && f[x][y] == f[x-1][y] + s[x+y-1] * p10[n-x]) o[x+y-1] = 'H', --x;
		else o[x+y-1] = 'M', --y;
	}
	puts(o);

	return 0;
}