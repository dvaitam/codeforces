#pragma GCC optimize("unroll-loops")

#pragma GCC optimize("O3")



#include<bits/stdc++.h>

#include<ext/pb_ds/tree_policy.hpp>

#include<ext/pb_ds/assoc_container.hpp>



#define ull unsigned long long

#define ll long long

#define ld long double

#define pb push_back

#define fi first

#define se second

#define endl '\n'

#define pw(x) (1LL << x)

#define pii pair <int, int>

#define fob find_by_order

#define ork order_of_key

#define int int64_t



using namespace std;

using namespace __gnu_pbds;



typedef tree<pii, null_type, less <pii>, rb_tree_tag,

        tree_order_statistics_node_update> oset;

mt19937 gen(time(0));

const int mod = 1e9 + 7;

const int MOD = 998244353;

const int dx[4] = {0, 1, 0, -1};

const int dy[4] = {1, 0, -1, 0};



int32_t main() {

    ios_base::sync_with_stdio(0);

    cin.tie(0);

    cout.tie(0);

#ifdef LOCAL

    freopen("input.txt", "r", stdin);

    freopen("output.txt", "w", stdout);

#else

//    freopen("distance.in", "r", stdin);

//    freopen("distance.out", "w", stdout);

#endif // LOCAL

    int q;

    cin >> q;

    while (q --) {

        int n;

        cin >> n;

        int a[n], s = 0;

        for(int i = 0; i < n; ++i) {

            cin >> a[i];

            s += a[i];

        }

        if(!s) {

            cout << "NO" << endl;

            continue;

        }

        cout << "YES" << endl;

        if(s > 0) {

            for(int i = 0; i < n; ++i) {

                if(a[i] > 0) {

                    cout << a[i] << " ";

                }

            }

            for(int i = 0; i < n; ++i) {

                if(a[i] <= 0) {

                    cout << a[i] << " ";

                }

            }

            cout << endl;

        } else {

            for(int i = 0; i < n; ++i) {

                if(a[i] < 0) {

                    cout << a[i] << " ";

                }

            }

            for(int i = 0; i < n; ++i) {

                if(a[i] >= 0) {

                    cout << a[i] << " ";

                }

            }

            cout << endl;

        }

    }

    return 0;

}