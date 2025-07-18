#include <bits/stdc++.h>
using namespace std;
#define all(arr) arr.begin(), arr.end()
#define forn(i, n) for (ll i = 0; i < n; i++)
#define print(arr) forn(i, arr.size()) { cout << arr[i] << ' '; }                     cout << endl;
#define scan(arr, n) forn(i, n) { long long int x; cin >> x; arr.push_back(x); }
#define ll long long 
#define pb push_back
#include <ext/pb_ds/assoc_container.hpp>
#include <ext/pb_ds/tree_policy.hpp>
using namespace __gnu_pbds;
template <class T> using ordered_set = tree<T, null_type, less<T>, rb_tree_tag, tree_order_statistics_node_update>;
const int M = 1e9 + 7;
double pi = 3.14159265358979323846;
template<typename T, bool maximum_mode = false>
struct RMQ {
    int n = 0;
    vector<vector<T>> range_min;
    RMQ(const vector<T> &values = {}) {
        if (!values.empty())
            build(values);
    }
    static int highest_bit(int x) {
        return x == 0 ? -1 : 31 - __builtin_clz(x);
    }
    static T better(T a, T b) {
        return maximum_mode ? max(a, b) : min(a, b);
    }
    void build(const vector<T> &values) {
        n = int(values.size());
        int levels = highest_bit(n) + 1;
        range_min.resize(levels);
        for (int k = 0; k < levels; k++)
            range_min[k].resize(n - (1 << k) + 1);
        if (n > 0)
            range_min[0] = values;
        for (int k = 1; k < levels; k++)
            for (int i = 0; i <= n - (1 << k); i++)
                range_min[k][i] = better(range_min[k - 1][i], range_min[k - 1][i + (1 << (k - 1))]);
    }
    void show() const {
        for (int i = 0; i < int(range_min.size()); i++) {
            for (int j = 0; j < int(range_min[0].size()); j++) {
                cout << range_min[i][j] <<endl[j + 1 == range_min[0].size()];
            }
        }
    }
    T query_value(int from, int to) const {
        assert(0 <= from && from <= to && to <= n - 1);
        int lg = highest_bit(to - from + 1);
        return better(range_min[lg][from], range_min[lg][to - (1 << lg) + 1]);
    }
};
ll int mul(ll int x,ll int y){
    return (x*1ll*y)%M;
}
ll int power(ll int x,ll int y){
   ll int ans=1;
    while(y>0){
        if(y&1) ans=mul(ans,x);
        x=mul(x,x);
        y=y>>1;
    }
    return ans;
}
ll int divide(ll int x,ll int y){
    return mul(x,power(y,M-2));
}
vector<ll int> getDivisors(ll int n)
{
    vector<ll int> ret; 
    for (ll int i=1; i<=sqrt(n); i++)
    {
       if (n%i == 0)
       {
          if (n/i == i)
              ret.pb(i);
          else {
              ret.pb(i);
              ret.pb(n/i);
          }
       }
    }
return ret;
}
int main(){
int t;
cin >> t;
while (t--){
    int n;cin>>n;
    map<int,int>mp,mp1;
    map<pair<int,int>,int> mp2;
    vector<int> vis(n+1,0);
    vector<int> adj[n+1];
    forn(i,n-1){
        int x,y;cin>>x>>y;
        adj[x].pb(y);adj[y].pb(x);
        mp[x]++;mp[y]++;
        int mi=min(x,y),ma=max(x,y);
        mp2[{mi,ma}]=i+1;
    }
    if(n%3!=0){cout<<-1<<endl;continue;}
    queue<int> q;
    forn(i,n){
        if(mp[i+1]==1) {q.push(i+1);}
        mp1[i+1]=1;
    }
    //mp.clear();
    set<pair<int,int>> st;
    int check=0;
    while(!q.empty()){
        int si=q.size();
        while(si--){
            auto it=q.front();q.pop();
            if(mp1[it]>3){check=1;}
            vis[it]=1;
            for(auto&child:adj[it]){
                if(!vis[child]){
                    mp[child]--;
                    if(mp1[it]==3){
                        int mi=min(child,it),ma=max(it,child);
                        st.insert({mi,ma});
                    }
                    if(mp1[it]!=3)mp1[child]+=mp1[it];
                    if(mp[child]==1) q.push(child);
                }
            }

        }
        if(check==1) break;
    }
    if(check==1) cout<<-1<<endl;
    else{
        cout<<st.size()<<endl;
        if(st.size()==0){cout<<""<<endl;continue;}
        // for(auto&it:st){
        //     cout<<it.first<<" "<<it.second<<endl;
        // }
        for(auto&it:st){
            cout<<mp2[{it.first,it.second}]<<" ";
        }
        cout<<endl;
    }
}
}