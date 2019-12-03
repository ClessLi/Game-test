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
	Height        float32
	Width         float32
	walls         [][]*Wall
	spikes        []*Spike
	heightWallNum int
	widthWallNum  int
}

//一个简单的测试用的游戏地图生成函数
func NewGameMap(width, height float32, mapFile string) *GameMap {
	heightWallNum := int(math.Ceil(float64(height / WallHeight)))
	widthWallNum := int(math.Ceil(float64(width / WallWidth)))
	grounds := heightWallNum / 4
	xGrounds := widthWallNum / 4
	fmt.Println("map wall size:", heightWallNum, widthWallNum)
	walls := make([][]*Wall, heightWallNum)
	widthSpikeNum := int(math.Ceil(float64(width/SpikeWidth))) / 2
	spikes := make([]*Spike, widthSpikeNum)
	for i := 0; i < heightWallNum; i++ {
		rowWalls := make([]*Wall, widthWallNum)
		if i < grounds || i > grounds*3 {
			for j := 0; j < widthWallNum; j++ {
				if j < xGrounds || j > xGrounds*3 || i == 0 {
					gameObj := NewGameObj(resource.GetTexture("block"), float32(j)*WallWidth, float32(i)*WallHeight, &mgl32.Vec2{WallWidth, WallHeight}, 0, &mgl32.Vec3{1, 1, 1})
					rowWalls[j] = &Wall{GameObj: *gameObj}
				}
			}
		}
		walls[i] = rowWalls
		if i == grounds {
			for s := 0; s < widthSpikeNum; s++ {
				spikesObj := NewGameObj(resource.GetTexture("spike"), float32(s+xGrounds)*SpikeWidth, float32(i)*SpikeHeight, &mgl32.Vec2{SpikeWidth, SpikeHeight}, 0, &mgl32.Vec3{1, 1, 1})
				spikes[s] = &Spike{GameObj: *spikesObj}
			}
		}
	}
	return &GameMap{
		Height:        height,
		Width:         width,
		walls:         walls,
		spikes:        spikes,
		heightWallNum: heightWallNum,
		widthWallNum:  widthWallNum,
	}
}

//检测一个物体是否与地图中的方块或尖刺发生碰撞
func (gameMap *GameMap) IsColl(gameObj GameObj, shift mgl32.Vec2) (bool, mgl32.Vec2) {
	position := gameObj.GetPosition()
	size := gameObj.GetSize()
	startX, endX, startY, endY := gameMap.FetchBox(mgl32.Vec2{position[0], position[1]}, mgl32.Vec2{size[0], size[1]})
	for i := startX; i <= endX; i++ {
		for j := startY; j < endY; j++ {
			wall := gameMap.walls[i][j]
			if wall != nil {
				isCol, position := physic.ColldingAABBPlace(gameObj, wall, shift)
				if isCol {
					return isCol, position
				}
			}
		}
	}
	return false, gameObj.GetPosition()
}

//将一个物体坐标转换为地图格子坐标范围
func (gameMap *GameMap) FetchBox(position, size mgl32.Vec2) (int, int, int, int) {
	startY := int(math.Floor(float64(position[0]/gameMap.Width*float32(gameMap.widthWallNum)))) - 1
	if startY <= 0 {
		startY = 0
	}
	endY := int(math.Ceil(float64((position[0]+size[0])/gameMap.Width*float32(gameMap.widthWallNum)))) + 1
	if endY >= gameMap.widthWallNum {
		endY = gameMap.widthWallNum - 1
	}
	startX := int(math.Floor(float64((position[1])/gameMap.Height*float32(gameMap.heightWallNum)))) - 1
	if startX < 0 {
		startX = 0
	}
	endX := int(math.Ceil(float64((position[1]+size[1])/gameMap.Height*float32(gameMap.heightWallNum)))) + 1
	if endX >= gameMap.heightWallNum {
		endX = gameMap.heightWallNum - 1
	}
	return startX, endX, startY, endY
}

//渲染地图
func (gameMap *GameMap) Draw(position mgl32.Vec2, zoom mgl32.Vec2, renderer *sprite.SpriteRenderer) {
	startX, endX, startY, endY := gameMap.FetchBox(position, zoom)
	for i := startX; i <= endX; i++ {
		for j := startY; j < endY; j++ {
			wall := gameMap.walls[i][j]
			if wall != nil {
				wall.Draw(renderer)
			}
		}
	}
}
