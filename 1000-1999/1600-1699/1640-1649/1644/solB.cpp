#include <bits/stdc++.h>

using namespace std;

#define go ios::sync_with_stdio (0); cin.tie(0); cout.tie(0);

typedef long long ll;

typedef vector<int> vi;

typedef vector<ll> vl;

typedef vector<vector <int>> vii;

typedef pair<int, int> pi;

#define F first

#define S second

#define PB push_back

#define reb(i,a,b) for (int i = a; i < b; i++)

#define endl  "\n"

#define NO cout <<"NO\n"

#define YES cout << "YES\n"

#define nax 100005

#define LFT p<<1, L, (L+R)>>1

#define RGT p<<1|1, ((L+R)>>1)+1, R

#define all(x) 		x.begin(), x.end()

#define rall(v)	v.rbegin(), v.rend()

#define dmid  ll mid = L + ((R - L ) >> 1);

#define T int t; cin >> t; while (t--)



ll M = 998244353;

ll mod(ll x) {

	return ((x % M + M) % M);

}

ll mul(ll a, ll b) {

	return mod(mod(a) * mod(b));

}

ll add(ll a, ll b) {

	return mod(mod(a) + mod(b));

}

ll dec(ll a, ll b) {

	return mod(mod(a) - mod(b));

}

int main() {

	//freopen("input.txt", "r", stdin);

	//freopen("output.txt", "w", stdout);

	//cout.flush();

	//memset(&a[0], 0, sizeof(a[0]) * a.size());

	//cout << fixed << setprecision(9);

	go;

	int t; cin >> t;

	while (t--) {

		int n; cin >> n; 

		vi v(n); int k = n;

		vi idx(n);

		if (n != 3) {

			for (int i = 0; i < n; i++) {

				v[i] = k;

				k--;

				idx[i] = i;

			}

			for (int i = 0; i < n; i++) {

				for (int j = 0; j < n; j++) {

					cout << v[idx[j]] << " ";

					idx[j] %= n;

					idx[j]--;

					if (idx[j] < 0) idx[j] = n - 1;

				}

				cout << endl;

			}

		}

		else {

			cout << 3 << ' ' << 2 << " " << 1 << endl;

			cout << 1 << ' ' << 3 << " " << 2 << endl;

			cout << 2 << ' ' << 3 << " " << 1 << endl;

		}

	}



return 0;

}