#include <bits/stdc++.h>
 
using namespace std;
 
typedef long long ll;
typedef unsigned long long ull;
 
#define MASK(i) (1ULL << (i))
#define GETBIT(mask, i) (((mask) >> (i)) & 1)
#define ALL(v) (v).begin(), (v).end()
 
ll max(ll a, ll b){return (a > b) ? a : b;}
ll min(ll a, ll b){return (a < b) ? a : b;}
ll gcd(ll a, ll b){return __gcd(a, b);}
ll lcm(ll a, ll b){return a / gcd(a, b) * b;}
 
ll LASTBIT(ll mask){return (mask) & (-mask);}
int pop_cnt(ull mask){return __builtin_popcountll(mask);}
int ctz(ull mask){return __builtin_ctzll(mask);}
int logOf(ull mask){return 63 - __builtin_clzll(mask);}
 
mt19937_64 rng(chrono::high_resolution_clock::now().time_since_epoch().count());
ll rngesus(ll l, ll r){return l + (ull) rng() % (r - l + 1);}
double rngesus_d(double l, double r){
    double cur = rngesus(0, MASK(60) - 1);
    cur /= MASK(60) - 1;
    return l + cur * (r - l);
}
 
template <class T1, class T2>
    bool maximize(T1 &a, T2 b){
        if (a < b) {a = b; return true;}
        return false;
    }
 
template <class T1, class T2>
    bool minimize(T1 &a, T2 b){
        if (a > b) {a = b; return true;}
        return false;
    }
 
template <class T>
    void printArr(T container, string separator = " ", string finish = "\n", ostream &out = cout){
        for(auto item: container) out << item << separator;
        out << finish;
    }
 
template <class T>
    void remove_dup(vector<T> &a){
        sort(ALL(a));
        a.resize(unique(ALL(a)) - a.begin());
    }

const int N = 2008;

vector<int> solve_gauss(vector<bitset<N>> a, int n, int m){
    vector<bool> banned(n);
    for(int j = 0; j < m; ++j){
        int idx = -1;
        for(int i = 0; i < n; ++i) if (!banned[i] && a[i][j]) 
            idx = i;
        if (idx == -1) continue;
        banned[idx] = true;
        for(int i = 0; i < n; ++i) if (!banned[i] && a[i][j]) 
            a[i] ^= a[idx];
    }

    vector<pair<int, int>> deg;
    for(int i = 0; i < n; ++i){
        int j = a[i]._Find_first();
        if (j == m) assert(false);
        if (j > m) continue;
        deg.push_back({j, i});
    }
    sort(ALL(deg), greater<pair<int, int>>());

    vector<int> ans(m);
    for(pair<int, int> item: deg){
        int i = item.second, j = item.first;
        int cur = a[i][m];
        for(int k = j+1; k < n; ++k) cur ^= a[i][k] * ans[k];

        ans[j] = cur;
    }

    return ans;
}

void solve(){
    int n, m; cin>> n >> m;
    vector<bitset<N>> a(n);
    for(int i = 0; i < m; ++i){
        int u, v; cin >> u >> v;
        u--; v--;
        a[u][n] = !a[u][n]; a[v][n] = !a[v][n];
        a[u][u] = !a[u][u]; a[v][u] = !a[v][u];
        a[u][v] = !a[u][v]; a[v][v] = !a[v][v];
    } 
    bool check = true;
    for(int i = 0; i < n; ++i){
        if (a[i][n]){
            check = false;
        }
    }

    if (check){
        cout << 1 << "\n";
        for(int i = 1; i <= n; ++i) cout << 1 << " ";
        cout << "\n";
    }
    else{
        vector<int> ans = solve_gauss(a, n, n);
        for(int &i: ans) i++;
        cout << 2 << "\n";
        printArr(ans);
    }
}

int main(void){
    ios::sync_with_stdio(0);cin.tie(0); cout.tie(0);

    clock_t start = clock();

    int t; cin >> t;
    while(t--) solve();

    cerr << "Time elapsed: " << clock() - start << " ms\n";
    return 0;
}