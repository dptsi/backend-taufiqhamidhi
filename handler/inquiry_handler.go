package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"math/big"
	"myits-gate-api/entity"
	"myits-gate-api/repository"
	"net/http"
	"strconv"
	"time"
	"unicode/utf8"
)

type InquiryRequest struct {
	GateId     string `json:"gate_id" binding:"required"`
	CardId     string `json:"card_id" binding:"required"`
	AccessType string `json:"access_type" binding:"required"`
}

type InquiryResponse struct {
	ResponseCode    string `json:"response_code"`
	ResponseMessage string `json:"response_message"`
	Allowed         bool   `json:"allowed"`
}

type InquiryHandler struct {
	gateAccessRepo repository.GateAccessRepository
	gateLogRepo    repository.GateLogRepository
}

func NewInquiryHandler(
	gateAccessRepo repository.GateAccessRepository,
	gateLogRepo repository.GateLogRepository,
) *InquiryHandler {

	return &InquiryHandler{
		gateAccessRepo: gateAccessRepo,
		gateLogRepo:    gateLogRepo,
	}
}

func (i *InquiryHandler) PostInquiry(c *gin.Context) {
	i.InquiryProcess(c, false)
}

func (i *InquiryHandler) PostInquiryComplete(c *gin.Context) {
	i.InquiryProcess(c, true)
}

func (i *InquiryHandler) InquiryProcess(c *gin.Context, isComplete bool) {

	var request InquiryRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, NewErrorResponse(invalidRequestCode, invalidRequestMessage))
		return
	}

	if !i.isRequestValid(c, request) {
		return
	}

	cardIdDecimal := hexToDecimal(request.CardId)
	//cardIdDecimalInt, err := strconv.Atoi(cardIdDecimal)
	cardIdDecimal = fixDecimalValue(cardIdDecimal)
	//log.Println(cardIdDecimal)
	hasRows, err := i.gateAccessRepo.GetGate(request.GateId)
	if err != nil {
		log.Fatal(err)
	}

	if hasRows != nil {
		log.Println(hasRows.IsGateUmum)
		if hasRows.IsGateUmum {
			access, err := i.gateAccessRepo.FindByUmum(cardIdDecimal)
			log.Println(access)

			if err != nil {
				c.JSON(http.StatusInternalServerError, NewErrorResponse(dataLinkErrorCode, dataLinkErrorMessage))
				return
			}

			//var allowed = false
			var allowed = 0

			if access == nil {
				if isComplete {
					response := InquiryResponse{
						ResponseCode:    notAllowed,
						ResponseMessage: notAllowedMessage,
						Allowed:         false,
					}
					c.JSON(http.StatusOK, response)
					return
				}

				c.JSON(http.StatusOK, allowed)
				return
			}

			log := entity.GateLog{
				LogId:      uuid.New(),
				CardId:     cardIdDecimal,
				GateId:     request.GateId,
				AccessDate: time.Now(),
				AccessType: request.AccessType,
			}

			err = i.gateLogRepo.Save(log)

			if err != nil {
				c.JSON(http.StatusInternalServerError, NewErrorResponse(insertFailedCode, insertFailedMessage))
				return
			}

			//allowed = true
			allowed = 1

			if isComplete {
				response := InquiryResponse{
					ResponseCode:    successCode,
					ResponseMessage: successMessage,
					Allowed:         true,
				}
				c.JSON(http.StatusOK, response)
				return
			}

			c.JSON(http.StatusOK, allowed)
			return
		}

		access, err := i.gateAccessRepo.FindBy(request.GateId, cardIdDecimal)
		log.Println(access)

		if err != nil {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(dataLinkErrorCode, dataLinkErrorMessage))
			return
		}

		//var allowed = false
		var allowed = 0

		if access == nil {
			if isComplete {
				response := InquiryResponse{
					ResponseCode:    notAllowed,
					ResponseMessage: notAllowedMessage,
					Allowed:         false,
				}
				c.JSON(http.StatusOK, response)
				return
			}

			c.JSON(http.StatusOK, allowed)
			return
		}

		log := entity.GateLog{
			LogId:      uuid.New(),
			CardId:     cardIdDecimal,
			GateId:     request.GateId,
			AccessDate: time.Now(),
			AccessType: request.AccessType,
		}

		err = i.gateLogRepo.Save(log)

		if err != nil {
			c.JSON(http.StatusInternalServerError, NewErrorResponse(insertFailedCode, insertFailedMessage))
			return
		}

		//allowed = true
		allowed = 1

		if isComplete {
			response := InquiryResponse{
				ResponseCode:    successCode,
				ResponseMessage: successMessage,
				Allowed:         true,
			}
			c.JSON(http.StatusOK, response)
			return
		}

		c.JSON(http.StatusOK, allowed)
		return

	} else {
		c.JSON(http.StatusBadRequest, NewErrorResponse(gateNotFoundCode, gateNotFoundMessage))
		return
	}

}

func (i *InquiryHandler) isRequestValid(c *gin.Context, r InquiryRequest) bool {

	if r.AccessType != "I" && r.AccessType != "O" {
		c.JSON(http.StatusBadRequest, NewErrorResponse(invalidAccessTypeCode, invalidAccessTypeMessage))
		return false
	}

	return true
}

func hexToDecimal(hex string) string {
	decimal := big.NewInt(0)
	_, success := decimal.SetString(hex, 16)
	if !success {
		return ""
	}
	return decimal.String()
}

func cutString(s string, maxLength int) string {
	// Check if the string is already within the desired length
	if utf8.RuneCountInString(s) <= maxLength {
		return s
	}

	// Cut the string at the desired length
	runes := []rune(s)
	return string(runes[:maxLength])
}

func fixDecimalValue(cardIdDecimal string) string {
	cardIdDecimalInt, err := strconv.Atoi(cardIdDecimal)

	decimalLen := len([]rune(cardIdDecimal))
	if decimalLen < 10 {
		cardIdDecimal = fmt.Sprintf("%0*d", 10, cardIdDecimalInt)
	}

	if err != nil {
		return "error"
	}

	return cardIdDecimal
}
