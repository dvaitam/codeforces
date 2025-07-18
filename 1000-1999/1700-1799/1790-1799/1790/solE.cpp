#include <bits/stdc++.h>
 
using namespace std; 
 
//-------------------------------------------------------------
// Custom set for finding kth element

// #include <ext/pb_ds/assoc_container.hpp>
// #include <ext/pb_ds/tree_policy.hpp>
// using namespace __gnu_pbds;
// typedef tree<int, null_type, less<int>, rb_tree_tag, tree_order_statistics_node_update> ordered_set;

// s.insert(x);
// s.erase(x);
// s.count(x);
// s.find_by_order(k);
// k-th biggest element in set
// s.order_of_key(x);
// how many elements in set are less than x

//-------------------------------------------------------------

const long long INF = 100000000000;
const long long large_prime = 1e9 +7;
using ll = long long;
using ld = long double;
#define ff first
#define ss second
#define pb push_back
#define mp make_pair
#define endl "\n"
#define all(v) (v).begin(),(v).end()
#define set_bits __builtin_popcountll

   
void fastIO(){ 
    ios_base::sync_with_stdio(false);
    cin.tie(NULL);} 
    
bool sortbysecdesc(const pair<pair<ll,ll>,ll> &a, const pair<pair<ll,ll>,ll> &b) {
       if(a.first.second==b.first.second)
            return b.first.first>a.first.first;
       return b.first.second > a.first.second;}
       
vector<ll> sieve(ll n) {ll*arr = new ll[n + 1](); vector<ll> vect; for (ll i = 2; i <= n; i++)if (arr[i] == 0) {vect.push_back(i); for (ll j = (ll(i) * ll(i)); j <= n; j += i)arr[j] = 1;} return vect;}
 


#ifndef ONLINE_JUDGE
#define dbg(a) cerr << endl; cerr << #a << ": "; _print(a); cerr << endl << endl;
#else
#define dbg(a)
#endif
void _print(ll t) {cerr << t;}
void _print(int t) {cerr << t;}
void _print(string t) {cerr << t;}
void _print(char t) {cerr << t;}
void _print(double t) {cerr << t;}
template <class T, class V> void _print(pair <T, V> p);
template <class T> void _print(vector <T> v);
template <class T> void _print(set <T> v);
template <class T, class V> void _print(map <T, V> v);
template <class T> void _print(multiset <T> v);
template <class T, class V> void _print(pair <T, V> p) {cerr << "{"; _print(p.ff); cerr << ","; _print(p.ss); cerr << "}";}
template <class T> void _print(vector <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}
template <class T> void _print(set <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}
template <class T> void _print(multiset <T> v) {cerr << "[ "; for (T i : v) {_print(i); cerr << " ";} cerr << "]";}
template <class T, class V> void _print(map <T, V> v) {cerr << "[ "; for (auto i : v) {_print(i); cerr << " ";} cerr << "]";}
template <class T, class V> void _print(unordered_map <T, V> umap) {cerr << "{ "; for(auto i : umap) {cerr << "{";_print(i.fs); cerr << ", "; _print(i.sn); cerr << "} ";} cerr << "}";}
 
 

// int p[100005];
// int sizes[100005];


// int get(int a ){
//     return p[a]=(p[a]==a ? a : get(p[a]));}

// void merge(int a, int b) {1
//     a = get(a);
//     b = get(b);
//     if (a != b) {
//         if (sizes[a] < sizes[b])
//             swap(a, b);
//         p[b] = a;
//         sizes[a] += sizes[b];
     
//     }
// }



// dp[i][j] -> ith studet aur k candies use ho gye hai 

// dp[i][j] = dp[i-1][j] + dp[i-1][j-1] ... dp[i-1][a[i]]


// bool check (int mid , vector<int> v){
    
    
    
    
    
// }


// 

void makeCombiUtil(vector<vector<int> >& ans, vector<int>& tmp, int n, int left, int k , vector<int>&dis)
{
    
    if (k == 0) {
        ans.push_back(tmp);
        return;
    }
 
    // i iterates from left to n. First time
    // left will be 1
    for (int i = left; i < n; ++i)
    {
        tmp.push_back(dis[i]);
        makeCombiUtil(ans, tmp, n, i + 1, k - 1,dis);
 
        // Popping out last inserted element
        // from the vector
        tmp.pop_back();
    }
}


vector<vector<int> > makeCombi(int n, int k,vector<int>&dis)
{
    vector<vector<int> > ans;
    vector<int> tmp;
    makeCombiUtil(ans, tmp, n, 0, k,dis);
    return ans;
}


ll ans(string a , string b){
    
    ll ans =0LL;
    int n = a.size();
    ll len =0LL;
    
    for(int i=0;i<n;i++){
        
        if(a[i] == b[i]){
            len++;
        }
        
        else{
            ans+= len*(len+1)/2;
            len = 0 ;
        }
        
    }
    ans+= len*(len+1)/2;
    
    return ans;    
}




void solve(int t){
    fastIO();
    
    ll x;
    cin>>x;
    
    if(x%2==1){cout<<-1<<endl;}
        
    else{
        
        if( ( (x/2) ^ (2*x - (x/2 )) ) == x ){
            cout<<(x/2)<< " "<<(2*x - (x/2 ))<<endl;
        }
        else{
            cout<<-1<<endl;
        }
        
    }
  
    
  }
  
    
    

int main(){ 
    fastIO();
    int t =1;
    cin>>t;
    
    for(int i=1;i<=t;i++){
        
        solve(i);
   
    }
}