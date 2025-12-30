#include<cstdio>
const int vx[4]={-1,0,0,1},vy[4]={0,-1,1,0};
/*
 0
1 2
 3
*/
int n,m;
struct Map{
	int a[55][55],k,x1[10800],y1[10800],x2[10800],y2[10800];
	void read(){
		for(int i=1,x,y;i<=m;i++)scanf("%d%d",&x,&y),a[x][y]=i;
	}
	void move(int x, int y, int d){
		a[x+vx[d]][y+vy[d]]=a[x][y];
		a[x][y]=0;
		x1[k]=x;
		y1[k]=y;
		x2[k]=x+vx[d];
		y2[k++]=y+vy[d];
	}
	void move_r(int x, int y){
		if(y+1<=n&&a[x][y+1])move_r(x,y+1);
		move(x,y,2);
	}
	void play(){
		int pos=0;
		for(int i=1;i<=n;i++){
			for(int j=1;j<=n;j++)if(a[i][j]){
				int x=j;++pos;
				while(x<pos)move_r(i,x++);
				while(x>pos)move(i,x--,1);
				for(x=i;x>1;x--)move(x,pos,0);
			}
		}
	}
	void fix(int*b){
		if(n==2&&a[1][1]!=b[1]){
			move(1,1,3);move(2,1,2);
			move(1,2,1);move(2,2,0);
		}
		else if(n>3){
			for(int i=1;i<=m;i++){
				int j=1,x=i;while(b[j]!=a[1][i])j++;
				move(1,i,3);
				while(x<j)move(2,x++,2);
				while(x>j)move(2,x--,1);
				move(2,x,3);
			}
			for(int i=1;i<=m;i++)move(3,i,0),move(2,i,0);
		}
	}
	void print(bool r){
		if(r)for(int i=k;i--;)printf("%d %d %d %d\n",x2[i],y2[i],x1[i],y1[i]);
		else for(int i=0;i<k;i++)printf("%d %d %d %d\n",x1[i],y1[i],x2[i],y2[i]);
	}
}M1,M2;
int main(){
	scanf("%d%d",&n,&m);
	M1.read();
	M2.read();
	M1.play();
	M2.play();
	M1.fix(M2.a[1]);
	printf("%d\n",M1.k+M2.k);
	M1.print(0);
	M2.print(1);
}
