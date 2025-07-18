#include <iostream>
#include <cstdlib>
#include <cstring>
#include <cstdio>
#include <vector>
#include <cmath>
#include <algorithm>
#include <map>
#include <queue>
#include <set>
#include <sstream>
//#include <priority_queue>
using namespace std;
#define ll long long
#define x first
#define y second
#define pii pair<int, int>
#define pdd pair<double, double>
#define L(s) (int)(s).size()
#define VI vector<int>
#define all(s) (s).begin(), (s).end()
#define pb push_back
#define mp make_pair
#define inf 1000000000
int n, m;
int a[2222], b[2222];
ll sa, sb;
ll best;
int best_cnt;
int best_idx[4];
int bsum[2000004];
int asum[2000004];
inline ll solve() {
    sa = 0;
    sb = 0;
    for(int i = 0; i < n; ++i) {
        sa += a[i];
    }
    for(int i = 0; i < m; ++i) {
        sb += b[i];
    }

    best = abs(sa - sb);

    best_cnt = 0;

    for(int i = 0; i < n; ++i) {
        for(int j = 0; j < m; ++j) {
            if (abs(sa - sb - 2 * a[i] + 2 * b[j]) < best) {
                best = abs(sa - sb - 2 * a[i] + 2 * b[j]);
                best_cnt = 1;
                best_idx[0] = i;
                best_idx[1] = j;
            }
        }
    }

    if (n < 2 || m < 2) return best;

    int ca = 0, cb = 0;
    for(int i = 0; i < n; ++i) {
        for(int j = i + 1; j < n; ++j) {
            asum[ca++] = a[i] + a[j];
        }
    }

    for(int i = 0; i < m; ++i) {
        for(int j = i + 1; j < m; ++j) {
            bsum[cb++] = b[i] + b[j];
        }
    }

    sort(asum, asum + ca);
    sort(bsum, bsum + cb);

    int ptr = 0;
    for(int i = 0; i < ca; ++i) {
        while(ptr < cb - 1 && abs(sa - sb - 2LL * asum[i] + 2LL * bsum[ptr]) >= abs(sa - sb - 2LL * asum[i] + 2LL * bsum[ptr + 1])) {
            ++ptr;
        }

        if (abs(sa - sb - 2LL * asum[i] + 2LL * bsum[ptr]) < best) {
            best = abs(sa - sb - 2LL * asum[i] + 2LL * bsum[ptr]);
            best_cnt = 2;
            best_idx[0] = i;
            best_idx[1] = ptr;
        }
    }

    if (best_cnt == 2) {
        for(int i = 0; i < n; ++i) {
            for(int j = i + 1; j < n; ++j) {
                if (a[i] + a[j] == asum[best_idx[0]]) {
                    best_idx[0] = i;
                    best_idx[2] = j;
                    i = j = n + 1;
                    break;
                }
            }
        }

        for(int i = 0; i < m; ++i) {
            for(int j = i + 1; j < m; ++j) {
                if (b[i] + b[j] == bsum[best_idx[1]]) {
                    best_idx[1] = i;
                    best_idx[3] = j;
                    i = j = m + 1;
                    break;
                }
            }
        }
    }
    return best;
}

inline ll naive() {

    sa = 0; sb = 0;
    for(int i = 0; i < n; ++i) {
        sa += a[i];
    }
    for(int i = 0; i < m; ++i) {
        sb += b[i];
    }

    ll ans = abs(sb - sa);

    for(int i = 0; i < n; ++i) {
        for(int j = 0; j < m; ++j) {
            for(int k = 0; k < n; ++k) {
                for(int l = 0; l < m; ++l) {

                    ll ca = sa - a[i] + b[j];
                    ll cb = sb - b[j] + a[i];
                    swap(a[i], b[j]);

                    ans = min(ans, abs(ca - cb));

                    ca += b[l] - a[k];
                    cb += a[k] - b[l];


                    ans = min(ans, abs(ca - cb));

                    swap(a[i], b[j]);
                }
            }
        }
    }

    return ans;

}

int main() {
//
//    for(int iter = 0; iter < 100000; ++iter) {
//        n = 10;
//        m = 10;
//        for(int i = 0; i < n; ++i) {
//            a[i] = 1 + rand() % 20;
//        }
//        for(int i = 0; i < m; ++i) {
//            b[i] = 1 + rand() % 20;
//        }
//
////
////        for(int i = 0; i < n; ++i) {
////            cout << a[i] << " ";
////        }
////        cout << endl;
////
////        for(int i = 0; i < m; ++i) {
////            cout << b[i] << " ";
////        }
////        cout << endl;
////
//        ll v2 = naive();
//        ll v1 = solve();
//        if (v1 != v2) {
//            cout << "ERRROR\n";
//            cout << v1 << " " << v2 << endl;
//
//            exit(0);
//        }
//        cout << "Ok " << iter << endl;
//    }
//    exit(0);
//

    scanf("%d", &n);
    for(int i = 0; i < n; ++i) {
        scanf("%d", &a[i]);
    }
    scanf("%d", &m);
    for(int i = 0; i < m; ++i) {
        scanf("%d", &b[i]);
    }
    solve();
    cout << best << endl;
    cout << best_cnt << endl;
    for(int i = 0; i < best_cnt; ++i) {
        cout << best_idx[2 * i] + 1 << " " << best_idx[2 * i + 1] + 1 << endl;
    }

}