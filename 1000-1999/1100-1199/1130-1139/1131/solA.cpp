#pragma comment(linker, "/stack:20000000")
#pragma GCC optimize("Ofast")
#pragma GCC target("sse,sse2,sse3,ssse3,sse4,popcnt,abm,mmx")

#include <stdio.h>
#include <bits/stdc++.h>
#include <x86intrin.h>
#include <ext/pb_ds/assoc_container.hpp>
using namespace __gnu_pbds;
template <typename T> using ordered_set = tree<T, null_type, std::less<T>, rb_tree_tag, tree_order_statistics_node_update>;
template <typename T1, typename T2> using ordered_map = tree<T1, T2, std::less<T1>, rb_tree_tag, tree_order_statistics_node_update>;

using namespace std;


#define pb push_back
#define ppb pop_back
#define mp make_pair
#define fs first
#define sc second
#define abs(a) ((a) < 0 ? -(a) : (a))
#define sqr(a) ((a) * (a))

typedef unsigned int uint;
typedef long long ll;
typedef unsigned long long ull;

inline pair<pair<bool, bool>, pair<string, bool> > useinout();

#ifdef SOL
double starttime;
#endif

void initialization() {
#ifdef SOL
    starttime = 1000. * clock() / CLOCKS_PER_SEC;
    if (useinout().fs.fs)
        freopen("input.txt", "r", stdin);
    if (useinout().fs.sc)
        freopen("output.txt", "w", stdout);
#else
    srand(__rdtsc());
    const string file = useinout().sc.fs;
    if (!file.empty()) {
        freopen((file + ".in").c_str(), "r", stdin);
        freopen((file + ".out").c_str(), "w", stdout);
    } else
    if(useinout().sc.sc) {
        freopen("input.txt", "r", stdin);
        freopen("output.txt", "w", stdout);
    }
#endif
}

int solution();

void finish(int exitval) {
    fflush(stdout);
#ifdef SOL
    fprintf(stderr, "\n-----------------\nexit value: %d\ntime: %.3lf ms\n-----------------\n", exitval, 1000. * clock() / CLOCKS_PER_SEC - starttime);
#endif
}

int main() {
    initialization();
    finish(solution());
    return (0);
}

const double eps = 1e-9;
const int mod = (int) 1e+9 + 7;
const double pi = acos(-1.);
const int maxn = 100100;



int solution() {

    int a, b, c, d;
    scanf("%d%d%d%d", &a, &b, &c, &d);
    printf("%d\n", (a + 2) * 2 + (b + 2) * 2 - 4 + (c + 2) * 2 + (d + 2) * 2 - 4 - min(c + 2, a + 2) * 2);

    return 0;
}

inline pair<pair<bool, bool>, pair<string, bool> > useinout() {
    return (mp(mp(0, 0), mp("", 0)));
}

//by Andrey Kim