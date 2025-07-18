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
const int maxm=1<<20;

int n,a[maxn+10],odd[maxm+10],even[maxm+10];
long long ans;

int main()
{
  n=read();
  for(int i=1; i<=n; ++i)
    {
      a[i]=read();
    }
  for(int i=1; i<=n; ++i)
    {
      a[i]^=a[i-1];
    }
  for(int i=1; i<=n; i+=2)
    {
      ans+=odd[a[i]]++;
    }
  for(int i=0; i<=n; i+=2)
    {
      ans+=even[a[i]]++;
    }
  printf("%I64d\n",ans);
  return 0;
}