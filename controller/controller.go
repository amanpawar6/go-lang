package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"example.com/common"
	"example.com/common/mailer"
	database "example.com/connectivity"
	"example.com/models"
	"example.com/response"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/gofiber/fiber/v2"
)

func HomePage(c *fiber.Ctx) error {
	fmt.Print("Welcome to the HomePage!\n")
	fmt.Println("Endpoint Hit: homePage")
	fmt.Println(time.Now().Clock())
	obj := map[string]string{}
	obj["message"] = "Welcome to the HomePage!"
	response.OnSuccess(c, "Welcome to the HomePage!", 200, obj)
	return nil
}

func TestApi(c *fiber.Ctx) error {
	obj := map[string]string{}
	err := c.BodyParser(&obj)
	if err != nil {
		response.OnError(c, err.Error(), 500, err)
		return nil
	}
	response.OnSuccess(c, "Body returned", 200, &obj)
	return nil
}

func UsersDetails(c *fiber.Ctx) error {
	rows, err := database.DB.Query("Select * from user")
	if err != nil {
		response.OnError(c, err.Error(), 500, err)
		return nil
	}
	defer rows.Close()
	result := models.Users{}
	for rows.Next() {
		User := models.User{}
		err := rows.Scan(&User.ID, &User.First_Name, &User.Last_Name, &User.Password, &User.Email)
		// Exit if we get an error
		if err != nil {
			response.OnError(c, err.Error(), 500, err)
			return nil
		}
		result.Users = append(result.Users, User)
	}
	response.OnSuccess(c, "List of all Users", 200, result)
	return nil
}

func SingleUserDetails(c *fiber.Ctx) error {
	id := map[string]int{}
	err := c.BodyParser(&id)
	log.Println(id)
	if err != nil {
		response.OnError(c, err.Error(), 500, err)
		return nil
	}
	User := models.User{}
	rows, err := database.DB.Query("Select * from user where id = ?", id["id"])
	if err != nil {
		response.OnError(c, err.Error(), 500, err)
		return nil
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&User.ID, &User.First_Name, &User.Last_Name, &User.Password, &User.Email)
		if err != nil {
			response.OnError(c, err.Error(), 500, err)
			return nil
		}
	}
	// Return User in JSON format
	response.OnSuccess(c, "All Users returned successfully", 200, User)
	return nil
}

// Create User
func Createuser(c *fiber.Ctx) error {
	body := new(models.User) //Defines Body Model
	if err := c.BodyParser(body); err != nil {
		response.OnError(c, err.Error(), 400, err)
		return nil
	}
	rows, err := database.DB.Query("Select * from user where email=?", body.Email)
	if err != nil {
		response.OnError(c, err.Error(), 400, err)
		return nil
	}
	defer rows.Close()
	User := models.User{}
	for rows.Next() {
		err := rows.Scan(&User.ID, &User.First_Name, &User.Last_Name, &User.Password, &User.Email)
		if err != nil {
			response.OnError(c, err.Error(), 400, err)
			return nil
		}
	}
	// log.Println(User)
	if User.Email != "" {
		response.OnSuccess(c, "User Already Exists", 200, User)
		return nil
	}
	hashedPassword, hasherr := common.HashPassword(body.Password)
	if hasherr != nil {
		response.OnError(c, hasherr.Error(), 400, hasherr)
		return nil
	}
	body.Password = hashedPassword
	res, err := database.DB.Exec("INSERT INTO user (id, first_name, last_name, email, password) VALUES (?, ?, ?, ?, ?)", nil, body.First_Name, body.Last_Name, body.Email, body.Password)
	if err != nil {
		response.OnError(c, err.Error(), 400, err)
		return nil
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		response.OnError(c, err.Error(), 400, err)
		return nil
	}
	body.ID = uint(lastId)
	response.OnSuccess(c, "User successfully created", 201, body)
	return nil
}

