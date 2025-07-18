#include <algorithm>

#include <bitset>

#include <cassert>

#include <chrono>

#include <climits>

#include <cstdint>

#include <cstdio>

#include <cstdlib>

#include <cstring>

#include <ctime>

#include <deque>

#include <fstream>

#include <functional>

#include <iostream>

#include <iomanip>

#include <map>

#include <numeric>

#include <queue>

#include <random>

#include <set>

#include <stack>

#include <sstream>

#include <tuple>

#include <vector>

#include<bits/stdc++.h>

#define B break

#define int long long int

#define pp pair<int,int>,greater<pair<int,int>>

#define mp(a,b) make_pair(a,b)

#define IIT(a) if(!a)

#define in(a) insert(a)

    #define E empty()

#define ee(a) erase(a)

#define all(v) v.begin(),v.end()

#define w(a) while(a>0)

#define py cout<<"YES"<<"\n"

#define pn cout<<"NO"<<"\n"

#define f(i,a,b) for(int i=a;i<b;i++)

#define ff(i,a,b) for(int i=a;i<=b;i++)

#define in(a) insert(a)

#define fg(i,a,n) for(int i=a;i*i<n;i++)

#define pb(b) push_back(b)

#define fff(i,a,b) for(int i=a;i>=b;i--)

#define I(a,b) if(a==b)

#define IT(a,b) if(a<b)

#define T(a,b) if(a<=b)

#define F first

#define S second



using namespace std;

using namespace chrono;



#ifdef DEBUG

	//#define LOCAL_INPUT_FILE

#else

	//#define USE_FILE_IO

#endif



#ifdef USE_FILE_IO

	#define INPUT_FILE "input.txt"

	#define OUTPUT_FILE "output.txt"

	#define cin ____cin

	#define cout ____cout

	ifstream cin(INPUT_FILE);

	ofstream cout(OUTPUT_FILE);

#else

	#ifdef LOCAL_INPUT_FILE

		#define cin ____cin

		ifstream cin("input.txt");

	#endif

#endif



const int infinity = (int)1e9 + 42;

const int64_t llInfinity = (int64_t)1e18 + 256;

const int module = (int)1e9 + 7; 

const long double eps = 1e-8;



mt19937_64 randGen(system_clock().now().time_since_epoch().count());



inline void raiseError(string errorCode) {

	cerr << "Error : " << errorCode << endl;

	exit(42);

}



inline int64_t gilbertOrder(int x, int y, int pow, int rotate) {

	if (pow == 0) {

		return 0;

	}

	int hpow = 1 << (pow-1);

	int seg = (x < hpow) ? (

		(y < hpow) ? 0 : 3

	) : (

		(y < hpow) ? 1 : 2

	);

	seg = (seg + rotate) & 3;

	const int rotateDelta[4] = {3, 0, 0, 1};

	int nx = x & (x ^ hpow), ny = y & (y ^ hpow);

	int nrot = (rotate + rotateDelta[seg]) & 3;

	int64_t subSquareSize = int64_t(1) << (2*pow - 2);

	int64_t ans = seg * subSquareSize;

	int64_t add = gilbertOrder(nx, ny, pow-1, nrot);

	ans += (seg == 1 || seg == 2) ? add : (subSquareSize - add - 1);

	return ans;

}



struct Query {

	int l, r, idx;

	int64_t ord;

	

	inline void calcOrder() {

		ord = gilbertOrder(l, r, 21, 0);

	}

};



inline bool operator<(const Query &a, const Query &b) {

	return a.ord < b.ord;

}



// https://github.com/kth-competitive-programming/kactl/blob/main/content/data-structures/FenwickTree.h



typedef long long ll;



struct FT {

	int n;

	vector<ll> s;

	FT(int n) : s(n), n(n) {}

	void update(int pos, ll dif) { // a[pos] += dif

		for (; pos < n; pos |= pos + 1) s[pos] += dif;

	}

	ll query(int pos) { // sum of values in [0, pos)

		ll res = 0;

		for (; pos > 0; pos &= pos - 1) res += s[pos-1];

		return res;

	}

	int lower_bound(ll sum) {// min pos st sum of [0, pos] >= sum

		// Returns n if no sum is >= sum, or -1 if empty sum is.

		if (sum <= 0) return -1;

		int pos = 0;

		for (int pw = 1 << 25; pw; pw >>= 1) {

			if (pos + pw <= n && s[pos + pw-1] < sum)

				pos += pw, sum -= s[pos-1];

		}

		return pos;

	}

};

















using namespace std;





int Gcd(int a,int b){

    if(b==0)

    return a;

    else

    return Gcd(b,a%b);

}









void printst(string s){

     cout<<s<<" ";

}



void prints(int s){

     cout<<s<<" ";

}





void endLine(){

    

    cout<<"\n";

}







void sorty(int *a,int n){

    sort(a,a+n);

}





void soort(vector<int>&z){

      sort(z.begin(),z.end());

}





struct custom_comparator {

    bool operator()(const pair<int, int>& a,

                    const pair<int, int>& b) const

    {

        return less_comparator(minmax(a.first, a.second),

                              minmax(b.first, b.second));

    }



    less<pair<int, int>> less_comparator;

};





int k= 1000000007;





int fact(int v){

    int f=1;

    f(i,2,v+1)

     f=(f*i)%k;

     return f%k;

}







int gcd(int a,int b) {

  int temp;

  while(b > 0) {

      temp = b;

      b = a % b;

      a = temp;

  }

  return a;

}













void primeFactors(int n,vector<int>&v)

{

    while (n % 2 == 0)

    {

        v.push_back(2);

        // cout << 2 << " ";

        n = n/2;

    }

 

  

    for (int i = 3; i <= sqrt(n); i = i + 2)

    {

        while (n % i == 0)

        {

            v.push_back(i);

            // cout << i << " ";

            n = n/i;

        }

    }



    if (n > 2)

        // cout << n << " ";

        v.push_back(n);

    

}







void solve(){

  int n,m,x,y,f=0,s;

  cin>>n;

  vector<int>v(n),p(n),q(n),pair(n+1),vis(n+1,0),c(n+1,0);

  map<int,int>mp;

  f(i,0,n){

     cin>>v[i];

  }

  priority_queue<int>pq;

  f(i,0,n)

    c[v[i]]++;

  f(i,1,n+1){

      if(!c[i])

         pq.push(i);

  }

  fff(i,n,1){

      if(c[i]==1){

          pair[i]=i;

      }else if(c[i]==2){

          if(pq.top()<i){

              pair[i]=pq.top();

              pair[pq.top()]=i;

              pq.pop();

          }else{

              f=1;

              break;

          }

      }else if(c[i]>2){

          f=1;

          break;

      }

  }

  if(f)

    printst("NO");

    else{

        printst("YES");

        endLine();

        f(i,0,n){

            if(!vis[v[i]]){

                p[i]=v[i];

                q[i]=pair[v[i]];

                vis[v[i]]=1;

            }else{

                p[i]=pair[v[i]];

                q[i]=v[i];

            }

        }

        f(i,0,n)

          prints(p[i]);

        endLine();

        f(i,0,n)

          prints(q[i]);

    }

  endLine();

}









signed main(){

    ios_base::sync_with_stdio(false);

    cin.tie(NULL);

    cout.tie(NULL);

    int test;

    test=1;

    cin>>test;

    

    while(test--)

    solve();

    return 0;

}