#include <cstdio>
#include <vector>
#include <cmath>
#include <algorithm>
using namespace std;
int main(){
    int n,k,c = 0;
    scanf("%d%d",&n,&k);
    for (int i = 0; i<k; i++) printf("%c",i+97);
    for (int i = k; i<n; i++) if (i%2 == k%2) printf("%c",k-2+97); else printf("%c",k-1+97);
    return 0;
}