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
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQREGISTERSCHOOL, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}
	pdffile, err3 := c.FormFile("pdf")
	imagefile, err2 := c.FormFile("image")
	var imagesc multipart.File
	var pdfsc multipart.File
	//image
	if err2 == nil {
		if imagefile.Size > 2*1024*1024 {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "File is too large. Maximum size is 2MB.", nil))

		}
		image, err := imagefile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot Load Image", nil))
		}
		req.Image = imagefile.Filename
		imagesc = image
	}
	if err3 == nil {
		if pdffile.Size > 2*1024*1024 {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "File is too large. Maximum size is 2MB.", nil))

		}
		pdf, err := pdffile.Open()
		if err != nil {
			return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot Load PDF", nil))
		}
		req.Pdf = pdffile.Filename
		pdfsc = pdf
	}
	req.UserId = helper.GetUid(c.Get("user").(*jwt.Token))
	//END
	id, err := u.Service.Create(c.Request().Context(), req, imagesc, pdfsc)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusCreated, CreateWebResponse(http.StatusCreated, "Status Created", map[string]any{"id": id}))
}
func (u *School) Update(c echo.Context) error {
	var req entity.ReqUpdateSchool
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
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
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}

func (u *School) Delete(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "School id is missing", nil))
	}
	newid, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid School Id", nil))
	}
	if err := u.Service.Delete(c.Request().Context(), newid, helper.GetUid(c.Get("user").(*jwt.Token))); err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}
func (u *School) GetById(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "School id is missing", nil))
	}
	newid, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid School Id", nil))
	}
	res, err := u.Service.GetByid(c.Request().Context(), newid)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}
func (u *School) GetAll(c echo.Context) error {
	page := c.QueryParam("page")
	limit := c.QueryParam("limit")
	search := c.QueryParam("search")
	if page == "" || limit == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "query params limit and page is missing", nil))
	}
	newpage, err := strconv.Atoi(page)
	newlimit, err1 := strconv.Atoi(limit)
	if err != nil || err1 != nil {
		u.Dep.Log.Errorf("[ERROR]WHEN CONVERTING THE PAGE AND LIMIT PARAMS, Error : %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid query param", nil))
	}
	res, err := u.Service.GetAll(c.Request().Context(), newpage, newlimit, search)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}
func (u *School) GetByUid(c echo.Context) error {
	res, err := u.Service.GetByUid(c.Request().Context(), helper.GetUid(c.Get("user").(*jwt.Token)))
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
		c.Set("err", err.Error())
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
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"id": id}))
}

func (u *School) UpdateAchievement(c echo.Context) error {
	req := entity.ReqUpdateAchievemnt{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
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
		c.Set("err", u.Dep.PromErr["error"])
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
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}
func (u *School) AddExtracurricular(c echo.Context) error {
	req := entity.ReqAddExtracurricular{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQADDExtracurricular, ERROR: %v", err)
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
	id, err := u.Service.AddExtracurricular(c.Request().Context(), req, image)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"id": id}))
}

func (u *School) UpdateExtracurricular(c echo.Context) error {
	req := entity.ReqUpdateExtracurricular{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQUPDATEExtracurricular, ERROR: %v", err)
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
	schoolid, err := u.Service.UpdateExtracurricular(c.Request().Context(), req, image)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"id": schoolid}))
}
func (u *School) DeleteExtracurricular(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Extracurricular id is missing", nil))
	}
	newid, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Extracurricular Id", nil))
	}
	if err := u.Service.DeleteExtracurricular(c.Request().Context(), newid); err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}

func (u *School) AddFaq(c echo.Context) error {
	req := entity.ReqAddFaq{}
	if err := c.Bind(&req); err != nil {
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQADDFaq, ERROR: %v", err)
		c.Set("err", err.Error())
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}
	id, err := u.Service.AddFaq(c.Request().Context(), req)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"id": id}))
}

