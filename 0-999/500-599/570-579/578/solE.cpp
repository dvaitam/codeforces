#include <bits/stdc++.h>
#define fi first
#define se second
#define mk make_pair
#define pb push_back
#define CH (ch=getchar())
#define Exit(...)    printf(__VA_ARGS__),exit(0)
#define dprintf(...) fprintf(stderr,__VA_ARGS__)
#define rep(i,V)     for(__typeof(*V.begin()) i:  V)
#define For(i,a,b)   for(int i=(int)a;i<=(int)b;i++)
#define Rep(i,a,b)   for(int i=(int)a;i< (int)b;i++)
#define Forn(i,a,b)  for(int i=(int)a;i>=(int)b;i--)
#define pend(x)      ((x)=='\n'||(x)=='\r'||(x)=='\t'||(x)==' ')
using namespace std;
typedef double  db;
typedef long long ll;
typedef pair<int,int> PII;
const int N=100005;
const ll Inf=(ll)1e10;
const int inf=(int)1e9;
const int mo=ll(1e9+7);

inline int IN(){
 char ch;CH; int f=0,x=0;
 for(;pend(ch);CH); if(ch=='-')f=1,CH;
 for(;!pend(ch);CH) x=x*10+ch-'0';
 return (f)?(-x):(x);
}

int Pow(int x,int y,int p){
 int A=1;
 for(;y;y>>=1,x=(ll)x*x%p) if(y&1) A=(ll)A*x%p;
 return A;
}

char s[N];
int A[2][N],n;
int pre[N],Next[N];
bool vis[N];
vector<PII> seq;

void add(int x,int y){
 pre[y]=x;
 Next[x]=y;
}

bool check(){
 if(seq.size()<2)return 0;
 int st=seq.back().fi,ed=seq.back().se;
 seq.pop_back();
 int st2=seq.back().fi,ed2=seq.back().se;
 seq.pb(mk(st,ed));
 if(s[st]!=s[ed2]){
  add(ed2,st);
  seq.pop_back();
  seq.pop_back();
  seq.pb(mk(st2,ed));
  return 1;
 }
 if(s[st2]!=s[ed]){
  add(ed,st2);
  seq.pop_back();
  seq.pop_back();
  seq.pb(mk(st,ed2));
  return 1;
 }
 if(s[st]!=s[ed]){
  if(ed<ed2){
   int z=pre[ed2];
   add(ed,ed2);
   add(ed2,st2);
   Next[z]=0;
   seq.pop_back();
   seq.pop_back();
   seq.pb(mk(st,z));
  }else{
   int z=pre[ed];
   add(ed2,ed);
   add(ed,st);
   Next[z]=0;
   seq.pop_back();
   seq.pop_back();
   seq.pb(mk(st2,z));
  }
  return 1;
 }
 return 0;
}

int main(){
 scanf("%s",s+1);
 n=strlen(s+1);
 For(i,1,n){
  int p=s[i]=='R';
  if(*A[p^1]){
   vis[i]=1;
   add(A[p^1][*A[p^1]],i);
   --*A[p^1];
  }
  A[p][++*A[p]]=i;
 }
 int cnt=-1;
 For(i,1,n){
  if(vis[i])continue;
  ++cnt;
  int j=i;
  while(Next[j]) j=Next[j];
  seq.pb(mk(i,j));
  while(check());
 }
 printf("%d\n",cnt);
 rep(t,seq){
  int j=t.fi;
  for(;j!=t.se;j=Next[j]) printf("%d ",j);
  printf("%d ",j);
 }
 return 0;
}