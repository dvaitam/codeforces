#include <cstdio>
#define N 200010
#define R 1000000007

int n, l, r, a[3];
int d[N][3] = {{1, 0, 0}};

int main()
{
	scanf("%d %d %d", &n, &l, &r);

	a[0] = r/3 - (l-1)/3;
	a[1] = (r+2)/3 - (l+1)/3;
	a[2] = (r+1)/3 - l/3;

	for(int i=1; i<=n; ++i)
	{
		d[i][0] = (1LL*d[i-1][0]*a[0] + 1LL*d[i-1][1]*a[2] + 1LL*d[i-1][2]*a[1])%R;
		d[i][1] = (1LL*d[i-1][0]*a[1] + 1LL*d[i-1][1]*a[0] + 1LL*d[i-1][2]*a[2])%R;
		d[i][2] = (1LL*d[i-1][0]*a[2] + 1LL*d[i-1][1]*a[1] + 1LL*d[i-1][2]*a[0])%R;
	}

	printf("%d", d[n][0]);
	return 0;
}