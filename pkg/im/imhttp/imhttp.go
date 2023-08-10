package imhttp

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/CanftIn/gothafoss/pkg/cache"
	"github.com/CanftIn/gothafoss/pkg/log"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserRole string

const (
	Admin      UserRole = "admin"
	SuperAdmin UserRole = "superAdmin"
)

type IMHttp struct {
	engine *gin.Engine
	pool   sync.Pool
}

func New() *IMHttp {
	imHttp := &IMHttp{
		engine: gin.New(),
		pool:   sync.Pool{},
	}
	imHttp.engine.Use(gin.Recovery())
	imHttp.pool.New = func() interface{} {
		return allocateContext()
	}
	return imHttp
}

func allocateContext() *Context {
	return &Context{Context: nil, log: log.NewTLog("context")}
}

type Context struct {
	*gin.Context
	log log.Log
}

func (c *Context) reset() {
	c.Context = nil
}

func (c *Context) ResponseError(err error) {
	c.JSON(http.StatusBadRequest, gin.H{
		"msg":    err.Error(),
		"status": http.StatusBadRequest,
	})
}

func (c *Context) ResponseErrorf(msg string, err error) {
	if err != nil {
		c.log.Error(msg, zap.Error(err), zap.String("path", c.FullPath()))
	}
	c.JSON(http.StatusBadRequest, gin.H{
		"msg":    msg,
		"status": http.StatusBadRequest,
	})
}

func (c *Context) ResponseErrorWithStatus(err error, status int) {
	c.JSON(http.StatusBadRequest, gin.H{
		"msg":    err.Error(),
		"status": status,
	})
}

func (c *Context) GetPage() (pageIndex int64, pageSize int64) {
	pageIndex, _ = strconv.ParseInt(c.Query("page_index"), 10, 64)
	pageSize, _ = strconv.ParseInt(c.Query("page_size"), 10, 64)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 15
	}
	return
}

func (c *Context) ResponseOK() {
	c.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func (c *Context) ResponseWithStatus(status int, data interface{}) {
	c.JSON(status, data)
}

func (c *Context) GetLoginUID() string {
	return c.MustGet("uid").(string)
}

func (c *Context) GetAppID() string {
	return c.GetHeader("appid")
}

func (c *Context) GetLoginName() string {
	return c.MustGet("name").(string)
}

func (c *Context) GetLoginRole() string {
	return c.GetString("role")
}

func (c *Context) CheckLoginRole() error {
	role := c.GetLoginRole()
	if role == "" {
		return errors.New("login user role error!")
	}
	if role != string(Admin) && role != string(SuperAdmin) {
		return errors.New("this user has no right to perform this operation!")
	}
	return nil
}

type HandlerFunc func(c *Context)

func (imHttp *IMHttp) IMHttpHandler(handlerFunc HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		hc := imHttp.pool.Get().(*Context)
		hc.reset()
		hc.Context = c
		defer imHttp.pool.Put(hc)

		handlerFunc(hc)
	}
}

func (imHttp *IMHttp) Run(addr ...string) error {
	return imHttp.engine.Run(addr...)
}

func (imHttp *IMHttp) RunTLS(addr, certFile, keyFile string) error {
	return imHttp.engine.RunTLS(addr, certFile, keyFile)
}

func (imHttp *IMHttp) POST(relativePath string, handlers ...HandlerFunc) {
	imHttp.engine.POST(relativePath, imHttp.handlersToGinHandleFuncs(handlers)...)
}

func (imHttp *IMHttp) GET(relativePath string, handlers ...HandlerFunc) {
	imHttp.engine.GET(relativePath, imHttp.handlersToGinHandleFuncs(handlers)...)
}

func (imHttp *IMHttp) Any(relativePath string, handlers ...HandlerFunc) {
	imHttp.engine.Any(relativePath, imHttp.handlersToGinHandleFuncs(handlers)...)
}

func (imHttp *IMHttp) Static(relativePath string, root string) {
	imHttp.engine.Static(relativePath, root)
}

// LoadHTMLGlob LoadHTMLGlob
func (imHttp *IMHttp) LoadHTMLGlob(pattern string) {
	imHttp.engine.LoadHTMLGlob(pattern)
}

// UseGin UseGin
func (imHttp *IMHttp) UseGin(handlers ...gin.HandlerFunc) {
	imHttp.engine.Use(handlers...)
}

// Use Use
func (imHttp *IMHttp) Use(handlers ...HandlerFunc) {
	imHttp.engine.Use(imHttp.handlersToGinHandleFuncs(handlers)...)
}

// ServeHTTP ServeHTTP
func (imHttp *IMHttp) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	imHttp.engine.ServeHTTP(w, req)
}

// HandleContext HandleContext
func (imHttp *IMHttp) HandleContext(c *Context) {
	imHttp.engine.HandleContext(c.Context)
}

func (imHttp *IMHttp) handlersToGinHandleFuncs(handlers []HandlerFunc) []gin.HandlerFunc {
	newHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	for _, handler := range handlers {
		newHandlers = append(newHandlers, imHttp.IMHttpHandler(handler))
	}
	return newHandlers
}

// AuthMiddleware 认证中间件
func (imHttp *IMHttp) AuthMiddleware(cache cache.Cache, tokenPrefix string) HandlerFunc {

	return func(c *Context) {
		token := c.GetHeader("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "token不能为空，请先登录！",
			})
			return
		}
		uidAndName := GetLoginUID(token, tokenPrefix, cache)
		if strings.TrimSpace(uidAndName) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "请先登录！",
			})
			return
		}
		uidAndNames := strings.Split(uidAndName, "@")
		if len(uidAndNames) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "token有误！",
			})
			return
		}
		c.Set("uid", uidAndNames[0])
		c.Set("name", uidAndNames[1])
		if len(uidAndNames) > 2 {
			c.Set("role", uidAndNames[2])
		}
		c.Next()
	}
}

// GetLoginUID GetLoginUID
func GetLoginUID(token string, tokenPrefix string, cache cache.Cache) string {
	uid, err := cache.Get(tokenPrefix + token)
	if err != nil {
		return ""
	}
	return uid
}

// RouterGroup RouterGroup
type RouterGroup struct {
	*gin.RouterGroup
	L *IMHttp
}

func newRouterGroup(g *gin.RouterGroup, l *IMHttp) *RouterGroup {
	return &RouterGroup{RouterGroup: g, L: l}
}

// POST POST
func (r *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.POST(relativePath, r.L.handlersToGinHandleFuncs(handlers)...)
}

// GET GET
func (r *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.GET(relativePath, r.L.handlersToGinHandleFuncs(handlers)...)
}

// DELETE DELETE
func (r *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.DELETE(relativePath, r.L.handlersToGinHandleFuncs(handlers)...)
}

// PUT PUT
func (r *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) {
	r.RouterGroup.PUT(relativePath, r.L.handlersToGinHandleFuncs(handlers)...)
}

// CORSMiddleware 跨域
func CORSMiddleware() HandlerFunc {
	return func(c *Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, token, accept, origin, Cache-Control, X-Requested-With, appid, noncestr, sign, timestamp")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT,DELETE,PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
