#include <bits/stdc++.h>

using namespace std;
using ll = long long;
using ull = unsigned long long;
using ld = long double;
constexpr char LF = '\n';
constexpr char SP = ' ';

int get_inv(int v) {
    int lg = __lg(v) + 1;
    int inv = (1 << lg) - v;
    if (inv == v) {
        inv = 0;
    }
    return inv;
}

int main() {
    ios_base::sync_with_stdio(false);
    cin.tie(0);

    int n;
    cin >> n;
    vector<int> a(n);
    for (auto &v: a)
        cin >> v;

    vector<vector<int>> a_p(n);
    for (int i = 0; i < n; i++) {
        a_p[i].push_back(a[i]);
        while (a_p[i].back() != 0) {
            a_p[i].push_back(get_inv(a_p[i].back()));
        }
        reverse(a_p[i].begin(), a_p[i].end());
    }
    
    int from = 0;
    int max_dist = -1;
    int max_dist_v = -1;
    for (int i = 0; i < n; i++) {
        if (i == from) 
            continue;

        int shared = 0;
        while (shared < min(a_p[from].size(), a_p[i].size()) && a_p[from][shared] == a_p[i][shared]) {
            shared++;
        }
        int dist = a_p[from].size() + a_p[i].size() - shared*2;
        // cout << "dist(" << a[from] << ", " << a[i] << "): " << dist << LF;
        
        if (dist > max_dist) {
            max_dist = dist;
            max_dist_v = i;
        }
    }

    from = max_dist_v;
    max_dist_v = -1;
    max_dist = -1;

    for (int i = 0; i < n; i++) {
        if (i == from) 
            continue;

        int shared = 0;
        while (shared < min(a_p[from].size(), a_p[i].size()) && a_p[from][shared] == a_p[i][shared]) {
            shared++;
        }
        int dist = a_p[from].size() + a_p[i].size() - shared*2;
        
        if (dist > max_dist) {
            max_dist = dist;
            max_dist_v = i;
        }
    }

    from++; max_dist_v++;
    cout << max_dist_v << SP << from << SP << max_dist << LF;
}