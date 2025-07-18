#include<iostream>
#include<cstdio>
#include<cstring>
#include<string>
#include<algorithm>
#include<cmath>
#include<set>
using namespace std;

int read(){
	int ans=0,w =1 ;
	char ch = ' ';
	while(ch<'0'||ch>'9'){
		w = ch=='-'?-1:w;
		ch = getchar();
	}
	while(ch>='0'&&ch<='9'){
		ans = ans*10 + ch - 48;
		ch = getchar();
	}
	return ans*w;
}
struct inti{
	int val,index;
};
bool cmp(inti a,inti b){
	return a.val<b.val;
}
int main(){
	int n;
	n = read();
	inti a[n+10];
	for(register int i=1;i<=n;++i){
		a[i].val =  read();
		a[i].index = i;
	}
	if(n<=3){
		printf("1\n");
		return 0;
	}
	sort(a+1,a+1+n,cmp);
	int flag = 0;
	int dui = 1,duidui=1,pos=0;
	int cnt = 0;
	if((a[n].val-a[1].val)%(n-2)==0){
		//printf("check division\n");
		int d = (a[n].val-a[1].val)/(n-2);
		for(register int i=0;i<n-1;++i){
			if(a[i+1].val!=a[1].val+cnt*d){
				if(dui){
					dui = 0;
					pos = a[i+1].index;
					continue;
				}
				else{
					duidui = 0;
					break;
				}
			}
			++ cnt;
		}
	}
	else duidui = 0;
	if(duidui&&pos){
		printf("%d\n",pos);
		return 0;
	}
	dui = 1;
	if((a[n-1].val - a[1].val)%(n-2)==0){
		int d = (a[n-1].val-a[1].val)/(n-2);
		//printf("check diff: %d\n",d);
		for(register int i=0;i<n-1;++i){
			if(a[i+1].val!=a[1].val+i*d){
				dui = 0;
				break;
			}
		}
	}
	else dui = 0;
	if(dui){
		printf("%d\n",a[n].index);
		return 0;
	}
	dui = 1;
	if((a[n].val-a[2].val)%(n-2)==0){
		int d = (a[n].val-a[2].val)/(n-2);
		for(register int i=0;i<n-2;++i){
			if(a[i+2].val!=a[2].val+i*d){
				dui = 0;
				break;
			}
		}
	}
	else dui = 0;
	if(dui){
		printf("%d\n",a[1].index);
	}
	else printf("-1\n");
	return 0;
}