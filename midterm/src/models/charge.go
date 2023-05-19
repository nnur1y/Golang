package models

type Charge struct {
	Amount         string `json:"amount"`
	RecipeId       string `json:"r_id"`
	Token          string `json:"token"`
	CardholderName string `json:"cardholderName"`
}

func (c *Charge) TableName() string {
	return "charge"

}
