package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"noto/database"
	"noto/model"
	"os"
	"sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type Manga model.Manga
type Ranobe model.Ranobe
type RequestedList model.RequestedList
type AnimeShort model.AnimeShort
type Listdata model.Listdata

type ById []Listdata

func (a ById) Len() int           { return len(a) }
func (a ById) Less(i, j int) bool { return a[i].TitleID < a[j].TitleID }
func (a ById) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

// получение пользовательского списка
func RetrieveList(c *fiber.Ctx) error {
	db := database.DB
	uname := c.Params("login") //берем из URLa юзернейм
	user := User{}
	query := User{Login: uname}         //формируем запрос
	err := db.Take(&user, &query).Error //записываем в user найденного пользователя (полностью)
	if err == gorm.ErrRecordNotFound {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "UserNotFound",
		})
	}
	entity := c.Params("entity")
	if Contains(model.EntityTypes, entity) == false {
		return c.JSON(fiber.Map{"code": 400, "message": "EntityUnknown"})
	}
	file, err := os.OpenFile("storage/lists/"+entity+"/id"+strconv.Itoa(user.ID)+"_"+entity+"s.json", os.O_RDWR|os.O_CREATE, 0644) //цепляем из записанного юзера его ID и ретривим файл со списком
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "UserlistFileNotFound",
		})
	}
	defer file.Close()

	// Read the existing data from the file
	existingData, err := ioutil.ReadAll(file)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "EmptyUserlist",
		})
	}

	// If the file is empty, write the new data as an array and return
	if len(existingData) == 0 {
		return c.JSON(fiber.Map{
			"code":    404,
			"message": "EmptyUserlist",
		})
	}

	// If the file is not empty, unmarshal the existing data
	var existingDataArray []Listdata
	err = json.Unmarshal(existingData, &existingDataArray)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "JSONUnmarshalError",
		})
	}

	sort.Sort(ById(existingDataArray))

	userListQty := len(existingDataArray)
	listedQuery := make([]RequestedList, userListQty)

	var idquery []int //слайс с выбранными айдишниками
	for d := 0; d < userListQty; d++ {
		idquery = append(idquery, existingDataArray[d].TitleID)
	}
	switch entity {
	case "anime":
		currentDBQuery := make([]Anime, userListQty)

		db.Model(&model.Anime{}).Order("id asc").Find(&currentDBQuery, idquery)
		for d := 0; d < userListQty; d++ {
		    listedQuery[d].Slug = currentDBQuery[d].Slug
			listedQuery[d].Rating = currentDBQuery[d].Rating
			listedQuery[d].Type = currentDBQuery[d].Type
			listedQuery[d].Title = currentDBQuery[d].Title
			listedQuery[d].Cover = currentDBQuery[d].Cover
			listedQuery[d].Altt = currentDBQuery[d].Altt
			listedQuery[d].Genres = currentDBQuery[d].Genres
			listedQuery[d].Themes = currentDBQuery[d].Themes
			listedQuery[d].AiredOn = currentDBQuery[d].AiredOn
		}
	case "manga":
		currentDBQuery := make([]Manga, userListQty)
		db.Model(&model.Manga{}).Order("id asc").Find(&currentDBQuery, idquery)
		for d := 0; d < userListQty; d++ {
		    listedQuery[d].Slug = currentDBQuery[d].Slug
			listedQuery[d].Rating = currentDBQuery[d].Rating
			listedQuery[d].Type = currentDBQuery[d].Type
			listedQuery[d].Title = currentDBQuery[d].Title
			listedQuery[d].Cover = currentDBQuery[d].Cover
			listedQuery[d].Altt = currentDBQuery[d].Altt
			listedQuery[d].Genres = currentDBQuery[d].Genres
			listedQuery[d].Themes = currentDBQuery[d].Themes
			listedQuery[d].AiredOn = currentDBQuery[d].AiredOn
		}
	case "ranobe":
		currentDBQuery := make([]Ranobe, userListQty)
		db.Model(&model.Ranobe{}).Order("id asc").Find(&currentDBQuery, idquery)
		for d := 0; d < userListQty; d++ {
		    listedQuery[d].Slug = currentDBQuery[d].Slug
			listedQuery[d].Rating = currentDBQuery[d].Rating
			listedQuery[d].Type = currentDBQuery[d].Type
			listedQuery[d].Title = currentDBQuery[d].Title
			listedQuery[d].Cover = currentDBQuery[d].Cover
			listedQuery[d].Altt = currentDBQuery[d].Altt
			listedQuery[d].Genres = currentDBQuery[d].Genres
			listedQuery[d].Themes = currentDBQuery[d].Themes
			listedQuery[d].AiredOn = currentDBQuery[d].AiredOn
		}
	default:
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "UnknownEntity",
		})
	}

	for d := 0; d < userListQty; d++ {
		listedQuery[d].TitleID = existingDataArray[d].TitleID
		listedQuery[d].Mark = existingDataArray[d].Mark
		listedQuery[d].Progress = existingDataArray[d].Progress
		listedQuery[d].Status = existingDataArray[d].Status
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
		"data":    listedQuery,
	})

}

