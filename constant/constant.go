package constant

type Direction int

const (
	UP    Direction = iota // 摄像机移动状态:上
	DOWN                   // 下
	LEFT                   // 左
	RIGHT                  // 右
)
