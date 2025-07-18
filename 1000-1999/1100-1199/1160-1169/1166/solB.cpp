#include<bits/stdc++.h>
using namespace std;
#define all(x) (x).begin(), (x).end()
typedef vector <long long> vi;
int main(){
    long long n,j,y=1,i,x=0;
    cin>>n;
    char ar[5]={'a','e','i','o','u'};
    if(n<25){
        cout<<-1;
        return 0;
    }
    for(i=5;i*i<=n;i++){
        if(n%i==0){
            y=i;
            break;
        }
    }
    if(n/y<5 || y==1){
        cout<<-1;
        return 0;
    }
    if(y==5){
        for(i=0;i<n/5;i++){
            if(i%5==0) cout<<"aeiou";
           else if(i%5==1) cout<<"eioua";
            else    if(i%5==2) cout<<"iouae";
            else if(i%5==3) cout<<"ouaei";
            else cout<<"uaeio";
        }
    }
    else{for(j=0;j<n/y;j++){for(i=0;i<y;i++){cout<<ar[x%5];x++;}}}
	return 0;
}