#include <bits/stdc++.h>

using namespace std;

#define fin stdin
#define fout stdout
//FILE *fin = fopen("x.in", "r"), *fout = fopen("x.out", "w");

typedef long long i64;

const int nmax = 2000;

struct pct {
    int x, y;

    bool operator < (const pct &shp) const {
        if (x != shp.x)
            return x < shp.x;
        return y < shp.y;
    }
} v[nmax + 1];

struct lin {
    int x, y;
    pair<int, int> p;

    bool operator < (const lin &shp) const {
        return 1LL * p.first * shp.p.second < 1LL * p.second * shp.p.first;
    }
};
vector<lin> wlin;

int ord[nmax + 1], pos[nmax + 1];

pair<int, int> panta;
inline bool cmp (int a, int b) {
    return 1LL * (v[a].y - v[b].y) * panta.second < 1LL * panta.first * (v[a].x - v[b].x);
}

inline i64 C2 (int x) {
    if (x < 2)
        return 0;

    return 1LL * x * (x - 1) / 2;
}

int main() {
    int n;
    fscanf(fin, "%d", &n);

    for (int i = 1; i <= n; ++ i) {
        fscanf(fin, "%d%d", &v[i].x, &v[i].y);
    }
    sort(v + 1, v + n + 1);

    long long ans = 0;

    for (int i = 1; i <= n; ++ i) {
        for (int j = i + 1; j <= n; ++ j) {
            if (v[i].x == v[j].x) {
                ans += C2(i - 1) * C2(n - j);
            } else {
                lin aux;
                aux.x = i, aux.y = j;
                aux.p.first = v[i].y - v[j].y;
                aux.p.second = v[i].x - v[j].x;

                if (aux.p.second < 0)
                    aux.p.second *= -1, aux.p.first *= -1;

                wlin.push_back(aux);
            }
        }
    }


    for (int i = 1; i <= n; ++ i)
        pos[i] = i;

    panta = {-1e9 - 1, 1};
    sort(pos + 1, pos + n + 1, cmp);
    for (int i = 1; i <= n; ++ i) {
        ord[pos[i]] = i;
    }

    sort(wlin.begin(), wlin.end());
    for (auto i : wlin) {
        assert(abs(ord[i.x] - ord[i.y]) == 1);

        swap(ord[i.x], ord[i.y]);

        int mx = max(ord[i.x], ord[i.y]);
        ans += C2(mx - 2) * C2(n - mx);
    }

    fprintf(fout, "%I64d\n", ans);

    return 0;
}