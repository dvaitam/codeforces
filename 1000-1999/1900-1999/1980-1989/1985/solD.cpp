#include<bits/stdc++.h>
using namespace std;

int main(){
int t; cin>>t;
while(t--){
int a,b; cin>>a>>b;
int x=0,y=0,n=0;char c;
for(int i=0;i<a;i++){
    for(int j=0;j<b;j++){
        cin>>c;if(c=='#'){x+=i+1;y+=j+1;n++;}}
}
cout<<x/n<<" "<<y/n<<endl;
}
}
