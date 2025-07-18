#pragma GCC optimize("Ofast")
#pragma GCC optimize("unroll-loops")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4")
#include <bits/stdc++.h>

using namespace std;

typedef long long ll;
typedef pair<int, int> pii;
typedef pair<ll, ll> pll;

#ifndef _MSC_VER
#define __builtin_popcont (int)__popcnt
#define $ system("pause")
#else
#define $
#endif
#define MOD ll(1e9 + 7)
#define MAXN int(4000)
#define inf (ll)1e17

template <typename S, typename T>
inline ostream &operator<<(ostream &out, const pair<S, T> &p)
{
	return out << "( " << p.first << " , " << p.second << " )";
}
template <typename S>
inline ostream &operator<<(ostream &out, const vector<S> &p)
{
	for (auto &e : p)
		out << e << ' ';
	return out;
}
template <typename T, typename S>
inline T smin(T &a, const S &b) { return a > b ? a = b : a; }
template <typename T, typename S>
inline T smax(T &a, const S &b) { return a < b ? a = b : a; }
ll po(ll v, ll u) { return u ? (po(v * v % MOD, u >> 1) * (u & 1 ? v : 1) % MOD) : 1; }
inline void add(ll &l, const ll &r) { l = (1ll * l + r) % MOD; }
ll gcd(ll v, ll u) { return u ? gcd(u, v % u) : v; }
//mt19937 rnd(chrono::steady_clock::now().time_since_epoch().count());

ll n, m;
int main()
{
	ios::sync_with_stdio(0);
	cin.tie(0);
	cout.tie(0);
	cin >> n >> m;
	if (n > m)
		swap(n, m);
	if (n == 1)
	{
		ll res = m;
		m %= 6;
		if (m <= 3)
			res -= m;
		if (m == 4)
			res -= 2;
		if (m == 5)
			res -= 1;
		return cout << res << endl, 0;
	}
	if (n == 2)
	{
		if (m == 2)
			return cout << 0 << endl, 0;
		if (m == 3 || m == 7)
			return cout << 2 * m - 2 << endl, 0;
		return cout << 2 * m << endl, 0;
	}
	cout << n * m - (n & m & 1) << endl;
	return 0;
}