//go:build integration

package tests

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/Eiji224/go_ozon/hw2/cart/internal/dto"
	"github.com/Eiji224/go_ozon/hw2/cart/internal/models"
	"github.com/Eiji224/go_ozon/hw2/cart/internal/router"
	"github.com/stretchr/testify/require"
)

var url = "http://localhost:8080"

func TestGetCart_WrongUserId(t *testing.T) {
	tests := []struct {
		name         string
		userId       any
		expectedCode int
		expectedBody map[string]string
	}{
		{"cart not found", 99999, http.StatusNotFound, map[string]string{"error": router.ErrCartNotFound}},
		{"userID is a string", "asd", http.StatusBadRequest, map[string]string{"error": router.ErrUserIdLessThanZero}},
		{"cart not found", -1, http.StatusBadRequest, map[string]string{"error": router.ErrUserIdLessThanZero}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filledUrl := url + "/user/"
			strUserId, ok := tt.userId.(string)
			if ok {
				filledUrl += strUserId + "/cart"
			} else {
				intUserId := tt.userId.(int)
				filledUrl += strconv.Itoa(intUserId) + "/cart"
			}

			gotCode, gotBody, err := sendAndGetContent(t, filledUrl, http.MethodGet, []byte{})
			if err != nil {
				require.Failf(t, "", "got error: %v", err)
			}
			parsedBody := map[string]string{}
			if err := json.Unmarshal([]byte(gotBody), &parsedBody); err != nil {
				require.Failf(t, "", "got error: %v", err)
			}

			require.Equal(t, tt.expectedCode, gotCode)
			require.Equal(t, tt.expectedBody, parsedBody)
		})
	}
}

func TestAddProduct_Success(t *testing.T) {
	req := map[string]int{"count": 1}
	body, _ := json.Marshal(req)
	expectedResp := &dto.GetCartResponse{
		Items: []*models.Product{{
			Sku:   models.Sku(1),
			Name:  "Теория нравственных чувств | Смит Адам",
			Count: 1,
			Price: 3379,
		}},
		TotalPrice: 3379,
	}

	gotCode, _, gotError := sendAndGetContent(t, url+"/user/1/cart/1", http.MethodPost, body)
	_, result, err := sendAndGetContent(t, url+"/user/1/cart", http.MethodGet, []byte{})
	if err != nil {
		require.Fail(t, "failed get result")
	}
	gotResp := &dto.GetCartResponse{}
	err = json.Unmarshal([]byte(result), gotResp)
	if err != nil {
		require.Fail(t, "failed parse result")
	}

	require.Equal(t, http.StatusOK, gotCode)
	require.Equal(t, nil, gotError)
	require.Equal(t, expectedResp, gotResp)
}

func TestAddMultipleProducts_Success(t *testing.T) {
	expectedResp := &dto.GetCartResponse{
		Items: []*models.Product{
			{
				Sku:   models.Sku(1),
				Name:  "Теория нравственных чувств | Смит Адам",
				Count: 1,
				Price: 3379,
			},
			{
				Sku:   models.Sku(2),
				Name:  "Богатство народов | Смит Адам",
				Count: 2,
				Price: 4250,
			},
			{
				Sku:   models.Sku(3),
				Name:  "Мышление, быстрое и медленное | Канеман Даниэль",
				Count: 3,
				Price: 2890,
			},
		},
		TotalPrice: 21329,
	}

	for i := 1; i <= 3; i++ {
		req := map[string]int{"count": i}
		body, _ := json.Marshal(req)

		gotCode, _, gotError := sendAndGetContent(t, url+"/user/1/cart/"+strconv.Itoa(i), http.MethodPost, body)
		require.Equal(t, http.StatusOK, gotCode)
		require.Equal(t, nil, gotError)
	}
	_, result, err := sendAndGetContent(t, url+"/user/1/cart", http.MethodGet, []byte{})
	if err != nil {
		require.Fail(t, "failed get result")
	}
	gotResp := &dto.GetCartResponse{}
	err = json.Unmarshal([]byte(result), gotResp)
	if err != nil {
		require.Fail(t, "failed parse result")
	}

	require.Equal(t, expectedResp, gotResp)
}

func TestAddProduct_WrongParams(t *testing.T) {
	tests := []struct {
		name         string
		userId       int
		skuId        int
		count        any
		expectedCode int
		expectedBody map[string]string
	}{
		{"userId < 0", -1, 1, 1, http.StatusBadRequest, map[string]string{"error": router.ErrUserIdLessThanZero}},
		{"sku doesn't exists", 1, 99999999, 1, http.StatusPreconditionFailed, map[string]string{"error": router.ErrProductNotFound}},
		{"sku < 0", 1, -1, 1, http.StatusBadRequest, map[string]string{"error": router.ErrSkuLessThanZero}},
		{"count is a string", 1, 1, "1", http.StatusBadRequest, map[string]string{"error": router.ErrCountDigitGreaterThanZero}},
		{"count < 0", 1, 1, -1, http.StatusBadRequest, map[string]string{"error": router.ErrCountDigitGreaterThanZero}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := map[string]any{"count": tt.count}
			body, _ := json.Marshal(req)

			filledUrl := url + "/user/" + strconv.Itoa(tt.userId) + "/cart/" + strconv.Itoa(tt.skuId)
			gotCode, gotResp, err := sendAndGetContent(t, filledUrl, http.MethodPost, body)
			if err != nil {
				require.Failf(t, "", "got error: %v", err)
			}
			parsedResp := map[string]string{}
			if err := json.Unmarshal([]byte(gotResp), &parsedResp); err != nil {
				require.Failf(t, "", "got error: %v", err)
			}
			require.Equal(t, tt.expectedCode, gotCode)
			require.Equal(t, tt.expectedBody, parsedResp)
		})
	}
}
