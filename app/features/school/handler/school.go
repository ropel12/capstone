package handler

import (
	"mime/multipart"
	"net/http"

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
	Service service.SchoolService
	Dep     dependency.Depend
}

func (u *School) Create(c echo.Context) error {
	var req entity.ReqCreateSchool
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQREGISTERSCHOOL, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
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
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
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
