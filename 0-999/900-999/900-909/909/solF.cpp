#include <bits/stdc++.h>



#define inf 0x3f3f3f3f

#define MOD 1000000007

#define pb push_back

#define sq(x) ((x)*(x))

#define x first

#define y second

#define eps 1e-9

#define bpt(x) (__builtin_popcount(x))

#define bptl(x) (__builtin_popcountll(x))

#define bit(x, y) (((x)>>(y))&1)

#define bclz(x) (__builtin_clz(x))

#define bclzl(x) (__builtin_clzll(x))

#define bctz(x) (__builtin_ctz(x))

#define bctzl(x) (__builtin_ctzll(x))

#define sgn(x) (((x)>eps)-((x)<-eps))

#define PQ priority_queue<pii, vector<pii>, greater<pii> >



using namespace std;

typedef long long INT;

typedef vector<int> VI;

typedef pair<int, int> pii;

typedef pair<pii, int> pi3;

typedef double DO;



template<typename T, typename U> inline void smin(T &a, U b) {if (a>b) a=b;}

template<typename T, typename U> inline void smax(T &a, U b) {if (a<b) a=b;}



template<class T>inline void gn(T &x) {char c, sg=0; while (c=getchar(), (c>'9' || c<'0') && c!='-');for((c=='-'?sg=1, c=getchar():0), x=0; c>='0' && c<='9'; c=getchar()) x=(x<<1)+(x<<3)+c-'0';if (sg) x=-x;}

template<class T, class T1>inline void gn(T &x, T1 &y) {gn(x); gn(y);}

template<class T, class T1, class T2>inline void gn(T &x, T1 &y, T2 &z) {gn(x); gn(y); gn(z);}

template<class T>inline void print(T x) {if (x<0) {putchar('-');return print(-x);}if (x<10) {putchar('0'+x);return;} print(x/10);putchar(x%10+'0');}

template<class T>inline void printsp(T x) {print(x); putchar(' ');}

template<class T>inline void println(T x) {print(x); putchar('\n');}

template<class T, class U>inline void print(T x, U y) {printsp(x); println(y);}

template<class T, class U, class V>inline void print(T x, U y, V z) {printsp(x); printsp(y); println(z);}



int power(int a, int b, int m, int ans=1) {

 for (; b; b>>=1, a=1LL*a*a%m) if (b&1) ans=1LL*ans*a%m;

 return ans;

}



#define NN 101010



int p[NN];



int main() {

#ifndef ONLINE_JUDGE

 freopen("in.in", "r", stdin);

 freopen("out.out", "w", stdout);

#endif

 

 int n, t;

 gn(n);t=n;

 if (n&1) puts("NO\n");

 else {

  puts("YES");

  int nn=n;

  while (n) {

   int m=1<<(31-bclz(n));

   n-=m;

   for (int i=0; i<=n; i++) p[m+i]=m-1-i, p[m-1-i]=m+i;

   n=m-n-2;

  }

  for (int i=1; i<=nn; i++) printsp(p[i]);puts("");

 }



 n=t;

 memset(p, 0, sizeof p);

 if (n<6 || bpt(n)==1) puts("NO");

 else {

  puts("YES");

  if (n==7) {

   puts("5 3 2 6 1 7 4");

   return 0;

  }

  int m=1<<(31-bclz(n)), nn=n;

  p[m]=n; p[n]=m;

  n-=m+1;

  if (n&1) {

   for (int i=1; i<=n; i++) p[i]=m+i, p[i+m]=i;

   for (int i=n+1; i<m-1; i+=2) p[i]=i+1, p[i+1]=i;

  } else {

   for (int i=1; i<=n; i+=2) p[i+m]=i+m+1, p[i+m+1]=i+m;

   p[1]=m-1; p[m-1]=1;

   for (int i=2; i<m-4; i+=2) p[i]=i+1, p[i+1]=i;

   p[m-2]=m-3; p[m-3]=m-4; p[m-4]=m-2;

  }

  for (int i=1; i<=nn; i++) printsp(p[i]);puts("");

  for (int i=1; i<=nn; i++) assert(p[i]&i);

 }



 return 0;

}