#include <bits/stdc++.h>


using namespace std;
#define int long long
mt19937 ran;
vector<vector<pair<int, int>>> g;
vector<vector<int>> up, dob, puti;
vector<pair<int, int>> tree;
vector<int> p, dp, ha, h, am, tin, tout, ad, path, nach;
unordered_map<int, int> la, vis, nizh;
vector<bool> used, iscol;
int t, m, ans;
pair<int, int> res;


void dfs(int v, int pr) {
    used[v] = true;
    tin[v] = t++;
    for (auto u: g[v]) {
        if (u.first == pr) continue;
        if (used[u.first] && tin[u.first] < tin[v]) {
            am[v]++, am[u.first]--, ha[v] += h[u.second], ha[u.first] -= h[u.second];
        } else if (!used[u.first]) {
            dfs(u.first, v);
            ha[v] += ha[u.first], am[v] += am[u.first];
        }
    }
    if (!nizh.count(ha[v])) nizh[ha[v]] = v;
}


bool is_anc(int v, int u) {
    return tin[v] <= tin[u] && tout[u] <= tout[v];
}


int lca(int v, int u) {
    if (is_anc(v, u)) return v;
    for (int i = 20; i >= 0; --i) {
        if (!is_anc(up[v][i], u)) v = up[v][i];
    }
    return up[v][0];
}


void dfs4(int v, int pr) {
    used[v] = true;
    tin[v] = t++;
    up[v][0] = pr;
    for (int i = 1; i < 21; ++i) {
        up[v][i] = up[up[v][i - 1]][i - 1];
    }
    for (auto u: g[v]) {
        if (used[u.first]) continue;
        dfs4(u.first, v);
        p[u.first] = u.second;
    }
    tout[v] = t++;
}


void dfs5(int v) {
    used[v] = true;
    for (auto u: g[v]) {
        if (used[u.first]) continue;
        dfs5(u.first);
        dp[v] += dp[u.first];
    }
}


void push(int v, int l, int r) {
    tree[v].first += ad[v];
    if (r - l != 1) {
        ad[2 * v + 1] += ad[v], ad[2 * v + 2] += ad[v];
    }
    ad[v] = 0;
}


void build(int v, int l, int r) {
    if (r - l == 1) {
        tree[v] = {0, l};
        return;
    }
    int mid = (l + r) / 2;
    build(2 * v + 1, l, mid);
    build(2 * v + 2, mid, r);
    tree[v] = max(tree[2 * v + 1], tree[2 * v + 2]);
}


pair<int, int> get(int v, int l, int r, int ql, int qr) {
    push(v, l, r);
    if (l >= qr || ql >= r) return {-1e18, -1e18};
    if (ql <= l && r <= qr) return tree[v];
    int mid = (l + r) / 2;
    return max(get(2 * v + 1, l, mid, ql, qr), get(2 * v + 2, mid, r, ql, qr));
}


void upd(int v, int l, int r, int ql, int qr, int val) {
    if (ql == qr) return;
    push(v, l, r);
    if (l >= qr || ql >= r) return;
    if (ql <= l && r <= qr) {
        ad[v] += val;
        push(v, l, r);
        return;
    }
    int mid = (l + r) / 2;
    upd(2 * v + 1, l, mid, ql, qr, val);
    upd(2 * v + 2, mid, r, ql, qr, val);
    tree[v] = max(tree[2 * v + 1], tree[2 * v + 2]);
}


void dfs6(int v, int ind) {
    used[v] = true;
    for (auto pref: dob[v]) {
        upd(0, 0, m, 0, pref, -2);
    }
    for (auto u: g[v]) {
        if (used[u.first]) continue;
        int l = -1;
        if (la.count(ha[u.first])) upd(0, 0, m, la[ha[u.first]] + 1, ind, -1e9), l = la[ha[u.first]];
        la[ha[u.first]] = ind;
        int num = nach[u.first];
        upd(0, 0, m, 0, ind, dp[v] - dp[u.first] + 2 * num);
        if (vis.count(ha[u.first]) && ha[u.first]) {
            auto anss = get(0, 0, m, vis[ha[u.first]], ind);
            if (ans < anss.first) {
                ans = anss.first, res.first = path[anss.second], res.second = u.second;
            }
        } else vis[ha[u.first]] = ind;
        path.push_back(u.second);
        dfs6(u.first, ind + 1);
        path.pop_back();
        upd(0, 0, m, 0, ind, -dp[v] + dp[u.first] - 2 * num);
        if (l != -1) upd(0, 0, m, l + 1, ind, 1e9);
    }
    for (auto pref: dob[v]) {
        upd(0, 0, m, 0, pref, 2);
    }
}


