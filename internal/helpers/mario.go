package helpers

import "strings"

// GeneratePyramid 生成右对齐金字塔（用于 mario-less）
// 示例 (height=4):
//
//	   #
//	  ##
//	 ###
//	####
func GeneratePyramid(height int) string {
	var result strings.Builder
	for row := 1; row <= height; row++ {
		for s := 0; s < height-row; s++ {
			result.WriteString(" ")
		}
		for h := 0; h < row; h++ {
			result.WriteString("#")
		}
		result.WriteString("\n")
	}
	return result.String()
}

// GenerateDoublePyramid 生成双金字塔（用于 mario-more）
// 示例 (height=4):
//
//	   #  #
//	  ##  ##
//	 ###  ###
//	####  ####
func GenerateDoublePyramid(height int) string {
	var result strings.Builder
	for row := 1; row <= height; row++ {
		for s := 0; s < height-row; s++ {
			result.WriteString(" ")
		}
		for h := 0; h < row; h++ {
			result.WriteString("#")
		}
		result.WriteString("  ")
		for h := 0; h < row; h++ {
			result.WriteString("#")
		}
		result.WriteString("\n")
	}
	return result.String()
}
