#include <bits/stdc++.h>



#define x first

#define y second

#define eps 1e-9

#define pb push_back

#define inf 0x3f3f3f3f

#define mod 1000000007

#define sq(_) ((_)*(_))

#define cb(_) ((_)*(_)*(_))

#define bit(x, y) (((x)>>(y))&1)

#define bctz(x) (__builtin_ctz(x))

#define bclz(x) (__builtin_clz(x))

#define bctzl(x) (__builtin_ctzll(x))

#define bclzl(x) (__builtin_clzll(x))

#define bpt(x) (__builtin_popcount(x))

#define bptl(x) (__builtin_popcountll(x))

#define PQ priority_queue<pii, vector<pii>, greater<pii> >



using namespace std;



typedef double DO;

typedef long long INT;

typedef vector<int> VI;

typedef pair<int, int> pii;

typedef pair<pii, int> pi3;



template<typename T, typename U> inline void smin(T &a, U b) {if (a>b) a=b;}

template<typename T, typename U> inline void smax(T &a, U b) {if (a<b) a=b;}



template<class T> inline void gn(T &x) {char c, sg=0; while (c=getchar(), (c>'9' || c<'0') && c!='-');for((c=='-'?sg=1, c=getchar():0), x=0; c>='0' && c<='9'; c=getchar()) x=(x<<1)+(x<<3)+c-'0';if (sg) x=-x;}

template<class T1, class T2> inline void gn(T1 &x1, T2 &x2) {gn(x1), gn(x2);}

template<class T1, class T2, class T3> inline void gn(T1 &x1, T2 &x2, T3 &x3) {gn(x1), gn(x2), gn(x3);}

template<class T1, class T2, class T3, class T4> inline void gn(T1 &x1, T2 &x2, T3 &x3, T4 &x4) {gn(x1), gn(x2), gn(x3), gn(x4);}



template<class T> inline void print(T x) {if (x<0) {putchar('-');return print(-x);}if (x<10) {putchar('0'+x);return;} print(x/10);putchar(x%10+'0');}

template<class T> inline void printsp(T x) {print(x); putchar(' ');}

template<class T> inline void println(T x) {print(x); putchar('\n');}

template<class T1, class T2> inline void print(T1 x1, T2 x2){printsp(x1), println(x2);}

template<class T1, class T2, class T3> inline void print(T1 x1, T2 x2, T3 x3){printsp(x1), printsp(x2), println(x3);}

template<class T1, class T2, class T3, class T4> inline void print(T1 x1, T2 x2, T3 x3, T4 x4){printsp(x1), printsp(x2), printsp(x3), println(x4);}



int power(int a, int b, int m, int ans=1) {

 for (; b; b>>=1, a=1LL*a*a%m) if (b&1) ans=1LL*ans*a%m;

 return ans;

}



#ifndef ONLINE_JUDGE

void debug() {cout << endl;}

template<typename T, typename... Args>void debug(T&& head, Args&&... tail) {cout << head << ' ';debug(tail...);}

#endif



#define NN 200010



char s[NN];

int nxt[NN], x[NN];

int a0[NN], a1[NN], n0, n1;

int vst[NN];



int main() {

#ifndef ONLINE_JUDGE

 freopen("in.in", "r", stdin);

 freopen("out.out", "w", stdout);

#endif



 scanf("%s", s+1);

 int n=strlen(s+1), m=0;

 

 for(int i=1; i<=n; i++) {

  if(s[i]=='0') {

   if(n1) {

    n1--;

    nxt[a1[n1]]=i;

    a0[n0++]=i;

   }

   else {

    a0[n0++]=i;

    m++;

   }

  }

  else {

   if(n0) {

    n0--;

    nxt[a0[n0]]=i;

    a1[n1++]=i;

   }

   else {

    puts("-1");

    return 0;

   }

  }

 }

 

 if(n1) {

  puts("-1");

  return 0;

 }

 

 println(m);

 for(int i=1; i<=n; i++) if(vst[i]==0) {

  x[0]=i;

  int nn=1, u=i;

  while(nxt[u]) {

   x[nn++]=nxt[u];

   u=nxt[u];

   vst[u]=1;

  }

  printsp(nn);

  for(int j=0; j<nn; j++) {

   printsp(x[j]);

  }

  puts("");

 }

}