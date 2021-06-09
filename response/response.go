package response

import (
	"encoding/json"
	"fmt"
	"time"

	database "example.com/connectivity"
	"github.com/gofiber/fiber/v2"
)

type status struct {
	Success bool
	Message string
	Err     error
	Result  interface{}
}

func logger(c *fiber.Ctx, resData interface{}, log error) error {
	var payload interface{}
	buffer := c.Body()
	json.Unmarshal(buffer, &payload)
	// fmt.Println(payload)
	if payload != nil {
		reqData := payload.(map[string]interface{})
		// fmt.Println(len(reqData))
		reqjson, err := json.Marshal(reqData)
		if err != nil {
			return err
		}
		resjson, err := json.Marshal(resData)
		if err != nil {
			return err
		}
		res := string(resjson)
		req := string(reqjson)
		// fmt.Println(res)
		// fmt.Println(req)
		_, er := database.DB.Query("INSERT INTO log (log_id, request, response, error, created_on, updated_on) VALUES (?, ?, ?, ?, ?, ?)", nil, string(req), string(res), nil, time.Now(), nil)
		if er != nil {
			fmt.Println(er)
		}
		return nil
	}
	resjson, err := json.Marshal(resData)
	if err != nil {
		return err
	}
	res := string(resjson)
	// fmt.Println(log)
	var errMessage string
	if log == nil {
		errMessage = ""
	} else {
		errMessage = log.Error()
	}
	_, er := database.DB.Query("INSERT INTO log (log_id, request, response, error, created_on, updated_on) VALUES (?, ?, ?, ?, ?, ?)", nil, nil, string(res), errMessage, time.Now(), nil)
	if er != nil {
		fmt.Println(er)
	}
	return nil
}

func OnError(c *fiber.Ctx, message string, code int, err error) error {
	errRes := status{
		Success: false,
		Err:     err,
		Message: message,
	}
	logger(c, nil, err)
	return c.Status(code).JSON(errRes)
}

func OnSuccess(c *fiber.Ctx, message string, code int, body interface{}) error {
	sucRes := status{
		Success: true,
		Message: message,
		Result:  body,
	}
	logger(c, body, nil)
	return c.Status(code).JSON(sucRes)
}
