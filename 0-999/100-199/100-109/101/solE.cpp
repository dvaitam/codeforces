#include<time.h>

#include<stdlib.h>

#include<assert.h>

#include<cmath>

#include<cstring>

#include<cstdio>

#include<set>

#include<map>

#include<queue>

#include<bitset>

#include<vector>

#include<iostream>

#include<algorithm>

#include<iomanip>

using namespace std;

typedef long long ll;

typedef unsigned long long ul;

typedef vector<int> vi;

typedef pair<int, int> pii;

#define rep(i,l,r) for(int i=l;i<(r);++i)

#define per(i,l,r) for(int i=r-1;i>=(l);--i)

#define sz(x) ((int)((x).size()))

#define sqr(x) ((x)*(x))

#define all(x) (x).begin(),(x).end()

#define mp make_pair

#define pb push_back

#define fi first

#define se second

#define de(x) cout << #x << " = " << x << endl;

#define debug(x) freopen(x".in", "r", stdin);

#define setIO(x) freopen(x".in", "r", stdin);freopen(x".out", "w", stdout);

const ll LINF = 1e17 + 7;

const ul BASE = 33;

const int N = 2e4 + 7;

const int INF = 1e9 + 7;

const int MOD = 1e9 + 7;

const double Pi = acos(-1.);

const double EPS = 1e-8;

ll kpow(ll a, ll b) {

	ll ret = 1;

	for (; b; b >>= 1, a = a * a)

		if (b & 1)

			ret = ret * a;

	return ret;

}

//--------------head--------------

int n, m, p, X[N], Y[N], F[N], G[N];

ll sum;

vector<char> ans;

inline int f(int a, int b) {

	return (X[a] + Y[b]) % p;

}

inline void mxv(int &x, int y) {

	if (x < y)

		x = y;

}

void solve(int x0, int y0, int x1, int y1) {

//	printf("%d %d %d %d\n", x0, y0, x1, y1);

	if (x0 == x1) {

		sum += f(x0, y0);

		rep(j, y0 + 1, y1 + 1)

			sum += f(x0, j), ans.pb('S');

		return;

	}

	rep(j, y0, y1 + 1)

		F[j] = G[j] = 0;

	int mid = (x0 + x1) >> 1;

	rep(i, x0, mid + 1)

	{

		F[y0] += f(i, y0);

		rep(j, y0 + 1, y1 + 1)

			mxv(F[j], F[j - 1]), F[j] += f(i, j);

	}

	per(i, mid + 1, x1 + 1)

	{

		G[y1] += f(i, y1);

		per(j, y0, y1)

			mxv(G[j], G[j + 1]), G[j] += f(i, j);

	}

	int bst = y1;

	rep(j, y0, y1)

		if (F[j] + G[j] > F[bst] + G[bst])

			bst = j;

	solve(x0, y0, mid, bst);

	ans.pb('C');

	solve(mid + 1, bst, x1, y1);

}

int main() {

	scanf("%d%d%d", &n, &m, &p);

	rep(i, 0, n)

		scanf("%d", &X[i]);

	rep(i, 0, m)

		scanf("%d", &Y[i]);

	solve(0, 0, n - 1, m - 1);

	printf("%I64d\n", sum);

	rep(i, 0, sz(ans))

		putchar(ans[i]);

	return 0;

}