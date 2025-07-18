#include<bits/stdc++.h>
using namespace std;
int main(){
    
    long long n,p=0;
    cin>>n;
    long long a[2*n];
    for(int i=0;i<2*n;i++)
     cin>>a[i];
     sort(a,a+2*n);
     for(int i=0;i<2*n-1;i++)
     {
         if(a[i]!=a[i+1])
         {p=1;break;}
     }
     if(!p)
     cout<<-1;
     else
     {
         for(int x:a)
          cout<<x<<" ";
     }
    return 0;
}