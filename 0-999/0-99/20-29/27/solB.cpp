#include <iostream>
#include <algorithm>
using namespace std;
int h[111],i,j,n,m,k,a[111],b[111];
pair <int,int> p[111];
int main()
{
    int x,y;
    cin>>n;
    for (i=1;i<=n;i++) h[i]=n-1,a[i]=50;
    for (i=1;i<n*(n-1)/2;i++){
      cin>>x>>y;
      h[x]--;h[y]--;
      a[x]++;a[y]--;
    }
  x=0;y=0;
    for (i=1;i<=n;i++)
    if (h[i])
     {
      if (x==0) x=i;else y=i;
     } else
     b[a[i]]=1;

   if (b[a[x]-1]==0 && b[a[y]+1]==0) cout<<y<<' '<<x;else cout<<x<<' '<<y;


}