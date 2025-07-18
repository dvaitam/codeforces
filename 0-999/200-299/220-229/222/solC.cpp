/******************************************

*    AUTHOR:         BHUVNESH JAIN        *

*    INSTITUITION:   BITS PILANI, PILANI  *

******************************************/

#include <bits/stdc++.h>

using namespace std;

 

typedef long long LL; 

typedef long double LD;



const int MAX	= 1e7 + 1;

const int LIM 	= 1e5 + 5;

const int MOD	= 1e9 + 7;

const LD  EPS	= 1e-10;

const double PI = acos(-1.0);



#define fastio			ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0)

#define Assert(x)		{if(!(x)){cerr<<"Assertion failed at line "<<__LINE__<<": "<<#x<<" = "<<(x)<<"\n";exit(1);}}

#define inchar			getchar//_unlocked

#define outchar(x)		putchar(x)//_unlocked(x)



double tick(){static clock_t oldt;clock_t newt=clock();double diff=1.0*(newt-oldt)/CLOCKS_PER_SEC;oldt = newt;return diff;}

template<typename T> void inpos(T &x){x=0;register T c=inchar();while(((c<48)||(c>57))&&(c!='-'))c=inchar();bool neg=false;\

if(c=='-')neg=true;for(;c<48||c>57;c=inchar());for(;c>47&&c<58;c=inchar())x=(x<<3)+(x<<1)+(c&15);if(neg)x=-x;}

template<typename T> void outpos(T n){if(n<0){outchar('-');n*=-1;}char snum[65];int i=0;do {snum[i++]=n%10+'0';n/=10;}\

while(n);i=i-1;while(i>=0)outchar(snum[i--]);outchar(' ');}

inline void instr(char *str){register char c=0;register int i=0;while(c<33)c=inchar();while (c!='\n'&&c!=' '&&c!=EOF){\

str[i]=c;c=inchar();++i;}str[i]='\0';}

template<typename T> T gcd(T a, T b){return (b?__gcd(a,b):a);}

template<typename T> T lcm(T a, T b){return (a*(b/gcd(a,b)));}

template<typename T> T mod(T a, T b) {return (a<b ? a : a%b);}

template<typename T> T add(T a, T b, T c){T x=a+b;return (x>=c ? x-c : x);}

template<typename T> T mod_neg(T a, T b) {a=mod(a, b);if(a<0){a+=b;}return a;}

LL mulmod(LL a,LL b, LL m){LL q=(LL)(((LD)a*(LD)b)/(LD)m);LL r=a*b-q*m;if(r>m)r%=m;if(r<0)r+=m;return r;}

template<typename T> T expo(T e, T n){T x=1,p=e;while(n){if(n&1)x=x*p;p=p*p;n>>=1;}return x;}

template<typename T> T power(T e, T n, T m){T x=1%m,p=e;while(n){if(n&1)x=mod(x*p,m);p=mod(p*p,m);n>>=1;}return x;}

template<typename T> T extended_euclid(T a, T b, T &x, T &y){T xx=0,yy=1;y=0;x=1;while(b){T q=a/b,t=b;b=a%b;a=t;t=xx;\

xx=x-q*xx;x=t;t=yy;yy=y-q*yy;y=t;}return a;}

template<typename T> T mod_inverse(T a, T n){T x,y;T d = extended_euclid(a, n, x, y);return (d>1?-1:mod_neg(x,n));}



int lp[MAX];

int factors[MAX];

int factors2[MAX];

vector<int> primes;

int a[LIM], b[LIM];



void factor_sieve() {

	lp[1] = 1;

	for (int i=2; i<MAX; ++i) {

		if (lp[i] == 0) {

			lp[i] = i;

			primes.push_back (i);

		}

		for (int j=0; j<primes.size() && primes[j]<=lp[i]; ++j) {

			int x = i * primes[j];

			if (x >= MAX) break;

			lp[x] = primes[j];

		}

	}

}



int main() {

	#ifndef ONLINE_JUDGE

		freopen("inp.txt", "r", stdin);

	#endif

	factor_sieve();

	int n, m, x, y;

	inpos(n), inpos(m);

	for(int i=0; i<n; ++i) {

		inpos(x);

		a[i] = x;

		while(x != 1) {

			y = lp[x];

			while(x % y == 0) {

				x /= y;

				factors[y] += 1;

			}

		}

	}

	for(int i=0; i<m; ++i) {

		inpos(x);

		b[i] = x;

		while(x != 1) {

			y = lp[x];

			while(x % y == 0) {

				x /= y;

				factors2[y] += 1;

			}

		}

	}

	outpos(n), outpos(m), outchar('\n');

	for(int i=0; i<n; ++i) {

		x = a[i];

		while(x != 1) {

			y = lp[x];

			while(x % y == 0) {

				x /= y;

				if (factors2[y] > 0) {

					a[i] /= y;

					factors2[y] -= 1;

				}

			}

		}

		outpos(a[i]);

	}

	outchar('\n');

	for(int i=0; i<m; ++i) {

		x = b[i];

		while(x != 1) {

			y = lp[x];

			while(x % y == 0) {

				x /= y;

				if (factors[y] > 0) {

					b[i] /= y;

					factors[y] -= 1;

				}

			}

		}

		outpos(b[i]);

	}

	// cout<<"Execution time : "<<tick()<<"\n";

	return 0;

}