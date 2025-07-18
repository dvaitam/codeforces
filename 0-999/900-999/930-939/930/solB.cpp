#include <iostream>
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <algorithm>
#include <vector>

using namespace std;
#define MAXN 5000
#define rint register int
#define gc() getchar()
inline int read(rint ans = 0, rint sgn = ' ', rint ch = gc())
{
 for(; ch < '0' || ch > '9'; sgn = ch, ch = gc());
 for(; ch >='0' && ch <='9';(ans*=10)+=ch-'0', ch = gc());
 return sgn-'-'?ans:-ans;
}
char s[MAXN+5]; int book[26]; int n, Ans, Res, Max; vector<int> c[26]; inline char A(int i){return i<=n?s[i]:s[i-n];} 
int main()
{
 scanf("%s",s+1), n = strlen(s+1); for(rint i = 1; i <= n; c[s[i]-'a'].push_back(i), i++);
 for(rint i = 0, j, l; i < 26; i++)
  if(c[i].size())
  {
   for(Max = 0, l = 1; l < n; l++)
   {
    memset(book,0,sizeof book), Res = 0;
    for(j = 0; j < (int)c[i].size(); ++book[A(c[i][j]+l)-'a'], j++);
    for(j = 0; j < 26; j++) if(book[j]==1) ++Res; Max = max(Max,Res);
   } Ans += Max;
  } printf("%.12lf\n",(double)Ans/n); return 0;
}