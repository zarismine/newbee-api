package request

type SaveOrder struct {
	AddressId   int   `json:"addressId" form:"addressId"`
	CartItemIds []int `json:"cartItemIds" form:"cartItemIds"`
}