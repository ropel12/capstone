package handler

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	entity "github.com/education-hub/BE/app/entities"
	"github.com/education-hub/BE/app/features/school/service"
	"github.com/education-hub/BE/config/dependency"
	"github.com/education-hub/BE/helper"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"go.uber.org/dig"
)

type School struct {
	dig.In
	Service   service.SchoolService
	Dep       dependency.Depend
	Gmeetsess map[int]bool
}

func (u *School) Create(c echo.Context) error {
	var req entity.ReqCreateSchool
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQREGISTERSCHOOL, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}

	pdffile, err3 := c.FormFile("pdf")
	imagefile, err2 := c.FormFile("image")
	if err2 != nil || err3 != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Missing Image or PDF", nil))
	}
	if pdffile.Size > 2*1024*1024 || imagefile.Size > 2*1024*1024 {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "File is too large. Maximum size is 2MB.", nil))

	}
	//load image
	image, err := imagefile.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot Load Image", nil))
	}
	req.Image = imagefile.Filename
	//load pdf
	pdf, err := pdffile.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot Load PDF", nil))
	}
	req.Pdf = pdffile.Filename
	req.UserId = helper.GetUid(c.Get("user").(*jwt.Token))
	//END
	id, err := u.Service.Create(c.Request().Context(), req, image, pdf)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusCreated, CreateWebResponse(http.StatusCreated, "Status Created", map[string]any{"id": id}))
}
func (u *School) Update(c echo.Context) error {
	var req entity.ReqUpdateSchool
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQUPDATESCHOOL, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}
	var image multipart.File
	var pdf multipart.File
	pdffile, err := c.FormFile("pdf")
	if err == nil {
		if pdffile.Size > 2*1024*1024 {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "File is too large. Maximum size is 2MB.", nil))

		}
		fpdf, err := pdffile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot Load PDF", nil))
		}
		req.Pdf = pdffile.Filename
		pdf = fpdf
	}
	imagefile, err := c.FormFile("image")
	if err == nil {
		fimage, err := imagefile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot Load Image", nil))
		}
		req.Image = imagefile.Filename
		image = fimage
	}
	res, err := u.Service.Update(c.Request().Context(), req, image, pdf)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}

func (u *School) Search(c echo.Context) error {
	searchval := c.QueryParam("q")
	if searchval == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "School is missing in the query param", nil))
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", u.Service.Search(searchval)))
}

func (u *School) AddAchievement(c echo.Context) error {
	req := entity.ReqAddAchievemnt{}
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQADDACHIEVEMENT, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}
	filehead, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Image is Missing", nil))
	}
	req.Image = filehead.Filename
	image, err := filehead.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot laod image", nil))
	}
	if filehead.Size > 2*1024*1024 {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "File is too large. Maximum size is 2MB.", nil))
	}
	id, err := u.Service.AddAchievement(c.Request().Context(), req, image)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"id": id}))
}

func (u *School) UpdateAchievement(c echo.Context) error {
	req := entity.ReqUpdateAchievemnt{}
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQUPDATEACHIEVEMENT, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}
	filehead, err := c.FormFile("image")
	var image multipart.File
	if err == nil {
		multipart, err := filehead.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot Load Image", nil))
		}
		req.Image = filehead.Filename
		image = multipart
	}
	schoolid, err := u.Service.UpdateAchievement(c.Request().Context(), req, image)
	if err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"id": schoolid}))
}
func (u *School) DeleteAchievement(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Achievement id is missing", nil))
	}
	newid, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Achievement Id", nil))
	}
	if err := u.Service.DeleteAchievement(c.Request().Context(), newid); err != nil {
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}

func (u *School) GenerateUrl(c echo.Context) error {
	req := entity.ReqCreateGmeet{}
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR]WHEN BINDING REQGMEET, Err : %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	if err := c.Validate(req); err != nil {
		return CreateErrorResponse(err, c)
	}
	start_date := fmt.Sprintf("%s+07:00", req.StartDate)
	end_date := fmt.Sprintf("%s+07:00", req.EndDate)
	u.Gmeetsess[req.SchoolId] = true
	url := u.Dep.Calendar.GenerateUrl(start_date, end_date, req.SchoolId)
	return c.Redirect(http.StatusFound, url)
}

func (u *School) CreateGmeet(c echo.Context) error {
	err := c.QueryParam("error")
	auth := c.QueryParam("code")
	state := c.QueryParam("state")
	data := strings.Split(state, "?")
	if err != "" || len(data) < 3 || len(data) > 4 || auth == "" {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}
	schoolid, err1 := strconv.Atoi(data[3])
	if err1 != nil {
		u.Dep.Log.Errorf("[ERROR]WHEN CONVERT SCHOOL ID , Error: %v", err)
	}
	if !u.Gmeetsess[schoolid] {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}
	schooldata, err1 := u.Service.GetByid(c.Request().Context(), schoolid)
	gmeetlink := u.Dep.Calendar.NewService(auth).Create(strings.Replace(data[1], ":00+07", "", 1), strings.Replace(data[2], ":00+07", "", 1), schooldata.Name)
	u.Service.Update(c.Request().Context(), entity.ReqUpdateSchool{Id: schooldata.Id, Gmeet: gmeetlink}, nil, nil)
	delete(u.Gmeetsess, schoolid)
	return c.Redirect(http.StatusFound, URLFRONTENDSUCCESSCREATEDGMEET)
}
