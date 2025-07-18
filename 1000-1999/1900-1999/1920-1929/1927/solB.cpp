#include <bits/stdc++.h>
#include <string>
#include <algorithm>
using namespace std;


                                                                                              
int main() {
  int t;
  cin>>t;
  while(t--){
    int n;
    cin>>n;
    int arr[n];
     string s = "";
    int k = 97;
    map<int,int> m;
    for(int i=0;i<n;i++){
      cin>>arr[i];
      if(arr[i]==0){
        m[k]+=1;
        s+=char(k);
        k++;
        continue;
      }
      for(auto &it:m){
        if(it.second==arr[i]){
          s+=char(it.first);
          it.second++;
          break;
        }
      }
    }
   cout<<s<<endl;

  }
return 0;
}