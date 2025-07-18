/*



_/      _/       _/_/_/      _/      _/    _/           _/_/_/_/_/

 _/    _/      _/      _/     _/    _/     _/           _/

  _/  _/      _/               _/  _/      _/           _/

   _/_/       _/                 _/        _/           _/_/_/_/

  _/  _/      _/                 _/        _/           _/

 _/    _/      _/      _/        _/        _/           _/

_/      _/       _/_/_/          _/        _/_/_/_/_/   _/_/_/_/_/



*/

#include <bits/stdc++.h>

#define ll long long

#define lc(x) ((x) << 1)

#define rc(x) ((x) << 1 | 1)

#define ru(i, l, r) for (int i = (l); i <= (r); i++)

#define rd(i, r, l) for (int i = (r); i >= (l); i--)

#define mid ((l + r) >> 1)

#define pii pair<int, int>

#define mp make_pair

#define fi first

#define se second

#define sz(s) (int)s.size()

#define maxn 300005

using namespace std;

inline int read() {

	int x = 0, w = 0; char ch = getchar();

	while(!isdigit(ch)) {w |= ch == '-'; ch = getchar();}

	while(isdigit(ch)) {x = x * 10 + ch - '0'; ch = getchar();}

	return w ? -x : x;

}

int n, a[maxn], p[maxn];

int buc[maxn], cnt[maxn];

void solve() {

	n = read();

	ru(i, 1, n) cnt[i] = buc[i] = 0;

	ru(i, 1, n) cnt[++buc[a[i] = read()]]++;

	int mx = n; while(!cnt[mx]) mx--;

	if(mx > (n + 1) / 2) {

		printf("NO\n");

		return;

	}

	int l = 0, r = 0;

	ru(i, 1, n) {

		if(mx == (n - l + 1) / 2) {

			l--;

			ru(j, 1, n) if(buc[j] == mx) {

				ru(k, i, n) {

					if(a[k] == j) p[l += 2] = k;

					else p[r += 2] = k;

				}

			}

			break;

		}

		if(a[p[l]] == a[i]) p[r += 2] = i;

		else {

			p[++l] = i;

			if(l < r) l++;

			else r = l;

		}

		cnt[buc[a[i]]--]--;

		while(!cnt[mx]) mx--;

	}

	printf("YES\n");

	ru(i, 1, n) printf("%d ", p[i]); printf("\n");

}

int main() {

	int T = read();

	while(T--) solve();

	return 0;

}