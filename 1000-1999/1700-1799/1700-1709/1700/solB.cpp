#include<bits/stdc++.h>

using namespace std;

 

int main(){

int t;

cin>>t;

while(t--){

int n;

cin>>n;

 

string s;

cin>>s;

vector<int>v(n);

for(int i=0;i<n;i++){

    v[i]=s[i]-'0';

}

 

 

for(int i=0;i<n;i++){

        v[i]=9-v[i];

    }

 

 

 

if(v[0]==0){

    int carry=0;

    int d=v[n-1]+2;

    carry=d/10;

    v[n-1]=d%10;

    for(int i=n-2;i>=0;i--){

         d=v[i]+3+carry;

         

         v[i]=d%10;

         carry=d/10;

        

    }

    

    

}

for(int i=0;i<n;i++){

    cout<<v[i];

}

cout<<endl;

}}