#include<bits/stdc++.h>

using namespace std;



#define sz(x) (int((x).size()))

#define endl '\n'

#define fi first

#define se second

#define pb push_back

#define pob pop_back

#define YES cout << "YES" << endl;

#define NO cout << "NO" << endl;

#define F1(x) function<void(int)> x = [&] 

#define F2(x) function<void(int, int)> x = [&] 

#define F3(x) function<void(int, int, int)> x = [&]

 

template<typename hd, typename tl> void chkmin(hd& a, tl b) { if(b < a) a = b; }

template<typename hd, typename tl> void chkmax(hd& a, tl b) { if(a < b) a = b; }

using LL = long long;

using PI = pair<int, int>;

using VI = vector<int>;

using VPI = vector<PI>;

const int N = 1000005;

const LL oo = 1000000000000000005LL;

const int P1 = 998244353;

const int P2 = 1000000007;

const int P = P1;



//assume -1 < x, y < P

inline int add(int x, int y) {

	return (x + y) % P;

}



inline int mul(int x, int y) {

	return 1LL * x * y % P;

}



inline int sub(int x, int y) {

	return (x + P - y) % P;

}



string itos(int x) {

	string str = "";

	while (x) {

		str += (x % 10 + '0');

		x /= 10;

	}

	reverse(str.begin(), str.end());

	return str;

}



int a[N], b[N], c[N], d[N];

int res[N];

bool vis[N];

int num[N];



int main() {

#ifndef ONLINE_JUDGE

	freopen("in.txt", "r", stdin);

#endif

	ios::sync_with_stdio(0);

	cin.tie(0); cout.tie(0);

	int T; cin >> T;

	while (T --) {

		int n; cin >> n;

		for (int i = 1; i <= n; i ++) c[i] = d[i] = 0;

		for (int i = 1; i <= n; i++) {

			cin >> a[i];

			if (!c[a[i]]) c[a[i]] = i;

			else d[a[i]] = i;

		}

		for (int i = 1; i <= n; i++) {

			cin >> b[i];

			if (!c[b[i]]) c[b[i]] = i;

			else d[b[i]] = i;

		}

		for (int i = 1; i <= n; i ++) num[i] = 0;

		for (int i = 1; i <= n; i ++) {

			num[a[i]] ++;

			num[b[i]] ++;

		}

		bool flg = 0;

		for (int i = 1; i <= n; i ++) if (num[i] != 2) flg = 1;

		if (flg) {

			cout << -1 << endl;

			continue;

		}

		for (int i = 1; i <= n; i ++) {

			res[i] = 0;

			vis[i] = 0;

		}

		for (int i = 1; i <= n; i ++) {

			if (vis[i]) continue;

			int x = c[i], cur = i, op = 1, cnt = 0;

			VI buf;

			buf.pb(x);

			if (a[x] != i) {

				res[x] = 1;

				cnt ++;

			} else res[x] = 0;

			while (1) {

				vis[cur] = 1;

				int nxt = a[x] + b[x] - cur;

				x = c[nxt] + d[nxt] - x, cur = nxt;

				if (vis[cur]) break;

				buf.pb(x);

				if (op == 1 && b[x] == nxt) {

					res[x] = 1;

					cnt ++;

				}

			}

			if (cnt > buf.size() - cnt) {

				for (int x : buf) res[x] ^= 1;

			}

		}

		int rlt = 0;

		for (int i = 1; i <= n; i ++) {

			if (res[i] == 1) rlt ++;

		}

		cout << rlt << endl;

		for (int i = 1; i <= n; i ++) {

			if (res[i] == 1) cout << i << ' ';

		} cout << endl;

	}

	return 0;

}