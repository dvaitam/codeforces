#include<cstdio>
#include<cstring>
#include<algorithm>
using namespace std;
#define LL long long
void fg(int ax,int ay,int bx,int by,int &x,int &y)
{
	if(ay<bx||by<ax){x=-1;y=-1;return;}
	if(ax<=bx&&bx<=ay&&ay<=by){x=bx;y=ay;return;}
	if(ax<=bx&&by<=ay){x=bx;y=by;return;}
	if(bx<=ax&&ay<=by){x=ax;y=ay;return;}
	if(bx<=ax&&ax<=by&&by<=ay){x=ax;y=by;return;}
}
void gx(int x,int y,int z,int w,LL &bla,LL &whi)
{
	int h=z-x+1,l=w-y+1;
	bla=1ll*h*l/2;whi=1ll*h*l/2;
	if((h&1)&&(l&1)){
		if((x+y)&1)bla++;
		else whi++;
	}
}
LL S(int x,int y,int z,int w)
{
	return 1ll*(z-x+1)*(w-y+1);
}
int main()
{
	int n,m,ax1,ax2,ay1,ay2,bx1,bx2,by1,by2,cx1,cx2,cy1,cy2,T;
	LL bla,whi,abla,awhi,bbla,bwhi,cbla,cwhi;
	scanf("%d",&T);
	while(T--){
		scanf("%d%d",&n,&m);
		gx(1,1,n,m,bla,whi);
		scanf("%d%d%d%d",&ax1,&ay1,&ax2,&ay2);
		scanf("%d%d%d%d",&bx1,&by1,&bx2,&by2);
		fg(ax1,ax2,bx1,bx2,cx1,cx2);
		fg(ay1,ay2,by1,by2,cy1,cy2);
		gx(ax1,ay1,ax2,ay2,abla,awhi);
		gx(bx1,by1,bx2,by2,bbla,bwhi);
		bla-=abla;whi-=awhi;
		bla-=bbla;whi-=bwhi;
		if(cx1==-1||cx2==-1||cy1==-1||cy2==-1){
			bla=1ll*bla+1ll*S(bx1,by1,bx2,by2);
			whi=1ll*whi+1ll*S(ax1,ay1,ax2,ay2);
			printf("%I64d %I64d\n",whi,bla);
			continue;
		}
		gx(cx1,cy1,cx2,cy2,cbla,cwhi);
		bla+=cbla;whi+=cwhi;
		bla=1ll*bla+1ll*S(bx1,by1,bx2,by2);
		whi=1ll*whi+1ll*S(ax1,ay1,ax2,ay2)-1ll*S(cx1,cy1,cx2,cy2);
		printf("%I64d %I64d\n",whi,bla);
	}
}