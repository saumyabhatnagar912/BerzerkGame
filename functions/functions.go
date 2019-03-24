package functions

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image"
	_ "image/png"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)
//Variables and types ------------------------------------------------------------------------------
var currentLevel = 1
var s1RandomSource = rand.NewSource(time.Now().UnixNano())
var r1Random = rand.New(s1RandomSource)
var bulletLastFired = time.Now()
type HeroType struct{
	xLoc      float64
	yLoc      float64
	lives     int
	score     int
	direction string
	dead      bool
	bullet    BulletType
}
var hero HeroType

type BadGuyType struct {
	xLoc            float64
	yLoc            float64
	touch           bool
	dead            bool
	direction		string
	bulletXLoc		float64
	bulletYLoc 		float64
	bulletDirection string
	bulletTouch	 	bool
	bulletFired 	bool
	bulletFiredTime	time.Time
	bulletFlag	bool
	bulletSprite *pixel.Sprite
}
var badGuys []BadGuyType

type BulletType struct{
	xLoc float64
	yLoc float64
	direction string
	touch bool
}
//--------------------------------------------------------------------------------------------------

//Run function from faiface Pixel library acts a main for the program.
//The main functionality and flow of the Game is defined in Run.
func Run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Berzerk!",
		Bounds: pixel.R(100, 0, 1000, 668),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	heroPic, err := LoadPicture("hero.png")
	if err != nil {
		panic(err)
	}

	badGuyPic, err := LoadPicture("badGuy.png")
	if err != nil {
		panic(err)
	}
	bulletPic,err:= LoadPicture("bullet.png")
	if err != nil {
		panic(err)
	}
	//blankPic,err:= LoadPicture("blank.png")
	//if err != nil {
	//	panic(err)
	//}
	spriteHeroBullet:=pixel.NewSprite(bulletPic,bulletPic.Bounds())
	//spriteBadGuyBullet:=pixel.NewSprite(bulletPic,bulletPic.Bounds())

	InitializeBadGuys()

	hero= NewHero()
	spriteHero := pixel.NewSprite(heroPic,heroPic.Bounds())
	spriteBadGuy := pixel.NewSprite(badGuyPic,badGuyPic.Bounds())



	for !win.Closed() {
		win.Clear(colornames.Black)
		imd := DrawGame(currentLevel)
		if hero.lives<=0{
			win.Clear(colornames.Black)
			spriteHeroBullet=pixel.NewSprite(bulletPic,bulletPic.Bounds())
			spriteHeroBullet.Draw(win,pixel.IM.Moved(pixel.V(0,0)))
			//	spriteBadGuyBullet=pixel.NewSprite(bulletPic,bulletPic.Bounds())
			//	spriteBadGuyBullet.Draw(win,pixel.IM.Moved(pixel.V(0,0)))
			badGuys=[]BadGuyType{}
			DisplayGameOver(win)
		}else {
			imd.Draw(win)
		}

		hero,spriteHero= CheckHeroStats(win,hero,spriteHeroBullet,spriteHero)

		for x:=range badGuys{
			if !badGuys[x].dead {
				badGuys[x]= BadGuyMovesTowardsHero(hero.xLoc,hero.yLoc,badGuys[x])
				badGuys[x].touch= CheckIfTouchedWall(currentLevel,badGuys[x].xLoc,badGuys[x].yLoc)
				if badGuys[x].touch{
					hero.score += 10
					PlaySound("badGuyDies.mp3")
					badGuys[x].dead=true
				}
			}
			if !badGuys[x].dead {
				hero,spriteHero = BadGuyTouchHero(badGuys[x].xLoc, badGuys[x].yLoc, hero,spriteHero)
			}
			if !badGuys[x].dead {
				badGuys[x].dead,hero.bullet.touch,hero.score= BulletHitBadGuy(badGuys[x].xLoc,badGuys[x].yLoc,hero.bullet.xLoc,hero.bullet.yLoc,hero.bullet.touch,hero.score)
			}

		}
		for range badGuys{
			y:=r1Random.Intn(len(badGuys))
			badGuys[y],spriteHero,hero= BadGuyShootsBullet(badGuys[y],hero,spriteHero)
			badGuys[y].bulletTouch = CheckIfTouchedWall(currentLevel,badGuys[y].bulletXLoc,badGuys[y].bulletYLoc)
			if badGuys[y].bulletTouch{
				badGuys[y].bulletXLoc = -500.0
				badGuys[y].bulletYLoc = -500.0
				//	badGuys[y].bulletSprite=pixel.NewSprite(blankPic,blankPic.Bounds())
			}

			if !badGuys[y].bulletTouch && !badGuys[y].dead {
				switch badGuys[y].bulletDirection {
				case "rightDown":
					badGuys[y].bulletXLoc += 2.0
					badGuys[y].bulletYLoc -= 2.0
				case "rightUp":
					badGuys[y].bulletXLoc += 2.0
					badGuys[y].bulletYLoc += 2.0
				case "leftDown":
					badGuys[y].bulletXLoc -= 2.0
					badGuys[y].bulletYLoc -= 2.0
				case "leftUp":
					badGuys[y].bulletXLoc -= 2.0
					badGuys[y].bulletYLoc += 2.0
				case "left":
					badGuys[y].bulletXLoc -= 2.0
				case "right":
					badGuys[y].bulletXLoc += 2.0
				case "up":
					badGuys[y].bulletYLoc += 2.0
				case "down":
					badGuys[y].bulletYLoc -= 2.0

				}
				badGuys[y].bulletSprite.Draw(win, pixel.IM.Moved(pixel.V(badGuys[y].bulletXLoc, badGuys[y].bulletYLoc)))
				badGuys[y].bulletFired = false
			}
		}

		for x:=range badGuys{
			if !badGuys[x].dead {
				spriteBadGuy.Draw(win, pixel.IM.Moved(pixel.V(badGuys[x].xLoc, badGuys[x].yLoc)))
			}

		}
		DisplayScore(win,hero.score)
		DisplayLives(win,hero.lives)
		DisplayLevel(win)

		if !hero.dead {
			spriteHero.Draw(win, pixel.IM.Moved(pixel.V(hero.xLoc, hero.yLoc)))
		}

		if !hero.bullet.touch {
			spriteHeroBullet.Draw(win, pixel.IM.Moved(pixel.V(hero.bullet.xLoc, hero.bullet.yLoc)))
		}


		if CheckAllBadGuysDead(badGuys){
			if currentLevel == 1 {
				if hero.yLoc >= 568 || hero.yLoc <= 100 {
					hero.yLoc = 300.0
					hero= ChangeLevel(hero)
				}
			}else if currentLevel ==2{
				if hero.yLoc >= 568 ||hero.xLoc<=150||hero.xLoc>=950{
					hero.xLoc = 250.0
					hero.yLoc = 300.0
					hero= ChangeLevel(hero)
				}
			}else if currentLevel== 3{
				if hero.yLoc >= 568|| hero.yLoc <=100{
					hero= ChangeLevel(hero)
				}

			}else if currentLevel ==4{
				if hero.yLoc >= 568|| hero.yLoc <=100 {
					PlaySound("heroWin.mp3")
					DisplayHeroWin(win)
				}
			}
		}
		win.Update()
	}
}
//LoadPicture function from faiface is a helper function to load pictures to be used for sprites.
//Available at https://github.com/faiface/pixel/wiki/Drawing-a-Sprite
func LoadPicture(path string)(pixel.Picture, error)  {
	file,err:=os.Open(path)
	if err!=nil{
		return nil,err
	}
	defer  file.Close()
	img,_,err:=image.Decode(file)
	if err!=nil{
		return nil,err
	}
	return pixel.PictureDataFromImage(img),nil
}
//DrawGame function is used to draw the game background depending on the current level.
//The input is an integer representing the current level.
//The output of this function is an IMDraw object.
func DrawGame(currentLevel int)*imdraw.IMDraw{
	if currentLevel == 1 {
		imd := imdraw.New(nil)
		imd.Color = colornames.Darkgray
		imd.EndShape = imdraw.SharpEndShape
		imd.Push(pixel.V(150, 100), pixel.V(150, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(950, 100), pixel.V(950, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(150, 100), pixel.V(475, 100))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(150, 568), pixel.V(475, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(675, 100), pixel.V(950, 100))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(675, 568), pixel.V(950, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(250, 334), pixel.V(850, 334))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(250, 200), pixel.V(250, 450))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(850, 200), pixel.V(850, 450))
		imd.Rectangle(10.0)

		return imd
	}else if currentLevel == 2 {
		imd := imdraw.New(nil)
		imd.Color = colornames.Darkgray
		imd.EndShape = imdraw.SharpEndShape
		imd.Push(pixel.V(150, 95), pixel.V(950, 95))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(150, 95), pixel.V(150, 170))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(410, 95), pixel.V(410, 170))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(682, 95), pixel.V(682, 170))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(950, 95), pixel.V(950, 170))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(150, 350), pixel.V(150, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(150, 568), pixel.V(410, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(410, 568), pixel.V(410, 350))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(682, 568), pixel.V(682, 350))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(682, 568), pixel.V(950, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(950, 568), pixel.V(950, 350))
		imd.Rectangle(10.0)
		return imd

	}else if currentLevel == 3 {
		imd := imdraw.New(nil)
		imd.Color = colornames.Darkgray
		imd.EndShape = imdraw.SharpEndShape
		imd.Push(pixel.V(140, 100), pixel.V(410, 100))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(682, 100), pixel.V(960, 100))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(140, 100), pixel.V(140, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(960, 100), pixel.V(960, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(140, 350), pixel.V(410, 350))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(682, 350), pixel.V(960, 350))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(140, 568), pixel.V(410, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(682, 568), pixel.V(960, 568))
		imd.Rectangle(10.0)
		return imd
	}else if currentLevel == 4 {
		imd := imdraw.New(nil)
		imd.Color = colornames.Darkgray
		imd.EndShape = imdraw.SharpEndShape
		imd.Push(pixel.V(150, 100), pixel.V(410, 100))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(682, 100), pixel.V(950, 100))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(150, 100), pixel.V(150, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(950, 100), pixel.V(950, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(150, 568), pixel.V(410, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(682, 568), pixel.V(950, 568))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(250, 450), pixel.V(800, 450))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(520, 450), pixel.V(520, 250))
		imd.Rectangle(10.0)
		imd.Push(pixel.V(250, 250), pixel.V(800, 250))
		imd.Rectangle(10.0)
		return imd
	}
	return nil
}
//SelectLegalSpot function contains set of legal x and y positions.
//The inputs to this function are two empty slices created for x and y positions.
//The output return random x and y location for bad Guys.
func SelectLegalSpotForBadGuy()(float64,float64){
	var x []float64
	var y []float64
	x=append(x,
		230.0, 260.0, 270.0, 280.0, 290.0, 300.0, 310.0, 320.0, 330.0, 340.0, 350.0, 360.0, 370.0, 380.0, 430.0,
		440.0, 450.0, 460.0, 470.0, 480.0, 490.0, 700.0, 710.0, 720.0, 730.0, 740.0, 750.0, 760.0, 770.0, 780.0, 820.0, 830.0, 840.0,
		860.0, 870.0, 880.0, 890.0)
	y = append(y, 150.0, 160.0, 180.0, 190.0, 460.0, 470.0, 480.0, 490.0, 500.0,
		510.0, 520.0, 530.0)

	return x[r1Random.Intn(34)],y[r1Random.Intn(11)]
}
//The inputs to this function are two empty slices created for x and y positions.
//The output return random x and y location for hero.
func SelectLegalSpotForHero()(float64,float64){
	var x []float64
	var y []float64
	x=append(x, 500.0, 540.0, 550.0, 560.0, 570.0, 580.0, 590.0, 600.0, 610.0, 620.0, 630.0,
		640.0, 650.0, 660.0,)
	y = append(y, 150.0, 160.0, 180.0, 190.0, 460.0, 470.0, 480.0, 490.0, 500.0,
		510.0, 520.0, 530.0)

	return x[r1Random.Intn(13)],y[r1Random.Intn(11)]
}
//InitializeBullet function is used to initialize bullet.
//This function returns a BulletType.
func InitializeBullet() BulletType {
	return BulletType{0,0,"",false}
}
//NewHero function is used to initialize hero.
//It returns HeroType with random x and y locations.
func NewHero() HeroType {
	x,y:= SelectLegalSpotForHero()
	return HeroType{x,y,3,0,"",false, InitializeBullet()}
}
//NewBadGuy function is used to initialize badGuy instance.
//It returns BadGuyType with random x and y locations.
func NewBadGuy() BadGuyType {
	x,y:= SelectLegalSpotForBadGuy()
	picBullet,_:=LoadPicture("bullet.png")
	return BadGuyType{x,y,false,false,"",0, 0, "", false, false,time.Now(),false,pixel.NewSprite(picBullet,picBullet.Bounds())}
}
//InitializeBadGuys function is used to create an initialized slice of 4+(2*currentLevel) badGuys.
func InitializeBadGuys()  {
	for x := 1; x <= 4+(2*currentLevel); x++ {
		newBad := NewBadGuy()
		badGuys=append(badGuys, newBad)
	}
}
//CheckHeroStats function checks for key strokes to move the hero sprite around and fire bullet.
//This function returns a HeroType with new hero location and bullet direction, and hero Sprite.
func CheckHeroStats(win *pixelgl.Window, heroStatus HeroType,spriteBullet *pixel.Sprite, spriteHero *pixel.Sprite)(HeroType,*pixel.Sprite){
	heroPic, err := LoadPicture("hero.png")
	if err != nil {
		panic(err)
	}
	if win.Pressed(pixelgl.KeyLeft) {
		heroStatus.direction = "left"
		spriteHero = pixel.NewSprite(heroPic,heroPic.Bounds())
		heroStatus.xLoc -= 2
	}
	if win.Pressed(pixelgl.KeyRight) {
		heroStatus.direction = "right"
		spriteHero = pixel.NewSprite(heroPic,heroPic.Bounds())
		heroStatus.xLoc += 2
	}
	if win.Pressed(pixelgl.KeyUp) {
		heroStatus.direction = "up"
		spriteHero = pixel.NewSprite(heroPic,heroPic.Bounds())
		heroStatus.yLoc += 2
	}
	if win.Pressed(pixelgl.KeyDown) {
		heroStatus.direction = "down"
		spriteHero = pixel.NewSprite(heroPic,heroPic.Bounds())
		heroStatus.yLoc -= 2
	}
	if win.Pressed(pixelgl.KeySpace){
		heroStatus.bullet.xLoc= heroStatus.xLoc
		heroStatus.bullet.yLoc= heroStatus.yLoc
		heroStatus.bullet.direction= heroStatus.direction
		heroStatus.bullet.touch=false
		spriteBullet.Draw(win,pixel.IM.Moved(pixel.V(heroStatus.bullet.xLoc, heroStatus.bullet.yLoc)))
		PlaySound("heroShoot.mp3")
	}

	switch heroStatus.bullet.direction {
	case "left":
		heroStatus.bullet.xLoc -= 2.0
		if !heroStatus.bullet.touch {
			heroStatus.bullet.touch = CheckIfTouchedWall(currentLevel, heroStatus.bullet.xLoc, heroStatus.bullet.yLoc)
		}
	case "right":
		heroStatus.bullet.xLoc += 2.0
		if !heroStatus.bullet.touch {
			heroStatus.bullet.touch = CheckIfTouchedWall(currentLevel, heroStatus.bullet.xLoc, heroStatus.bullet.yLoc)
		}
	case "up":
		heroStatus.bullet.yLoc += 2.0
		if !heroStatus.bullet.touch {
			heroStatus.bullet.touch = CheckIfTouchedWall(currentLevel, heroStatus.bullet.xLoc, heroStatus.bullet.yLoc)
		}
	case "down":
		heroStatus.bullet.yLoc -= 2.0
		if !heroStatus.bullet.touch {
			heroStatus.bullet.touch = CheckIfTouchedWall(currentLevel, heroStatus.bullet.xLoc, heroStatus.bullet.yLoc)
		}
	}

	heroStatus.dead = CheckIfTouchedWall(currentLevel, heroStatus.xLoc, heroStatus.yLoc)
	if heroStatus.dead{
		heroPicDead, err := LoadPicture("hero1.png")
		if err != nil {
			panic(err)
		}
		spriteHero = pixel.NewSprite(heroPicDead,heroPicDead.Bounds())
		PlaySound("heroLoseLife.mp3")
		fmt.Println(hero.dead)
		heroStatus.lives-=1
		tempLives:= heroStatus.lives
		tempScore:= heroStatus.score
		heroStatus = NewHero()
		heroStatus.lives=tempLives
		heroStatus.score=tempScore
		heroStatus.dead = true
	}
	return heroStatus,spriteHero
}
//badGuyBulletTouchHero function checks if the bullet from bad guy hits hero.
//The input to this function is herotype and badGuy type.
//The function returns hero.dead values for badG.bulletTouch
func BadGuyBulletTouchHero(heroStat HeroType,badGuy BadGuyType)(bool,bool){
	if ((heroStat.xLoc+10.0>=badGuy.bulletXLoc)&&(heroStat.xLoc-10.0<=badGuy.bulletXLoc)) &&
		((heroStat.yLoc+10.0>=badGuy.bulletYLoc)&&(heroStat.yLoc-10.0<=badGuy.bulletYLoc)){
		return true,true
	}
	return false,false
}
//CheckIfTouchedWall function checks for any sprite hero, badGuy or bullet touching the walls.
//The input to this function is current level, a and y position of sprite.
//It returns a bool . True is returned when sprite touches any wall.
func CheckIfTouchedWall(currentLevel int,xLoc, yLoc float64)bool  {
	if currentLevel==1 {
		if xLoc-16.0 <= 150 || xLoc+16.0 >= 950 {
			return true
		}
		if (yLoc-23.0 <= 100 && xLoc-16.0 <= 475) || (yLoc-23.0 <= 100 && xLoc+16.0 >= 674) ||
			(yLoc+23.0 >= 568 && xLoc-16.0 <= 475) || (yLoc+23.0 >= 568 && xLoc-16.0 >= 674) {
			return true
		}
		if ((yLoc+23.0 >= 200 && yLoc-23.0 <= 450) && (xLoc+16.0 >= 248 && xLoc-16.0 <= 252)) || ((xLoc+16.0 >= 250 && xLoc-16.0 <= 850) && (yLoc+23.0 >= 333 && yLoc-23.0 <= 335)) ||
			((xLoc+16.0 >= 848 && xLoc-16.0 <= 852) && (yLoc+23.0 >= 200 && yLoc-23.0 <= 450)) {
			return true
		}
		return false
	}else if currentLevel==2{
		if (xLoc-16.0 <= 150 && (yLoc-23.0<=170&&yLoc+23.0>=350)) || (xLoc+16.0 >= 950 && (yLoc-23.0<=170&&yLoc+23.0>=350)){
			return true
		}
		if yLoc-23.0 <= 95 || ((yLoc+23.0 >= 568 && xLoc-16.0 <= 410) || (yLoc+23.0 >= 568 && xLoc+16.0>=682)) {
			return true
		}
		if (xLoc+16.0 >=409&&xLoc-16.0<=411)&&(yLoc+23.0>=95&&yLoc-23.0<=170) {
			return true
		}
		if(xLoc+16.0 >=681&&xLoc-16.0<=683)&&(yLoc+23.0>=95&&yLoc-23.0<=170) {
			return true
		}
		if (xLoc+16.0 >=409&&xLoc-16.0<=411)&&(yLoc+23.0>=350&&yLoc-23.0<=568) {
			return true
		}
		if (xLoc+16.0 >=681&&xLoc-16.0<=683)&&(yLoc+23.0>=350&&yLoc-23.0<=568) {
			return true
		}
		return false
	}else if currentLevel == 3 {

		if (xLoc+16.0>=140&&xLoc-16.0<=410)&&(yLoc+23.0>=99&&yLoc-23.0<=101){
			return true
		}
		if (xLoc+16.0>=139&&xLoc-16.0<=141)&&(yLoc+23.0>=100&&yLoc-23.0<=568){
			return true
		}
		if (xLoc+16.0>=140&&xLoc-16.0<=410)&&(yLoc+23.0>=567&&yLoc-23.0<=569){
			return true
		}
		if (xLoc+16.0>=140&&xLoc-16.0<=410)&&(yLoc+23.0>=349&&yLoc-23.0<=351){
			return true
		}
		if (xLoc+16.0>=682&&xLoc-16.0<=960)&&(yLoc+23.0>=567&&yLoc-23.0<=569){
			return true
		}
		if (xLoc+16.0>=682&&xLoc-16.0<=960)&&(yLoc+23.0>=349&&yLoc-23.0<=351){
			return true
		}
		if (xLoc+16.0>=959&&xLoc-16.0<=960)&&(yLoc+23.0>=100&&yLoc-23.0<=568){
			return true
		}
		if (xLoc+16.0>=682&&xLoc-16.0<=950)&&(yLoc+23.0>=99&&yLoc-23.0<=101){
			return true
		}
	}else if currentLevel==4{
		if (xLoc+16.0>=150&&xLoc-16.0<=410)&&(yLoc+23.0>=99&&yLoc-23.0<=100){
			return true
		}
		if (xLoc+16.0>=682&&xLoc-16.0<=950)&&(yLoc+23.0>=99&&yLoc-23.0<=100){
			return true
		}
		if (xLoc+16.0>=150&&xLoc-16.0<=410)&&(yLoc+23.0>=567&&yLoc-23.0<=569){
			return true
		}
		if (xLoc+16.0>=800&&xLoc-16.0<=950)&&(yLoc+23.0>=567&&yLoc-23.0<=569){
			return true
		}
		if (xLoc+16.0>=149&&xLoc-16.0<=151)&&(yLoc+23.0>=100&&yLoc-23.0<=568){
			return true
		}
		if (xLoc+16.0>=949&&xLoc-16.0<=951)&&(yLoc+23.0>=100&&yLoc-23.0<=568){
			return true
		}
		if (xLoc+16.0>=250&&xLoc-16.0<=800)&&(yLoc+23.0>=449&&yLoc-23.0<=451){
			return true
		}
		if (xLoc+16.0>=250&&xLoc-16.0<=800)&&(yLoc+23.0>=249&&yLoc-23.0<=251){
			return true
		}
		if (xLoc+16.0>=519&&xLoc-16.0<=521)&&(yLoc+23.0>=250&&yLoc-23.0<=450){
			return true
		}
	}
	return false
}
//BadGuyMovesTowardsHero function checks the badGuy position with hero position and makes the badGuy move towards hero.
//This function returns the BadGuyType location and direction of badGuy.
func BadGuyMovesTowardsHero(HeroXLoc, HeroYLoc float64,badGuy BadGuyType)(BadGuyType)  {
	if badGuy.xLoc > HeroXLoc{
		if badGuy.yLoc > HeroYLoc{
			badGuy.xLoc -= 0.15
			badGuy.yLoc -= 0.15
			badGuy.direction = "left"
		}else {
			badGuy.xLoc -= 0.15
			badGuy.yLoc+= 0.15
			badGuy.direction = "up"
		}
	}else {
		if badGuy.yLoc > HeroYLoc {
			badGuy.xLoc += 0.15
			badGuy.yLoc -= 0.15
			badGuy.direction = "down"
		}else {
			badGuy.xLoc += 0.15
			badGuy.yLoc += 0.15
			badGuy.direction = "right"
		}
	}
	return badGuy
}
//BadGuyTouchHero function compares hero and badGuy location to see if they match.
//If the x,y position for badGuy and hero match, hero dies and if number of lives i greter than 0, the hero spawns at random location.
//This function returns updated HeroType.
func BadGuyTouchHero(badGuyXLoc,badGuyYLoc float64,heroStatus HeroType,spriteHero *pixel.Sprite)(HeroType,*pixel.Sprite){
	if (((badGuyXLoc+15.0 >= heroStatus.xLoc-15.0)&&(badGuyXLoc+15.0 <= heroStatus.xLoc+15.0))||((badGuyXLoc-15.0<= heroStatus.xLoc+15.0)&&(badGuyXLoc-15.0>= heroStatus.xLoc-15.0)))&&
		(((badGuyYLoc+15.0>= heroStatus.yLoc-15.0)&&(badGuyYLoc+15.0<= heroStatus.yLoc+15.0))||(badGuyYLoc-15.0<= heroStatus.yLoc+15.0)&&(badGuyYLoc-15.0>= heroStatus.yLoc-15.0)){
		heroStatus.dead = true
		heroPicDead, err := LoadPicture("hero1.png")
		if err != nil {
			panic(err)
		}
		spriteHero = pixel.NewSprite(heroPicDead,heroPicDead.Bounds())
		heroStatus.lives -= 1
		PlaySound("heroLoseLife.mp3")
		x,y:= SelectLegalSpotForHero()
		heroStatus.xLoc = x
		heroStatus.yLoc = y
	}
	return heroStatus,spriteHero

}
//BulletHitBadGuy function checks if hero's bullet hit the badGuy.
//If the bullet touches badGuy, hero's score is increased by 10.
//This function returns the value of badGuy.dead, if bullet touched badGuy and hero score.
func BulletHitBadGuy(badGuyXLoc, badGuyYLoc,bulletXLoc,bulletYLoc float64,bulletTouch bool,heroScore int) (bool,bool,int) {
	if !bulletTouch {
		if (bulletXLoc <= badGuyXLoc+20.0 && bulletXLoc >= badGuyXLoc-20.0) &&
			(bulletYLoc <= badGuyYLoc+20.0 && bulletYLoc >= badGuyYLoc-20.0) {
			bulletTouch = true
			heroScore += 10
			PlaySound("badGuyDies.mp3")
			return true, bulletTouch,heroScore
		}
	}
	return false,bulletTouch,heroScore
}
//CheckBadGuyBulletDirection function checks the position of bad buy with respect to hero.
//This function returns the direction in which bd guy will shot bullet.
func CheckBadGuyBulletDirection(badGuy BadGuyType,heroStat HeroType)string{
	if (heroStat.xLoc>=badGuy.xLoc)&&(heroStat.yLoc<=badGuy.yLoc){
		if (heroStat.xLoc>=badGuy.xLoc)&&((heroStat.yLoc+10.0>=badGuy.yLoc-10.0)&&(heroStat.yLoc-10.0<=badGuy.yLoc+10.0)){
			badGuy.bulletDirection = "right"
		}else{
			badGuy.bulletDirection = "rightDown"
		}

	} else if (heroStat.xLoc>=badGuy.xLoc)&&(heroStat.yLoc>=badGuy.yLoc) {
		if ((heroStat.xLoc-10.0>=badGuy.xLoc+10.0)&&(heroStat.xLoc+10.0<=badGuy.xLoc-10.0))&&(heroStat.yLoc>=badGuy.yLoc){
			badGuy.bulletDirection = "up"
		}else {
			badGuy.bulletDirection = "rightUp"
		}

	}else if (heroStat.xLoc<=badGuy.xLoc)&&(heroStat.yLoc<=badGuy.yLoc) {
		if (heroStat.xLoc <= badGuy.xLoc) && ((heroStat.yLoc+10.0 >= badGuy.yLoc-10.0) && (heroStat.yLoc-10.0 <= badGuy.yLoc+10.0)) {
			badGuy.bulletDirection = "left"
		} else {
			badGuy.bulletDirection = "leftDown"
		}
	}else if (heroStat.xLoc<=badGuy.xLoc)&&(heroStat.yLoc>=badGuy.yLoc) {
		if ((heroStat.xLoc-10.0>=badGuy.xLoc+10.0)&&(heroStat.xLoc+10.0<=badGuy.xLoc-10.0))&&(heroStat.yLoc<=badGuy.yLoc){
			badGuy.bulletDirection = "down"
		}else {
			badGuy.bulletDirection = "leftUp"
		}
	}
	return badGuy.bulletDirection
}
//BadGuyShootsBullet function makes badGuys randomly fire bullet in hero's direction.
func BadGuyShootsBullet(badG BadGuyType,heroStat HeroType,spriteHero *pixel.Sprite) (BadGuyType,*pixel.Sprite, HeroType) {
	timeSeconds:=1
	switch currentLevel {
	case 1:timeSeconds=4
	case 2:timeSeconds=3
	case 3:timeSeconds=2
	case 4:timeSeconds=2
	}
	if time.Since(bulletLastFired).Seconds()>float64(timeSeconds){

		if !badG.dead && !badG.bulletFired{
			badG.bulletTouch = false
			badG.bulletXLoc = badG.xLoc
			badG.bulletYLoc = badG.yLoc
			badG.bulletFired = true
			PlaySound("badGuyShoot.mp3")
			bulletLastFired = time.Now()
			badG.bulletFlag = true
		}

	}
	if badG.bulletFlag {
		badG.bulletFlag = false
		badG.bulletDirection= CheckBadGuyBulletDirection(badG,hero)
	}
	hero.dead,badG.bulletTouch = BadGuyBulletTouchHero(hero,badG)
	if badG.bulletTouch{
		blnk,_:=LoadPicture("blank.png")
		badG.bulletXLoc =0.0
		badG.bulletYLoc =0.0
		badG.bulletSprite=pixel.NewSprite(blnk,blnk.Bounds())
	}
	if !badG.dead && hero.dead {
		heroPicDead, err := LoadPicture("hero1.png")
		if err != nil {
			panic(err)
		}
		spriteHero = pixel.NewSprite(heroPicDead, heroPicDead.Bounds())
		PlaySound("heroLoseLife.mp3")
		fmt.Println("-",hero.dead)
		heroStat.lives -= 1
		tempLives := heroStat.lives
		tempScore := heroStat.score
		heroStat = NewHero()
		heroStat.lives = tempLives
		heroStat.score = tempScore
		heroStat.dead = true
	}
	return badG,spriteHero,heroStat
}
//DisplayScore function is used to display hero score.
func DisplayScore(win  *pixelgl.Window,heroScore int)  {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	scoreText := text.New(pixel.V(150, 80), basicAtlas)
	fmt.Fprintf(scoreText, "Score: %s", strconv.Itoa(heroScore))
	scoreText.Draw(win, pixel.IM)
}
//DisplayLives is used to display hero lives.
func DisplayLives(win *pixelgl.Window,heroLives int)  {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	livesText := text.New(pixel.V(900, 80), basicAtlas)
	fmt.Fprintf(livesText, "Lives: %s", strconv.Itoa(heroLives))
	livesText.Draw(win,pixel.IM)
}
//DisplayLevel function is used to display current level.
func DisplayLevel(win *pixelgl.Window)  {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	levelText := text.New(pixel.V(500, 80), basicAtlas)
	fmt.Fprintf(levelText, "Level: %s", strconv.Itoa(currentLevel))
	levelText.Draw(win,pixel.IM)
}
//PlaySound function from faiface beep is used to play sounds.
//Available at https://github.com/faiface/beep
func PlaySound(name string)  {
	//https://github.com/faiface/beep
	f, err := os.Open(name)
	if err != nil {
		log.Fatal(err)
	}
	s, format, _ := mp3.Decode(f)
	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	playing := make(chan struct{})
	speaker.Play(beep.Seq(s, beep.Callback(func() {
		close(playing)
	})))
	<-playing
}
//DisplayGameOver function is used to display "GAME OVER" when hero loses all the lives.
func DisplayGameOver(win *pixelgl.Window)  {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	gameOverText:= text.New(pixel.V(500, 600), basicAtlas)
	fmt.Fprintf(gameOverText, "GAME OVER")
	gameOverText.Draw(win,pixel.IM)
	PlaySound("gameOver.mp3")
}
//DisplayHeroWin function is used to display "YOU WIN" when hero wins the game.
func DisplayHeroWin(win *pixelgl.Window)  {
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	gameOverText:= text.New(pixel.V(500, 600), basicAtlas)
	fmt.Fprintf(gameOverText, "!!!!!YOU WIN!!!!!")
	gameOverText.Draw(win,pixel.IM)
	PlaySound("heroWin.mp3")
}
//CheckAllBadGuysDead function is used to check if all the badGuys in the level are ded.
func CheckAllBadGuysDead(badGuysSlice []BadGuyType) bool {
	for x:=range badGuysSlice{
		if badGuysSlice[x].dead == false{
			return false
		}
	}
	return true
}
//ChangeLevel function changes level when all badGuys are dead and hero reaches the exit location.
//This function keeps the score and lives of the hero and initialize rest of hero parameters.
func ChangeLevel(heroStatus HeroType) HeroType {
	currentLevel += 1
	tempScore := heroStatus.score
	tempLife := heroStatus.lives
	hero = NewHero()
	heroStatus.score = tempScore
	heroStatus.lives = tempLife
	InitializeBadGuys()
	PlaySound("levelComplete.mp3")
	return heroStatus
}
