#include <bits/stdc++.h>



/*

#pragma comment(linker, "/stack:200000000")

#pragma GCC optimize("Ofast")

#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx,avx,tune=native")

*/



using namespace std;

using ll = long long;

using ld = long double;

#define int long long

using pii = pair<int, int>;

using vpii = vector<pii>;

using vi = vector<int>;

#define F first

#define S second

#define mp make_pair

#define pb push_back

#define pob pop_back

#define _bp __builtin_popcount

#define all(arr) (arr).begin(), (arr.end())

#define sz(arr) int((arr).size())

#define forn(i, a, b) for(int i = a; i < b; i++)



void run() {

    ios::sync_with_stdio(false);

    cin.tie(NULL);

    cout.tie(NULL);

}



const int N =(1 << 21);

const int INF = 9223372036854775807;

const int mod = 1e9 + 7;



int binpow(int a, int b, int mod) {

    a %= mod;

    int res = 1;

    while (b > 0) {

        if (b & 1) res = (res * a) % mod;

        a = (a * a) % mod;

        b >>= 1;

    }

    return res;

}



mt19937 rnd(117);

int modInverse(int n, int mod) {

    return binpow(n, mod - 2, mod);

}



int n, k;

vector<char> prime (10000+1, true);

void solve() {

    cin >> n;

    string s;

    cin >> s;

    vi a;

    vi cnt(10, 0);

    forn(i, 0, n) {

        a.pb(s[i] - '0');

        cnt[s[i] - '0']++;

    }

    bool ok = false;

    if(cnt[1] > 0) ok = true;

    if(cnt[4] > 0 || cnt[6] > 0 || cnt[8] > 0 || cnt[9] > 0) ok = true;

    if(ok) {

        cout << 1 << '\n';

        if(cnt[1] > 0) cout << 1 << '\n';

        else {

            if(cnt[4] > 0) cout << 4 << '\n';

            else {

                if(cnt[6] > 0) cout << 6 << '\n';

                else {

                    if(cnt[8] > 0) cout << 8 << '\n';

                    else {

                        if(cnt[9] > 0) cout << 9 << '\n';

                    }

                }

            }



        }

        return;

    }

    ok = false;

    forn(i, 0, n) {

        forn(j, i + 1, n) {

            if(!prime[a[i] * 10 + a[j]]) {

                cout << 2 << '\n';

                cout << a[i] << a[j] << '\n';

                return;

            }

        }

    }

    forn(i, 0, n) {

        forn(j, i + 1, n) {

            forn(y, j + 1, n) {

                if(!prime[a[i] * 100 + a[j] * 10 + a[y]]) {

                    cout << 3 << '\n';

                    cout << a[i] << a[j] << a[y] << '\n';

                    return;

                }

            }

        }

    }



}



signed main() {

    run();

prime[0] = prime[1] = false;

for (int i=2; i<=10000; ++i)

	if (prime[i])

		if (i * 1ll * i <= 10000)

			for (int j=i*i; j<=10000; j+=i)

				prime[j] = false;

    int tt = 1;

    cin >> tt;

    while (tt--) solve();

    return 0;

}