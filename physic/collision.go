package physic

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

//检测两个矩形是否发生碰撞
func IsCollidingAABB(thisGameObj, anotherObj React) bool {
	tPosition := thisGameObj.GetPosition()
	tSize := thisGameObj.GetSize()

	aPosition := anotherObj.GetPosition()
	aSize := anotherObj.GetSize()
	return isCollidingReact(tPosition, tSize, aPosition, aSize)
}

type React interface {
	GetPosition() mgl32.Vec2
	GetSize() mgl32.Vec2
}

func isCollidingReact(position1, size1, position2, size2 mgl32.Vec2) bool {
	// x轴方向碰撞？
	collisionX := position1[0]+size1[0] >= position2[0] && position2[0]+size2[0] >= position1[0]
	// y轴方向碰撞？
	collisionY := position1[1]+size1[1] >= position2[1] && position2[1]+size2[1] >= position1[1]
	if collisionX && collisionY {
		//fmt.Println("x: ", position1[0], position2[0], ". x_size: ", size1[0], size2[0])
		//fmt.Println("y: ", position1[1], position2[1], ". y_size: ", size1[1], size2[1])
	}
	return collisionX && collisionY
}

//检测两个矩形运动后是否会发生碰撞
func WillCollidingAABB(thisGameObj, anotherObj React, dt mgl32.Vec2) bool {
	tPosition := thisGameObj.GetPosition().Add(dt)
	tSize := thisGameObj.GetSize()
	aPosition := anotherObj.GetPosition()
	aSize := anotherObj.GetSize()
	return isCollidingReact(tPosition, tSize, aPosition, aSize)
}

//检测两个矩形的碰撞，并获取碰撞位置
func ColldingAABBPlace(thisGameObj, anotherObj React, shift mgl32.Vec2) (bool, mgl32.Vec2) {
	position := thisGameObj.GetPosition()
	if shift[0] == 0 && shift[1] == 0 {
		return false, position
	}
	colldingShift := mgl32.Vec2{0, 0}
	colldingDt := shift.Normalize()
	//fmt.Println(colldingShift, shift)
	for math.Abs(float64(colldingShift[0])) <= math.Abs(float64(shift[0])) && math.Abs(float64(colldingShift[1])) <= math.Abs(float64(shift[1])) {
		//tempColldingShift := colldingShift.Sub(colldingDt)
		tempColldingShift := colldingShift.Add(colldingDt)
		//fmt.Println(colldingShift, shift)
		if WillCollidingAABB(thisGameObj, anotherObj, tempColldingShift) {
			//fmt.Println("sX:", thisGameObj.GetPosition()[0], "sY:", thisGameObj.GetPosition()[1], "dX:", anotherObj.GetPosition()[0], "dY:", anotherObj.GetPosition()[1])
			//fmt.Println("colldingDt:", colldingDt, "tempColldingShift:", tempColldingShift, "shift:", shift)
			//return true, thisGameObj.GetPosition().Sub(colldingShift)
			return true, thisGameObj.GetPosition().Add(colldingShift)
		}
		colldingShift = tempColldingShift
	}
	return false, thisGameObj.GetPosition()
}