func AddEntityOnList(c *fiber.Ctx) error {
	jsonx := new(Listdata)
	if err := c.BodyParser(jsonx); err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "InvalidJSON",
		})
	}
	locid := EntityIDFromSlug(c.Params("slug"), c.Params("entity"))
	entity := c.Params("entity")
	if Contains(model.EntityTypes, entity) == false {
		return c.JSON(fiber.Map{"code": 400, "message": "EntityUnknown"})
	}
	if locid == 0 {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "TitleNotFound",
		})
	}
	newData := Listdata{
		TitleID:  locid,
		Status:   jsonx.Status,
		Progress: jsonx.Progress,
		Mark:     jsonx.Mark,
	}
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})

	claims := token.Claims.(*jwt.StandardClaims)

	file, err := os.OpenFile("storage/lists/"+entity+"/id"+claims.Id+"_"+entity+"s.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorFileOpen",
		})
	}
	defer file.Close()

	// Read the existing data from the file
	existingData, err := ioutil.ReadAll(file)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorFileRead",
		})
	}

	// If the file is empty, write the new data as an array and return
	if len(existingData) == 0 {
		newDataArray := []Listdata{newData}
		newDataJSON, err := json.Marshal(newDataArray)
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "JSONMarshalError",
			})
		}
		_, err = file.Write(newDataJSON)
		if err != nil {
			return c.JSON(fiber.Map{
				"code":    400,
				"message": "ErrorDataWrite",
			})
		}
	}

	// If the file is not empty, unmarshal the existing data
	var existingDataArray []Listdata
	err = json.Unmarshal(existingData, &existingDataArray)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "JSONUnmarshalError",
		})
	}

	// Check if the new data already exists in the file
	for i, d := range existingDataArray {
		if d.TitleID == newData.TitleID {
			//проверка текущей оценки с новой
			if existingDataArray[i].Mark != newData.Mark {
				db := database.DB
				switch entity {
				case "anime":
					query := Anime{ID: locid}
					anime := Anime{}
					db.Take(&anime, &query)
					found := Anime{}
					found.Marks = anime.Marks
					switch existingDataArray[i].Mark {
					case 5:
						found.Marks.Mark5--
					case 4:
						found.Marks.Mark4--
					case 3:
						found.Marks.Mark3--
					case 2:
						found.Marks.Mark2--
					case 1:
						found.Marks.Mark1--
					}

					switch newData.Mark {
					case 5:
						found.Marks.Mark5++
					case 4:
						found.Marks.Mark4++
					case 3:
						found.Marks.Mark3++
					case 2:
						found.Marks.Mark2++
					case 1:
						found.Marks.Mark1++
					}
					divider := (found.Marks.Mark1 + found.Marks.Mark2 + found.Marks.Mark3 + found.Marks.Mark4 + found.Marks.Mark5)
					if divider == 0 {
						divider = 1
					}
					var rateComp float32 = float32((found.Marks.Mark5*5)+(found.Marks.Mark4*4)+(found.Marks.Mark3*3)+(found.Marks.Mark2*2)+(found.Marks.Mark1)) / float32(divider)
					found.MarkMean = rateComp
					db.Model(&Anime{}).Where("id = ?", locid).Updates(Anime{MarkMean: found.MarkMean, Marks: found.Marks})
				case "manga":
					query := Manga{ID: locid}
					manga := Manga{}
					db.Take(&manga, &query)
					found := Manga{}
					found.Marks = manga.Marks
					switch existingDataArray[i].Mark {
					case 5:
						found.Marks.Mark5--
					case 4:
						found.Marks.Mark4--
					case 3:
						found.Marks.Mark3--
					case 2:
						found.Marks.Mark2--
					case 1:
						found.Marks.Mark1--
					}

					switch newData.Mark {
					case 5:
						found.Marks.Mark5++
					case 4:
						found.Marks.Mark4++
					case 3:
						found.Marks.Mark3++
					case 2:
						found.Marks.Mark2++
					case 1:
						found.Marks.Mark1++
					}
					divider := (found.Marks.Mark1 + found.Marks.Mark2 + found.Marks.Mark3 + found.Marks.Mark4 + found.Marks.Mark5)
					if divider == 0 {
						divider = 1
					}
					var rateComp float32 = float32((found.Marks.Mark5*5)+(found.Marks.Mark4*4)+(found.Marks.Mark3*3)+(found.Marks.Mark2*2)+(found.Marks.Mark1)) / float32(divider)
					found.MarkMean = rateComp
					db.Model(&Manga{}).Where("id = ?", locid).Updates(Manga{MarkMean: found.MarkMean, Marks: found.Marks})
				case "ranobe":
					query := Ranobe{ID: locid}
					ranobe := Ranobe{}
					db.Take(&ranobe, &query)
					found := Ranobe{}
					found.Marks = ranobe.Marks
					switch existingDataArray[i].Mark {
					case 5:
						found.Marks.Mark5--
					case 4:
						found.Marks.Mark4--
					case 3:
						found.Marks.Mark3--
					case 2:
						found.Marks.Mark2--
					case 1:
						found.Marks.Mark1--
					}

					switch newData.Mark {
					case 5:
						found.Marks.Mark5++
					case 4:
						found.Marks.Mark4++
					case 3:
						found.Marks.Mark3++
					case 2:
						found.Marks.Mark2++
					case 1:
						found.Marks.Mark1++
					}
					divider := (found.Marks.Mark1 + found.Marks.Mark2 + found.Marks.Mark3 + found.Marks.Mark4 + found.Marks.Mark5)
					if divider == 0 {
						divider = 1
					}
					var rateComp float32 = float32((found.Marks.Mark5*5)+(found.Marks.Mark4*4)+(found.Marks.Mark3*3)+(found.Marks.Mark2*2)+(found.Marks.Mark1)) / float32(divider)
					found.MarkMean = rateComp
					db.Model(&Ranobe{}).Where("id = ?", locid).Updates(Ranobe{MarkMean: found.MarkMean, Marks: found.Marks})
				default:
					return c.JSON(fiber.Map{
						"code":    400,
						"message": "UnknownEntity",
					})
				}

			}
			// If the new data exists, update it and write the updated data back to the file
			existingDataArray[i] = newData
			updatedDataJSON, err := json.Marshal(existingDataArray)
			if err != nil {
				return c.JSON(fiber.Map{
					"code":    400,
					"message": "NewJSONUnmarshalError",
				})
			}
			err = ioutil.WriteFile("storage/lists/"+entity+"/id"+claims.Id+"_"+entity+"s.json", updatedDataJSON, 0644)
			if err != nil {
				return c.JSON(fiber.Map{
					"code":    400,
					"message": "NewErrorFileWrite",
				})
			}
			return c.JSON(fiber.Map{
				"code":    200,
				"message": "success",
			})
		}
	}

	// If the new data does not exist, append it to the existing data and write the updated data back to the file
	db := database.DB
	switch entity {
	case "anime":
		query := Anime{ID: locid}
		anime := Anime{}
		db.Take(&anime, &query)
		found := Anime{}
		found.Marks = anime.Marks

		switch newData.Mark {
		case 5:
			found.Marks.Mark5++
		case 4:
			found.Marks.Mark4++
		case 3:
			found.Marks.Mark3++
		case 2:
			found.Marks.Mark2++
		case 1:
			found.Marks.Mark1++
		}
		fmt.Println(found.Marks)
		divider := (found.Marks.Mark1 + found.Marks.Mark2 + found.Marks.Mark3 + found.Marks.Mark4 + found.Marks.Mark5)
		if divider == 0 {
			divider = 1
		}
		found.MarkMean = float32(((found.Marks.Mark5 * 5) + (found.Marks.Mark4 * 4) + (found.Marks.Mark3 * 3) + (found.Marks.Mark2 * 2) + (found.Marks.Mark1)) / divider)
		db.Model(&Anime{}).Where("id = ?", locid).Updates(Anime{MarkMean: found.MarkMean, Marks: found.Marks})
	case "manga":
		query := Manga{ID: locid}
		manga := Manga{}
		db.Take(&manga, &query)
		found := Manga{}
		found.Marks = manga.Marks

		switch newData.Mark {
		case 5:
			found.Marks.Mark5++
		case 4:
			found.Marks.Mark4++
		case 3:
			found.Marks.Mark3++
		case 2:
			found.Marks.Mark2++
		case 1:
			found.Marks.Mark1++
		}
		fmt.Println(found.Marks)
		divider := (found.Marks.Mark1 + found.Marks.Mark2 + found.Marks.Mark3 + found.Marks.Mark4 + found.Marks.Mark5)
		if divider == 0 {
			divider = 1
		}
		found.MarkMean = float32(((found.Marks.Mark5 * 5) + (found.Marks.Mark4 * 4) + (found.Marks.Mark3 * 3) + (found.Marks.Mark2 * 2) + (found.Marks.Mark1)) / divider)
		db.Model(&Manga{}).Where("id = ?", locid).Updates(Manga{MarkMean: found.MarkMean, Marks: found.Marks})
	case "ranobe":
		query := Ranobe{ID: locid}
		ranobe := Ranobe{}
		db.Take(&ranobe, &query)
		found := Anime{}
		found.Marks = ranobe.Marks

		switch newData.Mark {
		case 5:
			found.Marks.Mark5++
		case 4:
			found.Marks.Mark4++
		case 3:
			found.Marks.Mark3++
		case 2:
			found.Marks.Mark2++
		case 1:
			found.Marks.Mark1++
		}
		fmt.Println(found.Marks)
		divider := (found.Marks.Mark1 + found.Marks.Mark2 + found.Marks.Mark3 + found.Marks.Mark4 + found.Marks.Mark5)
		if divider == 0 {
			divider = 1
		}
		found.MarkMean = float32(((found.Marks.Mark5 * 5) + (found.Marks.Mark4 * 4) + (found.Marks.Mark3 * 3) + (found.Marks.Mark2 * 2) + (found.Marks.Mark1)) / divider)
		db.Model(&Ranobe{}).Where("id = ?", locid).Updates(Ranobe{MarkMean: found.MarkMean, Marks: found.Marks})
	default:
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "UnknownEntity",
		})
	}

	existingDataArray = append(existingDataArray, newData)
	updatedDataJSON, err := json.Marshal(existingDataArray)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "NewJSONMarshalError",
		})
	}
	err = ioutil.WriteFile("storage/lists/"+entity+"/id"+claims.Id+"_"+entity+"s.json", updatedDataJSON, 0644)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorDataWrite",
		})
	}

	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
	})
}

