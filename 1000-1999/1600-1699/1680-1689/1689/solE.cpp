// #define FRR

// #define DBG

// #define OIO



#include <bits/stdc++.h>

using std::cin;

using std::cout;

using std::cerr;

#include <ext/pb_ds/assoc_container.hpp>

#include <ext/pb_ds/hash_policy.hpp>

#include <ext/pb_ds/priority_queue.hpp>

#include <ext/pb_ds/tree_policy.hpp>

using namespace __gnu_pbds;



#ifdef FRR

#define frr freopen("1.in","r",stdin)

#else

#define frr

#endif



#ifdef OIO

#define oio

#else

#define oio std::ios::sync_with_stdio(false), cin.tie(nullptr)

#endif



#ifdef DBG

#define dbg(x) std::cerr << #x << ": " << (x) << "  "

#define dbgl std::cerr << "\n"

#else

#define dbg(x)

#define dbgl

#endif



typedef tree<std::pair<int, int>, null_type, std::less<std::pair<int, int>>, rb_tree_tag, tree_order_statistics_node_update> rbTree;

using i64 = long long;

#define int i64

#define double long double



struct DSU {

    std::vector<int> pa, siz;

    DSU() : pa(), siz() {}

    DSU(int n) : pa(n), siz(n, 1) {

        iota(pa.begin(), pa.end(), 0);

    };

    int leader(int x) {

        while (x != pa[x]) x = pa[x] = pa[pa[x]];

        return x;

    }

    bool same(int x, int y) {

        return leader(x) == leader(y);

    }

    bool mergeto(int x, int y) { // x -> y

        x = leader(x);

        y = leader(y);

        if (x == y) return false;

        siz[y] += siz[x];

        pa[x] = y;

        return true;

    }

    int size(int x) {

        return siz[leader(x)];

    }

};



void solve() {

    int n;

    cin >> n;

    std::vector<int> a(n);

    int ans = 0;

    for (int i = 0; i < n; ++i) {

        cin >> a[i];

        if (a[i] == 0) ++a[i], ++ans;

    }



    auto check = [&]() {

        DSU dsu(n);

        for (int k = 0; k <= 30; ++k) {

            int fir = -1;

            for (int i = 0; i < n; ++i) {

                if (a[i] >> k & 1) {

                    if (fir == -1) fir = i;

                    else dsu.mergeto(i, fir);

                }

            }

        }

        if (dsu.size(0) == n) return true; else return false;

    };



    auto show = [&]() {

        cout << ans << "\n";

        for (int i = 0; i < n; ++i) cout << a[i] << " ";

        cout << "\n";

        return;

    };



    if (check()) {

        show();

        return;

    }



    for (int i = 0; i < n; ++i) {

        ++a[i], ++ans;

        if (check()) {

            show();

            return;

        }

        --a[i], --ans;

    }



    for (int i = 0; i < n; ++i) {

        --a[i], ++ans;

        if (check()) {

            show();

            return;

        }

        ++a[i], --ans;

    }



    int cntmx = 0, mx = 0;

    for (int i = 0; i < n; ++i) {

        if ((a[i] & -a[i]) > mx) mx = (a[i] & -a[i]), cntmx = 1;

        else if ((a[i] & -a[i]) == mx) ++cntmx;

    }



    {

        for (int i = 0; i < n; ++i) {

            if ((a[i] & -a[i]) == mx) {

                --a[i], ++ans;

                break;

            }

        }

        for (int i = n - 1; i >= 0; --i) {

            if ((a[i] & -a[i]) == mx) {

                ++a[i], ++ans;

                break;

            }

        }

        show();

    }



}





signed main() {

    oio;

    frr;



    int t;

    cin >> t;

    while (t--) {

        solve();

    }

    return 0;

}