#include <bits/stdc++.h>

using namespace std;



#define loop           	for (int i = 0; i < n; i++)

#define lop(i, n)      	for (int i = 0; i < n; i++)

#define lp(i, k, n)    	for (int i=k;k < n?i < n: i>n;k < n? i+=1: i-=1)

#define trav(a) 		for (auto it = a.begin();  it != a.end();  it++)

#define yesno(x)      	cout<<(x?"YES\n":"NO\n")

#define yes				{cout<< "YES\n"; return;}

#define no             	{cout<< "NO\n"; return;}

#define all(x)         	x.begin(), x.end()

#define travauto(a)		for (auto& it:a)

#define sortall(x)    	sort(all(x))

#define ll             	long long

#define pb             	push_back

#define ss             	second

#define ff             	first

#define endl           	"\n"

typedef pair<int, int> 	pi;

typedef pair<ll, ll>   	pl;

typedef vector<int>    	vi;

typedef vector<ll>     	vl;

typedef vector<pi>     	vpi;

typedef vector<vi>     	vvi;

const ll mod  = 1000000007;

const ll inf  =	1e9;

const ll linf =	1e18;





void __print(int x)    	        {cerr << x;}

void __print(long x)   	        {cerr << x;}

void __print(float x)  	        {cerr << x;}

void __print(double x) 	        {cerr << x;}

void __print(unsigned x)       	{cerr << x;}

void __print(long long x)      	{cerr << x;}

void __print(long double x)    	{cerr << x;}

void __print(unsigned ll x)    	{cerr << x;}

void __print(unsigned long x)  	{cerr << x;}

void __print(const char *x)    	{cerr << '"' << x << '"';}

void __print(const string &x)  	{cerr << '"' << x << '"';}

void __print(char x)           	{cerr << '\''<< x <<'\'';}

void __print(bool x)           	{cerr <<(x?"true":"false");}



template<typename T, typename V>

	void __print(const pair<T, V> &x)

		{cerr << '{'; __print(x.first);cerr<<',';__print(x.second);cerr<<'}';}

template<typename T>

	void __print(const T &x) 

		{int f = 0; cerr << '{'; for (auto &i: x) cerr << (f++ ? "," : ""), __print(i); cerr << "}";}

	void _print()

		{cerr << "]\n";}

template <typename T, typename... V>

	void _print(T t, V... v)

		{__print(t); if (sizeof...(v)) cerr << ", "; _print(v...);}



#ifndef ONLINE_JUDGE

#define debug(x...) cerr << "[" << #x << "] = ["; _print(x)

#else

#define debug(x...)

#endif



#define int long long

bool comp(const pi x, const pi y){

	return x.ss-x.ff < y.ss-y.ff;

}

void Solve()

{

	int n; cin>>n;

	vpi a(n); loop cin>>a[i].ff>>a[i].ss;

	bool flag[n+1] = {0};

	sort(all(a),comp);

	debug(a);

	vi ans(n);

	loop {

		for(int k = a[i].ff;k<=a[i].ss;k++){

			if(!flag[k]) {

				ans[i] = k;

				flag [k] = 1;

				break;

			}

		}

	}	

	loop cout<<a[i].ff<<' '<<a[i].ss<<' '<<ans[i]<<endl;

	cout<<endl;

}



signed main()

{

	ios_base::sync_with_stdio(0), cin.tie(0), cout.tie(0);

	srand(chrono::high_resolution_clock::now().time_since_epoch().count());



	int Testcase = 1;

	cin>>Testcase;



	while (Testcase--) Solve();

	return 0;

}



//by stunnerhash