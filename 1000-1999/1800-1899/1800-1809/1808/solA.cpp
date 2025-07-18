// Online C++ compiler to run C++ program online
// ANMOL_GOEL_2214018
//HAR_HAR_MAHADEV
//श्री शिवाय नमस्तुभ्यं
    #include <bits/stdc++.h>
using namespace std;
 void solve(){
     long long  int l,r,g=0,o=0;
    cin>>l>>r;
    int cnt=l,ans=INT_MIN;
    
    for(int i=l;i<=r;i++){
        int maxi=-1,mini=10;
        string s=to_string(i);
        for(int j=0;j<s.size();j++){
            maxi= max(maxi,s[j]-'0');
            mini = min(mini,s[j]-'0');
        }
        if(maxi-mini==9){
               cout<<i<<endl;
            return ;
        }
        if(maxi-mini>ans){
            ans=maxi-mini;
            cnt=i;
        }
    }
   
        cout<<cnt<<endl;
    
 }
int main() {
    // Write C++ code here
   // std::cout << "Hello world!";
int t;
cin>>t;
while(t--){
   solve();
    
}
}