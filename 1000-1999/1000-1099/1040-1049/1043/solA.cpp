#include<bits/stdc++.h>
#define lld long long int
using namespace std;

int max(int a, int b)
{
  return (a>b)?a:b;
}

int main()
{
  int n,i,maxx=INT_MIN;
  cin >> n;
  float sum=0;
  int a[n];
  for(i=0;i<n;i++)
  {
    cin >> a[i];
    sum += a[i];
    if(a[i]>maxx)
      maxx = a[i];
  }
  float y=n;
  float x = ceil((2*sum)/y);
  int z = (2*sum)/n;
  if(z==x)
    cout << max(z+1,maxx);
  else
    cout << max(x,maxx);
  return 0;
}