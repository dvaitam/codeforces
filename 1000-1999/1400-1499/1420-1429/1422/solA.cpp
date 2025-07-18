#include <iostream>

#include <cstdio>

#include <map>

#include <set>

#include <queue>

#include <vector>

#include <cmath>

#include <limits.h>

#include <algorithm>

#include <cstdlib>



using namespace std;

#define left_son node<<1

#define right_son node<<1|1

#define lson left_son,start,mid

#define rson right_son,mid+1,end

#define left_len (mid-start+1)

#define right_len (end-mid)

#define all_len (end-start+1)

#define IOS  std::ios::sync_with_stdio(0);std::cin.tie(0);std::cout.tie(0);system("color 0c")

#define Rep(i,n) for(int i=1;i<=n;i++)

#define rep(i,n) for(int i=0;i<n;i++)

#define mes(n) memset(n,0,sizeof(n))

#define all(x) x.begin(),x.end()

#define E cout<<'\n'

#define F cout<<"above is ok"<<endl;

#define endl '\n'

#define YES cout<<"YES"<<'\n'

#define NO cout<<"NO"<<'\n'

#define vec vector<int>

#define mk make_pair

#define debug(x) cout<<#x<<" = "<<x<<endl

#define over cout<<ans<<endl

#define lowbit(x) ((-x)&x)

#define vi vector<int>

#define P pair<int,int>

#define pll pair<long long,long long>

#define vl vector<long long>

#define mk make_pair

#define pb push_back

#define fi first

#define se second

#define ep(x) emplace_back(x)

#define endll endl

#define ed endl



typedef double D;

typedef long long ll;



char ch;

int n, m, x, k, y, dx, dy, ans;



const int N = 1e6 + 10;

const int inf = 0x3fffffff;

const int maxn = 1e2 + 10;

const int mod = 998244353;



inline ll gcd(ll a, ll b) { return b > 0 ? gcd(b, a % b) : a; }

inline ll lcm(ll a, ll b) { return a * b / gcd(a, b); }

string s;



//for(int i = 1;i<=n;i++){

//		sum[i] = sum[i-1]+v[i-1];

//}





template <typename T> void inline read(T& x) { T f = 1; x = 0; char s = getchar(); while (s < '0' || s > '9') { if (s == '-') f = -1; s = getchar(); }while (s <= '9' && s >= '0') x = (x << 3) + (x << 1) + (s ^ 48), s = getchar(); x *= f; }

template<typename T> void pr(T x) { if (x < 0) putchar('-'), x = -x; if (x > 9) pr(x / 10); putchar(x % 10 + '0'); }







inline bool isPrime(int J) { if (J < 2)return 0; if (J == 2 || J == 3)return 1; if (J % 6 != 1 && J % 6 != 5)return 0; for (int i = 5; i * i <= J; i += 6) { if (J % i == 0 || J % (i + 2) == 0)return 0; }return 1; }





template <typename T>

void in(vector<T>&v){

	for(auto &x:v)

		cin>>x;

}

template <typename T>

void out(vector<T>&v){

	for(auto &x:v){

		cout<<x<<' ';

	}E;

}









void test() {

	

	

}



























































void before(){

	

}









void test_case(int test_num) {//std::cin>>t;

	ll a,b,c,d;

	cin>>a>>b>>c;

	cout<<a+b+c-1<<ed;

}





int main() {

	IOS;

	int t = 1;

	before();

	std::cin>>t;

	int case_num = t;

	while (t--) {

		test_case(case_num - t);

	}

	//cerr << '\n' << (double)clock() / CLOCKS_PER_SEC * 1000 << "ms" << endl;

	return 0;

}