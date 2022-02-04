package controller

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func GetFOAASMessage(ctx *gin.Context) {
	userId := ctx.GetHeader("userId")

	ctx.JSON(http.StatusOK, getMessageFromService(userId))
}

func getMessageFromService(userId string) interface{} {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", fmt.Sprintf("https://www.foaas.com/awesome/%v", userId), nil)
	req.Header.Add("Accept", "application/json")

	resp, _ := client.Do(req)
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)

	var responseObject interface{}
	json.Unmarshal(bodyBytes, &responseObject)
	return responseObject
}
