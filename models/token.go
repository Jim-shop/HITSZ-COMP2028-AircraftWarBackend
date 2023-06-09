/**
 * Copyright (c) [2023] [Jim-shop]
 * [AircraftWarBackend] is licensed under Mulan PubL v2.
 * You can use this software according to the terms and conditions of the Mulan PubL v2.
 * You may obtain a copy of Mulan PubL v2 at:
 *          http://license.coscl.org.cn/MulanPubL-2.0
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PubL v2 for more details.
 */

package models

import (
	"crypto/rand"
	"crypto/sha512"
	"fmt"
	"imshit/aircraftwar/db"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Token struct {
	Token string `form:"token" json:"token" uri:"token" xml:"token" binding:"required"`
}

func NewToken(user *User, c *gin.Context) (*Token, error) {
	randInt, err := rand.Int(rand.Reader, big.NewInt((1<<63)-1))
	if err != nil {
		log.Printf("Rand int generate error: %v\n", err)
		return nil, err
	}
	info := fmt.Sprintf("%s%s%s%s",
		c.ClientIP(),
		user.Name,
		time.Now(),
		randInt,
	)
	token := fmt.Sprintf("%x",
		sha512.Sum512([]byte(info)),
	)
	db.GetRedis().Set(
		token,
		user.ID,
		viper.GetDuration("token.timeout"),
	)
	return &Token{token}, nil
}

func ValidateToken(token *Token) bool {
	redis := db.GetRedis()
	// 查验Token是否有效
	_, err := redis.Get(token.Token).Result()
	if err != nil {
		return false
	}
	// 延期Token
	success, err := redis.Expire(token.Token, viper.GetDuration("token.timeout")).Result()
	if err != nil {
		log.Printf("Token expire error: %v\n", err)
		return false
	}
	return success
}

func GetUserIDByToken(token *Token) (int, error) {
	_id, err := db.GetRedis().Get(token.Token).Result()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(_id)
}
