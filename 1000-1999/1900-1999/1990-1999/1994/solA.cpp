#ifndef LOCAL
#pragma GCC optimize("O3,unroll-loops")
#endif
#include <bits/stdc++.h>
using namespace std;

#include <ext/pb_ds/assoc_container.hpp>
using namespace __gnu_pbds;
template<class T> using oset = tree<T, null_type, less<T>, rb_tree_tag, tree_order_statistics_node_update>; 

#define ar array
#define sz(v) int(v.size())
#define all(v) v.begin(), v.end()
#define FOR(i,a,b) for (int i=a; i<b; i++)
#define ROF(i,a,b) for (int i=b-1; i>=a; i--)
typedef long long ll;
typedef vector<int> vi;
typedef pair<int, int> pi;

const int MOD=998244353;//1e9+7;

inline namespace MD {
int addself(int &a, int b) { if ((a+=b)>=MOD) a-=MOD; return a; }
int add(int a, int b) { return addself(a,b); }
int subself(int &a, int b) { if ((a-=b)<0) a+=MOD; return a; }
int sub(int a, int b) { return subself(a,b); }
int mul(int a, int b) { return (ll)a*b%MOD; }
int inv(int a) { return a==1?1:mul(inv(MOD%a),MOD-MOD/a); }
auto gen_comb(int N) {
    vi fac(N+1), ifac(N+1), vv(N+1);
    vv[1]=1; FOR(i,2,N+1) vv[i]=mul(vv[MOD%i],MOD-MOD/i);
    fac[0]=ifac[0]=1; FOR(i,1,N+1) fac[i]=mul(fac[i-1],i), ifac[i]=mul(ifac[i-1],vv[i]);
    return make_tuple(fac,ifac,vv);
}
}
/*
https://github.com/kth-competitive-programming/kactl/blob/main/content/data-structures/FenwickTree.h
*/
struct FT {
    vector<ll> s;
    FT(int n) : s(n) {}
    void update(int pos, ll dif) { // a[pos] += dif
        for (; pos < sz(s); pos |= pos + 1) s[pos] += dif;
    }
    ll query(int pos) { // sum of values in [0, pos]
        ll res = 0;
        for (pos++; pos > 0; pos &= pos - 1) res += s[pos-1];
        return res;
    }
    int lower_bound(ll sum) {// min pos st sum of [0, pos] >= sum
        // Returns n if no sum is >= sum, or -1 if empty sum is.
        if (sum <= 0) return -1;
        int pos = 0;
        for (int pw = 1 << 25; pw; pw >>= 1) {
            if (pos + pw <= sz(s) && s[pos + pw-1] < sum)
                pos += pw, sum -= s[pos-1];
        }
        return pos;
    }
};
struct segt {
    int l, r; segt *lc, *rc;
    ll x, lz;
    segt(int l, int r) : l(l), r(r) {
        x=0, lz=0;
        if (l<r) {
            int m=(l+r)/2;
            lc=new segt(l,m), rc=new segt(m+1,r);
            pul();
        }
    }
    void pul() {
        x=min(lc->x,rc->x);
    }
    void put(ll v) {
        x+=v, lz+=v;
    }
    void pus() {
        if (lz) {
            lc->put(lz), rc->put(lz);
            lz=0;
        }
    }
    void upd(int i, ll v) {
        assert(i>=l&&i<=r);
        if (l==r) x=v;
        else pus(), (i<=lc->r?lc:rc)->upd(i,v), pul();
    }
    void add(int ql, int qr, ll v) {
        assert(ql<=r&&qr>=l);
        if (ql<=l&&qr>=r) return (void)put(v);
        else {
            pus();
            if (ql<=lc->r) lc->add(ql,qr,v);
            if (qr>=rc->l) rc->add(ql,qr,v);
            pul();
        }
    }
    ll qry(int ql, int qr) {
        assert(ql<=r&&qr>=l);
        if (ql<=l&&qr>=r) return x;
        pus();
        if (qr<=lc->r) return lc->qry(ql,qr);
        if (ql>=rc->l) return rc->qry(ql,qr);
        return lc->qry(ql,qr)+rc->qry(ql,qr);
    }
};
/*
https://github.com/bqi343/cp-notebook/blob/master/Implementations/content/data-structures/Static%20Range%20Queries%20(9.1)/RMQ%20(9.1).h
*/
template<class T> struct RMQ { // floor(log_2(x))
    int level(int x) { return 31-__builtin_clz(x); }
    vector<T> v; vector<vi> jmp;
    int cmb(int a, int b) {
        return v[a]==v[b]?min(a,b):(v[a]<v[b]?a:b); }
    void init(const vector<T>& _v) {
        v = _v; jmp = {vi(sz(v))};
        iota(all(jmp[0]),0);
        for (int j = 1; 1<<j <= sz(v); ++j) {
            jmp.push_back(vi(sz(v)-(1<<j)+1));
            FOR(i,0,sz(jmp[j])) jmp[j][i] = cmb(jmp[j-1][i],
                jmp[j-1][i+(1<<(j-1))]);
        }
    }
    int index(int l, int r) {
        assert(l <= r); int d = level(r-l+1);
        return cmb(jmp[d][l],jmp[d][r-(1<<d)+1]); }
    T query(int l, int r) { return v[index(l,r)]; }
};
/*
https://codeforces.com/contest/1983/submission/269275960
*/
struct HLD {
    int n;
    std::vector<int> siz, top, dep, parent, in, out, seq;
    std::vector<std::vector<int>> adj;
    int cur;
    
