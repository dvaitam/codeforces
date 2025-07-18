#include <bits/stdc++.h>

using namespace std;
#define int long long
#define endl '\n'
typedef pair<int, int> PII;
#define fr(i, a, n) for(int i = a; i <= n; i ++)
#define rf(i, a, n) for(int i = n; i >= a; i --)

const int N = 2e5 + 1;
const int mod = 998244353;
int T, t, n, m, a, b, idx, flag, a1, b1, a2, b2, k, x, y, mx, mn;
int l, r, mid, c, mi, d;
int ans, sum, cnt;
int q[500][500], p[N], dp[N], h[N], mp[N];
char path[5][N], st[5][N];
string s, ss, s1, s2;

int dx[4] = {0, 1, 0, -1};
int dy[4] = {1, 0, -1, 0};

struct node1
{
	int a, b;
}str[N];

struct node
{
	int l, r, sum, ans;
}tr[N * 4];

int cmp(node1 x, node1 y)
{
	return x.a > y.a;
}

vector<int> g[N];

void solve()
{
	cin >> n >> m;
	//	1000 1001 1004 1005 1008
	//	1002 1003 1006 1007 1010
	//  2000 2001 
	//  2002 2003
	//	100
	//	a = 2000 ^ 2001 ^ 1002 ^ 1003;
	//	b = 1001 ^ 1004 ^ 1003 ^ 1006;
	//	c = 1004 ^ 1005 ^ 1006 ^ 1007;
	//	cout << a << ' ' << b << ' ' << c << endl;
	q[1][1] = 0;
	ans = n * 2;
	sum = 0;
	while(ans >= 1)
	{
		sum ++;
		ans /= 2;
	}
	sum = pow(2, sum);
	//	cout << sum << endl;
	fr(i, 1, n)
	{
		fr(j, 1, m)
		{
			if(i == 1)
			{
				if(j == 1) continue;
				if(j % 2 == 1)
				{
					q[i][j] = q[i][j - 1] + 3;
				}else
				{
					q[i][j] = q[i][j - 1] + 1;
				}
			}else if(i % 2 == 0)
			{
				q[i][j] = q[i - 1][j] + 2;
			}else if(i % 2 == 1)
			{
				q[i][j] = q[i - 2][j] + pow(2, 9);
			}
		}
	}
	cout << n * m << endl;
	fr(i, 1, n)
	{
		fr(j, 1, m)
		{
			cout << q[i][j] << ' ';
		}
		cout << endl;
	}
	fr(i, 1, n - 3)
	{
		fr(j, 1, m - 3)
		{
			a = q[i][j] ^ q[i][j + 1] ^ q[i + 1][j] ^ q[i + 1][j + 1];
			b = q[i + 2][j + 2] ^ q[i + 2][j + 3] ^ q[i + 3][j + 2] ^ q[i + 3][j + 3];
			d = q[i + 2][j] ^ q[i + 2][j + 1] ^ q[i + 3][j] ^ q[i + 3][j + 1];
			c = q[i][j + 2] ^ q[i][j + 3] ^ q[i + 1][j + 2] ^ q[i + 1][j + 3];
			//			cout << a << ' ' << b << ' ' << c << ' ' << d << endl;
		}
	}
	//	ans = 1006 ^ 1009 ^ 1005 ^ 1008;
	//	cout << ans << endl;
}

signed main()
{
	ios::sync_with_stdio(false);
	cin.tie(0), cout.tie(0);
	cin >> T;
	while(T -- )
		solve();
	return 0;
}