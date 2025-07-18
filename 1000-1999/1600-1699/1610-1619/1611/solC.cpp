#include<bits/stdc++.h>

using namespace std;





typedef long long int ll;

typedef unsigned long long ull;

typedef long double lld;

typedef string str;



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

#ifndef ONLINE_JUDGE

#define debug(x...) cerr << "[" << #x << "] = ["; _print(x)

#else

#define debug(x...)

#endif





#define pb push_back

#define all(x) (x).begin(), (x).end()

#define print(x) cout << x << "\n"

#define fr first

#define sc second



ll cdiv(ll a, ll b) { return a/b+((a^b)>0&&a%b); } // divide a by b rounded up

ll fdiv(ll a, ll b) { return a/b-((a^b)<0&&a%b); } // divide a by b rounded down



void init()

{

    #ifndef ONLINE_JUDGE

    freopen("input.txt", "r", stdin);

    freopen("output.txt", "w", stdout);

    freopen("error.txt", "w", stderr);

    #endif //

}



/*

	maximum element has been added at the end

	the position of the maximum element has to be 1st or last

	

	4 1 2 3

	

	3 2 1 4

	

	1 2 3

	



*/



void solve()

{

    int n; cin >> n;

    vector<int> a(n);

    for(int i = 0; i<n; i++) cin >> a[i];

    if(a[0]!=n && a[n-1]!=n)

    {

    	print(-1);

    	return;

    }

    if(a[0] == n)

    {

    	for(int i = n-1; i>=0; i--) cout << a[i] << " ";

    }

	else

	{

		for(int i = n-2; i>=0; i--) cout << a[i] << " ";

		cout << n;

	}

	print("");

}







signed main()

{

    ios::sync_with_stdio(false);

    cin.tie(0);

    init();

    int T = 1;

    cin>>T;

    while(T--)

        solve();

    return 0;

}