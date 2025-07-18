#include <bits/stdc++.h>
using namespace std;

#ifdef ILIKEGENTOO
#   define Eo(x) { cerr << #x << " = " << (x) << endl; }
#   define E(x) { cerr << #x << " = " << (x) << "   "; }
#   define FREOPEN(x)
#else
#   define Eo(x)
#   define E(x)
#   define FREOPEN(x) (void)freopen(x ".in", "r", stdin);(void)freopen(x ".out", "w", stdout);
#endif
#define EO(x) Eo(x)
#define Sz(x) (int((x).size()))
#define All(x) (x).begin(),(x).end()

template<class A, class B> ostream &operator<<(ostream &os, const pair<A, B>& p) { return os << '(' << p.first << ", " << p.second << ')'; }

template<class C> void operator<<(vector<C> &v, const C &x){v.push_back(x);}
template<class D> void operator>>(vector<D> &v, D &x){assert(!v.empty()); x=v.back(); v.pop_back();}
template<class E> void operator<<(set<E> &v, const E &x){v.insert(x);}
template<class F> void operator<<(queue<F> &c, const F& v){v.push(v);}
template<class G> void operator>>(queue<G> &c, const G& v){const G r=v.front();v.pop();return r;}

typedef double flt;
typedef long long int64;
typedef unsigned long long uint64;
typedef pair<int, int> pii;

const int inf = 0x3f3f3f3f;
const int64 inf64 = 0x3f3f3f3f3f3f3f3fLL;
const flt eps = 1e-8;
const flt pi = acos(-1.0);
const int dir[4][2] = { {0, 1}, {1, 0}, {0, -1}, {-1, 0} };

int n;
int64 w;
const int maxn = 100500;
//int a[maxn], b[maxn];

const int maxg = 303;
const int maxb = maxn / maxg + 2;

vector<int> user_by_b[maxn];

pii ab[maxn];
bool have[maxn];

struct Block {
    int64 ans;
    int pans;
    int rightcnt;
    int nextrcnt;

    Block() {
        ans = 0;
        pans = 0;
        rightcnt = 0;
        nextrcnt = inf;
    }

    int ld, ldx;
    void get_leader(const int from, const int to) {
        int ourcnt = 0;
        ans = pans = 0;
        for (int i = to - 1; i >= from; --i) if (have[i]) {
            ++ourcnt;
            int64 curans = int64(rightcnt + ourcnt) * ab[i].first;
            if (curans > ans) {
                ans = curans;
                pans = ab[i].first;
                ld = i;
                ldx = rightcnt + ourcnt;
            }
        }
    }

    void recalc(const int from, const int to) {
        get_leader(from, to);
        int ourcnt = 0;
        nextrcnt = inf;
        for (int i = to - 1; i >= from; --i) if (have[i]) {
            ++ourcnt;
            if (i == ld) continue;
            if (ab[i].first <= ab[ld].first) continue;

            int curx = rightcnt + ourcnt;
            int64 lp = (int64(ab[ld].first) * ldx - int64(ab[i].first) * curx);
            assert(lp >= 0);
            int64 rp = ab[i].first - ab[ld].first;
            int64 t = lp / rp; //(lp + rp - 1) / rp;
            ++t;
            if (t < inf) nextrcnt = min(nextrcnt, rightcnt + int(t));
        }
        E(from); E(to); E(ld); E(ldx); E(ans); Eo(nextrcnt);
    }

    void update(const int from, const int to) {
        ++rightcnt;
        ans += pans;
        if (rightcnt >= nextrcnt)
            recalc(from, to);
    }

} bs[maxb];

int main() {
    // static_assert(sizeof(long) == 8, "32-bit !!! :'(");
    // FREOPEN("f");
    ios_base::sync_with_stdio(false);

    cin >> n >> w;
    int mb = 0;
    for (int i = 0; i < n; ++i) {
        int aa, bb; cin >> aa >> bb;
        ab[i] = pii(aa, bb);
        mb = max(mb, bb);
    }
    sort(ab, ab + n);
    for (int i = 0; i < n; ++i) {
        user_by_b[ab[i].second] << i;
    }

    int watch_ads = n;
    for (int c = 0; c <= mb + 1; ++c) {
        if (c) {
            for (int user : user_by_b[c - 1]) {
                --watch_ads;
                if (ab[user].first == 0) continue;
                have[user] = true;

                int gid = user / maxg;
                // E(user); Eo(gid);
                bs[gid].recalc(gid * maxg, min(n, (gid + 1) * maxg));
                for (int i = 0; i < gid; ++i)
                    bs[i].update(i * maxg, (i + 1) * maxg);
            }
        }
        E(c); Eo(watch_ads);

        int64 ans = 0;
        int pans = 0;
        for (int i = 0; i < maxb; ++i) if (bs[i].ans > ans) {
            ans = bs[i].ans;
            pans = bs[i].pans;
            E(i); E(ans); Eo(pans);
        }
        ans += w * c * watch_ads;
        cout << ans << ' ' << pans << '\n';
    }

    return 0;
}