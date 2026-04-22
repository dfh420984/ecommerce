package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

// 创建纯色占位图片
func createPlaceholderImage(filename string, width, height int, r, g, b uint8) {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充颜色
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.RGBA{r, g, b, 255})
		}
	}

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		log.Printf("创建文件失败 %s: %v", filename, err)
		return
	}
	defer file.Close()

	// 根据文件扩展名选择格式
	if len(filename) >= 4 && filename[len(filename)-4:] == ".png" {
		err = png.Encode(file, img)
	} else {
		err = jpeg.Encode(file, img, &jpeg.Options{Quality: 80})
	}

	if err != nil {
		log.Printf("保存图片失败 %s: %v", filename, err)
		return
	}

	log.Printf("成功创建图片: %s (%dx%d)", filename, width, height)
}

func main() {
	log.Println("开始生成占位图片...")

	// 轮播图 (800x400)
	createPlaceholderImage("uploads/banner1.jpg", 800, 400, 102, 126, 234)   // 蓝色
	createPlaceholderImage("uploads/banner2.jpg", 800, 400, 118, 75, 162)    // 紫色
	createPlaceholderImage("uploads/banner3.jpg", 800, 400, 255, 107, 107)   // 红色

	// 商品图 (400x400)
	createPlaceholderImage("uploads/product1.jpg", 400, 400, 255, 159, 67)   // 橙色
	createPlaceholderImage("uploads/product2.jpg", 400, 400, 255, 202, 40)   // 黄色
	createPlaceholderImage("uploads/product3.jpg", 400, 400, 46, 213, 115)   // 绿色

	// Logo (200x200)
	createPlaceholderImage("uploads/logo.png", 200, 200, 102, 126, 234)      // 蓝色

	log.Println("所有占位图片生成完成！")
}
