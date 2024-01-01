package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	ratingService      string = "http://localhost:8050"
	libraryService     string = "http://localhost:8060"
	reservationService string = "http://localhost:8070"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type LibraryResponse struct {
	Library_uid string `json:"library_uid"`
	Name        string `json:"name"`
	City        string `json:"city"`
	Address     string `json:"address"`
}

type LibrariesLimited struct {
	Page          int               `json:"page"`
	PageSize      int               `json:"pageSize"`
	TotalElements int               `json:"totalElements"`
	Items         []LibraryResponse `json:"items"`
}

type BookResponse struct {
	Book_uid        string `json:"book_uid"`
	Name            string `json:"name"`
	Author          string `json:"author"`
	Genre           string `json:"genre"`
	Condition       string `json:"condition"`
	Available_count int    `json:"available_count"`
}

type BookToUserResponse struct {
	Book_uid string `json:"book_uid"`
	Name     string `json:"name"`
	Author   string `json:"author"`
	Genre    string `json:"genre"`
}

type BookLimited struct {
	Page          int            `json:"page"`
	PageSize      int            `json:"pageSize"`
	TotalElements int            `json:"totalElements"`
	Items         []BookResponse `json:"items"`
}

type RatingResponse struct {
	Stars int `json:"stars"`
}

type ReservationResponse struct {
	Reservation_uid string `json:"reservation_uid"`
	Username        string `json:"username"`
	Book_uid        string `json:"book_uid"`
	Library_uid     string `json:"library_uid"`
	Status          string `json:"status"`
	Start_date      string `json:"start_date"`
	Till_date       string `json:"till_date"`
}

type ReservationToUserResponse struct {
	Reservation_uid string             `json:"reservation_uid"`
	Status          string             `json:"status"`
	Start_date      string             `json:"start_date"`
	Till_date       string             `json:"till_date"`
	Book            BookToUserResponse `json:"book"`
	Library         LibraryResponse    `json:"library"`
}

type TakeBookResponse struct {
	Reservation_uid string             `json:"reservation_uid"`
	Status          string             `json:"status"`
	Start_date      string             `json:"start_date"`
	Till_date       string             `json:"till_date"`
	Book            BookToUserResponse `json:"book"`
	Library         LibraryResponse    `json:"library"`
	Rating          RatingResponse     `json:"rating"`
}

type CreateReservationRequest struct {
	BookUid    string `json:"bookUid"`
	LibraryUid string `json:"libraryUid"`
	TillDate   string `json:"tillDate"`
}

