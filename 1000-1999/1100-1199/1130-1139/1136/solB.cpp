#include<bits/stdc++.h>
#define ll long long
using namespace std;
int main(){
   ll n,p,m=0;
   cin>>n>>p;
   if(p == 1 || p==n){
    cout<<6 + (n-2)*3;
   }
   else {
   if(p <= (n-p))
        cout<<(6+ p -1 +(n-2)*3);
   else
        cout<< (6 + (n-p-1)*3 + (n-p) + (p-1)*3);
   }
   return 0;
}