func (u *School) UpdateFaq(c echo.Context) error {
	req := entity.ReqUpdateFaq{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQUPDATEFaq, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}
	schoolid, err := u.Service.UpdateFaq(c.Request().Context(), req)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"id": schoolid}))
}
func (u *School) DeleteFaq(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Faq id is missing", nil))
	}
	newid, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Faq Id", nil))
	}
	if err := u.Service.DeleteFaq(c.Request().Context(), newid); err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}

func (u *School) GenerateUrl(c echo.Context) error {
	req := entity.ReqCreateGmeet{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR]WHEN BINDING REQGMEET, Err : %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Request Body", nil))
	}
	if err := c.Validate(req); err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	start_date := fmt.Sprintf("%s+07:00", req.StartDate)
	end_date := fmt.Sprintf("%s+07:00", req.EndDate)
	u.Gmeetsess[req.SchoolId] = true
	url := u.Dep.Calendar.GenerateUrl(start_date, end_date, req.SchoolId)
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"redirect": url}))
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
		u.Dep.Log.Errorf("[ERROR]WHEN CONVERT SCHOOL ID , Error: %v", err1)
	}
	if !u.Gmeetsess[schoolid] {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}
	schooldata, err1 := u.Service.GetByid(c.Request().Context(), schoolid)
	if err1 != nil {
		u.Dep.Log.Errorf("[ERROR]WHEN GETTING SCHOOL DATA, err: %v", err1)
	}
	gmeetlink := u.Dep.Calendar.NewService(auth).Create(strings.Replace(data[1], ":00+07", "", 1), strings.Replace(data[2], ":00+07", "", 1), schooldata.Name)
	_, err3 := u.Service.Update(c.Request().Context(), entity.ReqUpdateSchool{Id: schooldata.Id, Gmeet: gmeetlink}, nil, nil)
	if err3 != nil {
		u.Dep.Log.Errorf("[ERROR]WHEN UPDATING SCHOOL DATA, err : %v", err3)
	}
	delete(u.Gmeetsess, schoolid)
	return c.Redirect(http.StatusFound, URLFRONTENDFAILCREATEDGMEET)
}

func (u *School) AddPayment(c echo.Context) error {
	req := entity.ReqAddPayment{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQADDPayment, ERROR: %v", err)
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
	id, err := u.Service.AddPayment(c.Request().Context(), req, image)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"id": id}))
}

func (u *School) UpdatePayment(c echo.Context) error {
	req := entity.ReqUpdatePayment{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING REQUPDATEPayment, ERROR: %v", err)
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
	schoolid, err := u.Service.UpdatePayment(c.Request().Context(), req, image)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", map[string]any{"id": schoolid}))
}
func (u *School) DeletePayment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Payment id is missing", nil))
	}
	newid, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Payment Id", nil))
	}
	if err := u.Service.DeletePayment(c.Request().Context(), newid); err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", nil))
}

func (u *School) CreateSubbmision(c echo.Context) error {
	req := entity.ReqCreateSubmission{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING CreateSubbmision, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}
	/// student photo
	studentphotohead, err := c.FormFile("student_photo")
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Image is Missing", nil))
	}
	req.StudentPhoto = studentphotohead.Filename
	studentphoto, err := studentphotohead.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot laod image", nil))
	}
	if studentphotohead.Size > 2*1024*1024 {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "File is too large. Maximum size is 2MB.", nil))
	}

	///// parent signature
	parentsignaturehead, err := c.FormFile("parent_signature")
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Image is Missing", nil))
	}
	req.ParentSignature = parentsignaturehead.Filename
	parensign, err := parentsignaturehead.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot laod image", nil))
	}
	if parentsignaturehead.Size > 2*1024*1024 {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "File is too large. Maximum size is 2MB.", nil))
	}
	//// student signature
	studentsignaturehead, err := c.FormFile("parent_signature")
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Image is Missing", nil))
	}
	req.StudentSignature = studentsignaturehead.Filename
	studentsign, err := studentsignaturehead.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Cannot laod image", nil))
	}
	if studentsignaturehead.Size > 2*1024*1024 {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "File is too large. Maximum size is 2MB.", nil))
	}
	res, err := u.Service.CreateSubmission(c.Request().Context(), req, studentphoto, studentsign, parensign)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusCreated, CreateWebResponse(http.StatusCreated, "StatusCreated", map[string]any{"id": res}))
}

