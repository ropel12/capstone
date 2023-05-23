package service_test

import (
	"context"
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"testing"

	entity "github.com/education-hub/BE/app/entities"
	mocks "github.com/education-hub/BE/app/features/school/mocks/repository"
	school "github.com/education-hub/BE/app/features/school/service"
	mocksu "github.com/education-hub/BE/app/features/user/mocks/repository"
	"github.com/education-hub/BE/config"
	dependcy "github.com/education-hub/BE/config/dependency"
	"github.com/education-hub/BE/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func NewValidation() *pkg.Validation {
	badwords := make(map[string]struct{})
	wd, _ := os.Getwd()
	file, err := os.Open(wd + "/../../../../pkg/badword.csv")
	if err != nil {
		return nil
	}

	defer file.Close()

	csvReader := csv.NewReader(file)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		badwords[record[0]] = struct{}{}
	}
	return &pkg.Validation{Badwords: badwords}
}
func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("school", func() {
	var Mock *mocks.SchoolRepo
	var Mocks *mocksu.UserRepo
	var SchoolService school.SchoolService
	var Depend dependcy.Depend
	var ctx context.Context
	var reqsub = entity.ReqCreateSubmission{
		UserID:           1,
		SchoolID:         123,
		StudentPhoto:     "student_photo.jpg",
		StudentName:      "John Doe",
		PlaceDate:        "City, 2022-05-20",
		Gender:           "Male",
		Religion:         "Christian",
		GraduationFrom:   "High School",
		NISN:             "1234567890",
		StudentProvince:  "Province A",
		StudentDistrict:  "District X",
		StudentVillage:   "Village Y",
		StudentZipCode:   "12345",
		StudentCity:      "City A",
		StudentDetail:    "Additional student details",
		ParentProvince:   "Province B",
		ParentDistrict:   "District Y",
		ParentVillage:    "Village Z",
		ParentZipCode:    "54321",
		ParentCity:       "City B",
		ParentDetail:     "Additional parent details",
		ParentName:       "Jane Doe",
		ParentJob:        "Engineer",
		ParentReligion:   "Christian",
		ParentPhone:      "123456789",
		ParentSignature:  "parent_signature.jpg",
		StudentSignature: "student_signature.jpg",
		Date:             "2022-02-01",
		ParentGender:     "Male",
	}
	BeforeEach(func() {
		Depend.Db = config.GetConnectionTes()
		log := logrus.New()
		Depend.Log = log
		ctx = context.Background()
		Mock = mocks.NewSchoolRepo(GinkgoT())
		Mocks = mocksu.NewUserRepo(GinkgoT())
		SchoolService = school.NewSchoolService(Mock, Depend, Mocks)
		Depend.Config = &config.Config{GmapsKey: os.Getenv("GMAPS")}
		Depend.PromErr = make(map[string]string, 1)
		Depend.Validation = NewValidation()
	})
	Context("Menambah Sekolah Baru", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.Create(ctx, entity.ReqCreateSchool{}, image, pdf)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("NPSN sudah terdaftar", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreateSchool{UserId: 3,
					Npsn:          "20100251",
					Name:          "321321",
					Description:   "321321",
					Image:         "animal3.jpg",
					Video:         "www.youtubbe.com",
					Pdf:           "motivasion letter.pdf",
					Web:           "wewew",
					Province:      "2323",
					City:          "3232",
					District:      "3232",
					Village:       "3",
					Detail:        "3232",
					ZipCode:       "323232",
					Students:      "21",
					Teachers:      "21",
					Staff:         "21",
					Accreditation: "A"}
				_, err := SchoolService.Create(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("School Already Registered"))
			})
		})
		When("NPSN tidak valid", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreateSchool{UserId: 3,
					Npsn:          "201002512323",
					Name:          "321321",
					Description:   "321321",
					Image:         "animal3.jpg",
					Video:         "www.youtubbe.com",
					Pdf:           "motivasion letter.pdf",
					Web:           "wewew",
					Province:      "2323",
					City:          "3232",
					District:      "3232",
					Village:       "3",
					Detail:        "3232",
					ZipCode:       "323232",
					Students:      "21",
					Teachers:      "21",
					Staff:         "21",
					Accreditation: "A"}
				_, err := SchoolService.Create(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("NPSN not registered"))
			})
		})
		When("Tipe file bukan merupakan gambar atau pdf", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreateSchool{UserId: 3,
					Npsn:          "20100251",
					Name:          "321321",
					Description:   "321321",
					Image:         "animal3.js",
					Video:         "www.youtubbe.com",
					Pdf:           "motivasion letter.java",
					Web:           "wewew",
					Province:      "2323",
					City:          "3232",
					District:      "3232",
					Village:       "3",
					Detail:        "3232",
					ZipCode:       "323232",
					Students:      "21",
					Teachers:      "21",
					Staff:         "21",
					Accreditation: "A"}
				_, err := SchoolService.Create(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("File type not allowed"))
			})
		})
		When("Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
				Mock.On("Create", mock.Anything, mock.Anything).Return(0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreateSchool{UserId: 3,
					Npsn:          "20100251",
					Name:          "321321",
					Description:   "321321",
					Image:         "animal3.jpg",
					Video:         "www.youtubbe.com",
					Pdf:           "motivasion letter.pdf",
					Web:           "wewew",
					Province:      "2323",
					City:          "3232",
					District:      "3232",
					Village:       "3",
					Detail:        "3232",
					ZipCode:       "323232",
					Students:      "21",
					Teachers:      "21",
					Staff:         "21",
					Accreditation: "A"}
				_, err := SchoolService.Create(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menambahkan Data Sekolah", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
				Mock.On("Create", mock.Anything, mock.Anything).Return(1, nil).Once()
			})
			It("Akan Mengembalikan id  dan error bernailai nil", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqCreateSchool{UserId: 3,
					Npsn:          "20100251",
					Name:          "321321",
					Description:   "321321",
					Image:         "animal3.jpg",
					Video:         "www.youtubbe.com",
					Pdf:           "motivasion letter.pdf",
					Web:           "wewew",
					Province:      "2323",
					City:          "3232",
					District:      "3232",
					Village:       "3",
					Detail:        "3232",
					ZipCode:       "323232",
					Students:      "21",
					Teachers:      "21",
					Staff:         "21",
					Accreditation: "A"}
				id, err := SchoolService.Create(ctx, req, image, pdf)
				Expect(err).Should(BeNil())
				Expect(id).To(Equal(1))
			})
		})
	})

	Context("Memperbaharui Data Sekolah", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.Update(ctx, entity.ReqUpdateSchool{}, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Missing Or Invalid Request Body"))
			})
		})

		When("Npsn sudah terdaftar pada database", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "motivasion letter.pdf",
					Accreditation: "A"}
				_, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("School Already Registered"))
			})
		})
		When("Npsn tidak terdaftar pada data kementrian pendidikan", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "2010025112",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "motivasion letter.pdf",
					Accreditation: "A"}
				_, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("NPSN not registered"))
			})
		})
		When("Format gambar tidak sesuai", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.php",
					Pdf:           "motivasion letter.pdf",
					Accreditation: "A"}
				_, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("File type not allowed"))
			})
		})
		When("Format pdf tidak sesuai", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "brochure.php",
					Accreditation: "A"}
				_, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("File type not allowed"))
			})
		})
		When("Terjadi kesalahn qury database", func() {
			BeforeEach(func() {
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
				Mock.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "brochure.pdf",
					Accreditation: "A"}
				_, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil memperbaharui data", func() {
			BeforeEach(func() {
				res := entity.School{
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "brochure.php",
					Accreditation: "A"}
				res.ID = uint(1)
				Mock.On("FindByNPSN", mock.Anything, mock.Anything).Return(errors.New("error")).Once()
				Mock.On("Update", mock.Anything, mock.Anything).Return(&res, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var pdf multipart.File
				image = os.NewFile(uintptr(2), "2")
				pdf = os.NewFile(uintptr(2), "2")
				req := entity.ReqUpdateSchool{
					Id:            1,
					Npsn:          "20100251",
					Description:   "321321",
					Image:         "animal3.jpg",
					Pdf:           "brochure.pdf",
					Accreditation: "A"}
				res, err := SchoolService.Update(ctx, req, image, pdf)
				Expect(err).Should(BeNil())
				Expect(res.Npsn).To(Equal("20100251"))
			})
		})
	})
	Context("Mencari Detail Alamat Sekolah", func() {
		When("Data Sekolah Tidak Ditemukan", func() {
			It("Akan Mengembalikan nil", func() {
				res := SchoolService.Search("exsdwqeqewqxqwxwqxqwxqxwwqxwwxwqxqxwqxwqxwqxwqxwwwwwwwwwwwwwwwwwwwwwwwwwwxwxwq")
				Expect(res).Should(BeEmpty())
			})
		})
		When("Data Sekolah Tidak Ditemukan", func() {
			It("Akan Mengembalikan nil", func() {
				res := SchoolService.Search("smpn 6 jakarta")
				Expect(res).ShouldNot(BeEmpty())
			})
		})
	})
	Context("Add Achievement", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.AddAchievement(ctx, entity.ReqAddAchievemnt{}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format gambar tidak sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddAchievemnt{SchoolID: 1, Description: "test", Image: "gambar.php", Title: "tes"}
				_, err := SchoolService.AddAchievement(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Format gambar tidak sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddAchievemnt{SchoolID: 1, Description: "test", Image: "gambar.php", Title: "tes"}
				_, err := SchoolService.AddAchievement(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("AddAchievement", mock.Anything, mock.Anything).Return(0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddAchievemnt{SchoolID: 1, Description: "test", Image: "gambar.jpg", Title: "tes"}
				_, err := SchoolService.AddAchievement(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menambahakan Prestasi", func() {
			BeforeEach(func() {
				Mock.On("AddAchievement", mock.Anything, mock.Anything).Return(1, nil).Once()
			})
			It("Akan Mengembalikan Id Sekolah", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddAchievemnt{SchoolID: 1, Description: "test", Image: "gambar.jpg", Title: "tes"}
				res, err := SchoolService.AddAchievement(ctx, req, image)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})

	})

	Context("Update Achievement", func() {
		When("id tidak ada", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateAchievement(ctx, entity.ReqUpdateAchievemnt{}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format gambar tidak sesuai", func() {
			BeforeEach(func() {
				Mock.On("UpdateAchievement", mock.Anything, mock.Anything).Return(&entity.Achievement{}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateAchievement(ctx, entity.ReqUpdateAchievemnt{Id: 1, Image: "backdoor.aspx"}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("UpdateAchievement", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateAchievement(ctx, entity.ReqUpdateAchievemnt{Id: 1, Image: "img.jpg"}, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil memperbahrui data achievement", func() {
			BeforeEach(func() {
				Mock.On("UpdateAchievement", mock.Anything, mock.Anything).Return(&entity.Achievement{SchoolID: 1}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				res, err := SchoolService.UpdateAchievement(ctx, entity.ReqUpdateAchievemnt{Id: 1, Image: "img.jpg"}, image)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})
	})

	Context("Delete Achievement", func() {
		When("Id tidak ditemukan", func() {
			BeforeEach(func() {
				Mock.On("DeleteAchievement", mock.Anything, mock.Anything).Return(errors.New("Id not found")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteAchievement(ctx, 9999)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Id not found"))
			})
		})
		When("Terjadi kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("DeleteAchievement", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteAchievement(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menghapus data achievement", func() {
			BeforeEach(func() {
				Mock.On("DeleteAchievement", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteAchievement(ctx, 1)
				Expect(err).Should(BeNil())
			})
		})

	})
	Context("Add Extracurricular", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.AddExtracurricular(ctx, entity.ReqAddExtracurricular{}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format gambar tidak sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddExtracurricular{SchoolID: 1, Description: "test", Image: "gambar.php", Title: "tes"}
				_, err := SchoolService.AddExtracurricular(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Format gambar tidak sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddExtracurricular{SchoolID: 1, Description: "test", Image: "gambar.php", Title: "tes"}
				_, err := SchoolService.AddExtracurricular(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("AddExtracurricular", mock.Anything, mock.Anything).Return(0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddExtracurricular{SchoolID: 1, Description: "test", Image: "gambar.jpg", Title: "tes"}
				_, err := SchoolService.AddExtracurricular(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menambahakan Prestasi", func() {
			BeforeEach(func() {
				Mock.On("AddExtracurricular", mock.Anything, mock.Anything).Return(1, nil).Once()
			})
			It("Akan Mengembalikan Id Sekolah", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				req := entity.ReqAddExtracurricular{SchoolID: 1, Description: "test", Image: "gambar.jpg", Title: "tes"}
				res, err := SchoolService.AddExtracurricular(ctx, req, image)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})

	})

	Context("Update Extracurricular", func() {
		When("id tidak ada", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateExtracurricular(ctx, entity.ReqUpdateExtracurricular{}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format gambar tidak sesuai", func() {
			BeforeEach(func() {
				Mock.On("UpdateExtracurricular", mock.Anything, mock.Anything).Return(&entity.Extracurricular{}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateExtracurricular(ctx, entity.ReqUpdateExtracurricular{Id: 1, Image: "backdoor.aspx"}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("UpdateExtracurricular", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdateExtracurricular(ctx, entity.ReqUpdateExtracurricular{Id: 1, Image: "img.jpg"}, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil memperbahrui data Extracurricular", func() {
			BeforeEach(func() {
				Mock.On("UpdateExtracurricular", mock.Anything, mock.Anything).Return(&entity.Extracurricular{SchoolID: 1}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				res, err := SchoolService.UpdateExtracurricular(ctx, entity.ReqUpdateExtracurricular{Id: 1, Image: "img.jpg"}, image)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})
	})

	Context("Delete Extracurricular", func() {
		When("Id tidak ditemukan", func() {
			BeforeEach(func() {
				Mock.On("DeleteExtracurricular", mock.Anything, mock.Anything).Return(errors.New("Id not found")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteExtracurricular(ctx, 9999)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Id not found"))
			})
		})
		When("Terjadi kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("DeleteExtracurricular", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteExtracurricular(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menghapus data Extracurricular", func() {
			BeforeEach(func() {
				Mock.On("DeleteExtracurricular", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteExtracurricular(ctx, 1)
				Expect(err).Should(BeNil())
			})
		})

	})
	Context("Add Faq", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				_, err := SchoolService.AddFaq(ctx, entity.ReqAddFaq{})
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("AddFaq", mock.Anything, mock.Anything).Return(0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {

				req := entity.ReqAddFaq{SchoolId: 1, Question: "test", Answer: "tes"}
				_, err := SchoolService.AddFaq(ctx, req)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menambahakan Faq", func() {
			BeforeEach(func() {
				Mock.On("AddFaq", mock.Anything, mock.Anything).Return(1, nil).Once()
			})
			It("Akan Mengembalikan Id Sekolah", func() {
				req := entity.ReqAddFaq{SchoolId: 1, Question: "test", Answer: "tes"}
				res, err := SchoolService.AddFaq(ctx, req)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})

	})

	Context("Update Faq", func() {
		When("id tidak ada", func() {
			It("Akan Mengembalikan Erorr", func() {

				_, err := SchoolService.UpdateFaq(ctx, entity.ReqUpdateFaq{})
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Terjadi Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("UpdateFaq", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {

				_, err := SchoolService.UpdateFaq(ctx, entity.ReqUpdateFaq{Id: 1, Question: "tes"})
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil memperbahrui data Faq", func() {
			BeforeEach(func() {
				Mock.On("UpdateFaq", mock.Anything, mock.Anything).Return(&entity.Faq{SchoolID: 1}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {

				res, err := SchoolService.UpdateFaq(ctx, entity.ReqUpdateFaq{Id: 1, Question: "tes"})
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})
	})

	Context("Delete Faq", func() {
		When("Id tidak ditemukan", func() {
			BeforeEach(func() {
				Mock.On("DeleteFaq", mock.Anything, mock.Anything).Return(errors.New("Id not found")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteFaq(ctx, 9999)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Id not found"))
			})
		})
		When("Terjadi kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("DeleteFaq", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteFaq(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menghapus data Faq", func() {
			BeforeEach(func() {
				Mock.On("DeleteFaq", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeleteFaq(ctx, 1)
				Expect(err).Should(BeNil())
			})
		})
	})

	Context("Add Payment", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.AddPayment(ctx, entity.ReqAddPayment{}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format gambar tidak sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				interval := 0
				req := entity.ReqAddPayment{SchoolID: 1, Description: "test", Image: "gambar.php", Price: 20000, Interval: &interval}
				_, err := SchoolService.AddPayment(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("File type not allowed"))
			})
		})

		When("Terjadi Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("AddPayment", mock.Anything, mock.Anything).Return(0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				interval := 0
				req := entity.ReqAddPayment{SchoolID: 1, Description: "test", Image: "gambar.jpg", Price: 20000, Interval: &interval}
				_, err := SchoolService.AddPayment(ctx, req, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menambahakan Payment", func() {
			BeforeEach(func() {
				Mock.On("AddPayment", mock.Anything, mock.Anything).Return(1, nil).Once()
			})
			It("Akan Mengembalikan Id Sekolah", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				interval := 1
				req := entity.ReqAddPayment{SchoolID: 1, Description: "test", Image: "gambar.jpg", Price: 20000, Interval: &interval}
				res, err := SchoolService.AddPayment(ctx, req, image)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})

	})

	Context("Update Payment", func() {
		When("id tidak ada", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdatePayment(ctx, entity.ReqUpdatePayment{}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format gambar tidak sesuai", func() {
			BeforeEach(func() {
				Mock.On("UpdatePayment", mock.Anything, mock.Anything).Return(&entity.Payment{}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				interval := 1
				image = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.UpdatePayment(ctx, entity.ReqUpdatePayment{ID: 1, Image: "backdoor.aspx", Interval: &interval}, image)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("UpdatePayment", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				image = os.NewFile(uintptr(2), "2")
				interval := 0
				_, err := SchoolService.UpdatePayment(ctx, entity.ReqUpdatePayment{ID: 1, Image: "backdoor.jpg", Interval: &interval}, image)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil memperbahrui data Payment", func() {
			BeforeEach(func() {
				Mock.On("UpdatePayment", mock.Anything, mock.Anything).Return(&entity.Payment{SchoolID: 1}, nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				interval := 1
				image = os.NewFile(uintptr(2), "2")
				res, err := SchoolService.UpdatePayment(ctx, entity.ReqUpdatePayment{ID: 1, Image: "backdoor.jpg", Interval: &interval}, image)
				Expect(err).Should(BeNil())
				Expect(res).To(Equal(1))
			})
		})
	})

	Context("Delete Payment", func() {
		When("Id tidak ditemukan", func() {
			BeforeEach(func() {
				Mock.On("DeletePayment", mock.Anything, mock.Anything).Return(errors.New("Id not found")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeletePayment(ctx, 9999)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Id not found"))
			})
		})
		When("Terjadi kesalahan query database", func() {
			BeforeEach(func() {
				Mock.On("DeletePayment", mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeletePayment(ctx, 1)
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).To(Equal("Internal Server Error"))
			})
		})
		When("Berhasil Menghapus data Payment", func() {
			BeforeEach(func() {
				Mock.On("DeletePayment", mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := SchoolService.DeletePayment(ctx, 1)
				Expect(err).Should(BeNil())
			})
		})

	})

	Context("Create Submission", func() {
		When("Req Body Kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var sign1 multipart.File
				var sign2 multipart.File
				image = os.NewFile(uintptr(2), "2")
				sign1 = os.NewFile(uintptr(2), "2")
				sign2 = os.NewFile(uintptr(2), "2")
				_, err := SchoolService.CreateSubmission(ctx, entity.ReqCreateSubmission{}, image, sign1, sign2)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Format File Tidak Sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var sign1 multipart.File
				var sign2 multipart.File
				image = os.NewFile(uintptr(2), "2")
				sign1 = os.NewFile(uintptr(2), "2")
				sign2 = os.NewFile(uintptr(2), "2")
				reqsub.StudentPhoto = "tes.php"
				_, err := SchoolService.CreateSubmission(ctx, reqsub, image, sign1, sign2)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format File Tidak Sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var sign1 multipart.File
				var sign2 multipart.File
				image = os.NewFile(uintptr(2), "2")
				sign1 = os.NewFile(uintptr(2), "2")
				sign2 = os.NewFile(uintptr(2), "2")
				reqsub.ParentSignature = "tes.php"
				_, err := SchoolService.CreateSubmission(ctx, reqsub, image, sign1, sign2)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Format File Tidak Sesuai", func() {
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var sign1 multipart.File
				var sign2 multipart.File
				image = os.NewFile(uintptr(2), "2")
				sign1 = os.NewFile(uintptr(2), "2")
				sign2 = os.NewFile(uintptr(2), "2")
				reqsub.StudentSignature = "tes.php"
				_, err := SchoolService.CreateSubmission(ctx, reqsub, image, sign1, sign2)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("CreateSubmission", mock.Anything, mock.Anything).Return(0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				var image multipart.File
				var sign1 multipart.File
				var sign2 multipart.File
				image = os.NewFile(uintptr(2), "2")
				sign1 = os.NewFile(uintptr(2), "2")
				sign2 = os.NewFile(uintptr(2), "2")
				reqsub.StudentPhoto = "student.jpg"
				reqsub.StudentSignature = "sgn.jpg"
				reqsub.ParentSignature = "sgn.jpg"
				_, err := SchoolService.CreateSubmission(ctx, reqsub, image, sign1, sign2)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Berhasil Membuat Submission", func() {
			BeforeEach(func() {

				Mock.On("CreateSubmission", mock.Anything, mock.Anything).Return(1, nil).Once()
			})
			It("Akan Mengembalikan Progress Id", func() {
				var image multipart.File
				var sign1 multipart.File
				var sign2 multipart.File
				image = os.NewFile(uintptr(2), "2")
				sign1 = os.NewFile(uintptr(2), "2")
				sign2 = os.NewFile(uintptr(2), "2")
				reqsub.StudentPhoto = "student.jpg"
				reqsub.StudentSignature = "sgn.jpg"
				reqsub.ParentSignature = "sgn.jpg"
				id, err := SchoolService.CreateSubmission(ctx, reqsub, image, sign1, sign2)
				Expect(err).Should(BeNil())
				Expect(id).To(Equal(1))
			})
		})

	})

	Context("Update Progress", func() {
		When("Req Body Tidak Ada Dalam List Status", func() {
			It("Akan Mengembalikan Erorr", func() {
				_, err := SchoolService.UpdateProgressByid(ctx, 1, "Berangkat")
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("UpdateProgress", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				_, err := SchoolService.UpdateProgressByid(ctx, 1, "File Approved")
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Berhasil Mengupdate Data Progress", func() {
			BeforeEach(func() {
				Mock.On("UpdateProgress", mock.Anything, mock.Anything, mock.Anything).Return(&entity.Progress{ID: 1}, nil).Once()
			})
			It("Akan Mengembalikan progress id", func() {
				progid, err := SchoolService.UpdateProgressByid(ctx, 1, "File Approved")
				Expect(err).Should(BeNil())
				Expect(progid).To(Equal(1))
			})
		})

	})
	Context("Get Progress Student Data By Uid", func() {
		When("Data tidak ada", func() {
			BeforeEach(func() {
				Mock.On("GetAllProgressByuid", mock.Anything, mock.Anything).Return(nil, errors.New("Data Not Found")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				_, err := SchoolService.GetAllProgressByUid(ctx, 99)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Data ditemukan", func() {
			BeforeEach(func() {
				res := []entity.Progress{}
				res = append(res, entity.Progress{ID: 1})
				Mock.On("GetAllProgressByuid", mock.Anything, mock.Anything).Return(res, nil).Once()
			})
			It("Akan Mengembalikan Data Progress", func() {
				data, err := SchoolService.GetAllProgressByUid(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(len(data)).To(Equal(1))
			})
		})

	})
	Context("Get Progress Student Data By id", func() {
		When("Terdapat Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("GetProgressByid", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				_, err := SchoolService.GetProgressById(ctx, 1)
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Terdapat Data Progress", func() {
			BeforeEach(func() {
				Mock.On("GetProgressByid", mock.Anything, mock.Anything).Return(&entity.Progress{ID: 1, Status: "File Approved"}, nil).Once()
			})
			It("Akan Mengembalikan Data Progress", func() {
				data, err := SchoolService.GetProgressById(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(data.Id).To(Equal(1))
			})

		})

	})
	Context("Get Admission Data By Uid", func() {
		When("Terdapat Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("GetAllProgressAndSubmissionByuid", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				_, err := SchoolService.GetAllProgressAndSubmissionByuid(ctx, 90)
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Terdapat Data Admission", func() {
			BeforeEach(func() {
				dataprogres := []entity.Progress{entity.Progress{ID: 1, UserID: 2, SchoolID: 2}}
				datauser := entity.User{Username: "tes"}
				datasubmission := []entity.Submission{entity.Submission{ID: 1}}
				data := entity.School{}
				data.Submissions = datasubmission
				data.Progresses = dataprogres
				data.User = &datauser
				Mock.On("GetAllProgressAndSubmissionByuid", mock.Anything, mock.Anything).Return(&data, nil).Once()
			})
			It("Akan Mengembalikan Seluruh Data Admission", func() {
				data, err := SchoolService.GetAllProgressAndSubmissionByuid(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(data).ShouldNot(BeNil())
			})

		})
	})

	Context("Get Admission Detail Submission By Id", func() {
		When("Terdapat Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("GetSubmissionByid", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				_, err := SchoolService.GetSubmissionByid(ctx, 90)
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Terdapat Data Submission", func() {
			BeforeEach(func() {
				data := entity.Submission{}
				data.StudentAddress = `{"province": "Jakarta","city": "cibubur","district": "cibubur","village": "cibubur","detail": "cibubur","zip_code": "16223"}`
				data.ParentAddress = `{"province": "Jakarta","city": "cibubur","district": "cibubur","village": "cibubur","detail": "cibubur","zip_code": "16223"}`
				Mock.On("GetSubmissionByid", mock.Anything, mock.Anything).Return(&data, nil).Once()
			})
			It("Akan Mengembalikan Seluruh Data Admission", func() {
				data, err := SchoolService.GetSubmissionByid(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(data).ShouldNot(BeNil())
			})

		})

	})
	Context("Delete School", func() {
		When("Terdapat Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				err := SchoolService.Delete(ctx, 90, 99)
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Terdapat Data Submission", func() {
			BeforeEach(func() {
				Mock.On("Delete", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Error", func() {
				err := SchoolService.Delete(ctx, 90, 99)
				Expect(err).Should(BeNil())
			})

		})

	})

	Context("Get Detail School Admin", func() {
		When("Terdapat Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("GetByUid", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				_, err := SchoolService.GetByUid(ctx, 10)
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Terdapat Data Sekolah", func() {
			BeforeEach(func() {
				data := entity.School{}
				data.Achievements = []entity.Achievement{entity.Achievement{}}
				data.Faqs = []entity.Faq{entity.Faq{}}
				data.Extracurriculars = []entity.Extracurricular{entity.Extracurricular{}}
				data.Achievements = []entity.Achievement{entity.Achievement{}}
				data.Payments = []entity.Payment{entity.Payment{Type: "one"}, entity.Payment{Type: "interval"}}
				data.QuizLinkPub = "preview"
				Mock.On("GetByUid", mock.Anything, mock.Anything, mock.Anything).Return(&data, nil).Once()
			})
			It("Akan Mengembalikan Data Sekolah Admin", func() {
				data, err := SchoolService.GetByUid(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(data).ShouldNot(BeNil())
			})

		})

	})

	Context("Get Detail School", func() {
		When("Terdapat Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, mock.Anything).Return(nil, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				_, err := SchoolService.GetByid(ctx, 10)
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Terdapat Data Sekolah", func() {
			BeforeEach(func() {
				data := entity.School{}
				data.Achievements = []entity.Achievement{entity.Achievement{}}
				data.Faqs = []entity.Faq{entity.Faq{}}
				data.Extracurriculars = []entity.Extracurricular{entity.Extracurricular{}}
				data.Achievements = []entity.Achievement{entity.Achievement{}}
				data.Reviews = []entity.Reviews{entity.Reviews{}}
				data.Payments = []entity.Payment{entity.Payment{Type: "one"}, entity.Payment{Type: "interval"}, entity.Payment{Type: "interval"}}
				Mock.On("GetById", mock.Anything, mock.Anything, mock.Anything).Return(&data, nil).Once()
			})
			It("Akan Mengembalikan Data Sekolah Admin", func() {
				data, err := SchoolService.GetByid(ctx, 1)
				Expect(err).Should(BeNil())
				Expect(data).ShouldNot(BeNil())
			})

		})

	})

	Context("Get All School", func() {
		When("Terdapat Kesalahn Query Database", func() {
			BeforeEach(func() {
				Mock.On("GetAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, 0, errors.New("Internal Server Error")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				_, err := SchoolService.GetAll(ctx, 20, 20, "")
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Terdapat Data Sekolah", func() {
			BeforeEach(func() {
				data := []entity.School{entity.School{User: &entity.User{}}}
				Mock.On("GetAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(data, 5, nil).Once()
			})
			It("Akan Mengembalikan Data Sekolah Admin", func() {
				data, err := SchoolService.GetAll(ctx, 1, 20, "")
				Expect(err).Should(BeNil())
				Expect(data).ShouldNot(BeNil())
			})

		})

	})

	Context("Add Review", func() {
		When("Req Body Kosong", func() {
			It("Akan Mengembalikan Error", func() {
				_, err := SchoolService.AddReview(ctx, entity.Reviews{})
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Mengandung Kata Kasar", func() {
			It("Akan Mengembalikan Error", func() {
				_, err := SchoolService.AddReview(ctx, entity.Reviews{SchoolID: 1, Review: "tai lah"})
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("AddReview", mock.Anything, mock.Anything).Return(0, errors.New("err")).Once()
			})
			It("Akan Mengembalikan Data Sekolah Admin", func() {
				id, err := SchoolService.AddReview(ctx, entity.Reviews{SchoolID: 1, Review: "mantap lah"})
				Expect(err).ShouldNot(BeNil())
				Expect(id).To(Equal(0))
			})

		})
		When("Berhasil Menambahakan Review", func() {
			BeforeEach(func() {
				Mock.On("AddReview", mock.Anything, mock.Anything).Return(1, nil).Once()
			})
			It("Akan Mengembalikan Data Sekolah Admin", func() {
				id, err := SchoolService.AddReview(ctx, entity.Reviews{SchoolID: 1, Review: "mantap lah"})
				Expect(err).Should(BeNil())
				Expect(id).To(Equal(1))
			})

		})

	})
	Context("Add Quiz", func() {
		When("Req Body Kosong", func() {
			It("Akan Mengembalikan Error", func() {
				req := []entity.ReqAddQuiz{}
				err := SchoolService.CreateQuiz(ctx, req)
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Tidak Terdapat Data Sekolah", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, mock.Anything).Return(nil, errors.New("err")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				req := []entity.ReqAddQuiz{entity.ReqAddQuiz{SchoolID: 1}}
				err := SchoolService.CreateQuiz(ctx, req)
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Kesalahan Query Database", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, mock.Anything).Return(nil, nil).Once()
				Mock.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("err")).Once()
			})
			It("Akan Mengembalikan error", func() {
				req := []entity.ReqAddQuiz{entity.ReqAddQuiz{SchoolID: 1}}
				err := SchoolService.CreateQuiz(ctx, req)
				Expect(err).ShouldNot(BeNil())
			})

		})
		When("Berhasil Menambahkan Quiz", func() {
			BeforeEach(func() {
				Mock.On("GetById", mock.Anything, mock.Anything).Return(nil, nil).Once()
				Mock.On("Update", mock.Anything, mock.Anything).Return(&entity.School{}, nil).Once()
			})
			It("Akan Mengembalikan nil error", func() {
				req := []entity.ReqAddQuiz{entity.ReqAddQuiz{SchoolID: 1}}
				err := SchoolService.CreateQuiz(ctx, req)
				Expect(err).Should(BeNil())

			})

		})

	})

	Context("Get Test Result", func() {
		When("Jika Tidak Terdapat Data Link", func() {
			BeforeEach(func() {
				Mock.On("GetByUid", mock.Anything, mock.Anything).Return(nil, errors.New("tes")).Once()
			})
			It("Akan Mengembalikan Error", func() {
				_, err := SchoolService.GetTestResult(ctx, 2)
				Expect(err).ShouldNot(BeNil())

			})

		})
		When("Terdapat Data Link", func() {
			BeforeEach(func() {
				Mock.On("GetByUid", mock.Anything, mock.Anything).Return(&entity.School{}, nil).Once()
			})
			It("Akan Mengembalikan Data Sekolah Admin", func() {
				res, err := SchoolService.GetTestResult(ctx, 2)
				Expect(err).Should(BeNil())
				Expect(len(res)).To(Equal(1))

			})

		})

	})
})
