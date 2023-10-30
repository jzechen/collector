/**
 * @Time: 2023/10/20 16:29
 * @Author: jzechen
 * @File: define.go
 * @Software: GoLand collector
 */

package dto

// NullRsp a base null define.
type NullRsp struct{}

type LoginReq struct {
	UserID   string `json:"userID" binding:"required,max=32"`
	Password string `json:"password" binding:"required,max=32"`
}

/// TODO: add common data struct here, like the protobuf files, json struct definition...
