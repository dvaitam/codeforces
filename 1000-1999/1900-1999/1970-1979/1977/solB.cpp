#include <bits/stdc++.h>
using namespace std;

typedef unsigned long long ull;
typedef map<int,int> mint;
typedef pair<int,int> pint;

#define sz(a) s.size();
#define is '\n'
#define yes "YES\n"
#define no "NO\n"
#define all(v) sort(v.begin(), v.end())
#define s(a) sort(a, a + n);
#define pb push_back
#define pf push_front
#define ll long long
#define PODVALMANSURABI ios_base::sync_with_stdio(0); cin.tie(0); cout.tie(0);

const int N=5e5+5;
const int MOD=998244353;
const int MOD1=1e9+7;
const int INF=1e9;
const int IOI=30;
const bool ISSA = true;
const int MAX = 1e9 + 7;

bool was[N];
vector<int> g[N];

void dfs(int u){
    was[u] = 1;
    for (int to : g[u]){
        if (was[to]){
            continue;
        }
        dfs(to);
    }
}
void solve(int x){
    vector<int> vc;
    while (x > 0) {
        if (x % 2 == 0) {
            vc.push_back(0);
        } else {
            if ((x % 4) == 1) {
                vc.push_back(1);
                x -= 1;
            } else {
                vc.push_back(-1);
                x += 1;
            }
        }
        x /= 2;
    }
    cout << vc.size() << is;
    for (int ai : vc) {
        cout << ai << " ";
    }
    cout << is;
}
int main(){
    //freopen("matrix.in", "r", stdin);
    //freopen("matrix.out", "w", stdout);
    PODVALMANSURABI;
    int t;
    cin >> t;
    while(t--){
        int x;
        cin >> x;
        solve(x);
    }
    return 0;
    cout << "IDI POKRAYKAY";
}