type ReservationAmount struct {
	Amount int `json:"amount"`
}

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) GetLibrariesByCity(c *gin.Context) {
	params := c.Request.URL.Query()
	requestURL := fmt.Sprintf("%s/api/v1/libraries/", libraryService)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	q := req.URL.Query()
	q.Add("city", c.Query("city"))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var libraries []LibraryResponse
	if json.Unmarshal(resBody, &libraries) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	fmt.Println(libraries)

	pageParam := params.Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	sizeParam := params.Get("size")
	if sizeParam == "" {
		sizeParam = "100"
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	right := page * size
	if len(libraries) < right {
		right = len(libraries)
	}

	librariesStripped := make([]LibraryResponse, 0)

	if (page-1)*size <= len(libraries) {
		librariesStripped = libraries[(page-1)*size : right]
	}

	data := LibrariesLimited{
		Page:          page,
		PageSize:      size,
		TotalElements: len(librariesStripped),
		Items:         librariesStripped,
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetBooksByLibraryUid(c *gin.Context) {
	params := c.Request.URL.Query()
	requestURL := fmt.Sprintf("%s/api/v1/libraries/%s/books/", libraryService, c.Param("uid"))

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	q := req.URL.Query()
	q.Add("showAll", c.Query("showAll"))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var books []BookResponse
	if json.Unmarshal(resBody, &books) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	fmt.Println(books)

	pageParam := params.Get("page")
	if pageParam == "" {
		pageParam = "1"
	}
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	sizeParam := params.Get("size")
	if sizeParam == "" {
		sizeParam = "100"
	}
	size, err := strconv.Atoi(sizeParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	right := page * size
	if len(books) < right {
		right = len(books)
	}

	booksStripped := make([]BookResponse, 0)

	if (page-1)*size <= len(books) {
		booksStripped = books[(page-1)*size : right]
	}

	data := BookLimited{
		Page:          page,
		PageSize:      size,
		TotalElements: len(booksStripped),
		Items:         booksStripped,
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) GetRating(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	requestURL := fmt.Sprintf("%s/api/v1/rating/", ratingService)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	req.Header.Set("X-User-Name", username)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var rating RatingResponse
	if json.Unmarshal(resBody, &rating) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, RatingResponse{
		Stars: rating.Stars,
	})
}

func (h *Handler) GetReservations(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	requestURL := fmt.Sprintf("%s/api/v1/reservations/", reservationService)

	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	req.Header.Set("X-User-Name", username)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var reservations []ReservationResponse
	if json.Unmarshal(resBody, &reservations) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	response := make([]ReservationToUserResponse, len(reservations))

	for i, reservation := range reservations {
		requestBookURL := fmt.Sprintf("%s/api/v1/books/%s/", libraryService, reservation.Book_uid)

		req, err := http.NewRequest(http.MethodGet, requestBookURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		var book BookToUserResponse
		if json.Unmarshal(resBody, &book) != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		requestLibraryURL := fmt.Sprintf("%s/api/v1/libraries/%s/", libraryService, reservation.Library_uid)

		reqLib, err := http.NewRequest(http.MethodGet, requestLibraryURL, nil)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		resLib, err := http.DefaultClient.Do(reqLib)
		if err != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		resLibBody, err := io.ReadAll(resLib.Body)
		if err != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		var library LibraryResponse
		if json.Unmarshal(resLibBody, &library) != nil {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: err.Error(),
			})
			return
		}

		response[i] = ReservationToUserResponse{
			Reservation_uid: reservation.Reservation_uid,
			Status:          reservation.Status,
			Start_date:      reservation.Start_date,
			Till_date:       reservation.Till_date,
			Book:            book,
			Library:         library,
		}

	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateReservation(c *gin.Context) {

	username := c.GetHeader("X-User-Name")

	if username == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "username must be given as X-User-Name Header",
		})
		return
	}

	var inputCreateBody CreateReservationRequest

	err := json.NewDecoder(c.Request.Body).Decode(&inputCreateBody)
	if err != nil {
		fmt.Printf("failed to decode body %s\n", err.Error())
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	//getting amount
	requestAmountURL := fmt.Sprintf("%s/api/v1/reservations/amount", reservationService)

	reqAmount, err := http.NewRequest(http.MethodGet, requestAmountURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	reqAmount.Header.Set("X-User-Name", username)

	resAmount, err := http.DefaultClient.Do(reqAmount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resBodyAmount, err := io.ReadAll(resAmount.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var reservationAmount ReservationAmount
	if json.Unmarshal(resBodyAmount, &reservationAmount) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	//getting a rating
	requestRatingURL := fmt.Sprintf("%s/api/v1/rating/", ratingService)

	reqRating, err := http.NewRequest(http.MethodGet, requestRatingURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	reqRating.Header.Set("X-User-Name", username)

	resRating, err := http.DefaultClient.Do(reqRating)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resBodyRating, err := io.ReadAll(resRating.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var rating RatingResponse
	if json.Unmarshal(resBodyRating, &rating) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if reservationAmount.Amount >= rating.Stars {
		c.JSON(http.StatusBadRequest, MessageResponse{
			Message: "user cannot take new book",
		})
		return
	}

	//create reservation
	requestCreateURL := fmt.Sprintf("%s/api/v1/reservations", reservationService)

	marshalled, err := json.Marshal(inputCreateBody)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
	}

	reqCreate, err := http.NewRequest(http.MethodPost, requestCreateURL, bytes.NewReader(marshalled))
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}
	reqCreate.Header.Set("X-User-Name", username)

	resCreate, err := http.DefaultClient.Do(reqCreate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resBodyCreate, err := io.ReadAll(resCreate.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var createReserv ReservationResponse
	if json.Unmarshal(resBodyCreate, &createReserv) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	//create response
	requestBookURL := fmt.Sprintf("%s/api/v1/books/%s/", libraryService, createReserv.Book_uid)

	reqBook, err := http.NewRequest(http.MethodGet, requestBookURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resBook, err := http.DefaultClient.Do(reqBook)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resBodyBook, err := io.ReadAll(resBook.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var book BookToUserResponse
	if json.Unmarshal(resBodyBook, &book) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	requestLibraryURL := fmt.Sprintf("%s/api/v1/libraries/%s/", libraryService, createReserv.Library_uid)

	reqLib, err := http.NewRequest(http.MethodGet, requestLibraryURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resLib, err := http.DefaultClient.Do(reqLib)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resLibBody, err := io.ReadAll(resLib.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var library LibraryResponse
	if json.Unmarshal(resLibBody, &library) != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	response := TakeBookResponse{
		Reservation_uid: createReserv.Reservation_uid,
		Status:          createReserv.Status,
		Start_date:      createReserv.Start_date,
		Till_date:       createReserv.Till_date,
		Book:            book,
		Library:         library,
		Rating:          rating,
	}

	requestUpdateCountURL := fmt.Sprintf("%s/api/v1/books/%s/count", libraryService, book.Book_uid)

	reqCount, err := http.NewRequest(http.MethodPut, requestUpdateCountURL, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	resCount, err := http.DefaultClient.Do(reqCount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Message: err.Error(),
		})
		return
	}

	if resCount.StatusCode != 200 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Message: "error while updating count",
		})
		return
	}

	c.JSON(http.StatusOK, response)

	//TO DO уменьшить available_count
}