    HLD() {}
    HLD(int n) {
        init(n);
    }
    void init(int n) {
        this->n = n;
        siz.resize(n);
        top.resize(n);
        dep.resize(n);
        parent.resize(n);
        in.resize(n);
        out.resize(n);
        seq.resize(n);
        cur = 0;
        adj.assign(n, {});
    }
    void addEdge(int u, int v) {
        adj[u].push_back(v);
        adj[v].push_back(u);
    }
    void work(int root = 0) {
        top[root] = root;
        dep[root] = 0;
        parent[root] = -1;
        dfs1(root);
        dfs2(root);
    }
    void dfs1(int u) {
        if (parent[u] != -1) {
            adj[u].erase(std::find(adj[u].begin(), adj[u].end(), parent[u]));
        }
        
        siz[u] = 1;
        for (auto &v : adj[u]) {
            parent[v] = u;
            dep[v] = dep[u] + 1;
            dfs1(v);
            siz[u] += siz[v];
            if (siz[v] > siz[adj[u][0]]) {
                std::swap(v, adj[u][0]);
            }
        }
    }
    void dfs2(int u) {
        in[u] = cur++;
        seq[in[u]] = u;
        for (auto v : adj[u]) {
            top[v] = v == adj[u][0] ? top[u] : v;
            dfs2(v);
        }
        out[u] = cur;
    }
    int lca(int u, int v) {
        while (top[u] != top[v]) {
            if (dep[top[u]] > dep[top[v]]) {
                u = parent[top[u]];
            } else {
                v = parent[top[v]];
            }
        }
        return dep[u] < dep[v] ? u : v;
    }
    
    int dist(int u, int v) {
        return dep[u] + dep[v] - 2 * dep[lca(u, v)];
    }
    
    int jump(int u, int k) {
        if (dep[u] < k) {
            return -1;
        }
        
        int d = dep[u] - k;
        
        while (dep[top[u]] > d) {
            u = parent[top[u]];
        }
        
        return seq[in[u] - dep[u] + d];
    }
    
    bool isAncester(int u, int v) {
        return in[u] <= in[v] && in[v] < out[u];
    }
    
    int rootedParent(int u, int v) {
        std::swap(u, v);
        if (u == v) {
            return u;
        }
        if (!isAncester(u, v)) {
            return parent[u];
        }
        auto it = std::upper_bound(adj[u].begin(), adj[u].end(), v, [&](int x, int y) {
            return in[x] < in[y];
        }) - 1;
        return *it;
    }
    
    int rootedSize(int u, int v) {
        if (u == v) {
            return n;
        }
        if (!isAncester(v, u)) {
            return siz[v];
        }
        return n - siz[rootedParent(u, v)];
    }
    
    int rootedLca(int a, int b, int c) {
        return lca(a, b) ^ lca(b, c) ^ lca(c, a);
    }
};

void solve_tc() {
    int n, m;
    cin>>n>>m;
    if (n==1 && m==1) {
        int a;cin>>a;
        cout<<"-1\n";
        return;
    }
    FOR(i,0,n) FOR(j,0,m) {
        int a;
        cin>>a,a--;
        cout<<(a+1)%(n*m)+1<<" \n"[j == m-1];
    }
}

int main() {
    ios::sync_with_stdio(0);
    cin.tie(0);
    int T=1;
    cin>>T;
    while (T--) {
        solve_tc();
    }
}