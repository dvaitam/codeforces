#include<iostream>
using namespace std;
int main()
{
   int n,i,k;
   cin>>n>>k;
   for(i=1;i<=2*n;i+=2)
   {
     if (k) cout<<i<<" "<<i+1<<" ", k--;
     else cout<<i+1<<" "<<i<<" ";
   }
   return 0;
}