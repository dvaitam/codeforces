#include <bits/stdc++.h>
#define fi first
#define se second
#define mp make_pair
#define pb push_back
#define all(v) v.begin(), v.end()
#define sz(x) ((int) (x).size())
using namespace std;
using ll = long long;
using pii = pair<int,int>;
using pll = pair<ll,ll>;
using vi = vector<int>;
using vp = vector<pii>;
using vl = vector<ll>;
using vvi = vector<vi>;
using vvl = vector<vl>;
using vb = vector<bool>;

template<typename A, typename B> ostream& operator<<(ostream &os, const pair<A, B> &p){return os << '(' << p.fi << ", " << p.se << ')';}
template<typename C, typename T = typename enable_if<!is_same<C, string>::value, typename C::value_type>::type>
ostream& operator<<(ostream &os, const C &v){string sep; for(const T &x : v) os << sep << x, sep = " "; return os;}
#define deb(...) logger(#__VA_ARGS__, __VA_ARGS__)
template<typename ...Args>
void logger(string vars, Args&&... values){
    cout << "[Debug]\n\t" << vars << " = ";
    string d = "[";
    (..., (cout << d << values, d = "] ["));
    cout << "]\n";
}

mt19937 rng(chrono::steady_clock::now().time_since_epoch().count());
int unif(int l, int r) { return uniform_int_distribution(l, r)(rng); }


const int inf = INT_MAX;
const ll linf = LLONG_MAX;

const int N = 2e5 + 3;
vi children[N];

struct HLD {
    int n, t;
    vi heavy, p, sub, h, tour, in, head;
    vector <vi> adj;

    HLD(int n_) : n(n_), t(0), heavy(n), p(n), sub(n),
        h(n), tour(n), in(n), head(n), adj(n) {}

    void add_edge(int u, int v) {
        adj[u].pb(v);
        adj[v].pb(u);
    }

    void dfs_prep(int u, int par) {
        sub[u] = 1;
        heavy[u] = -1;
        h[u] = h[par] + 1;
        p[u] = par;
        for(auto v : adj[u]) {
            if(v == par) continue;
            dfs_prep(v, u);
            sub[u] += sub[v];
            if(heavy[u] == -1 || sub[v] > sub[heavy[u]])
                heavy[u] = v;
        }
    }
    void dfs_hld(int u, int par) {
        in[u] = t;
        tour[t++] = u;
        head[u] = u != heavy[par] ? u : head[par];
        if(heavy[u] == -1) return;
        dfs_hld(heavy[u], u);
        children[u].pb(heavy[u]);
        for(auto v : adj[u]) {
            if(v == par || v == heavy[u]) continue;
            children[u].pb(v);
            dfs_hld(v, u);
        }
    }

    int lca(int u, int v) {
        while(head[u] != head[v]) {
            if(h[head[u]] < h[head[v]]) swap(u, v);
            u = p[head[u]];
        }
        if(h[u] > h[v]) swap(u, v);
        return u;
    }

    void init() {
        dfs_prep(0, 0);
        dfs_hld(0, 0);
    }
};


vvi adj;
int leaf1[N], leaf2[N];
vi leaf_to_time;

void rec1(int u, int l, int r) {
    while(sz(adj) <= u) adj.pb({}), leaf_to_time.pb(-1);
    if(u != 0) {
        int par = (u-1)/2;
        adj[u].pb(par);
    }
    //cout << u << " is [" << l << ", " << r << "] for left tree\n";
    if(l == r) {
        leaf_to_time[u] = l;
        leaf1[l] = u;
      //  cout << "leaf1(" << l << ") = " << leaf1[l] << "\n";
    }
    else {
        int m = (l + r)/2;
        rec1(2*u+1, l, m);
        rec1(2*u+2, m+1, r);
    }
}

int off;

void rec2(int u, int l, int r) {
  //  cout << u + off << " is [" << l << ", " << r << "]\n";
    while(sz(adj) <= u + off) adj.pb({}), leaf_to_time.pb(-1);
   // cout << u + off << " is [" << l << ", " << r << "] for right tree\n";
    if(u != 0) {
        int par = (u-1)/2;
        adj[par + off].pb(u + off);
    }
    if(l == r) {
        leaf2[l] = u + off;
       // leaf_to_time[u + off] = l;
      //  cout << "leaf2(" << l << ") = " << leaf2[l] << "\n";
    }
    else {
        int m = (l + r)/2;
        rec2(2*u+1, l, m);
        rec2(2*u+2, m+1, r);
    }
}

void add1(int u, int l, int r, int s, int e, int x) {
    if(r < s || e < l) return;
    if(s <= l && r <= e) {
        adj[u].pb(x);
      //  deb(u, x, "1");
      //cout << u << " -> " << x << " add 1\n";
        return;
    }
    int m = (l + r) / 2;
    add1(2*u+1, l, m, s, e, x);
    add1(2*u+2, m+1, r, s, e, x);
}

void add2(int u, int l, int r, int s, int e, int x) {
    if(r < s || e < l) return;
    if(s <= l && r <= e) {
        adj[x].pb(u + off);
       // deb(x, u + off, "2");
       //cout << x << " -> " << u+off << " add 2\n";
        return;
    }
    int m = (l + r) / 2;
    add2(2*u+1, l, m, s, e, x);
    add2(2*u+2, m+1, r, s, e, x);
}

