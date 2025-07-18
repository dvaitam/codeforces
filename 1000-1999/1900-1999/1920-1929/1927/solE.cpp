//******************************************************************
//						Author: Ritik Arora
#include<bits/stdc++.h>
using namespace std;
#define fastio() ios_base::sync_with_stdio(false);cin.tie(NULL);cout.tie(NULL)

// #include <ext/pb_ds/assoc_container.hpp> // Common file
// using namespace __gnu_pbds;
// typedef tree<ll, null_type, less<ll>, rb_tree_tag,tree_order_statistics_node_update> ordered_set; // find_by_order, order_of_key

typedef long long ll;
typedef unsigned long long ull;
typedef long double lld;

#define pp pair<int,int>
#define ppl pair<ll,ll>
#define ff first
#define ss second
#define set_bits __builtin_popcountll
#define nline '\n'
#define inf LLONG_MAX
#define vi vector<int>
#define vvi vector<vector<int>> 
#define vl vector<long long>
#define vvl vector<vector<long long>>
ll mod1=1e9+7;
ll mod2=998244353;

#ifndef ONLINE_JUDGE
#define debug(x) cerr << #x <<" "; _print(x); cerr << endl;
#else
#define debug(x)
#endif

void _print(ll t) {cerr << t;}
void _print(int t) {cerr << t;}
void _print(string t) {cerr << t;}
void _print(char t) {cerr << t;}
void _print(lld t) {cerr << t;}
void _print(double t) {cerr << t;}
void _print(ull t) {cerr << t;}

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

bool revsort(ll a, ll b) {return a > b;}
ll gcd(ll a, ll b) {if (b > a) {return gcd(b, a);} if (b == 0) {return a;} return gcd(b, a % b);}
ll expo(ll a, ll b, ll mod) {ll res = 1; while (b > 0) {if (b & 1)res = (res * a) % mod; a = (a * a) % mod; b = b >> 1;} return res;}

vector<ll> sieve(int n) {int*arr = new int[n + 1](); vector<ll> vect; for (int i = 2; i <= n; i++)if (arr[i] == 0) {vect.push_back(i); for (int j = 2 * i; j <= n; j += i)arr[j] = 1;} return vect;}

ll mminvprime(ll a, ll b) {return expo(a, b - 2, b);}
ll mod_add(ll a, ll b, ll m) {a = a % m; b = b % m; return (((a + b) % m) + m) % m;}
ll mod_mul(ll a, ll b, ll m) {a = a % m; b = b % m; return (((a * b) % m) + m) % m;}
ll mod_sub(ll a, ll b, ll m) {a = a % m; b = b % m; return (((a - b) % m) + m) % m;}
ll mod_div(ll a, ll b, ll m) {a = a % m; b = b % m; return (mod_mul(a, mminvprime(b, m), m) + m) % m;}  //only for prime m

// class StringHash {
    // public:
    // 	vector<long long>ps1,ps2;
    // 	long long Q1 = 271, Q2 = 277, M1 = 1000000007, M2 = 998244353;
    // 	StringHash(string s) {
    // 		ps1 = vector<long long>(s.size()+1); ps2 = vector<long long>(s.size()+1);
    // 		for (int i = 1; i <= s.size(); i++) {
    // 			long long c = s[i-1] + 1;
    // 			ps1[i] = ((Q1 * ps1[i-1]) + c)%M1;
    // 			ps2[i] = ((Q2 * ps2[i-1]) + c)%M2;
    // 		}
    // 	}
    // 	long long int powxy(long long int x, long long int y, long long M) {
    // 		if (y == 0) return 1;
    // 		if (y%2 == 1) return (x*powxy(x, y-1, M))%M;
    // 		long long int t = powxy(x, y/2, M);
    // 		return (t*t)%M;
    // 	}
    // 	long long substrHash1(int firstIndex, int lastIndex) {
    // 		long long rem = (powxy(Q1,lastIndex-firstIndex+1,M1) * ps1[firstIndex])%M1;
    // 		return (ps1[lastIndex+1] - rem + M1)%M1;
    // 	}
    // 	long long substrHash2(int firstIndex, int lastIndex) {
    // 		long long rem = (powxy(Q2,lastIndex-firstIndex+1,M2) * ps2[firstIndex])%M2;
    // 		return (ps2[lastIndex+1] - rem + M2)%M2;
    // 	}
    // 	pair<long long, long long> substrHash(int firstIndex, int lastIndex) {
    // 		return {substrHash1(firstIndex, lastIndex), substrHash2(firstIndex, lastIndex)};
    // 	}
// };

// class Dsu {
// public:
//     ll v;
//     vector<ll> parent, size;
//     Dsu(ll n) {
//         v=n;
//         parent.resize(v + 1);
//         size.resize(v + 1, 1);
//         for (int i = 0 ; i <= v ; i++) {
//             parent[i] = i;
//         }
//     }
//     int find(int node) {
//         return parent[node] = (parent[node]==node) ? node : find(parent[node]);
//     }
//     void Union(int a, int b) {
//         a = find(a);
//         b = find(b);
//         if (a==b) return;
//         if (size[a] >= size[b]) {
//             parent[b] = a;
//             size[a] += size[b];
//         }
//         else {
//             parent[a] = b;
//             size[b] += size[a];
//         }
//     }
// };

void solve()
{
    int n,k;
    cin>>n>>k;
    vector<int> a(n);
    int start=1;
    int end=n;
    for(int i=0;i<=k-1;i++)
    {
        for(int j=i;j<n;j+=k)
        {
            if(i&1)
            {
                a[j]=end--;
            }
            else{
                a[j]=start++;
            }
        }
        
    }
    for(auto i:a)
    {
        cout<<i<<" ";
    }
    cout<<endl;
}
int main()
{
    //StringHash hash(word);
    //hash.substrHash(idx1,idx2)==hash.substrHash(idx3,idx4);
    #ifndef ONLINE_JUDGE
        freopen("Error.txt", "w", stderr);
    #endif
    fastio();
    int t;
    cin>>t;
    while(t--)
    {
        solve();
    }
    return 0;
}