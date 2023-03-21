package service

import (
	"OJPlatform/define"
	"OJPlatform/helper"
	"OJPlatform/models"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

// GetProblemList
// @Tags 公共方法
// @Summary 问题列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param keyword query string false "keyword"
// @Param category_identity query string false "category_identity"
// @Success 200 {string} json "{"msg":"success", "data":""}"
// @Router /problem/list [get]
func GetProblemList(context *gin.Context) {
	// 接收QueryString参数, 用于分页参数
	size, err := strconv.Atoi(context.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("GetProblemList - Size Parse Error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetProblemList - Page Parse Error: " + err.Error(),
		})
		return
	}
	page, err := strconv.Atoi(context.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetProblemList - Page Parse Error:", err)
		return
	}
	page = page - 1

	keyword := context.Query("keyword")
	categoryIdentity := context.Query("category_identity")

	var count int64
	// 获取models层模糊查询的DB
	tx := models.GetProblemList(keyword, categoryIdentity)
	// 分页
	var pList []models.Problem
	err = tx.Count(&count).Offset(page * size).Limit(size).Find(&pList).Error
	// 如果content长度过大, 可以省略对它的模糊查询
	// err = tx.Omit("content").Count(&count).Offset(page * size).Limit(size).Find(&pList).Error
	if err != nil {
		log.Println("GetProblemList - Pagination Error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetProblemList - Pagination Error: " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"count": count,
			"list":  pList,
		},
	})
}

// GetProblemDetail
// @Tags 公共方法
// @Summary 问题详情
// @Param identity query string false "problem_identity"
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /problem/detail [get]
func GetProblemDetail(context *gin.Context) {
	identity := context.Query("identity")
	if identity == "" {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "问题唯一标识不能为空",
		})
		return
	}
	p, err := models.GetProblemDetail(identity)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "未找到相关问题",
			})
			return
		}
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetProblemDetail - " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": p,
	})
}

// ProblemCreate
// @Tags 管理员方法
// @Summary 创建问题
// @Param authorization header string true "authorization"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_memory formData int false "max_memory"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /admin/problem/create [post]
func ProblemCreate(context *gin.Context) {
	title := context.PostForm("title")
	content := context.PostForm("content")
	maxRuntime, _ := strconv.Atoi(context.PostForm("max_runtime"))
	maxMemory, _ := strconv.Atoi(context.PostForm("max_memory"))
	categoryIds := context.PostFormArray("category_ids")
	testCases := context.PostFormArray("test_cases")
	if title == "" || content == "" || len(categoryIds) == 0 || len(testCases) == 0 {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	problem := models.Problem{
		Identity:   helper.GetUUID(),
		Title:      title,
		Content:    content,
		MaxRuntime: maxRuntime,
		MaxMemory:  maxMemory,
	}
	// 处理CategoryIds数组
	pcList := make([]models.ProblemCategory, 0)
	for _, id := range categoryIds {
		categoryId, _ := strconv.Atoi(id)
		pcList = append(pcList, models.ProblemCategory{
			ProblemId:  problem.Model.ID,
			CategoryId: uint(categoryId),
		})
	}
	problem.ProblemCategories = pcList
	// 处理TestCases数组
	var tcList []models.TestCase
	fmt.Println(testCases)
	for _, tc := range testCases {
		var tcMap map[string]string
		err := json.Unmarshal([]byte(tc), &tcMap)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例格式错误",
			})
			return
		}
		if _, ok := tcMap["input"]; !ok {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例输入为空",
			})
			return
		}
		if _, ok := tcMap["output"]; !ok {
			context.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "测试用例输出为空",
			})
			return
		}
		newTC := models.TestCase{
			Identity:        helper.GetUUID(),
			ProblemIdentity: problem.Identity,
			Input:           tcMap["input"],
			Output:          tcMap["output"],
		}
		tcList = append(tcList, newTC)
	}
	problem.TestCases = tcList
	err := models.DB.Create(&problem).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "ProblemCreate - " + err.Error(),
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"identity": problem.Identity,
		},
	})
}

// ProblemModify
// @Tags 管理员方法
// @Summary 修改问题
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity"
// @Param title formData string true "title"
// @Param content formData string true "content"
// @Param max_runtime formData int false "max_runtime"
// @Param max_memory formData int false "max_memory"
// @Param category_ids formData []string false "category_ids" collectionFormat(multi)
// @Param test_cases formData []string true "test_cases" collectionFormat(multi)
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /admin/problem/modify [put]
func ProblemModify(context *gin.Context) {
	identity := context.PostForm("identity")
	title := context.PostForm("title")
	content := context.PostForm("content")
	maxRuntime, _ := strconv.Atoi(context.PostForm("max_runtime"))
	maxMemory, _ := strconv.Atoi(context.PostForm("max_memory"))
	categoryIds := context.PostFormArray("category_ids")
	testCases := context.PostFormArray("test_cases")
	if identity == "" || title == "" || content == "" || len(categoryIds) == 0 || len(testCases) == 0 {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "参数不能为空",
		})
		return
	}
	// 将以下三个步骤放到事务中, 同时成功或同时失败
	err := models.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Problem基本字段的更新
		var problem = models.Problem{
			Identity:   identity,
			Title:      title,
			Content:    content,
			MaxRuntime: maxRuntime,
			MaxMemory:  maxMemory,
		}
		err := tx.Where("identity = ?", identity).Updates(&problem).Error
		if err != nil {
			return err
		}
		// 2. 查询问题详情
		err = tx.Where("identity = ?", identity).Find(&problem).Error
		if err != nil {
			return err
		}
		// 3. Categories的更新
		// (1) [全部]删除 => 删除ProblemCategories表中的记录
		err = tx.Where("problem_id = ?", problem.ID).Delete(&models.ProblemCategory{}).Error
		if err != nil {
			return err
		}
		// (2) [全部]新增 => 向ProblemCategories表添加新记录
		var pcList []models.ProblemCategory
		for _, id := range categoryIds {
			categoryId, _ := strconv.Atoi(id)
			pc := models.ProblemCategory{
				ProblemId:  problem.ID,
				CategoryId: uint(categoryId),
			}
			pcList = append(pcList, pc)
		}
		err = tx.Create(&pcList).Error
		if err != nil {
			return err
		}
		// 4. TestCases的更新
		// (1) [全部]删除
		err = tx.Where("problem_identity = ?", problem.Identity).Delete(&models.TestCase{}).Error
		if err != nil {
			return err
		}
		// (2) [全部]新增
		var tcList []models.TestCase
		for _, tc := range testCases {
			var tcMap map[string]string
			err := json.Unmarshal([]byte(tc), &tcMap)
			if err != nil {
				return err
			}
			if tcMap["input"] == "" || tcMap["output"] == "" {
				return errors.New("testcase format error")
			}
			tcList = append(tcList, models.TestCase{
				Identity:        helper.GetUUID(),
				ProblemIdentity: problem.Identity,
				Input:           tcMap["input"],
				Output:          tcMap["output"],
			})
		}
		err = tx.Create(&tcList).Error
		return nil
	})
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "ProblemModify Error - " + err.Error(),
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
	})
}
