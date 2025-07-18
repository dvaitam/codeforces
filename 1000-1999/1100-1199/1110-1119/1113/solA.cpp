#include<bits/stdc++.h>
#include<string>
#include<cmath>
/*BISHOY MAGDY ADEEB*/
using namespace std;
int main(){
int z,a,b,m=1,j=2,sum=0;cin>>a>>b;
for(int i=1;i<=a;i++){
if(b>0){sum++;b--;}
else{sum+=j;j++;}
a--;
i--;
if(a==1){break;}
}


cout<<sum;







return 0;

}