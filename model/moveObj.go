package model

import (
	"fmt"
	"github.com/ClessLi/Game-test/constant"
	"github.com/ClessLi/Game-test/resource"
	"github.com/go-gl/mathgl/mgl32"
)

//可移动的游戏对象
type MoveObj struct {
	GameObj
	//在上下左右方向是否可移动
	stockUp, stockDown, stockLeft, stockRight bool
	//水平移动速度
	movementSpeed float32
	//飞行速度
	fallSpeed float32
	//下坠速度
	flySpeed float32
	//移动时的动画纹理
	moveTextures []*resource.Texture2D
	//静止时的纹理
	standTextures []*resource.Texture2D
	//当前静止帧
	standIndex int
	//静止帧之间的切换阈值
	standDelta float32
	//游戏地图
	gameMap *GameMap
	//当前运动帧
	moveIndex int
	//运动帧之间的切换阈值
	moveDelta float32
}

func NewMoveObject(gameObj GameObj, movementSpeed, flySpeed float32, moveTextures []*resource.Texture2D, standTextures []*resource.Texture2D, gameMap *GameMap) *MoveObj {
	moveObj := &MoveObj{
		GameObj:       gameObj,
		movementSpeed: movementSpeed,
		fallSpeed:     100,
		gameMap:       gameMap,
		moveTextures:  moveTextures,
		flySpeed:      flySpeed,
		moveIndex:     0,
		moveDelta:     0,
		standTextures: standTextures,
		standIndex:    0,
		standDelta:    0,
	}
	return moveObj
}

//恢复静止
func (moveObj *MoveObj) Stand(delta float32) {
	if moveObj.standIndex >= len(moveObj.standTextures) {
		moveObj.standIndex = 0
	}
	moveObj.standDelta += delta
	if moveObj.standDelta > 0.1 {
		moveObj.standDelta = 0
		moveObj.texture = moveObj.standTextures[moveObj.standIndex]
		moveObj.standIndex += 1
	}
}

//由用户主动发起的运动
func (moveObj *MoveObj) Move(delta float32, ml []constant.Direction) {
	shift := mgl32.Vec2{0, 0}
	if moveObj.moveIndex >= len(moveObj.moveTextures) {
		moveObj.moveIndex = 0
	}
	moveObj.moveDelta += delta
	if moveObj.moveDelta > 0.05 {
		moveObj.moveDelta = 0
		moveObj.texture = moveObj.moveTextures[moveObj.moveIndex]
		moveObj.moveIndex += 1
	}
	for i := 0; i < len(ml); i++ {
		direction := ml[i]
		switch direction {
		case constant.DOWN:
			fmt.Println("click down")
			if !moveObj.stockDown && moveObj.y+moveObj.size[1] < moveObj.gameMap.Height {
				shift[1] += moveObj.flySpeed * delta
			}
		case constant.UP:
			fmt.Println("click up")
			if !moveObj.stockUp && moveObj.y > 0 {
				shift[1] -= moveObj.flySpeed * delta
			}
		case constant.LEFT:
			fmt.Println("click left")
			moveObj.ForWardX()
			if !moveObj.stockLeft && moveObj.x > 0 {
				shift[0] -= moveObj.movementSpeed * delta
			}
		case constant.RIGHT:
			fmt.Println("click right")
			moveObj.ReverseX()
			if !moveObj.stockRight && moveObj.x+moveObj.size[0] < moveObj.gameMap.Width {
				shift[0] += moveObj.movementSpeed * delta
			}
		}
	}
	isCol, position := moveObj.gameMap.IsColl(moveObj.GameObj, shift)
	if isCol {
		moveObj.SetPosition(position)
	} else {
		moveObj.x += shift[0]
		moveObj.y += shift[1]
	}
}
