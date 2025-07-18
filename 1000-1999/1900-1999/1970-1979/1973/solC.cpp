#include <bits/stdc++.h>
 
#include<ext/pb_ds/assoc_container.hpp>  
#include<ext/pb_ds/tree_policy.hpp>
using namespace __gnu_pbds;
using namespace std;
 
template <typename T>
using ordered_set = tree<T, null_type, less<T>, rb_tree_tag, tree_order_statistics_node_update>;
 
template <typename T>
using ordered_multiset = tree<T, null_type, less_equal<T>, rb_tree_tag, tree_order_statistics_node_update>; 
 
#pragma GCC target("avx2,bmi,bmi2,popcnt,lzcnt")
#pragma GCC optimization ("O3")
#pragma GCC optimization ("unroll-loops")
 
#define M 1000000007
#define INFL 1000000000000000000
#define INF 1000000000
#define pi 3.141592653589793238462643383279
// #define maxn 1000000
#define mp make_pair
#define pb push_back
#define popb pop_back
#define GCD(x, y) __gcd(x, y)
#define all(x) x.begin(), x.end()
#define allr(x) x.rbegin(), x.rend()
#define fix(n) std::fixed<<std::setprecision(n)
#define endl "\n"
#define INPUT(a,n) for(int i = 0; i < n; i++) cin>>a[i];
// #define INPUT(a,n,m) for(int i = 0; i < n; i++) for(int j = 0; j < m; j++) cin>>a[i][j];
#define OUTPUT(a,n) for(int i = 0; i < n; i++) cout<<a[i]<<" ";
      #define fr(n) for(ll i=0;i<n;i++)
typedef long long int ll;
typedef unsigned long long u64;
typedef long double ld;
typedef vector<int> vi;
typedef vector<long long> vll;
typedef pair<int, int> pii;
typedef pair<long long, long long> pll;
typedef vector<pair<int, int>> vpii;
typedef vector<vector<int>> vvi;
typedef vector<vector<ll>> vvl;
typedef vector<pair<long long, long long>> vpll;
typedef map<ll,ll> mpll;
typedef map<int,int> mpii;
 
void __print(int x) {cerr<<x;}
void __print(long x) {cerr<<x;}
void __print(long long x) {cerr<<x;}
void __print(unsigned x) {cerr<<x;}
void __print(unsigned long x) {cerr<<x;}
void __print(unsigned long long x) {cerr<<x;}
void __print(float x) {cerr<<x;}
void __print(double x) {cerr<<x;}
void __print(long double x) {cerr<<x;}
void __print(char x) {cerr<<'\''<<x<<'\'';}
void __print(const char *x) {cerr<<'\"'<<x<<'\"';}
void __print(const string &x) {cerr<<'\"'<<x<<'\"';}
void __print(bool x) {cerr<<(x==1);}
template<typename T,typename V>     void __print(const pair<T,V> &x) {cerr<<'{';__print(x.first);cerr<<',';__print(x.second);cerr<<'}';}
template<typename T>    void __print(const T &x) {int f = 0;cerr<<'{';for(auto &i:x) cerr<<(f++?",":""),__print(i);cerr<<"}";}
void _print() {cerr<<"]\n";}
template <typename T,typename... V> void _print(T t,V... v) {__print(t);if (sizeof...(v)) cerr<<",";_print(v...);}
#ifndef ONLINE_JUDGE
#define debug(x...) cerr<<"["<<#x<<"] = ["; _print(x);
#else
#define debug(x...);
#endif

// const int mxN = 500+5;
// vector<vector<int>> g(2*mxN),h(2*mxN);
// vector<int> vis(2*mxN,0),tout(2*mxN,0),comp(2*mxN,0);
// int tt = 0,n,ct = 0;
// set<array<int,2>> ss;

// void dfs(int u){
//       vis[u] = 1;
//       for(auto v : g[u]){
//             if(!vis[v]){
//                   dfs(v);
//             }
//       }
//       tout[u] = ++tt;
// }

// void dfx(int u){
//       vis[u] = 1;
//       for(auto v : h[u]){
//             if(!vis[v]){
//                   dfx(v);
//             }
//       }
//       ss.erase({tout[u],u});
//       comp[u] = ct;
// }
// void scc(){

//       for(int i = 2; i <= 2*n+1; i++){
            
//             debug("HEREeeee");
//             if(!vis[i])
//                   dfs(i);
//       }
//       for(int i = 2; i<= 2*n+1; i++){
//             if(tout[i]){
//                   ss.insert({tout[i],i});
//             }
//             vis[i] = 0;
//             debug("ERERERERE");
//       }
//       while(!ss.empty()){
//             auto u = (*(--ss.end()))[1];
//             ++ct;
//             dfx(u);
//       }
// }
// void clear(){
//       tt = ct = 0;
//       for(int i = 1; i <= 2*n+2; i++){
//             vis[i] = comp[i] = tout[i] = 0;
//             g[i].clear();
//             h[i].clear();
//       }
// }
void solve(){

      int n;
      cin>>n;
      int a[n];
      int b[n];
      int k;
      vector<array<int,2>> ev,od;

      for(int i = 0; i < n; i++){
            cin>>a[i];
            if(i&1)
                  od.push_back({a[i],i});
            else ev.push_back({a[i],i});
            if(a[i] == 1)
                  k = i;
      }     
      sort(all(od));
      sort(all(ev)); 
      if(k&1){
            k = n;
            for(auto x : ev){
                  b[x[1]] = k;
                  k--;
            }

            for(auto x : od){
                  b[x[1]] = k;
                  k--;
            }
      }
      else{
            k = n;
            for(auto x : od){
                  b[x[1]] = k;
                  k--;
            }
            for(auto x : ev){
                  b[x[1]] = k;
                  k--;
            }
      }
      for(int i = 0; i < n; i++)
            cout<<b[i]<<" ";
      cout<<endl;

}
int main()
{
      #ifndef ONLINE_JUDGE
            freopen("input.txt", "r" , stdin);
            freopen("output.txt", "w", stdout);
      #endif
      ios_base::sync_with_stdio(false);
      cin.tie(NULL);
      cout.tie(NULL);
//-----------------------------------------//    
 
      int t = 1;
      cin >> t;
      // sieve();
 
      while (t--)
      {
            solve();     
            debug("endtc");
      }
 
      cerr << "time taken : " << (float)clock() / CLOCKS_PER_SEC << " secs" << "\n";
      return 0;
 
}
 
// set upper bound --> returns iterator to the number greater than n
// set lower bound --> returns iterator to the number greater than or equal to n
// careful with multisets
// careful while setting limits in binary search