// Update User
func UpdateUser(c *fiber.Ctx) error {
	body := new(models.User)
	if err := c.BodyParser(body); err != nil {
		response.OnError(c, err.Error(), 400, err)
		return nil
	}
	if body.ID != 0 {
		User := new(models.User)
		result, err := database.DB.Query("Select * from user where id = ?", body.ID)
		if err != nil {
			response.OnError(c, err.Error(), 500, err)
			return nil
		}
		defer result.Close()
		for result.Next() {
			err := result.Scan(&User.ID, &User.First_Name, &User.Last_Name, &User.Email, &User.Password)
			if err != nil {
				response.OnError(c, err.Error(), 500, err)
				return nil
			}
		}
		userDetails := make(map[string]interface{})
		data, err := json.Marshal(User) // Convert to a json string
		if err != nil {
			response.OnError(c, err.Error(), 500, err)
			return nil
		}
		json.Unmarshal(data, &userDetails) // Convert back into interface
		for items := range userDetails {
			if items == "first_name" {
				if body.First_Name == "" {
					data, ok := userDetails[items]
					if ok {
						body.First_Name = data.(string)
					}
				}
			} else if items == "last_name" {
				if body.Last_Name == "" {
					data, ok := userDetails[items]
					if ok {
						body.Last_Name = data.(string)
					}
				}
			} else if items == "email" {
				if body.Email == "" {
					data, ok := userDetails[items]
					if ok {
						body.Email = data.(string)
					}
				}
			} else if items == "password" {
				if body.Password != "" {
					if userDetails[items] != body.Password {
						hashedPassword, hasherr := common.HashPassword(body.Password)
						if hasherr != nil {
							response.OnError(c, hasherr.Error(), 400, hasherr)
							return nil
						}
						body.Password = hashedPassword
					}
				} else if body.Password == "" {
					data, ok := userDetails[items]
					if ok {
						body.Password = data.(string)
					}
				}
			}
		}

		_, DBerr := database.DB.Query("UPDATE user SET first_name=?, last_name=?, email=?, password=? where id=?", body.First_Name, body.Last_Name, body.Email, body.Password, body.ID)
		if DBerr != nil {
			response.OnError(c, DBerr.Error(), 400, DBerr)
			return nil
		}

	} else {
		response.OnSuccess(c, "User ID is Missing", 200, body)
		return nil
	}

	response.OnSuccess(c, "Success", 200, body)
	return nil
}

// Delete User
func DeleteUser(c *fiber.Ctx) error {
	id := map[string]int{}
	err := c.BodyParser(&id)
	log.Println(id)
	if err != nil {
		response.OnError(c, err.Error(), 500, err)
		return nil
	}
	_, DBerr := database.DB.Query("Delete from user where id=?", id["id"])
	if DBerr != nil {
		response.OnError(c, DBerr.Error(), 500, DBerr)
		return nil
	}
	response.OnSuccess(c, "User Deleted", 200, nil)
	return nil
}

// Login
func Userlogin(c *fiber.Ctx) error {
	var bodypayload interface{}
	buffer := c.Body()
	json.Unmarshal(buffer, &bodypayload)
	body := bodypayload.(map[string]interface{})
	result, err := database.DB.Query("Select * from user where email=?", body["email"])
	if err != nil {
		response.OnError(c, err.Error(), 500, err)
		return nil
	}
	defer result.Close()
	User := new(models.User)
	for result.Next() {
		err := result.Scan(&User.ID, &User.First_Name, &User.Last_Name, &User.Email, &User.Password)
		if err != nil {
			response.OnError(c, err.Error(), 500, err)
			return nil
		}
	}
	if User.ID == 0 {
		response.OnError(c, "Wrong Email", 401, nil)
	} else {
		password := body["password"].(string)
		CheckPasswordHash := common.CheckPasswordHash(password, User.Password)
		if CheckPasswordHash {
			response.OnSuccess(c, "Login Successful", 200, nil)
		} else {
			response.OnError(c, "Wrong Password", 401, nil)
		}
	}
	return nil
}

