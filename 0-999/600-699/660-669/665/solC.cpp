#include<cstdio>

#include<cstring>

#include<algorithm>

#include<cmath>

using namespace std;

int main(){

    char a[200005];

    int i,n;

    scanf("%s",a);

    n=strlen(a);

    for(i=1;i<n;i++){

        if(a[i]==a[i-1]){

            for(int j='a';j<='z';j++){

                if(j!=a[i]&&j!=a[i+1]){

                    a[i]=j;break;

                }

            }

        }

    }

    puts(a);

    return 0;

}