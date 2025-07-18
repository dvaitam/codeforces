#include<bits/stdc++.h>
using namespace std;
typedef long long ll;
typedef pair<int,int> pii;
typedef pair<ll,ll> pll;
typedef vector<ll> vll;
typedef vector<int> vi;
typedef vector<vector<int>> vvi;
typedef vector<vector<ll>> vvll;
#define fo(i,s,e) for(ll i = s; i < e + 1; i++)
#define revfo(i,e,s) for(ll i = e; i >= s; i--)
#define pb push_back
#define all(v) v.begin(), v.end()
#define all2(v) v.begin()+1, v.end()
#define endl '\n'
const ll N = 2005;
vvll adj(N);
vll vis(N), p(N);
ll n, m, target, f;
void reset(){
    fo(i, 1, n) adj[i].clear();
}
void dfs(ll parent, ll vertex){
    vis[vertex] = 1;
    p[vertex] = parent;
    for(auto neigh: adj[vertex]){
        if(neigh == target && target != parent){
            p[target] = vertex;
            f = 1;
            return;
        }
        if(!vis[neigh]){
            dfs(vertex, neigh);
        }
    }
}
int main(){
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    int t; cin >> t;
    fo(tt, 1, t){
        cin >> n >> m;
        reset();
        vector<pll> ufu;
        fo(i, 1, m){
            ll u, v; cin >> u >> v;
            ufu.pb({u, v});
            adj[u].pb(v);
            adj[v].pb(u);
        }
        // if(tt == 100){
        //     string s = to_string(n);
        //     s += "_";
        //     s += to_string(m);
        //     for(auto x: ufu){
        //         s += "_";
        //         s += to_string(x.first);
        //         s += "_";
        //         s += to_string(x.second);
        //     }
        //     cout << s << endl;
        //     continue;
        // }
        ll ansflag = 0;
        fo(i, 1, n){
            fo(i, 1, n) vis[i] = 0;
            if(adj[i].size() >= 4){
                f = 0;
                fo(i, 1, n) p[i] = 0;
                target = i;
                dfs(0, i);
                if(f == 1){
                    ansflag = 1;
                    break;
                }
            }
        }
        if(ansflag == 1){
            cout << "YES\n";
            ll curr = target;
            vector<pll> ans; 
            vll inans(n + 1);
            while(p[curr] != target){
                inans[curr] = 1;
                inans[p[curr]] = 1;
                ans.pb({curr, p[curr]});
                curr = p[curr];
            }
            ans.pb({curr, p[curr]});
            ll cnt = 0;
            for(auto neigh: adj[target]){
                if(!inans[neigh] && cnt < 2){
                    cnt++;
                    ans.pb({target, neigh});
                }
            }
            cout << ans.size() << endl;
            for(auto x: ans){
                cout << x.first << " " << x.second << endl;
            }
        }
        else cout << "NO\n";
    }
    return 0;
}