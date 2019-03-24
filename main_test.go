package main

import (
	"math/rand"
	"testing"
	"time"
	"github.com/faiface/pixel"
)


func TestBadGuyMovesTowardsHero(t *testing.T) {
	badGuy1 := NewBadGuy()
	badGuy2 := NewBadGuy()
	badGuy3 := NewBadGuy()
	badGuy1.xLoc = 350.0
	badGuy1.yLoc = 450.0
	badGuy1.direction="up"
	badGuy2.xLoc = 250.0
	badGuy2.yLoc = 650.0
	badGuy2.direction="down"
	badGuy3.xLoc = 350.0
	badGuy3.yLoc = 750.0
	badGuy3.direction="left"
	var tests = []struct{
		heroX float64
		heroY float64
		bad   BadGuyType
	}{
		{250.0,500.0,badGuy1},
		{500.0,300.0,badGuy2},
		{250.0,500.0,badGuy3},
	}
	for _,test:=range tests{
		badOut:= BadGuyMovesTowardsHero(test.heroX,test.heroY,test.bad)
		if badOut.direction != test.bad.direction{
			t.Error("error ",badOut)
		}
	}

}

func TestSelectLegalSpotForHero(t *testing.T){
	var s1RandomSource = rand.NewSource(time.Now().UnixNano())
	var r1Random = rand.New(s1RandomSource)
	var legalX []float64
	var legalY []float64
	for v:=0;v<=400;v++{
		legalX=append(legalX,float64(v))
	}
	for v:=200;v<=300;v++{
		legalY=append(legalY,float64(v))
	}
	var tests = []struct{
		x             float64
		y             float64
	}{

		{legalX[r1Random.Intn(400)],legalY[r1Random.Intn(100)]},
		{legalX[r1Random.Intn(400)],legalY[r1Random.Intn(100)]},
		{legalX[r1Random.Intn(400)],legalY[r1Random.Intn(100)]},
		{legalX[r1Random.Intn(400)],legalY[r1Random.Intn(100)]},
	}
	for _,test:=range tests{
		outX,outY:= SelectLegalSpotForHero()
		if outX == test.x {
			t.Error("error!  Output x:",outX,"Expected x: ",test.x)
		}
		if outY==test.y{
			t.Error("error!  Output y:",outY,"Expected y: ",test.y)
		}
	}




}

func TestSelectLegalSpotForBadGuy(t *testing.T){
	var s1RandomSource = rand.NewSource(time.Now().UnixNano())
	var r1Random = rand.New(s1RandomSource)
	var legalX []float64
	var legalY []float64
	for v:=500;v<=650;v++{
		legalX=append(legalX,float64(v))
	}
	for v:=200;v<=400;v++{
		legalY=append(legalY,float64(v))
	}
	var tests = []struct{
		x             float64
		y             float64
	}{

		{legalX[r1Random.Intn(150)],legalY[r1Random.Intn(200)]},
		{legalX[r1Random.Intn(150)],legalY[r1Random.Intn(200)]},
		{legalX[r1Random.Intn(150)],legalY[r1Random.Intn(200)]},
		{legalX[r1Random.Intn(150)],legalY[r1Random.Intn(200)]},
	}
	for _,test:=range tests{
		outX,outY:= SelectLegalSpotForBadGuy()
		if outX == test.x {
			t.Error("error!  Output x:",outX,"Expected x: ",test.x)
		}
		if outY==test.y{
			t.Error("error!  Output y:",outY,"Expected y: ",test.y)
		}
	}




}

func TestInitializeBullet(t *testing.T)  {
	type bulletType struct{
		xLoc float64
		yLoc float64
		direction string
		touch bool
	}

	expectedResult :=bulletType{0,0,"",false}
	bullet:= InitializeBullet()
	if bullet.xLoc != expectedResult.xLoc && bullet.yLoc!=expectedResult.yLoc{
		t.Error("error! Expected bullet x,y: ",expectedResult.xLoc,expectedResult.yLoc,"Actual x,y: ",bullet.xLoc,bullet.yLoc)
	}
	if bullet.direction!=expectedResult.direction{
		t.Error("bullet direction different")
	}
	if bullet.touch!=expectedResult.touch{
		t.Error("bullet touch different")
	}

}

