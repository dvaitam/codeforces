#include<bits/stdc++.h>

using namespace std;

typedef long long LL;

typedef unsigned long long ULL;

typedef pair<double, int> pdi;

ULL _ = chrono::steady_clock::now().time_since_epoch().count();

const int N = 50003, V = 200003;

const double PI = acos(-1);

int n, L, k, B, x[N], y[N], xx[N];

vector<int> vec[V];

double R;

int hah(int x, int y){return (((ULL)x << 32 | y) ^ _) % V;}

void build(int l, int r){

    if(B) for(int i = l;i <= r;++ i) vec[xx[i]].resize(0);

    B = R * 2 + 1;

    for(int i = l;i <= r;++ i) vec[xx[i] = hah(x[i] / B, y[i] / B)].push_back(i);

}

void remove(int p){vec[xx[p]].erase(vec[xx[p]].begin());}

void add(int p){vec[xx[p] = hah(x[p] / B, y[p] / B)].push_back(p);}

namespace SegTree {

    const int N = 1 << 16;

    int pre[N<<1], sum[N<<1];

    void pup(int p){

        sum[p] = sum[p<<1] + sum[p<<1|1];

        pre[p] = max(pre[p<<1], sum[p<<1] + pre[p<<1|1]);

    }

    void upd(int p, int v){

        p += N; sum[p] += v; pre[p] = max(sum[p], 0);

        while(p >>= 1) pup(p);

    }

    void clr(int p){p += N; while(p){sum[p] = pre[p] = 0; p >>= 1;}}

}

LL SQ(int x){return (LL)x * x;}

bool check(int p, double R){

    using namespace SegTree;

    auto mdf = [&](int q, int v){upd(q, v); if(q + L <= n) upd(q + L, -v);};

    vector<pdi> vc;

    int xx = x[p] / B, yy = y[p] / B;

    for(int dx = -1;dx <= 1;++ dx)

        for(int dy = -1;dy <= 1;++ dy)

            for(int j : vec[hah(xx + dx, yy + dy)]) if(j != p && abs(j - p) < L){

                double arg2 = sqrt(SQ(x[p] - x[j]) + SQ(y[p] - y[j])) / (2 * R);

                if(arg2 >= 1) continue;

                double arg1 = atan2(y[j] - y[p], x[j] - x[p]); arg2 = acos(arg2);

                double argl = arg1 - arg2, argr = arg1 + arg2;

                if(argl < -PI || argr > PI) mdf(j, 1);

                if(argl < -PI) argl += 2 * PI;

                if(argr > PI) argr -= 2 * PI;

                vc.emplace_back(argl, j); vc.emplace_back(argr, -j);

            }

    sort(vc.begin(), vc.end());

    bool flg = pre[1] >= k - 1;

    for(auto [$, j] : vc){

        if(flg) break;

        mdf(abs(j), j > 0 ? 1 : -1);

        flg |= pre[1] >= k - 1;

    }

    for(auto [$, j] : vc){clr(abs(j)); if(abs(j) + L <= n) clr(abs(j) + L);}

    return flg;

}

void solve(){

    cin >> n >> L >> k;

    R = 225675834 * sqrt((k-1.) / L);

    for(int i = 1;i <= n;++ i){

        cin >> x[i] >> y[i];

        x[i] += 1e8; y[i] += 1e8;

    }

    build(1, L);

    for(int i = 1;i <= n;++ i){

        if(i > L) remove(i - L);

        if(check(i, R)){

            double l = 0, r = R, md;

            while(r - l >= l * 1e-10)

                (check(i, md = (l + r) / 2) ? r : l) = md;

            R = l; build(max(i - L, 0) + 1, min(i + L - 1, n));

        }

        if(i + L <= n) add(i + L);

    }

    for(int i = n - L + 1;i <= n;++ i) vec[xx[i]].resize(0);

    printf("%.10lf\n", R);

}

int main(){

    ios::sync_with_stdio(false);

    int T; cin >> T; while(T --) solve();

}