void add_upper_helper(HLD &G, int u, int v, int c) {
    int z = leaf2[G.in[c]];
   // cout << z << ", which is leaf2 of " << c+1 << ", will be bigger than any stuff on path " << u+1 << " to " << v+1 << "\n";
   // cout << G.in[c] << " is in of c\n";
    while(G.head[u] != G.head[v]) {
        if(G.h[G.head[u]] < G.h[G.head[v]]) swap(u, v);
        assert(G.in[G.head[u]] <=  G.in[u]);
        add1(0, 0, G.n-1, G.in[G.head[u]], G.in[u], z);
        u = G.p[G.head[u]];
    }
    if(G.h[u] > G.h[v]) swap(u, v);
    assert(G.in[u] <= G.in[v]);
    add1(0, 0, G.n-1, G.in[u], G.in[v], z);
}

int first_ancestor(HLD &G, int anc, int dec) {
    int l = 1, r = sz(children[anc]);
    while(l != r) {
        int m = (l + r) / 2;
        if(m == sz(children[anc]) || G.in[children[anc][m]] > G.in[dec]) {
            r = m;
        }
        else {
            l = m + 1;
        }
    }
    return children[anc][r-1];
};



void add_lower_helper(HLD &G, int u, int v, int c) {
    int z = leaf1[G.in[c]];
   // cout << z << ", which is leaf1 of " << c+1 << ", will be smaller than any stuff on path " << u+1 << " to " << v+1 << "\n";
    while(G.head[u] != G.head[v]) {
        if(G.h[G.head[u]] < G.h[G.head[v]]) swap(u, v);
        assert(G.in[G.head[u]] <=  G.in[u]);
        add2(0, 0, G.n-1, G.in[G.head[u]], G.in[u], z);
        u = G.p[G.head[u]];
    }
    if(G.h[u] > G.h[v]) swap(u, v);
    assert(G.in[u] <= G.in[v]);
    add2(0, 0, G.n-1, G.in[u], G.in[v], z);
}


void add_upper(HLD &G, int a, int b, int c) {
    for(int u : {a, b}) {
        if(u != c) {
            if(G.in[c] <= G.in[u] && G.in[c] + G.sub[c] >= G.in[u] + G.sub[u]) {
                int v = first_ancestor(G, c, u);
                add_upper_helper(G, u, v, c);
            }
            else {
                add_upper_helper(G, u, G.p[c], c);
            }
        }
    }
}

void add_lower(HLD &G, int a, int b, int c) {
    for(int u : {a, b}) {
        if(u != c) {
            if(G.in[c] <= G.in[u] && G.in[c] + G.sub[c] >= G.in[u] + G.sub[u]) {
                int v = first_ancestor(G, c, u);
               add_lower_helper(G, u, v, c);
            }
            else {
                add_lower_helper(G, u, G.p[c], c);
            }
        }
    }
}


int main()
{
    #ifndef LOCAL
        ios_base::sync_with_stdio(0); cin.tie(0);
    #endif // LOCAL

    int n, m;
    cin >> n >> m;
    HLD G(n);
    for(int i = 1; i < n; i++) {
        int u, v;
        cin >> u >> v;
        --u; --v;
        G.add_edge(u, v);
    }
    G.init();
    rec1(0, 0, n-1);
    off = sz(adj);
    rec2(0, 0, n-1);
    for(int t = 0; t < n; t++) {
        adj[leaf2[t]].pb(leaf1[t]);
    }

//    for(int i = 0; i < n; i++) {
//        cout << "in(" << i << ") = " << G.in[i] << "\n";
//    }
    vector <array<int, 3>> queries;
    for(int i = 0; i < m; i++) {
        int t, a, b, c;

        cin >> t >> a >> b >> c;

        --a; --b; --c;

        queries.pb({a, b, c});

        if(t == 2) {
            add_upper(G, a, b, c);
        }
        else {
            add_lower(G, a, b, c);
        }
    }
    queue <int> q;
    int nodes = sz(adj);
    vi in(nodes);
    for(int u = 0; u < nodes; u++) {
       // sort(all(adj[u]));
      //  adj[u].resize(unique(all(adj[u])) - begin(adj[u])); //???
        for(int v : adj[u]) {
            in[v]++;
        }
    }
    for(int u = 0; u < nodes; u++) {
        if(in[u] == 0) {
            q.push(u);
        }
    }
    vi ord;
    while(!q.empty()) {
        int u = q.front();
        q.pop();
        ord.pb(u);
        for(int v : adj[u]) {
            --in[v];
            if(in[v] == 0) q.push(v);
        }
    }
    if(sz(ord) != nodes) cout << "-1\n";
    else {
        vi ans(n);
        int ptr = 0;
        vi time_to_node(n);
        for(int i = 0; i < n; i++) time_to_node[G.in[i]] = i;
        for(int u : ord) {
            if(~leaf_to_time[u]) {
                ans[time_to_node[leaf_to_time[u]]] = ++ptr;
            }
        }
        cout << ans << "\n";
    }

   return 0;
}