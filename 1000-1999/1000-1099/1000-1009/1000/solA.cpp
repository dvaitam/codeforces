#include <bits/stdc++.h>
#define ll long long int
#define lop(i,a,n) for(int i=a;i<n;i++)
using namespace std;

int main(){
 ios_base::sync_with_stdio(false);cin.tie(NULL);
 #ifndef ONLINE_JUDGE
 freopen("input.txt", "r", stdin);
 freopen("output.txt", "w", stdout);
 #endif
   ll n,k;
   cin>>n;
   vector <string> v[10],v1[10];
   string s;
   lop(i,0,n){
    cin>>s;
    v[s.size()].push_back(s);
   }
    lop(i,0,n){
    cin>>s;
    v1[s.size()].push_back(s);
   }
   int ans=0,ss=0;
   lop(i,1,5){
    sort(v[i].begin(),v[i].end());
    sort(v1[i].begin(),v1[i].end());
      ll ar[26];
      memset(ar,0,sizeof(ar));
    lop(j,0,v[i].size()){
      lop(k,0,v[i][j].size()){
        ar[v[i][j][k]-'a'+32]++;
        ar[v1[i][j][k]-'a'+32]--;
      }
    }
          lop(k,0,26) ans+=abs(ar[k]);
      ans/=2;
      ss+=ans;ans=0;
   }
   cout<<ss;
 return 0;
}