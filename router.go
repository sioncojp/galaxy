package galaxy

import (
	"net/http"
	"regexp"

	"gopkg.in/gin-gonic/gin.v1"
)

func Router(config *Config) *gin.Engine {
	r := gin.Default()

	// health_check ...For monitoring this http server
	r.GET("/health_check", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
	r.HEAD("/health_check", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	// /repository ...Git clone target repository
	r.POST("/repository", func(c *gin.Context) {
		err := config.GitClone()
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": "Success!! Create Repository.",
		})
	})

	// /container/list ...Show all running container
	r.GET("/container/list", func(c *gin.Context) {
		lines, err := ShowAllContainer()
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		re := regexp.MustCompile(`\s{2,}`)
		for _, line := range lines {
			str := re.Split(line, -1)
			c.JSON(200, gin.H{
				"conrainer_id": str[0],
				"image":        str[1],
				"command":      str[2],
				"created":      str[3],
				"status":       str[4],
				"names":        str[5],
			})
		}
	})

	// /container/:commit_number ...CRUD Container and commits table
	r.POST("/container/:commit_number", func(c *gin.Context) {
		cn := c.PostForm("commit_number")
		cn, err := config.GitCommitNumerTo40digit(cn)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		if err := config.CreateContainer(cn); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": "Success!! Create Commit Data",
			"value":   cn,
		})
	})

	r.DELETE("/container/:commit_number", func(c *gin.Context) {
		cn := c.PostForm("commit_number")
		cn, err := config.GitCommitNumerTo40digit(cn)
		if err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}

		if err := config.DeleteContainer(cn); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": "Success!! Delete Commit Data",
			"value":   cn,
		})
	})

	// /container/proxy ...CRUD a proxy container
	r.POST("/container_proxy", func(c *gin.Context) {
		if err := config.CreateContainerProxy(); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": "Success!! Create a proxy container",
		})
	})

	r.DELETE("/container_proxy", func(c *gin.Context) {
		if err := DeleteContainerProxy(); err != nil {
			c.JSON(400, gin.H{"message": err.Error()})
			return
		}
		c.JSON(200, gin.H{
			"message": "Success!! Delete a proxy container",
		})
	})

	return r
}
