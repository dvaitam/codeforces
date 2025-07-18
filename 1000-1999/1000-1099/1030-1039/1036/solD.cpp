#include <cstdio>

int read()
{
  int x=0,f=1;
  char ch=getchar();
  while((ch<'0')||(ch>'9'))
    {
      if(ch=='-')
        {
          f=-f;
        }
      ch=getchar();
    }
  while((ch>='0')&&(ch<='9'))
    {
      x=x*10+ch-'0';
      ch=getchar();
    }
  return x*f;
}

const int maxn=300000;

int n,m,a[maxn+10],b[maxn+10],ans;
long long suma,sumb;

int main()
{
  n=read();
  for(int i=1; i<=n; ++i)
    {
      a[i]=read();
      suma+=a[i];
    }
  m=read();
  for(int i=1; i<=m; ++i)
    {
      b[i]=read();
      sumb+=b[i];
    }
  if(suma!=sumb)
    {
      puts("-1");
      return 0;
    }
  suma=sumb=0;
  int l=1,r=1;
  while((l<=n)||(r<=m))
    {
      if(suma<=sumb)
        {
          suma+=a[l++];
        }
      else
        {
          sumb+=b[r++];
        }
      if(suma==sumb)
        {
          ++ans;
        }
    }
  printf("%d\n",ans);
  return 0;
}