func TestBadGuyBulletTouchHero(t *testing.T) {

	hero1:=NewHero()
	badGuy1 := NewBadGuy()
	hero2:=NewHero()
	badGuy2 := NewBadGuy()

	hero1.xLoc=345.0
	hero1.yLoc=445.0
	badGuy1.bulletXLoc = 350.0
	badGuy1.bulletYLoc = 450.0

	hero2.xLoc = 500.0
	hero2.yLoc = 500.0
	badGuy2.bulletXLoc = 350.0
	badGuy2.bulletYLoc = 450.0

	x,y := BadGuyBulletTouchHero(hero1,badGuy1)
	if !x && !y {
		t.Error("error! expected x: true, actual x: ",x,"expected y: true, actual y: ",y)
	}

	x1,y1:=BadGuyBulletTouchHero(hero2,badGuy1)
	if x1 && y1{
		t.Error("error! expected x: true, actual x: ",x1,"expected y: true, actual y: ",y1)
	}

}

func TestCheckAllBadGuysDead(t *testing.T) {
	level :=1
	var badGuys1 []BadGuyType
	var badGuys2 []BadGuyType
	for x := 1; x <= 4+(2*level); x++ {
		newBad := NewBadGuy()
		badGuys1=append(badGuys1, newBad)
	}
	badGuys1[0].dead=true
	badGuys1[2].dead=true
	badGuys1[3].dead=true
	badGuys1[4].dead=true
	badGuys1[5].dead=true
	badGuys1[1].dead=true
	check:=CheckAllBadGuysDead(badGuys1)
	if !check{
		t.Error("error: expected:true, actual:",check)
	}
	for x := 1; x <= 4+(2*level); x++ {
		newBad := NewBadGuy()
		badGuys2=append(badGuys2, newBad)
	}
	badGuys2[1].dead=true
	badGuys2[2].dead=true
	badGuys2[3].dead=false
	badGuys2[4].dead=true
	badGuys2[5].dead=true
	badGuys2[0].dead=true
	check2:=CheckAllBadGuysDead(badGuys2)
	if check2{
		t.Error("error: expected:false, actual:",check2)
	}
}

func TestChangeLevel(t *testing.T) {
	hero1 := NewHero()
	hero1.lives = 2
	hero1.score = 220
	hero1 = ChangeLevel(hero1)
	if hero1.lives!=2 && hero1.score!=220{
		t.Error("error, lives expected:2, actual lives:",hero1.lives,"score expected: 220, actual score:",hero1.score)
	}
}

func TestBadGuyTouchHero(t *testing.T) {
	hero1:=NewHero()
	hero1.xLoc = 500.0
	hero1.yLoc = 400.0
	badGuy1 := NewBadGuy()
	badGuy1.xLoc = 500.0
	badGuy1.yLoc = 400.0
	heroPic, err := LoadPicture("hero.png")
	if err != nil {
		panic(err)
	}
	herosprite:=pixel.NewSprite(heroPic,heroPic.Bounds())
	hero1,_=BadGuyTouchHero(badGuy1.xLoc,badGuy1.yLoc,hero1,herosprite)
	if !hero1.dead{
		t.Error("error, expected hero1.dead: true, actual hero1.dead: ",hero1.dead)
	}
	hero2:=NewHero()
	badGuy2 := NewBadGuy()
	hero2.xLoc=300.0
	hero2.yLoc=500.0
	badGuy2.xLoc=400.0
	badGuy2.yLoc=200.0
	hero2,_=BadGuyTouchHero(badGuy2.xLoc,badGuy2.yLoc,hero2,herosprite)
	if hero2.dead{
		t.Error("error, expected hero2.dead: false, actual hero2.dead: ",hero2.dead)
	}
}
