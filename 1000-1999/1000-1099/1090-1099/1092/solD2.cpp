#include <iostream>
#include <cstring>
#include <cmath>
#include <cstdio>
#include <cstdlib>
#include <algorithm>
using namespace std;
const int maxn=2e5+5;
int get_num(){
    int num = 0;
    char c;
    bool flag = false;
    while((c = getchar()) == ' ' || c == '\r' || c == '\n');
    if(c == '-')
        flag = true;
    else num = c - '0';
    while(isdigit(c = getchar()))
        num = num * 10 + c - '0';
    return (flag ? -1 : 1) * num;
}
int fa[maxn],a[maxn];
void work(){
	printf("NO");
}
int n;
int main(){
	n = get_num();
	for(int i = 1;i <= n;++i)
		a[i] = get_num();
	for(int i = 1;i <= n;++i)
		if(a[i] > a[0]) a[0] = a[i];
	a[++n]=a[0];
	for(int i = 1;i <= n;++i){
		int j = i-1;
		while(a[i] > a[j]){ 	
			if((i - fa[j])%2==0){
                work();
                return 0;
            }
			j = fa[j];
		}
		if(a[i] == a[j]) fa[i] = fa[j];
		else fa[i] = j;	
	}
	printf("YES");
	return 0;
}