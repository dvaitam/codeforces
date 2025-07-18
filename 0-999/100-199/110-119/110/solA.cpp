#include<iostream>
#include<algorithm>
#include<cstdio>
#include<vector>
#include<string>
#include<cstring>
#include<set>
#include<cmath>
#include<iomanip>
#include<list>
#include<stack>
#include<queue>
#include <sstream>
typedef long long int ll;
const double pi=2*cos(0);
using namespace std;
int digits(unsigned long long int a)
{
  int rem,s=0;
  while(a>0)
  {
    rem=a%10;
    if(rem==4||rem==7)
     s++;
    a/=10;
   }
  return s;
}
int main()
{
  //freopen("input.txt","r",stdin);
  //freopen("output.txt","w",stdout);
  unsigned long long int a;
  cin>>a;
  int n,i;
  n=digits(a);
  int r,t=0;
  if(n==0)
   printf("NO\n");
  else
  {
  while(n>0)
  {
   r=n%10;
   if(r!=4&&r!=7)
   {
     t=1;
     break;
    }
    n/=10;
   }
   if(t==0)
    printf("YES\n");
   else 
    printf("NO\n");
   }
   return 0;
}