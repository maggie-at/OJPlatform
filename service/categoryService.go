package service

import (
	"OJPlatform/define"
	"OJPlatform/helper"
	"OJPlatform/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// GetCategoryList
// @Tags 管理员方法
// @Summary 类别列表
// @Param authorization header string true "authorization"
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Success 200 {string} json "{"msg":"success", "data":""}"
// @Router /admin/category/list [get]
func GetCategoryList(context *gin.Context) {
	size, err := strconv.Atoi(context.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("GetCategoryList - Size Parse Error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetCategoryList - Size Parse Error: " + err.Error(),
		})
		return
	}
	page, err := strconv.Atoi(context.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetCategoryList - Page Parse Error:", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetCategoryList - Page Parse Error: " + err.Error(),
		})
		return
	}
	page = page - 1

	keyword := context.Query("keyword")
	var count int64
	// 获取models层查询结果
	tx := models.GetCategoryList(keyword)
	// 分页
	var cList []models.Category
	err = tx.Count(&count).Offset(page * size).Limit(size).Find(&cList).Error
	if err != nil {
		log.Println("GetCategoryList - Pagination Error:", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetSubmissionList - Pagination Error: " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"count": count,
			"list":  cList,
		},
	})
}

// CategoryCreate
// @Tags 管理员方法
// @Summary 创建类别
// @Param authorization header string true "authorization"
// @Param name formData string true "name"
// @Param parentId formData string false "parentId"
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /admin/category/create [post]
func CategoryCreate(context *gin.Context) {
	name := context.PostForm("name")
	parentId, _ := strconv.Atoi(context.PostForm("parentId"))
	identity := helper.GetUUID()
	category := models.Category{
		Identity: identity,
		Name:     name,
		ParentId: parentId,
	}
	err := models.DB.Create(&category).Error
	if err != nil {
		log.Println("CategoryCreate - ", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "CategoryCreate - " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"identity": identity,
		},
	})
}

// CategoryModify
// @Tags 管理员方法
// @Summary 修改类别
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param name formData string true "name"
// @Param parentId formData string false "parentId"
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /admin/category/modify [put]
func CategoryModify(context *gin.Context) {
	identity := context.PostForm("identity")
	name := context.PostForm("name")
	parentId, _ := strconv.Atoi(context.PostForm("parentId"))
	if name == "" || identity == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "输入参数错误",
		})
		return
	}
	category := models.Category{
		Identity: identity,
		Name:     name,
		ParentId: parentId,
	}
	err := models.DB.Model(&models.Category{}).Where("identity = ?", identity).Updates(&category).Error
	if err != nil {
		log.Println("CategoryModify - ", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "修改失败",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"identity": category.Identity,
		},
	})
}

// CategoryDelete
// @Tags 管理员方法
// @Summary 删除类别
// @Param authorization header string true "authorization"
// @Param identity query string true "identity"
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /admin/category/delete [delete]
func CategoryDelete(context *gin.Context) {
	// Delete应该用query进行参数获取
	identity := context.Query("identity")
	if identity == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不正确",
		})
		return
	}
	var pCnt int64
	err := models.DB.Model(&models.ProblemCategory{}).Where("category_id = (SELECT id FROM categories where identity = ? LIMIT 1)", identity).Count(&pCnt).Error
	if err != nil {
		log.Println("CategoryDelete - ", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询该分类对应的问题列表失败",
		})
	}
	fmt.Println(pCnt)
	if pCnt > 0 {
		log.Println("CategoryDelete - 该类别存在关联问题")
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "该类别存在关联问题, 无法删除",
		})
		return
	}
	err = models.DB.Model(&models.Category{}).Where("identity=?", identity).Delete(&models.Category{}).Error
	if err != nil {
		log.Println("CategoryDelete - ", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "删除类别失败",
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
