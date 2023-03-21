package service

import (
	"OJPlatform/define"
	"OJPlatform/helper"
	"OJPlatform/models"
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"
)

// GetSubmissionList
// @Tags 公共方法
// @Summary 提交记录列表
// @Param page query int false "page"
// @Param size query int false "size"
// @Param problem_identity query string false "problem_identity"
// @Param user_identity query string false "user_identity"
// @Param status query string false "status"
// @Success 200 {string} json "{"msg":"success", "data":""}"
// @Router /submission/list [get]
func GetSubmissionList(context *gin.Context) {
	size, err := strconv.Atoi(context.DefaultQuery("size", define.DefaultSize))
	if err != nil {
		log.Println("GetSubmissionList - Size Parse Error:", err)
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetSubmissionList - Size Parse Error: " + err.Error(),
		})
		return
	}
	page, err := strconv.Atoi(context.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("GetSubmissionList - Page Parse Error:", err.Error())
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "GetSubmissionList - Page Parse Error: " + err.Error(),
		})
		return
	}
	page = page - 1

	problemIdentity := context.Query("problem_identity")
	userIdentity := context.Query("user_identity")
	status, _ := strconv.Atoi(context.Query("status"))

	var count int64
	tx := models.GetSubmissionList(userIdentity, problemIdentity, status)
	// 分页
	var sList []models.Submission
	err = tx.Count(&count).Offset(page * size).Limit(size).Find(&sList).Error
	if err != nil {
		log.Println("GetSubmissionList - Pagination Error:", err.Error())
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
			"list":  sList,
		},
	})
}

// SubmitAnswer
// @Tags 用户私有方法
// @Summary 代码提交
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "problem_identity"
// @Param code body string true "code"
// @Success 200 {string} json "{"code":200, "msg":"success", "data":""}"
// @Router /user/submit [post]
func SubmitAnswer(context *gin.Context) {
	problemIdentity := context.Query("problem_identity")
	code, err := ioutil.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "SubmitAnswer - Receive Code Error:" + err.Error(),
		})
		return
	}
	// 代码保存
	filePath, err := helper.SubmitCodeSave(code)
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "SubmitAnswer - File Save Error:" + err.Error(),
		})
		return
	}
	// 从context中获取用户登录信息UserClaims
	user, _ := context.Get("user")
	userClaims := user.(*helper.UserClaims)
	submission := models.Submission{
		Identity:        helper.GetUUID(),
		ProblemIdentity: problemIdentity,
		UserIdentity:    userClaims.Identity,
		Path:            filePath,
	}
	// todo: 进行代码结果判断 (结果正确性/占用内存/运行时间)
	var problem models.Problem
	// 加载测试用例
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCases").First(&problem).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "SubmitAnswer - Load Problem Error" + err.Error(),
		})
		return
	}
	// 处理所有TestCases
	var msg string
	WA := make(chan int)  // 答案错误的channel
	OOM := make(chan int) // 超内存错误的channel
	CE := make(chan int)  // 编译错误的channel
	passCount := 0        // 对通过的测试用例计数
	var lock sync.Mutex   // 加在passCount上的互斥锁
	for _, testCase := range problem.TestCases {
		// 执行测试
		go func() {
			var beginMem runtime.MemStats
			// 内存状态, Alloc: 已申请且仍在使用的字节
			runtime.ReadMemStats(&beginMem)
			fmt.Printf("KB: %v\n", beginMem.Alloc/1024)
			// go run code/user-submit/main.go
			cmd := exec.Command("go", "run", filePath)
			var out, stderr bytes.Buffer
			cmd.Stderr = &stderr
			cmd.Stdout = &out
			// 输入测试用例
			stdinPipe, err := cmd.StdinPipe()
			if err != nil {
				log.Fatalln(err)
			}
			io.WriteString(stdinPipe, testCase.Input)
			// 执行 go run xx/yy.go 这条命令
			runtime.ReadMemStats(&beginMem)
			err = cmd.Run()
			if err != nil {
				log.Println(err, stderr.String())
				if err.Error() == "exit status 2" {
					msg = stderr.String()
					CE <- 1
					return
				}
			}
			var endMem runtime.MemStats
			runtime.ReadMemStats(&endMem)
			// 答案错误
			if testCase.Output != out.String() {
				msg = "答案错误"
				WA <- 1
				return
			}
			if (endMem.Alloc/1024)-(beginMem.Alloc/1024) > uint64(problem.MaxMemory) {
				msg = "运行超内存"
				OOM <- 1
				return
			}
			lock.Lock()
			passCount += 1
			lock.Unlock()
			fmt.Println("Err: ", string(stderr.Bytes()))
			// 执行结果
			fmt.Println("运行输出: ", out.String())
			fmt.Println("执行结果: ", out.String() == testCase.Output)
		}()
	}
	select {
	// -1-待判断; 1-正确; 2-错误; 3-超时; 4-超内存
	case <-WA:
		submission.Status = 2
	case <-OOM:
		submission.Status = 4
	case <-CE:
		submission.Status = 5
	case <-time.After(time.Millisecond * time.Duration(problem.MaxRuntime)):
		if passCount == len(problem.TestCases) {
			submission.Status = 1
			msg = "答案正确"
		} else {
			submission.Status = 3
			msg = "运行超时"
		}
	}
	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		// (1) 创建提交记录
		err1 := tx.Create(&submission).Error
		if err1 != nil {
			return errors.New("Submission Save Fail: " + err1.Error())
		}
		mp := make(map[string]interface{})
		mp["submit_num"] = gorm.Expr("submit_num + ?", 1)
		if submission.Status == 1 {
			mp["pass_num"] = gorm.Expr("pass_num + ?", 1)
		}
		// (2) 更新用户信息(pass_num, submit_num)
		err2 := tx.Model(&models.User{}).Where("identity = ?", userClaims.Identity).Updates(&mp).Error
		if err2 != nil {
			return errors.New("User Modify Fail: " + err2.Error())
		}
		// (3) 更新Problem信息
		err3 := tx.Model(&models.Problem{}).Where("identity = ?", problemIdentity).Updates(&mp).Error
		if err3 != nil {
			return errors.New("Problem Modify Fail: " + err3.Error())
		}
		return nil
	}); err != nil {
		context.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "SubmitAnswer - Submit Error:" + err.Error(),
		})
	}
	context.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "success",
		"data": map[string]interface{}{
			"status": submission.Status,
			"msg":    msg,
		},
	})
}
