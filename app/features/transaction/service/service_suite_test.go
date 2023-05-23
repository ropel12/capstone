package service_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/education-hub/BE/app/entities"
	mocks "github.com/education-hub/BE/app/features/school/mocks/repository"
	mocksuu "github.com/education-hub/BE/app/features/transaction/mocks/repository"
	transaction "github.com/education-hub/BE/app/features/transaction/service"
	mocksu "github.com/education-hub/BE/app/features/user/mocks/repository"
	"github.com/education-hub/BE/config"
	dependcy "github.com/education-hub/BE/config/dependency"
	"github.com/education-hub/BE/config/dependency/container"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("transaction", func() {
	var Mock *mocks.SchoolRepo
	var Mocks *mocksu.UserRepo
	var Mockss *mocksuu.TransactionRepo
	var TransactionService transaction.TransactionService
	var Depend dependcy.Depend
	var ctx context.Context

	BeforeEach(func() {
		Depend.Db = config.GetConnectionTes()
		log := logrus.New()
		Depend.Log = log
		ctx = context.Background()
		Mock = mocks.NewSchoolRepo(GinkgoT())
		Mocks = mocksu.NewUserRepo(GinkgoT())
		Mockss = mocksuu.NewTransactionRepo(GinkgoT())
		TransactionService = transaction.NewTransactionService(Mockss, Depend, Mocks, Mock)
		Depend.Config = &config.Config{GmapsKey: os.Getenv("GMAPS")}
		Depend.PromErr = make(map[string]string, 1)
		Depend.Mds = container.NewMidtrans(&config.Config{Midtrans: config.MidtransConfig{ServerKey: "SB-Mid-server-TvgWB_Y9s81-rbMBH7zZ8BHW", ClientKey: "SB-Mid-client-nKsqvar5cn60u2Lv", Env: 1, ExpiryDuration: 1}})
	})

	Context("Create Transaction", func() {
		When("Request Body kosong", func() {
			It("Akan Mengembalikan Erorr", func() {
				_, err := TransactionService.CreateTransaction(ctx, entities.ReqCheckout{}, 1)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Request Type Tidak ada Di list", func() {
			It("Akan Mengembalikan Erorr", func() {
				_, err := TransactionService.CreateTransaction(ctx, entities.ReqCheckout{SchoolID: 1, Type: "regisss", PaymentMethod: "bca"}, 1)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Payment Method Tidak ada di daftar list", func() {
			It("Akan Mengembalikan Erorr", func() {
				_, err := TransactionService.CreateTransaction(ctx, entities.ReqCheckout{SchoolID: 1, Type: "registration", PaymentMethod: "bcaaaa"}, 1)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Terjadi kesalah query database pada saat mengambil data payment school", func() {
			BeforeEach(func() {
				Mockss.On("GetSchoolPayment", mock.Anything, mock.Anything).Return(nil, errors.New("Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				_, err := TransactionService.CreateTransaction(ctx, entities.ReqCheckout{SchoolID: 1, Type: "herregistration", PaymentMethod: "bca"}, 1)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("Terjadi kesalah query database pada memasukan data ke transaction", func() {
			BeforeEach(func() {
				data := []entities.Payment{entities.Payment{Description: "Tool", Price: 1000}}
				Mockss.On("GetSchoolPayment", mock.Anything, mock.Anything).Return(&entities.School{Payments: data}, nil).Once()
				Mockss.On("CreateTranscation", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Error")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				_, err := TransactionService.CreateTransaction(ctx, entities.ReqCheckout{SchoolID: 1, Type: "herregistration", PaymentMethod: "bca"}, 1)
				Expect(err).ShouldNot(BeNil())
			})
		})

		When("berhasil membuat transaction", func() {
			BeforeEach(func() {
				data := []entities.Payment{entities.Payment{Description: "Tool", Price: 1000}}
				Mockss.On("GetSchoolPayment", mock.Anything, mock.Anything).Return(&entities.School{Payments: data}, nil).Once()
				Mockss.On("CreateTranscation", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Data Transaksi", func() {
				res, err := TransactionService.CreateTransaction(ctx, entities.ReqCheckout{SchoolID: 1, Type: "herregistration", PaymentMethod: "indomaret"}, 1)
				Expect(err).Should(BeNil())
				Expect(res).ShouldNot(BeNil())
			})
		})
	})

	Context("GetAllTrasactionCart", func() {
		When("Data Cart Tidak ada", func() {
			BeforeEach(func() {
				Mockss.On("GetAllCartByuid", mock.Anything, mock.Anything).Return(nil, errors.New("Data Not Found")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				_, err := TransactionService.GetAllTrasactionCart(ctx, 1)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Terdapat Data", func() {
			BeforeEach(func() {
				data := []entities.Carts{entities.Carts{School: entities.School{}}}
				Mockss.On("GetAllCartByuid", mock.Anything, mock.Anything).Return(data, nil).Once()
			})
			It("Akan Mengembalikan Data Cart", func() {
				_, err := TransactionService.GetAllTrasactionCart(ctx, 1)
				Expect(err).Should(BeNil())
			})
		})

	})

	Context("GetDetailTransaction", func() {
		When("Terdapat Data Di Trasaction", func() {
			BeforeEach(func() {
				Mockss.On("GetTransaction", mock.Anything, mock.Anything, mock.Anything).Return(&entities.Transaction{}, nil).Once()
			})
			It("Akan Mengembalikan Data Trasaction", func() {
				_, err := TransactionService.GetDetailTransaction(ctx, 99, 1)
				Expect(err).Should(BeNil())
			})
		})
		When("Tidak Terdapat Data Di Cart Maupun Trx", func() {
			BeforeEach(func() {
				Mockss.On("GetTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Data Not Found")).Once()
				Mockss.On("GetCart", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Data Not Found")).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				_, err := TransactionService.GetDetailTransaction(ctx, 99, 1)
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Jika Terdapat Data Cart Dan tipenya Registration", func() {
			BeforeEach(func() {
				Mockss.On("GetTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Data Not Found")).Once()
				Mockss.On("GetCart", mock.Anything, mock.Anything, mock.Anything).Return(&entities.Carts{Type: "registration"}, nil).Once()
			})
			It("Akan Mengembalikan Data Cart Registrasi", func() {
				_, err := TransactionService.GetDetailTransaction(ctx, 2, 1)
				Expect(err).Should(BeNil())
			})
		})
		When("Jika Terdapat Data Cart Dan tipenya Her Registration", func() {
			BeforeEach(func() {
				Mockss.On("GetTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Data Not Found")).Once()
				Mockss.On("GetCart", mock.Anything, mock.Anything, mock.Anything).Return(&entities.Carts{Type: "herregistration", School: entities.School{Payments: []entities.Payment{entities.Payment{Type: "interval"}}}}, nil).Once()
			})
			It("Akan Mengembalikan Data Cart Her Registrasi", func() {
				_, err := TransactionService.GetDetailTransaction(ctx, 2, 2)
				Expect(err).Should(BeNil())
			})
		})
		When("Jika Terdapat Data Cart Dan tipenya Her Registration", func() {
			BeforeEach(func() {
				Mockss.On("GetTransaction", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("Data Not Found")).Once()
				Mockss.On("GetCart", mock.Anything, mock.Anything, mock.Anything).Return(&entities.Carts{Type: "herregistration", School: entities.School{Payments: []entities.Payment{entities.Payment{Type: "one"}}}}, nil).Once()
			})
			It("Akan Mengembalikan Data Cart Her Registrasi", func() {
				_, err := TransactionService.GetDetailTransaction(ctx, 2, 2)
				Expect(err).Should(BeNil())
			})
		})

	})

	Context("UpdateStatus", func() {
		When("Terdapat kesalahn qury db", func() {
			BeforeEach(func() {
				Mockss.On("UpdateStatus", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("Internal Server Error")).Once()
			})
			It("Akan MengembalikanEroor", func() {
				err := TransactionService.UpdateStatus(ctx, "paid", "INV-001222")
				Expect(err).ShouldNot(BeNil())
			})
		})
		When("Berhasil Mengupdate Transaction", func() {
			BeforeEach(func() {
				Mockss.On("UpdateStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
			})
			It("Akan Mengembalikan Erorr", func() {
				err := TransactionService.UpdateStatus(ctx, "paid", "INV-0013123")
				Expect(err).Should(BeNil())
			})
		})
	})

})
