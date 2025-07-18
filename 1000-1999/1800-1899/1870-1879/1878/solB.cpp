#include <bits/stdc++.h>
using namespace std;
#define loop(n) for(int i=0;i<n;i++)
#define input(a,n) for(int i=0;i<n;i++){cin>>a[i];}
#define loop2(n) for(int j=0;j<n;j++)
void solve()
{
 long long n;
 cin>>n;
 long long a[n];
 a[0]=3;
 a[1]=5;
 cout<<a[0]<<" "<<a[1]<<" ";
 for(int i=2;i<n;i++){
    a[i]=a[i-1]+1;
    cout<<a[i]<<" ";
 }
 cout<<endl;
}
int main () {
 int t=1;
cin>>t;
while(t--)
{
solve();
}
return 0 ;
}