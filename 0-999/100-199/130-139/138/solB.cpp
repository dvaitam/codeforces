#include <bits/stdc++.h>
using namespace std;

#define TRACE
#ifdef TRACE
#define TR(...) __f(#__VA_ARGS__, __VA_ARGS__)
template <typename Arg1>
void __f(const char* name, Arg1&& arg1){
  cerr << name << " : " << arg1 << std::endl;
}
template <typename Arg1, typename... Args>
void __f(const char* names, Arg1&& arg1, Args&&... args){
  const char* comma = strchr(names + 1, ',');cerr.write(names, comma - names) << " : " << arg1<<" | ";__f(comma+1, args...);
}
#else
#define TR(...)
#endif

typedef long long                LL;
typedef vector < int >           VI;
typedef pair < int,int >         II;
typedef vector < II >            VII;

#define MOD                      1000000007
#define EPS                      1e-12
#define N                        200100
#define PB                       push_back
#define MP                       make_pair
#define F                        first 
#define S                        second
#define ALL(v)                   v.begin(),v.end()
#define SZ(a)                    (int)a.size()
#define FILL(a,b)                memset(a,b,sizeof(a))
#define SI(n)                    scanf("%d",&n)
#define SLL(n)                   scanf("%lld",&n)
#define PLLN(n)                  printf("%lld\n",n)
#define PIN(n)                   printf("%d\n",n)
#define REP(i,j,n)               for(LL i=j;i<n;i++)
#define PER(i,j,n)               for(LL i=n-1;i>=j;i--)
#define endl                     '\n'
#define fast_io                  ios_base::sync_with_stdio(false);cin.tie(NULL)

int cnt[10] , c1[10] , c2[10];

inline int Find(int i , int j) {
  REP(x,0,10) c1[x] = cnt[x] , c2[x] = cnt[x];
  if(c1[i] == 0 || c2[j] == 0) return 0 ;
  c1[i] --; c2[j] --;
  int ret = 1;
  REP(i,0,10)
    REP(j,0,10)
    if(i + j == 9)
      ret += min(c1[i] , c2[j]);
  return ret;
}

int main() {
  fast_io;
  string n , a1 , a2; cin >> n;
  for(char i : n)
    cnt[i - '0'] ++;
  int mx = 0;
  while(cnt[0] > cnt[9]) {
    cnt[0] --;
    a1.PB('0');
    a2.PB('0');
  }
  int mi = -1 , mj = -1;
  REP(i,0,10) {
    REP(j,0,10) {
      if((i + j) != 10) continue;
      int x = Find(i , j);
      if(mx < x) {
        mx = x;
        mi = i , mj = j;
      }
    }
  }
  REP(i,0,10) c1[i] = c2[i] = cnt[i];
  if(mi != -1 && mj != -1) {
    a1.PB(char(mi + '0'));
    a2.PB(char(mj + '0'));
    c1[mi] --;
    c2[mj] --;
    REP(i,0,10) REP(j,0,10) {
      if(i + j == 9) {
        while(c1[i] > 0 && c2[j] > 0) {
          c1[i] --; c2[j] --;
          a1.PB(char(i + '0')); a2.PB(char(j + '0'));
        }
      }
    }
    REP(i,0,10)
      while(c1[i] > 0) {
        c1[i] --;
        a1.PB(char(i + '0'));
      }
    REP(i,0,10)
      while(c2[i] > 0) {
        c2[i] --;
        a2.PB(char(i + '0'));
      }
    reverse(ALL(a1));
    reverse(ALL(a2));
    cout << a1 << endl << a2 << endl;
  }
  else {
    if(a1.empty())
      cout << n << endl << n << endl;
    else {
      REP(i,0,10)
        while(c1[i] > 0) {
          c1[i] --;
          a1.PB(char(i + '0'));
        }
      REP(i,0,10)
        while(c2[i] > 0) {
          c2[i] --;
          a2.PB(char(i + '0'));
        }
      reverse(ALL(a1));
      reverse(ALL(a2));
      cout << a1 << endl << a2 << endl;
    }
  }
  return 0;
}