func RemoveEntityFromList(c *fiber.Ctx) error {
	locid := EntityIDFromSlug(c.Params("slug"), c.Params("entity"))
	entity := c.Params("entity")
	if locid == 0 {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "TitleNotFound",
		})
	}
	newData := Listdata{
		TitleID: locid,
	}
	cookie := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil //using the SecretKey which was generated in th Login function
	})

	claims := token.Claims.(*jwt.StandardClaims)
	// Open the file for reading and writing
	file, err := os.OpenFile("storage/lists/"+entity+"/id"+claims.Id+"_"+entity+"s.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorFileOpen",
		})
	}
	defer file.Close()
	// Read the existing data from the file
	existingData, err := ioutil.ReadAll(file)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorFileRead",
		})
	}

	// If the file is empty, return
	if len(existingData) == 0 {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorFileEmpty",
		})
	}

	// Unmarshal the existing data
	var existingDataArray []Listdata
	err = json.Unmarshal(existingData, &existingDataArray)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "JSONUnmarshalError",
		})
	}

	// Find the data with the specified ID and remove it from the array
	var newDataArray []Listdata
	found := false
	for _, d := range existingDataArray {
		if d.TitleID != newData.TitleID {
			newDataArray = append(newDataArray, d)
		} else {
			found = true
			db := database.DB
			switch entity {
			case "anime":
				query := Anime{ID: locid}
				anime := Anime{}
				db.Take(&anime, &query)
				found := Anime{}
				found.Marks = anime.Marks
				switch d.Mark {
				case 5:
					found.Marks.Mark5--
				case 4:
					found.Marks.Mark4--
				case 3:
					found.Marks.Mark3--
				case 2:
					found.Marks.Mark2--
				case 1:
					found.Marks.Mark1--
				}
				divider := (found.Marks.Mark1 + found.Marks.Mark2 + found.Marks.Mark3 + found.Marks.Mark4 + found.Marks.Mark5)
				if divider == 0 {
					divider = 1
				}
				var rateComp float32 = float32((found.Marks.Mark5*5)+(found.Marks.Mark4*4)+(found.Marks.Mark3*3)+(found.Marks.Mark2*2)+(found.Marks.Mark1)) / float32(divider)
				found.MarkMean = rateComp
				db.Model(&Anime{}).Where("id = ?", locid).Updates(Anime{MarkMean: found.MarkMean, Marks: found.Marks})
			case "manga":
				query := Manga{ID: locid}
				manga := Manga{}
				db.Take(&manga, &query)
				found := Manga{}
				found.Marks = manga.Marks
				switch d.Mark {
				case 5:
					found.Marks.Mark5--
				case 4:
					found.Marks.Mark4--
				case 3:
					found.Marks.Mark3--
				case 2:
					found.Marks.Mark2--
				case 1:
					found.Marks.Mark1--
				}
				divider := (found.Marks.Mark1 + found.Marks.Mark2 + found.Marks.Mark3 + found.Marks.Mark4 + found.Marks.Mark5)
				if divider == 0 {
					divider = 1
				}
				var rateComp float32 = float32((found.Marks.Mark5*5)+(found.Marks.Mark4*4)+(found.Marks.Mark3*3)+(found.Marks.Mark2*2)+(found.Marks.Mark1)) / float32(divider)
				found.MarkMean = rateComp
				db.Model(&Manga{}).Where("id = ?", locid).Updates(Manga{MarkMean: found.MarkMean, Marks: found.Marks})
			case "ranobe":
				query := Ranobe{ID: locid}
				ranobe := Ranobe{}
				db.Take(&ranobe, &query)
				found := Anime{}
				found.Marks = ranobe.Marks
				switch d.Mark {
				case 5:
					found.Marks.Mark5--
				case 4:
					found.Marks.Mark4--
				case 3:
					found.Marks.Mark3--
				case 2:
					found.Marks.Mark2--
				case 1:
					found.Marks.Mark1--
				}
				divider := (found.Marks.Mark1 + found.Marks.Mark2 + found.Marks.Mark3 + found.Marks.Mark4 + found.Marks.Mark5)
				if divider == 0 {
					divider = 1
				}
				var rateComp float32 = float32((found.Marks.Mark5*5)+(found.Marks.Mark4*4)+(found.Marks.Mark3*3)+(found.Marks.Mark2*2)+(found.Marks.Mark1)) / float32(divider)
				found.MarkMean = rateComp
				db.Model(&Ranobe{}).Where("id = ?", locid).Updates(Ranobe{MarkMean: found.MarkMean, Marks: found.Marks})
			default:
				return c.JSON(fiber.Map{
					"code":    400,
					"message": "UnknownEntity",
				})
			}
		}
	}

	// If the data with the specified ID was not found, return
	if !found {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "DataNotFound",
		})
	}

	// Marshal the updated data and write it back to the file
	updatedDataJSON, err := json.Marshal(newDataArray)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "NewJSONUnmarshalError",
		})
	}
	err = ioutil.WriteFile("storage/lists/"+entity+"/id"+claims.Id+"_"+entity+"s.json", updatedDataJSON, 0644)
	if err != nil {
		return c.JSON(fiber.Map{
			"code":    400,
			"message": "ErrorFileWrite",
		})
	}
	return c.JSON(fiber.Map{
		"code":    200,
		"message": "success",
	})
}
