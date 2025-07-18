#include <iostream>
#include <cstdio>
#include <cstring>
#include <climits>
#include <string>
#include <map>
#include <vector>
#include <set>
#include <list>
#include <cstdlib>
#include <cmath>
#include <algorithm>
#include <queue>
#include <stack>
#include <functional>
#include <complex>
#define mp make_pair
#define X first
#define Y second
#define MEMSET(a, b) memset(a, b, sizeof(a))
using namespace std;

typedef unsigned int ui;
typedef long long ll;
typedef unsigned long long ull;
typedef pair<int, int> pii;
typedef vector<int> vi;
typedef vi::iterator vi_it;
typedef map<int, int> mii;
typedef priority_queue<int> pqi;
typedef priority_queue<int, vector<int>, greater<int> > rpqi;
typedef priority_queue<pii> pqp;
typedef priority_queue<pii, vector<pii>, greater<pii> > rpqp;

const int MOD = 998244353;
ui dp[1001][2001][2];

int main(int argc, char *argv[])
{
	int n, k;
	cin >> n >> k;
	dp[1][2][0] = 1;
	dp[1][1][1] = 1;
	for (int i = 2; i <= n; ++i)
	{
		for (int j = 1; j <= k; ++j)
		{
			dp[i][j][0] = dp[i - 1][j][0] + (dp[i - 1][j - 1][1] << 1);
			if (j >= 2) dp[i][j][0] += dp[i - 1][j - 2][0];
			dp[i][j][1] = dp[i - 1][j][1] + dp[i - 1][j - 1][1] + (dp[i - 1][j][0] << 1);
			dp[i][j][0] %= MOD;
			dp[i][j][1] %= MOD;
		}
	}
	cout << (dp[n][k][0] + dp[n][k][1]) * 2 % MOD << endl;
	return 0;
}