// Depressed boy Bhaskor Roy

#include<bits/stdc++.h>



#define ll long long

#define ull unsigned long long

#define       forn(i,n)              for(int i=0;i<n;i++)

#define          all(v)              v.begin(), v.end()

#define         rall(v)              v.rbegin(),v.rend()



#define            pb                push_back

#define          sz(a)               (int)a.size()

using namespace std;

#define optimize() ios_base::sync_with_stdio(0);cin.tie(0);cout.tie(0);

#define endl '\n'



void s1(){

int n;

cin>>n;

int a[n];

int ev=0,od=0;

vector<int>oda,eva;

for(int i=0;i<n;i++){

    cin>>a[i];

    if(a[i]%2){

        od++;

     oda.push_back(i+1);

    }

    else{

        ev++;

        eva.push_back(i+1);



    }

}

if(n==1 && od==1){

    cout<<-1<<endl;

}



else if(ev>=1){

    cout<<1<<endl;

    cout<<eva[0]<<endl;

}

else{

    cout<<2<<endl;

    cout<<oda[0]<<" "<<oda[1]<<endl;

}



}



int main(){

    optimize();

int t;

cin>>t;

while(t--){

    s1();



}





return 0;

}