package handler

import (
	"fmt"
	"net/http"

	"github.com/coreos/go-iptables/iptables"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/sithukyaw666/iptablelb4/model"
	"github.com/sithukyaw666/iptablelb4/utils"
)

type iptableHandler struct {
	ipt *iptables.IPTables
}

func NewiptableHandler() (*iptableHandler, error) {
	ipt, err := iptables.New()
	if err != nil {
		return nil, err
	}
	return &iptableHandler{ipt: ipt}, nil
}

func (ipt *iptableHandler) HealthCheck(c *gin.Context) {
	if ipt.ipt != nil {
		c.IndentedJSON(http.StatusOK, gin.H{
			"ping": "pong",
		})
	}
}

func (ipt *iptableHandler) AddRule(c *gin.Context) {

	request := new(model.Rules)

	if err := c.ShouldBindJSON(request); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}

	ipLists, _ := utils.GetLocalIPs()
	fmt.Println(ipLists)
	upStreamLength := len(request.Upstreams)
	for i, server := range request.Upstreams {
		ingress, egress := utils.GenerateIptablerules(i, upStreamLength, ipLists[0], server.IpAddress, server.Port, request.Algorithm)
		logrus.Info(ingress)
		logrus.Info(egress)

	}

}

func (ipt *iptableHandler) ListRule(c *gin.Context) {
	rules, err := ipt.ipt.List("nat", "PREROUTING")

	if err != nil {
		logrus.Error(err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"Error": "Can't read the iptables list",
		})
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"Data": rules,
	})
}
func (ipt *iptableHandler) UpdateRule(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"Data": "Updated Data",
	})
}
func (ipt *iptableHandler) DeleteRule(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, gin.H{
		"Data": "Deleted Data",
	})
}