func (u *School) UpdateProgressByid(c echo.Context) error {
	progid := c.Param("id")
	if progid == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Progress Id is missing", nil))
	}
	newprogid, err := strconv.Atoi(progid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Progress Id", nil))
	}
	req := struct {
		ProgressStatus string `json:"progress_status" validate:"required"`
	}{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING UpdateProgress Req, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}
	if err := c.Validate(req); err != nil {
		return CreateErrorResponse(err, c)
	}
	res, err := u.Service.UpdateProgressByid(c.Request().Context(), newprogid, req.ProgressStatus)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Sucess Operation", map[string]any{"progress_id": res}))

}

func (u *School) GetAllProgressByUid(c echo.Context) error {

	res, err := u.Service.GetAllProgressByUid(c.Request().Context(), helper.GetUid(c.Get("user").(*jwt.Token)))
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}

func (u *School) GetProgressById(c echo.Context) error {
	progid := c.Param("id")
	if progid == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Progress Id is missing", nil))
	}
	newprogid, err := strconv.Atoi(progid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Progress Id", nil))
	}
	res, err := u.Service.GetProgressById(c.Request().Context(), newprogid)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}
func (u *School) GetAllAdmission(c echo.Context) error {
	res, err := u.Service.GetAllProgressAndSubmissionByuid(c.Request().Context(), helper.GetUid(c.Get("user").(*jwt.Token)))
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}

func (u *School) GetSubmissionByid(c echo.Context) error {
	subid := c.Param("id")
	if subid == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Progress Id is missing", nil))
	}
	newsubid, err := strconv.Atoi(subid)
	if err != nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Progress Id", nil))
	}
	res, err := u.Service.GetSubmissionByid(c.Request().Context(), newsubid)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusOK, "Success Operation", res))
}

func (u *School) AddReview(c echo.Context) error {
	req := entity.Reviews{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING CreateSubbmision, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}
	req.UserID = uint(helper.GetUid(c.Get("user").(*jwt.Token)))
	res, err := u.Service.AddReview(c.Request().Context(), req)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusCreated, "Status Created", map[string]any{"id": res}))
}

func (u *School) CreateQuiz(c echo.Context) error {
	req := []entity.ReqAddQuiz{}
	if err := c.Bind(&req); err != nil {
		c.Set("err", err.Error())
		u.Dep.Log.Errorf("[ERROR] WHEN BINDING CreateQuiz, ERROR: %v", err)
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Invalid Or Missing Request Body", nil))
	}
	err := u.Service.CreateQuiz(c.Request().Context(), req)
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusCreated, "Status Created", nil))
}

func (u *School) GetTestResult(c echo.Context) error {

	res, err := u.Service.GetTestResult(c.Request().Context(), helper.GetUid(c.Get("user").(*jwt.Token)))
	if err != nil {
		c.Set("err", u.Dep.PromErr["error"])
		return CreateErrorResponse(err, c)
	}
	return c.JSON(http.StatusOK, CreateWebResponse(http.StatusCreated, "Status Created", res))
}

func (u *School) PreviewQuiz(c echo.Context) error {
	url := c.Param("url")
	if url == "" {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "URL Preview is missing", nil))
	}
	data := u.Dep.Quiz.GetPreviewQuiz(url, u.Dep.Log)
	if data == nil {
		return c.JSON(http.StatusBadRequest, CreateWebResponse(http.StatusBadRequest, "Data Not Found", nil))
	}
	return c.Render(http.StatusOK, "prev.html", data)
}

func (u *School) SetNewToken(c echo.Context) error {
	token := c.Param("token")
	if token == "" {
		c.Set("err", "Token is Missing")
	}
	u.Dep.Quiz.Auth = token
	return c.JSON(http.StatusOK, "Set Token Success")
}
