#include<bits/stdc++.h>
using namespace std;
#define vec vector<long long>
#define ll long long int
#define ld long double 
#define f(n) for(int i=0;i<n;i++)
int main(){
    ios::sync_with_stdio(false);
    cin.tie(0);
    cout.tie(0);
  int t;
  cin>>t;
  while(t--){
       ll n,d,h;
       cin>>n>>d>>h;
       vec y(n);
       f(n) cin>>y[i];
       double area = double(d*h)/double(2.0);
       double totarea = n* area;
       for(int i=1;i<n;i++){
          double diff=(y[i]-y[i-1]);
          if(diff<h){
            double h2= (h-diff);
            double d2 = double(h2 * d)/double(h);

            double ext = double(d2*h2)/double(2.0);
            totarea-=ext;
          }
       }

       cout<<setprecision(7)<<totarea<<endl;
  }
}