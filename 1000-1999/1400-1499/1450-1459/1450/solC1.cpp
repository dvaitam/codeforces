#include <bits/stdc++.h>

/*#include <ext/pb_ds/detail/standard_policies.hpp>

#include <ext/pb_ds/tree_policy.hpp>

#include <ext/pb_ds/assoc_container.hpp>

using namespace __gnu_pbds;

template <class T> using ordered_set = tree<T, null_type, less<T>, rb_tree_tag,tree_order_statistics_node_update>;

*/

//#pragma GCC optimize("Ofast")

//#pragma GCC optimize("unroll-loops")

//#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,tune=native") //bad

//#pragma GCC target("avx,avx2")

#define ll long long

#define ull unsigned ll

#define ff first

#define ss second

#define int ll

#define all(v) v.begin(), v.end()

#define rall(v) v.rbegin(), v.rend()

#define pb push_back

#define pii pair <int, int>

#define pdd pair <double, double>

#define _size(a) (int)a.size()

using namespace std;

mt19937 rnd(time(nullptr));

const ll inf = 1e15;

//constexpr ll mod = 998244353;

constexpr ll mod = 1e9+7;

const int N = 3e5+1, B = 350;

int sqr(int a){ return a * a; }



void solve(){

	int n;

	cin >> n;

	vector < vector <char> > a(n, vector <char> (n));

	for (int i = 0; i < n; ++i)

		for (int j = 0; j < n; ++j) cin >> a[i][j];

	int cnt = 0;

	for (int i = 0; i < n; ++i)

		for (int j = 0; j < n; ++j) cnt += a[i][j] == 'X';

	

	for (int k = 0; k < 3; ++k){

		int c = 0;

		for (int i = 0; i < n; ++i)

			for (int j = 0; j < n; ++j)

				if ((i + j) % 3 == k && a[i][j] == 'X') ++c;

		if (c <= cnt / 3){

			for (int i = 0; i < n; ++i){

				for (int j = 0; j < n; ++j)

					if ((i + j) % 3 == k && a[i][j] == 'X') cout << "O";

					else cout << a[i][j];

				cout << '\n';

			}

			return;

		}

	}

}



signed main(){

    ios::sync_with_stdio(false); cin.tie(NULL); cout.tie(NULL);

    int tt = 1;

	cin >> tt;

    while (tt--){

        solve();

        cout << '\n';

    }

}