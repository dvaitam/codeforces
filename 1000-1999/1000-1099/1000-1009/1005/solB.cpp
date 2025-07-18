#include<bits/stdc++.h>
using namespace std;
char a[200006],b[200006];
int len,cnt,ma;
int main(){
	scanf("%s",a);
	scanf("%s",b);
		for(int i=strlen(a)-1,j=strlen(b)-1;i>=0,j>=0;i--,j--)
			if(a[i]==b[j])
				ma++;
			else
				break;
	printf("%d",strlen(a)+strlen(b)-2*ma);
	scanf("%d",cnt);
}