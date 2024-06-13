package utility

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

// this function appends a token query parameter to the specified url
func ParseUrl(url, param string) string {
	// remove any trailing slashes
	url = strings.TrimSuffix(url, "/")

	// check if there's already a query parameter on it, if it is, append your own to it
	if strings.Contains(url, "?") {
		/*
			i.e.
			url = https://sinechat.com/verification?random=hi
			new url = https://sinechat.com/verification?greet=hi&token=blahblah
		*/
		url = fmt.Sprintf("%s&token=%s", url, param)
	} else {
		/*
			i.e
			 url = https://sinechat.com/verification
			 new url = https://sinechat.com/verification?token=blahblah
		*/
		url = fmt.Sprintf("%s?token=%s", url, param)
	}

	return url
}

func GetHeader(c *gin.Context, key string) string {
	header := ""
	if c.GetHeader(key) != "" {
		header = c.GetHeader(key)
	} else if c.GetHeader(strings.ToLower(key)) != "" {
		header = c.GetHeader(strings.ToLower(key))
	} else if c.GetHeader(strings.ToUpper(key)) != "" {
		header = c.GetHeader(strings.ToUpper(key))
	} else if c.GetHeader(strings.Title(key)) != "" {
		header = c.GetHeader(strings.Title(key))
	}
	return header
}
