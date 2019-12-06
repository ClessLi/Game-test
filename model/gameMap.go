package model

import (
	"fmt"
	"github.com/ClessLi/Game-test/physic"
	"github.com/ClessLi/Game-test/resource"
	"github.com/ClessLi/Game-test/sprite"
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type GameMap struct {
	Height       float32
	Width        float32
	boxes        [][]*Box
	heightBoxNum int
	widthBoxNum  int
}

//一个简单的测试用的游戏地图生成函数
func NewGameMap(width, height float32, mapFile string) *GameMap {
	heightBoxNum := int(math.Ceil(float64(height / BoxHeight)))
	widthBoxNum := int(math.Ceil(float64(width / BoxWidth)))
	grounds := heightBoxNum / 4
	xGrounds := widthBoxNum / 4
	fmt.Println("map box size:", heightBoxNum, widthBoxNum)
	boxes := make([][]*Box, heightBoxNum)
	for i := 0; i < heightBoxNum; i++ {
		rowboxs := make([]*Box, widthBoxNum)
		if i < grounds || i > grounds*3 {
			for j := 0; j < widthBoxNum; j++ {
				if j < xGrounds || j > xGrounds*3 || i == 0 {
					rowboxs[j] = setWall(j, i)
				}
			}
			if i == 1 {
				for s := 0; s < widthBoxNum; s++ {
					if s >= xGrounds && s <= xGrounds*3 {
						rowboxs[s] = setSpike(s, i)
					}
				}
			}
		}
		boxes[i] = rowboxs
	}
	return &GameMap{
		Height:       height,
		Width:        width,
		boxes:        boxes,
		heightBoxNum: heightBoxNum,
		widthBoxNum:  widthBoxNum,
	}
}

func setSpike(x, y int) *Box {
	spike := NewGameObj(resource.GetTexture("spike"), float32(x)*BoxWidth, float32(y)*BoxHeight, &mgl32.Vec2{BoxWidth, BoxHeight}, 0, &mgl32.Vec3{1, 1, 1})
	return &Box{*spike}
}

func setWall(x, y int) *Box {
	wall := NewGameObj(resource.GetTexture("wall"), float32(x)*BoxWidth, float32(y)*BoxHeight, &mgl32.Vec2{BoxWidth, BoxHeight}, 0, &mgl32.Vec3{1, 1, 1})
	return &Box{*wall}
}

//检测一个物体是否与地图中的方块或尖刺发生碰撞
func (gameMap *GameMap) IsColl(gameObj GameObj, shift mgl32.Vec2) (bool, mgl32.Vec2) {
	position := gameObj.GetPosition()
	size := gameObj.GetSize()
	//startX, endX, startY, endY := gameMap.FetchBox(mgl32.Vec2{position[0], position[1]}, mgl32.Vec2{size[0], size[1]})
	startX, endX, startY, endY := gameMap.FetchBox(mgl32.Vec2{position[0], position[1]}, mgl32.Vec2{size[0], size[1]}, 1)
	positionMap := make(map[float64]mgl32.Vec2)
	for i := startY; i <= endY; i++ {
		for j := startX; j < endX; j++ {
			box := gameMap.boxes[i][j]
			if box != nil {
				isCol, p := physic.ColldingAABBPlace(gameObj, box, shift)
				if isCol {
					//fmt.Println(j, i)
					spacing := getSpacing(position, p)
					positionMap[spacing] = p
					//fmt.Println(box.x, box.y, box.size, box.texture.ID)
				}
			}
		}
	}
	if len(positionMap) > 0 {
		return true, getShortPositon(positionMap)
	}
	return false, gameObj.GetPosition()
}

func getShortPositon(pmap map[float64]mgl32.Vec2) mgl32.Vec2 {
	min := float64(^uint(0) >> 1)
	for s, _ := range pmap {
		//if isEqual(min, float64(0)) && !isEqual(s, float64(0)) {
		if min > s {
			min = s
		}
	}
	fmt.Println("pmap:", pmap, "min spacing:", min)
	return pmap[min]
}

func isEqual(f1 float64, f2 float64) bool {
	return math.Abs(f1-f2) == float64(0)
}

func getSpacing(s mgl32.Vec2, d mgl32.Vec2) float64 {
	w := d[0] - s[0]
	h := d[1] - s[1]
	//fmt.Println("s:", s, "d:", d, "w:", w, "h:", h, "spacing:", math.Sqrt(math.Pow(float64(w), 2) + math.Pow(float64(h), 2)))
	return math.Sqrt(math.Pow(float64(w), 2) + math.Pow(float64(h), 2))
}

// 将一个物体坐标转换为地图格子坐标范围
func (gameMap *GameMap) FetchBox(position, size mgl32.Vec2, extraboxnum int) (int, int, int, int) {
	startX := int(math.Floor(float64(position[0]/gameMap.Width*float32(gameMap.widthBoxNum-extraboxnum)))) - 1
	if startX <= 0 {
		startX = 0
	}
	endX := int(math.Ceil(float64((position[0]+size[0])/gameMap.Width*float32(gameMap.widthBoxNum+extraboxnum)))) + 1
	if endX >= gameMap.widthBoxNum {
		endX = gameMap.widthBoxNum - 1
	}
	startY := int(math.Floor(float64(position[1]/gameMap.Height*float32(gameMap.heightBoxNum-extraboxnum)))) - 1
	if startY < 0 {
		startY = 0
	}
	endY := int(math.Ceil(float64((position[1]+size[1])/gameMap.Height*float32(gameMap.heightBoxNum+extraboxnum)))) + 1
	if endY >= gameMap.heightBoxNum {
		endY = gameMap.heightBoxNum - 1
	}
	return startX, endX, startY, endY
}

//渲染地图
func (gameMap *GameMap) Draw(position mgl32.Vec2, zoom mgl32.Vec2, renderer *sprite.SpriteRenderer) {
	startX, endX, startY, endY := gameMap.FetchBox(position, zoom, 0)
	for i := startY; i <= endY; i++ {
		for j := startX; j < endX; j++ {
			box := gameMap.boxes[i][j]
			if box != nil {
				box.Draw(renderer)
			}
		}
	}
}
