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

const int maxn=100000;
const int mod=1000000007;

int quickpow(int a,int b,int m)
{
  int res=1;
  while(b)
    {
      if(b&1)
        {
          res=1ll*res*a%m;
        }
      a=1ll*a*a%m;
      b>>=1;
    }
  return res;
}

int n,q,sum[maxn+10],f[maxn+10],g[maxn+10];
char a[maxn+10];

int main()
{
  n=read();
  q=read();
  scanf("%s",a+1);
  for(int i=1; i<=n; ++i)
    {
      sum[i]=a[i]-'0'+sum[i-1];
    }
  g[0]=1;
  for(int i=1; i<=n; ++i)
    {
      g[i]=g[i-1]*2;
      if(g[i]>=mod)
        {
          g[i]-=mod;
        }
    }
  for(int i=1; i<=n; ++i)
    {
      f[i]=f[i-1]+g[i-1];
      if(f[i]>=mod)
        {
          f[i]-=mod;
        }
    }
  while(q--)
    {
      int l=read(),r=read(),cnt=sum[r]-sum[l-1];
      int ans=(f[cnt]+1ll*(g[r-l+1-cnt]-1+mod)*f[cnt])%mod;
      printf("%d\n",ans);
    }
  return 0;
}