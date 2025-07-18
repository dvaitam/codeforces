#include <bits/stdc++.h>

using namespace std;

int num[200005];
int cnt[200005];
bool flag[200005];

int main(int argc, char** argv) {
int n;
int m=0;
cin>>n;
int max1=0;
int xu=0;
long long ans=0;
for (int a=1;a<=n;a++)
{
   scanf("%d",&num[a]);
   ans+=num[a];
   if (num[a]>max1)
   {
      xu=a;
      max1=num[a];
   }
}
long long l=ans;
l-=max1;
for (int a=1;a<=n;a++)
{
   if (l-num[a]==max1&&a!=xu)
   {
      flag[a]=1;
      cnt[++m]=a;
   }
}
int max2=0;
int xu2=0;
for (int a=1;a<=n;a++)
   if (num[a]>max2&&a!=xu)
   {
      xu2=a;
      max2=num[a];
   }
l=ans;
l-=max2;
l-=max1;
if (l==max2&&flag[xu]==0)
   cnt[++m]=xu;
cout<<m<<endl;
for (int a=1;a<=m;a++)
   printf("%d ",cnt[a]);

	return 0;
}