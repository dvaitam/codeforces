#include <bits/stdc++.h>
using namespace std;
#define int long long
#define INF (int)1e18
#define f first
#define s second

mt19937_64 RNG(chrono::steady_clock::now().time_since_epoch().count());

void Solve() 
{
    int m;
    map <int, int> last;
    cin>>m;
    
    vector <int> adj[m];
    for (int i=0; i<m; i++){
        int n; cin>>n;
        
        for (int j=0; j<n; j++){
            int x; cin>>x;
            adj[i].push_back(x);
            last[x] = i + 1;
        }
    }
    
    vector <int> ans;
    for (int i=0; i<m; i++){
        int win = -1;
        for (auto x: adj[i]){
            if (last[x] == i + 1) win = x;
        }
        
        if (win == -1){
            cout<<-1<<"\n";
            return;
        }
        ans.push_back(win);
    }
    
    for (auto x: ans)
    cout<<x<<" ";
    cout<<"\n";
}

int32_t main() 
{
    auto begin = std::chrono::high_resolution_clock::now();
    ios_base::sync_with_stdio(0);
    cin.tie(0);
    int t = 1;
    cin >> t;
    for(int i = 1; i <= t; i++) 
    {
        //cout << "Case #" << i << ": ";
        Solve();
    }
    auto end = std::chrono::high_resolution_clock::now();
    auto elapsed = std::chrono::duration_cast<std::chrono::nanoseconds>(end - begin);
    cerr << "Time measured: " << elapsed.count() * 1e-9 << " seconds.\n"; 
    return 0;
}