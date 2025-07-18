#include <cmath>
#include <ctime>
#include <cctype>
#include <cstdio>
#include <cstring>
#include <cstdlib>
#include <cassert>
#include <climits>

#include <set>
#include <map>
#include <stack>
#include <queue>
#include <vector>
#include <bitset>
#include <complex>
#include <iostream>
#include <algorithm>

#define fi first
#define se second
#define pb push_back
#define lowbit(x) ((x) & -(x))
#define siz(x) ((int)(x).size())
#define all(x) (x).begin(),(x).end()
#define debug(x) cerr << #x << " = " << (x) << endl
#define rep(i, s, t) for (int i = (s), _t = (t); i < _t; ++i)
#define per(i, s, t) for (int i = (t) - 1, _s = (s); i >= _s; --i)

using namespace std;

typedef long long ll;
typedef unsigned long long ull;
typedef unsigned int ui;
typedef double db;
typedef pair<int,int> pii;
typedef pair<ll,ll> pll;
typedef vector<int> veci;

const int inf = (int)1e9;
const int mod = (int)1e9 + 7;
const int dxy[] = { -1, 0, 1, 0, -1 };
const ll INF = 1LL << 60;
const db pi = acos(-1.0);
const db eps = 1e-8;

template<class T>void rd(T &x)
{
	x = 0;
	char c;
	while (c = getchar(), c < '0');
	do
	{
		x = x * 10 + (c^'0');
	}
	while(c = getchar(), c >= '0');
}

// 只能输出整型非负数
template<class T>void pt(T x)
{
	if (!x)
	{
		putchar('0');
		return;
	}
	static char stk[65];
	int tp = 0;
	for(; x; x /= 10) stk[tp++] = x % 10;
	per(i, 0, tp) putchar(stk[i]^'0');
}

template<class T>inline void pts(T x)
{
	pt(x);
	putchar(' ');
}

template<class T>inline void ptn(T x)
{
	pt(x);
	putchar('\n');
}

template<class T>inline void Max(T &a,T b)
{
	if(b > a) a = b;
}

template<class T>inline void Min(T &a,T b)
{
	if(b < a) a = b;
}
// EOT


const int N = (int)5e5 + 5;

char L[N], R[N];

int main()
{
	int n, K;
	scanf("%d%d", &n, &K);
	scanf("%s%s", L+1, R+1);
	ll cnt = 0;
	ll cur = 1;
	ll ans;
	rep(i, 1, n+1)
	{
		cur <<=  1;
		if (L[i] == 'b') --cur;
		if (R[i] == 'a') --cur;
		if (cur > K)
		{
			ans = cnt + (ll)K * (n-i+1);
			break;
		}
		cnt += cur;
		ans = cnt + cur * (n-i);
	}
	cout << ans << endl;
	return 0;
}