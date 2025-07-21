#include<iostream>
using namespace std;
int main(){
int t;
cin>>t;
while(t--){
int s,b,c;
cin>>s>>b>>c;
for(int i=0;i<5;i++){
if(s<=b&&s<=c)s++;
else if(b<=c && b<=s)b++;
else c++;}
cout<<s*b*c<<endl;}}