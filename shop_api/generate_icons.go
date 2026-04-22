package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

// 创建简单的 tabBar 图标
func createTabBarIcon(filename string, iconType string, isActive bool) {
	// 图标尺寸 81x81
	width, height := 81, 81
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 定义颜色
	var mainColor color.RGBA
	var bgColor color.RGBA

	if isActive {
		// 激活状态 - 红色
		mainColor = color.RGBA{255, 107, 107, 255}
		bgColor = color.RGBA{255, 255, 255, 0}
	} else {
		// 未激活状态 - 灰色
		mainColor = color.RGBA{153, 153, 153, 255}
		bgColor = color.RGBA{255, 255, 255, 0}
	}

	// 填充透明背景
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, bgColor)
		}
	}

	// 绘制图标
	switch iconType {
	case "home":
		drawHomeIcon(img, mainColor)
	case "category":
		drawCategoryIcon(img, mainColor)
	case "cart":
		drawCartIcon(img, mainColor)
	case "user":
		drawUserIcon(img, mainColor)
	}

	// 保存文件
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("创建文件失败 %s: %v", filename, err)
		return
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		log.Printf("保存图片失败 %s: %v", filename, err)
		return
	}

	log.Printf("成功创建图标: %s", filename)
}

// 绘制首页图标（房子形状）
func drawHomeIcon(img *image.RGBA, c color.RGBA) {
	// 屋顶
	for x := 15; x <= 65; x++ {
		for y := 20; y <= 40; y++ {
			if x+y >= 55 && x+y <= 85 {
				img.Set(x, y, c)
			}
		}
	}

	// 房体
	for x := 25; x <= 55; x++ {
		for y := 40; y <= 65; y++ {
			img.Set(x, y, c)
		}
	}

	// 门
	for x := 35; x <= 45; x++ {
		for y := 50; y <= 65; y++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
}

// 绘制分类图标（网格形状）
func drawCategoryIcon(img *image.RGBA, c color.RGBA) {
	// 四个方块
	squares := [][4]int{
		{20, 20, 35, 35},
		{45, 20, 60, 35},
		{20, 45, 35, 60},
		{45, 45, 60, 60},
	}

	for _, sq := range squares {
		for x := sq[0]; x <= sq[2]; x++ {
			for y := sq[1]; y <= sq[3]; y++ {
				img.Set(x, y, c)
			}
		}
	}
}

// 绘制购物车图标
func drawCartIcon(img *image.RGBA, c color.RGBA) {
	// 车筐
	for x := 20; x <= 60; x++ {
		for y := 25; y <= 40; y++ {
			img.Set(x, y, c)
		}
	}

	// 车筐底部
	for x := 20; x <= 60; x++ {
		for y := 40; y <= 45; y++ {
			img.Set(x, y, c)
		}
	}

	// 车筐把手
	for x := 55; x <= 65; x++ {
		for y := 20; y <= 25; y++ {
			img.Set(x, y, c)
		}
	}

	// 车轮
	drawCircle(img, 30, 55, 6, c)
	drawCircle(img, 50, 55, 6, c)
}

// 绘制用户图标
func drawUserIcon(img *image.RGBA, c color.RGBA) {
	// 头部（圆形）
	drawCircle(img, 40, 30, 12, c)

	// 身体
	for x := 25; x <= 55; x++ {
		for y := 45; y <= 65; y++ {
			// 梯形形状
			if y == 45 {
				if x >= 30 && x <= 50 {
					img.Set(x, y, c)
				}
			} else if y == 55 {
				if x >= 27 && x <= 53 {
					img.Set(x, y, c)
				}
			} else {
				if x >= 25 && x <= 55 {
					img.Set(x, y, c)
				}
			}
		}
	}
}

// 绘制圆形
func drawCircle(img *image.RGBA, cx, cy, r int, c color.RGBA) {
	for y := cy - r; y <= cy+r; y++ {
		for x := cx - r; x <= cx+r; x++ {
			dx := x - cx
			dy := y - cy
			if dx*dx+dy*dy <= r*r {
				img.Set(x, y, c)
			}
		}
	}
}

func main() {
	log.Println("开始生成 tabBar 图标...")

	// 首页图标
	createTabBarIcon("shop_miniapp/static/tabbar/home.png", "home", false)
	createTabBarIcon("shop_miniapp/static/tabbar/home-active.png", "home", true)

	// 分类图标
	createTabBarIcon("shop_miniapp/static/tabbar/category.png", "category", false)
	createTabBarIcon("shop_miniapp/static/tabbar/category-active.png", "category", true)

	// 购物车图标
	createTabBarIcon("shop_miniapp/static/tabbar/cart.png", "cart", false)
	createTabBarIcon("shop_miniapp/static/tabbar/cart-active.png", "cart", true)

	// 我的图标
	createTabBarIcon("shop_miniapp/static/tabbar/user.png", "user", false)
	createTabBarIcon("shop_miniapp/static/tabbar/user-active.png", "user", true)

	log.Println("所有 tabBar 图标生成完成！")
}
