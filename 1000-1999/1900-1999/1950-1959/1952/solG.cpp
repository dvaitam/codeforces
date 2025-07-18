#include "bits/stdc++.h"
using namespace std;
void __print(int x) {cerr << x;}
void __print(long x) {cerr << x;}
void __print(long long x) {cerr << x;}
void __print(unsigned x) {cerr << x;}
void __print(unsigned long x) {cerr << x;}
void __print(unsigned long long x) {cerr << x;}
void __print(float x) {cerr << x;}
void __print(double x) {cerr << x;}
void __print(long double x) {cerr << x;}
void __print(char x) {cerr << '\'' << x << '\'';}
void __print(const char *x) {cerr << '\"' << x << '\"';}
void __print(const string &x) {cerr << '\"' << x << '\"';}
void __print(bool x) {cerr << (x ? "true" : "false");}
template<typename T, typename V>
void __print(const pair<T, V> &x) {cerr << '{'; __print(x.first); cerr << ','; __print(x.second); cerr << '}';}
template<typename T>
void __print(const T &x) {int f = 0; cerr << '{'; for (auto &i: x) cerr << (f++ ? "," : ""), __print(i); cerr << "}";}
void _print() {cerr << "]\n";}
template <typename T, typename... V>
void _print(T t, V... v) {__print(t); if (sizeof...(v)) cerr << ", "; _print(v...);}
template <typename T, typename V>
void mdebug(map<T,vector<V>>m){
    for(auto x:m){
        cerr << x.F << " : [ " ;
        for(auto c:x.S)
            cerr << c << " ";
        cerr << "]"<<endl ;
    }
}
#define debug(x...) cerr << "[" << #x << "] = ["; _print(x)
//#ifndef ONLINE_JUDGE
//#ifndef LOCAL
#define debug(x...) cerr << "[" << #x << "] = ["; _print(x)
//#else
//#define debug(x...)
//#endif
//#pragma GCC optimize "trapv"
#define int long long 
#define ld long double 
#define endl '\n' 
int Power(int Num,int Index){
    int r=1 ;
    for(int i=Index;i>0;i/=2,Num*=Num){
        if(i%2==1){
            r*=Num ;    
        }    
    }    
    return r ;
}

void solve(){
    ld x;cin>>x;
    cout << log2(x);
} 
signed main(){
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);
    int tc ;
    tc = 1 ;
    // cin >> tc ; 
    while(tc--){
        solve() ;
    }
    return 0;
}