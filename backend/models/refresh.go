package models

import (
	"github.com/beego/beego/v2/client/orm"
)

type RefreshToken struct {
	Id    int    `orm:"pk;auto"`
	Token string `orm:"unique"`
	User  *User  `orm:"rel(fk);on_delete(cascade)"`
}

func init() {
	orm.RegisterModel(new(RefreshToken))
}

func (r *RefreshToken) TableName() string {
	return "refresh_tokens"
}

func SaveRefreshToken(token string, userID int) error {
	o := orm.NewOrm()

	user, err := GetUserById(userID)
	if err != nil {
		return err
	}

	refreshToken := RefreshToken{
		Token: token,
		User:  user,
	}
	_, err = o.Insert(&refreshToken)
	return err
}

func ValidateRefreshTokenInDB(token string, userID int) (bool, error) {
	o := orm.NewOrm()
	var refreshToken RefreshToken
	err := o.QueryTable("refresh_tokens").Filter("token", token).Filter("user_id", userID).One(&refreshToken)

	if err == orm.ErrNoRows {
		return false, nil // Token not found → Invalid
	} else if err != nil {
		return false, err // DB error
	}

	return true, nil // Token is valid
}

func DeleteAllRefreshTokensForUser(userID int) error {
	o := orm.NewOrm()

	// Delete all refresh tokens associated with the given userId
	_, err := o.QueryTable("refresh_tokens").Filter("user_id", userID).Delete()
	return err
}
