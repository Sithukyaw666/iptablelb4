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
	if err := utils.EnableIPForwarding(); err != nil {
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

	request := new(model.Request)

	if err := c.ShouldBindJSON(request); err != nil {
		utils.StandardResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	chainName := request.ServerFarm

	if ok, _ := ipt.ipt.ChainExists("nat", chainName); ok {
		logrus.Error("Chain name %s already exists", chainName)
		response := fmt.Sprintf("Chain name: %s already exists", chainName)
		utils.StandardResponse(c, http.StatusBadRequest, "error", response, nil)
	}

	if err := ipt.ipt.NewChain("nat", chainName); err != nil {
		logrus.Error("Can't add new chain", err)
		response := fmt.Sprintf("Cannot add new chain: %s", chainName)
		utils.StandardResponse(c, http.StatusBadRequest, "error", response, nil)

	}
	if err := ipt.ipt.AppendUnique("nat", "PREROUTING", "-j", chainName); err != nil {
		logrus.Error("Can't append new chain", err)
		return
	}
	upStreamLength := len(request.Upstreams)
	for i, server := range request.Upstreams {
		ingress, egress := utils.GenerateIptablerules(i, upStreamLength, server.IpAddress, server.Port, request.Algorithm)
		if err := ipt.ipt.AppendUnique("nat", chainName, ingress...); err != nil {
			logrus.Error("Can't append ingress rule to iptables", err)
			return
		}
		if err := ipt.ipt.AppendUnique("nat", "POSTROUTING", egress...); err != nil {
			logrus.Error("Can't append egress rule to iptables", err)
			return
		}
	}
	response := fmt.Sprintf("Loadbalancing rule for %s configured successfully", chainName)
	utils.StandardResponse(c, http.StatusOK, "success", response, nil)
}

func (ipt *iptableHandler) ListFarm(c *gin.Context) {
	chains, err := ipt.ipt.ListChains("nat")
	var serverFarms []string
	if err != nil {
		logrus.Error(err)
		response := fmt.Sprintf("%v", err)
		utils.StandardResponse(c, http.StatusRequestTimeout, "error", response, nil)
	}
	for _, chain := range chains {
		if !utils.IsPredefinedChain(chain) {
			serverFarms = append(serverFarms, chain)
		}
	}
	utils.StandardResponse(c, http.StatusOK, "success", "Listed all the backend server farms", serverFarms)
}

func (ipt *iptableHandler) ListFarmByName(c *gin.Context) {
	chainName := c.Param("farm")
	rules, err := ipt.ipt.List("nat", chainName)
	if err != nil {
		logrus.Error("Can't read the chain in the iptables")
		return
	}
	response := new(model.Request)
	var algorithm string
	for _, rule := range rules[1:] {
		mode, ip, port, err := utils.ExtractModeAndDestination(rule)
		if err != nil {
			logrus.Error("Can't extract the required fields", err)
		}
		upstream := new(model.Upstreams)

		upstream.IpAddress = ip
		upstream.Port = port

		response.Upstreams = append(response.Upstreams, *upstream)

		if mode == "nth" {
			algorithm = "round-robin"
		} else {
			algorithm = mode
		}

	}

	response.ServerFarm = chainName
	response.Algorithm = algorithm

	utils.StandardResponse(c, http.StatusOK, "success", "Listed all the backend servers", response)
}

func (ipt *iptableHandler) UpdateRule(c *gin.Context) {

	request := new(model.Request)
	if err := c.ShouldBindJSON(request); err != nil {
		utils.StandardResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	chainName := request.ServerFarm

	if err := ipt.ipt.ClearChain("nat", chainName); err != nil {
		logrus.Error(err)
	}
	upStreamLength := len(request.Upstreams)
	for i, server := range request.Upstreams {
		ingress, egress := utils.GenerateIptablerules(i, upStreamLength, server.IpAddress, server.Port, request.Algorithm)
		if err := ipt.ipt.AppendUnique("nat", chainName, ingress...); err != nil {
			logrus.Error("Can't append ingress rule to iptables", err)
			return
		}
		if err := ipt.ipt.AppendUnique("nat", "POSTROUTING", egress...); err != nil {
			logrus.Error("Can't append egress rule to iptables", err)
			return
		}
	}
	response := fmt.Sprintf("Loadbalancing rule for %s updated successfully", chainName)
	utils.StandardResponse(c, http.StatusOK, "success", response, nil)
}
func (ipt *iptableHandler) DeleteRule(c *gin.Context) {

	chainName := c.Param("farm")

	if err := ipt.ipt.DeleteIfExists("nat", "PREROUTING", "-j", chainName); err != nil {
		logrus.Error(err)
	}
	if err := ipt.ipt.ClearChain("nat", chainName); err != nil {
		logrus.Error("can't clear chain", err)
		return
	}
	if err := ipt.ipt.DeleteChain("nat", chainName); err != nil {
		logrus.Error(err)
		return
	}
	response := fmt.Sprintf("Loadbalancing rule for %s deleted successfully", chainName)

	utils.StandardResponse(c, http.StatusOK, "success", response, nil)
}
