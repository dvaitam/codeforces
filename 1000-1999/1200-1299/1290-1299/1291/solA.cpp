// Siddharth Ruria
//Linkedin: https://www.linkedin.com/in/ruria-siddharth/
//Codeforces: https://codeforces.com/profile/Sw00sh
//Codechef: https://www.codechef.com/users/airurdis

#include <iostream>
#include <bits/stdc++.h>
// #include <sys/resource.h>
// #include <ext/pb_ds/assoc_container.hpp>
// #include <ext/pb_ds/tree_policy.hpp>
using namespace std;
// using namespace chrono;
// using namespace __gnu_pbds;

// #define ruria 1

//Speed
#define fast ios_base::sync_with_stdio(false);
#define input cin.tie(NULL);
#define output cout.tie(NULL);

//Aliases
using ll= long long;
using lld= long double;
using ull= unsigned long long;

//Constants
const lld pi= 3.141592653589793238;
const ll INF= LONG_LONG_MAX;
const ll mod=1e9+7;

//TypeDef
typedef vector<int> vi;
typedef vector<pair<int, int>> vpi;
typedef pair<ll, ll> pll;
typedef vector<ll> vll;
typedef vector<pll> vpll;
typedef vector<string> vs;
typedef unordered_map<ll,ll> umll;
typedef map<ll,ll> mll;

// Macros
#define F first
#define S second
#define pb push_back
#define mp make_pair
#define fl(i, n) for(int i= 0; i< n; ++i)
#define rtl(i, m, n) for(int i= n;i>= m; --i)
#define yes cout<<"YES\n"
#define no cout<<"NO\n"
#define minus cout << "-1\n"
#define all(v) v.begin(),v.end()

//Debug
#ifdef ruria
#define debug(x) cerr << #x << " "; cerr << x << " "; cerr << endl;
#else
#define debug(x);
#endif

// Operator overloads
template<typename T1, typename T2> // cin >> pair<T1, T2>
istream& operator>>(istream &istream, pair<T1, T2> &p) { return (istream >> p.first >> p.second); }
template<typename T> // cin >> vector<T>
istream& operator>>(istream &istream, vector<T> &v){ for (auto &it: v)cin >> it; return istream; }
template<typename T1, typename T2> // cout << pair<T1, T2>
ostream& operator<<(ostream &ostream, const pair<T1, T2> &p) { return (ostream << p.F << " " << p.S); }
template<typename T> // cout << vector<T>
ostream& operator<<(ostream &ostream, const vector<T> &c) { for (auto &it: c) cout << it << " "; return ostream; }

// Utility functions
template <typename T>
void print(T &&t)  { cout << t << "\n"; }
void printarr(int arr[], int n){ fl(i, n) cout << arr[i] << " "; cout << "\n"; }
template<typename T>
void printvec(vector<T>v){ ll n= v.size(); fl(i, n) cout << v[i] << " "; cout << "\n"; }
template<typename T>  
ll sumvec(vector<T>v){ ll n= v.size(); ll s= 0; fl(i, n) s+= v[i]; return s; }

// Mathematical functions
ll gcd(ll a, ll b){ if (b == 0) return a; return gcd(b, a % b); } //__gcd 
ll lcm(ll a, ll b){ return (a/gcd(a,b)*b); }
ll moduloMultiplication(ll a,ll b,ll mod){ ll res = 0; a %= mod; while (b){ if(b & 1) res = (res + a) % mod; b >>= 1; } return res; }
ll powermod(ll x, ll y, ll p){ ll res = 1; x = x % p; if (x == 0) return 0; while (y > 0){ if(y & 1) res = (res*x) % p; y = y>>1; x = (x*x) % p; } return res; }

//Graph-dfs
// bool gone[MN];
// vector<int> adj[MN];
// void dfs(int loc){
//     gone[loc]=true;
//     for(auto x:adj[loc])if(!gone[x])dfs(x);
// }

//Sorting
bool sorta(const pair<int,int> &a,const pair<int,int> &b){return (a.second < b.second);}
bool sortd(const pair<int,int> &a,const pair<int,int> &b){return (a.second > b.second);}

//Bits
string decToBinary(int n){string s="";int i = 0;while (n > 0) {s =to_string(n % 2)+s;n = n / 2;i++;}return s;}
ll binaryToDecimal(string n){string num = n;ll dec_value = 0;int base = 1;int len = num.length();for(int i = len - 1; i >= 0; i--){if (num[i] == '1')dec_value += base;base = base * 2;}return dec_value;}

//Check
bool isPrime(ll n){if(n<=1)return false;if(n<=3)return true;if(n%2==0||n%3==0)return false;for(int i=5;i*i<=n;i=i+6)if(n%i==0||n%(i+2)==0)return false;return true;}
bool isPowerOfTwo(ll n){if(n==0)return false;return (ceil(log2(n)) == floor(log2(n)));}
bool isPerfectSquare(ll x){if (x >= 0) {ll sr = sqrt(x);return (sr * sr == x);}return false;}

//Code
void solve() {
	int n;
	cin >> n;
	string s;
	cin >> s;
	if(n< 2) {
		cout << -1 << "\n";
		return;
	}
	bool check= false;
	int ct= 0;
	int fi= 0;
	int se= 0;
	for(int i= 0; i< n; ++i) {
		int val= (int)s[i]-'0';
		if(val&1 and fi== 0) {
			ct++;
			fi= val;
		} else if(val&1 and se== 0) {
			ct++;
			se= val;
		}
		if(ct== 2) {
			check= true;
			break;
		}
	}
	if(!check) cout << -1 << "\n";
	else cout << fi << se << "\n";
}	

//Main
int main() {

	fast
	input
	output

	int tt= 1;
	cin >> tt;
    while(tt--) {
    	solve();
    }

    //Uncomment for Kickstart
    // fl(i,t) {
    //     cout<<"Case #"<<i+1<<": ";
    //     solve();
    //     cout<<"\n";
    // }
}

// End