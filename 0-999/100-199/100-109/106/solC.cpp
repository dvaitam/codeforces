#include<cstdio>
#include<algorithm>

using namespace std;

int dp[1010];
int n, m, c0, d0;
int a, b, c, d;

int main()
{
    scanf("%d %d %d %d", &n, &m, &c0, &d0);
	for (int i = c0; i <= n; ++i)
	  dp[i] = max(dp[i - c0] + d0, dp[i]);
	for(int i=0;i<m;i++)
	{
		scanf("%d %d %d %d", &a, &b, &c, &d);
		for(int j = 0; j < a / b; ++j)
			for(int k = n; k >= c; --k)
				dp[k] = max(dp[k], dp[k - c] + d);
	}
	printf("%d\n", dp[n]);
	return 0;
}