// Update Password
func UpdatePassword(c *fiber.Ctx) error {
	var bodypayload interface{}
	buffer := c.Body()
	json.Unmarshal(buffer, &bodypayload)
	body := bodypayload.(map[string]interface{})
	if body["email"] == nil || body["newpassword"] == nil || body["confirmpassword"] == nil || body["oldpassword"] == nil {
		response.OnError(c, "Params Missing", 500, nil)
	} else {
		if body["newpassword"] != body["confirmpassword"] {
			response.OnError(c, "new password and confirm password not matched", 401, nil)
		} else {
			result, err := database.DB.Query("Select * from user where email=?", body["email"])
			if err != nil {
				response.OnError(c, err.Error(), 500, err)
				return nil
			}
			defer result.Close()
			User := new(models.User)
			for result.Next() {
				err := result.Scan(&User.ID, &User.First_Name, &User.Last_Name, &User.Email, &User.Password)
				if err != nil {
					response.OnError(c, err.Error(), 500, err)
					return nil
				}
			}
			// fmt.Println(User)
			oldpassword := body["oldpassword"].(string)
			CheckPasswordHash := common.CheckPasswordHash(oldpassword, User.Password)
			if CheckPasswordHash {
				newpassword := body["newpassword"].(string)
				hashedPassword, hasherr := common.HashPassword(newpassword)
				if hasherr != nil {
					response.OnError(c, hasherr.Error(), 400, hasherr)
					return nil
				}
				_, passErr := database.DB.Query("update user set password=? where id=?", hashedPassword, User.ID)
				if passErr != nil {
					response.OnError(c, passErr.Error(), 500, passErr)
					return nil
				}
				response.OnSuccess(c, "Password Updated", 200, nil)
			} else {
				response.OnError(c, "Old Password not matched", 401, nil)
			}
		}
	}
	return nil
}

// Forget Password
func ForgetPassword(c *fiber.Ctx) error {
	var bodypayload interface{}
	buffer := c.Body()
	json.Unmarshal(buffer, &bodypayload)
	body := bodypayload.(map[string]interface{})
	if body["email"] == nil {
		response.OnError(c, "Params Missing", 500, nil)
	} else {
		result, err := database.DB.Query("Select * from user where email=?", body["email"])
		if err != nil {
			response.OnError(c, err.Error(), 500, err)
			return nil
		}
		defer result.Close()
		User := new(models.User)
		for result.Next() {
			err := result.Scan(&User.ID, &User.First_Name, &User.Last_Name, &User.Email, &User.Password)
			if err != nil {
				response.OnError(c, err.Error(), 500, err)
				return nil
			}
		}
		newpassword := common.RandomPasswordGenerator()
		hashedPassword, hasherr := common.HashPassword(newpassword)
		if hasherr != nil {
			response.OnError(c, hasherr.Error(), 400, hasherr)
			return nil
		}
		_, passErr := database.DB.Query("update user set password=? where id=?", hashedPassword, User.ID)
		if passErr != nil {
			response.OnError(c, passErr.Error(), 500, passErr)
			return nil
		}
		mailer.Mailer(User.Email, User.First_Name, newpassword)
		fmt.Println(newpassword)
	}
	response.OnSuccess(c, "Password sends your email", 200, nil)
	return nil
}

func Bulkupload(c *fiber.Ctx) error {

	// This function returns the file path of the saved file or an error if it occurs
	file, err := common.FileUpload(c)
	if err != nil {
		response.OnError(c, err.Error(), 401, err)
		return nil
	}
	filepath := fmt.Sprintf("./upload/%s", file)

	f, err := excelize.OpenFile(filepath)
	if err != nil {
		response.OnError(c, "file not found", 401, err)
		return nil
	}
	sheetname := f.GetSheetMap()
	// fmt.Println(sheetname)
	rows, err := f.GetRows(sheetname[1])
	if err != nil {
		response.OnError(c, "file not found", 401, err)
		return nil
	}

	var wg sync.WaitGroup

	data1 := rows[:len(rows)/2]
	data2 := rows[(len(rows)/2)+1:]

	wg.Add(2)
	go Save(data1, &wg, 1)
	go Save(data2, &wg, 2)

	wg.Wait()

	fileremoveErr := os.Remove(filepath)
	if fileremoveErr != nil {
		response.OnError(c, fileremoveErr.Error(), 401, fileremoveErr)
		return nil
	}
	return nil
}

func Save(rows [][]string, wg *sync.WaitGroup, v int) {

	// fmt.Println(v)

	defer wg.Done()

	sqlStr := "INSERT INTO test(n1, n2, n3) VALUES "
	vals := []interface{}{}

	for _, row := range rows {
		sqlStr += "(?, ?, ?),"
		vals = append(vals, row[0], row[1], row[2])
	}

	//trim the last ,
	sqlStr = strings.TrimSuffix(sqlStr, ",")

	//Replacing ? with $n for postgres
	// sqlStr = ReplaceSQL(sqlStr, "?")

	//prepare the statement
	stmt, _ := database.DB.Prepare(sqlStr)

	//format all vals at once
	_, err := stmt.Exec(vals...)
	if err != nil {
		return
	}
}
