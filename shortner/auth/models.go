package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type forgotPasswordOtp struct {
	Otp     string    `json:"otp"`
	Created time.Time `json:"created"`
}
type User struct {
	Id                 primitive.ObjectID   `bson:"_id" json:"_id"`
	Name               string               `json:"name"`
	PetName            string               `json:"pet_name"`
	SentYear2021Review bool                 `json:"sent_year20_21_review"`
	Phone              string               `json:"phone"`
	Email              string               `json:"email"`
	Password           string               `json:"password"`
	Gender             string               `json:"gender"`
	OnesignalIds       string               `json:"onesignal_ids`
	ForgotPasswordOtp  forgotPasswordOtp    `json:"forgot_password_otp"`
	IsVerified         bool                 `json:"is_verified"`
	VerificationCode   string               `json:"verification_code"`
	ProfilePicture     string               `json:"profile_picture"`
	FbId               string               `json:"fb_id"`
	RandomKey          string               `json:"random_key"`
	PublicNoteCount    int                  `json:"public_note_count"`
	Favorites          []primitive.ObjectID `json:"favorites"`
	UserTags           []string             `json:"user_tags"`
	Country            string               `json:"country"`
	Blocked            bool                 `json:"blocked"`
	Flagged            bool                 `json:"flagged"`
	UserIp             string               `json:"user_ip"`
	TzOffset           int                  `json:"tz_offset"`
	CreatedAt          time.Time            `json:"CreatedAt"`
	UpdatedAt          time.Time            `json:"UpdatedAt"`
}

type RegisterInput struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	Phone          string `json:"phone"`
	Name           string `json:"name"`
	OTP            string `json:"otp"`
	PetName        string `json:"pet_name"`
	NotificationId string `json:"notification_id"`
	UerIp          string `json:"user_ip"`
	BAvatar        struct {
		ASet struct {
			V string `json:"v"`
			N string `json:"n"`
			C string `json:"c"`
		} `json:"a_set"`
	} `json:"b_avatar"`
	Blocked bool `json:"blocked"`
}
