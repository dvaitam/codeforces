#include<cstdio>
int n,sx,sy,ans,ansx,ansy,x1,x2,y1,y2;
int max(int a,int b){
	return a>b?a:b;
}
inline int read(){
	int x=0,f=1;
	char ch = getchar();
	while(ch<'0'||ch>'9'){
		if(ch=='-') f = -1;
		ch = getchar();
	}
	while('0'<=ch&&ch<='9'){
		x = x*10 + ch - '0';
		ch = getchar();
	}
	return x*f;
}
int main(){	
	n = read();
	sx = read(); sy = read();
	x1 = x2 = 0;
	y1 = y2 = 0;
	for(int i=1;i<=n;++i){
		int x = read(),y = read();
		if(x<sx) x1++;
		if(x>sx) x2++;
		if(y<sy) y1++;
		if(y>sy) y2++; 
	}
	ansx = sx;
	ansy = sy;
	ans = max(x1,max(x2,max(y1,y2)));
	if(ans==x1) ansx--;
	else if(ans==x2) ansx++;
	else if(ans==y1) ansy--;
	else if(ans==y2) ansy++;
	printf("%d\n",ans);
	printf("%d %d",ansx,ansy);	
	return 0;
}