void dfs7(int v, int he, int hh){
    used[v] = true;
    if (hh != 0) {
        for (auto pa: puti[v]) {
            dob[lca(pa, nizh[hh])].push_back(he - 1);
        }
    }
    if (hh == 0) hh = ha[v];
    for(auto u:g[v]){
        if (used[u.first]) continue;
        if (is_anc(u.first, nizh[hh])) dfs7(u.first, he + 1, hh);
        else dfs7(u.first, he + 1, 0);
    }
}


void solve() {
    int n, k;
    cin >> n >> m;
    t = 0, ans = 0;
    res = {0, 1};
    la.clear(), vis.clear(), h.clear(), nizh.clear(), iscol.clear();
    g.assign(n, {}), am.assign(n, 0), nach.assign(n, 0), dob.assign(n, {}), tree.assign(4 * m, {}), tin.assign(n, 0), ad.assign(4 * m, 0), used.assign(n, false), ha.assign(n, 0), puti.assign(n, {}), tin.assign(n, 0), tout.assign(n, 0), used.assign(n, false), up.assign(n, vector<int>(21, 0)), dp.assign(n, 0), p.assign(n, -1);
    build(0, 0, m);
    vector<pair<int, int>> ed;
    unordered_map<int, int> wh;
    for (int i = 0; i < m; ++i) {
        int v, u;
        cin >> v >> u;
        ed.push_back({v, u});
        v--, u--;
        g[v].push_back({u, i}), g[u].push_back({v, i});
    }
    for (int i = 0; i < m; ++i) {
        h.push_back(ran());
        wh[h[i]] = i;
    }
    dfs(0, 0);
    used.assign(n, false);
    cin >> k;
    dfs4(0, 0);
    for (int i = 0; i < k; ++i) {
        int v, u;
        cin >> v >> u;
        v--, u--;
        int lc = lca(v, u);
        dp[v]++, dp[u]++, dp[lc] -= 2;
        if (lc != v) {
            int v1 = v;
            for (int j = 20; j >= 0; --j) {
                if (!is_anc(up[v][j], lc)) v = up[v][j];
            }
            nach[v]++;
            puti[v].push_back(v1);
        }
        if (lc != u) {
            int u1 = u;
            for (int j = 20; j >= 0; --j) {
                if (!is_anc(up[u][j], lc)) u = up[u][j];
            }
            nach[u]++;
            puti[u].push_back(u1);
        }
    }
    used.assign(n, false);
    dfs5(0);
    vector<pair<int, int>> dp1;
    for (int i = 0; i < n; ++i) {
        if (am[i] == 1) {
            if (ans < dp[i]) ans = dp[i], res.first = p[i], res.second = wh[ha[i]];
        } else if (am[i] == 0) dp1.push_back({dp[i], p[i]});
    }
    sort(dp1.rbegin(), dp1.rend());
    if (dp1.size() > 1 && ans < dp1[0].first + dp1[1].first){
        ans = dp1[0].first + dp1[1].first;
        res = {dp1[0].second, dp1[1].second};
        res.second = max({res.second, 0ll, 1 - res.first});
    }
    used.assign(n, false);
    dfs7(0, 0, 0);
    used.assign(n, false);
    dfs6(0, 0);
    cout << ans << '\n' << ed[res.first].first << " " << ed[res.first].second << '\n' << ed[res.second].first << " " << ed[res.second].second << '\n';
}


signed main() {
    ios_base::sync_with_stdio(false);
    cin.tie(nullptr), cout.tie(nullptr);
    int y;
    cin >> y;
    while (y--) {
        solve();
    }
}