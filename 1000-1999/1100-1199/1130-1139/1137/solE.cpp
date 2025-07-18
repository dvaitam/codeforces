#include <cstdio>
#include <cstring>
#include <algorithm>
#include <cstdlib>
#include <ctime>
#include <cmath>
#include <vector>
#include <map>
#include <set>
#include <queue>
#include <chrono>
#include <random>
#define int long long
#define M 600010
#define PII pair <int, int>
#define x first
#define y second

using namespace std;

typedef long long ll;

int n, m;

int op[M], QA[M], QB[M];

ll getans(PII A, ll k, ll b) {
    return A.x * k + A.y + b;
}

ll calc(PII A, PII B) {
    ll dety = A.y - B.y;
    ll detx = B.x - A.x;
    if(dety % detx == 0) return dety / detx;
    return dety / detx + 1;
}

bool check(PII A, PII B, PII C) {
    return calc(A, B) > calc(B, C);
}

void tryadd(vector <PII> &V, PII x) {
    if(x.y >= V.back().y) return;
    while(V.size() > 1 && !check(V[(int) V.size() - 2], V.back(), x)) V.pop_back();
    V.push_back(x);
}

void Solve(int l, int r, ll initot) {
    vector <PII> V(0);
    V.push_back(PII(0, 0));
    ll tot = initot;
    ll adds = 0, addb = 0;
    for(int i = l; i <= r; i++) {if(op[i] == 2) {
        int k = QA[i];
        tryadd(V, PII(tot, -adds * tot - addb));
        tot += k;
    } else {
        adds += QB[i];
        addb += QA[i];
    }
        while(V.size() > 1 && getans(V[(int) V.size() - 2], adds, addb) <= getans(V.back(), adds, addb)) V.pop_back();
        printf("%lld %lld\n", V.back().x + 1, getans(V.back(), adds, addb));
    }
}

signed main() {
    scanf("%lld%lld", &n, &m);
    for(int i = 1; i <= m; i++) {
        scanf("%lld%lld", &op[i], &QA[i]);
        if(op[i] == 3) scanf("%lld", &QB[i]);
    }
    op[++m] = 1;
    ll tot = n;
    for(int i = 1, nxt; i <= m; i = nxt + 1) {
        ll tmptot = tot;
        for(nxt = i; op[nxt] != 1; nxt++) {
            if(op[nxt] == 2) tot += QA[nxt];
        }
        Solve(i, nxt - 1, tmptot);
        if(nxt < m) printf("1 0\n");
        tot += QA[nxt];
    }
    return 0;
}