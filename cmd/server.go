package cmd

import (
	"errors"
	"fmt"
	"tosilk/silk"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	serverPort int
)

func newRouter() (router *gin.Engine) {
	gin.SetMode(gin.ReleaseMode)

	router = gin.New()

	router.Use(gin.Recovery())

	router.POST("/tosilk", func(c *gin.Context) {
		payload := &struct {
			Base64 string `json:"base64"`
			URL    string `json:"url"`
		}{}
		response := &struct {
			Err    string `json:"err"`
			Result string `json:"result"`
		}{}
		if err := c.ShouldBindJSON(payload); err == nil {
			var b64 string
			{
				if payload.Base64 != "" {
					b64, err = silk.Base64ToSilkBase64(payload.Base64)
				} else if payload.URL != "" {
					b64, err = silk.URLToSilkBase64(payload.URL)
				} else {
					b64, err = "", errors.New("base64 和 url 必选一项")
				}
			}
			if err == nil {
				response.Result = b64
				c.JSON(200, response)
			} else {
				response.Err = err.Error()
				c.JSON(200, response)
			}
		} else {
			response.Err = err.Error()
			c.JSON(200, response)
		}
	})

	router.NoRoute(func(c *gin.Context) {
		c.Writer.Write([]byte("POST /tosilk"))
	})

	router.NoMethod(func(c *gin.Context) {
		c.Writer.Write([]byte("POST /tosilk"))
	})

	return
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "开启HTTP服务",
	Long: `通过--port/-p 指定HTTP监听的端口
向 /tosilk POST json格式的数据，其中base64和url两个字段二选一，如果都有，只选用base64
任何数据都没有做验证，请保证你的请求正确!
返回的数据格式(Json)为: { "err": "", "result": "" }
如果处理成功则 err 为空， result 为转换后的base64字符串
处理失败则 err 为错误信息`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Running on 0.0.0.0:%d\n", serverPort)
		newRouter().Run(fmt.Sprintf(":%d", serverPort))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 0, "HTTP服务的端口")
	serverCmd.MarkFlagRequired("port")
}
