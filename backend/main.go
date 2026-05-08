package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Field struct {
	ID         int      `json:"id"`
	Name       string   `json:"name"`
	Address    string   `json:"address"`
	District   string   `json:"district"`
	PriceRange string   `json:"price_range"`
	Phone      string   `json:"phone"`
	Lat        float64  `json:"lat"`
	Lng        float64  `json:"lng"`
	Facilities []string `json:"facilities"`
	Status     string   `json:"status"`
	OpenHours  string   `json:"open_hours"`
}

var fields = []Field{
	{
		ID:         1,
		Name:       "福田体育公园",
		Address:    "深圳市福田区福华路与彩田路交叉口",
		District:   "福田",
		PriceRange: "120-180元/小时",
		Phone:      "0755-82046666",
		Lat:        22.5431,
		Lng:        114.0549,
		Facilities: []string{"更衣室", "淋浴", "停车场", "灯光"},
		Status:     "available",
		OpenHours:  "08:00-22:00",
	},
	{
		ID:         2,
		Name:       "南山足球公园",
		Address:    "深圳市南山区科苑北路与朗山二路交叉口",
		District:   "南山",
		PriceRange: "100-150元/小时",
		Phone:      "0755-26552888",
		Lat:        22.5312,
		Lng:        113.9345,
		Facilities: []string{"更衣室", "淋浴", "停车场"},
		Status:     "busy",
		OpenHours:  "07:00-23:00",
	},
	{
		ID:         3,
		Name:       "罗湖体育中心",
		Address:    "深圳市罗湖区经二路与经三路交叉口",
		District:   "罗湖",
		PriceRange: "80-130元/小时",
		Phone:      "0755-25529988",
		Lat:        22.5513,
		Lng:        114.1297,
		Facilities: []string{"更衣室", "淋浴", "灯光"},
		Status:     "available",
		OpenHours:  "09:00-21:00",
	},
	{
		ID:         4,
		Name:       "宝安体育场",
		Address:    "深圳市宝安区裕安一路与新安一路交叉口",
		District:   "宝安",
		PriceRange: "90-140元/小时",
		Phone:      "0755-27888888",
		Lat:        22.5438,
		Lng:        113.8833,
		Facilities: []string{"更衣室", "淋浴", "停车场", "灯光", "wifi"},
		Status:     "full",
		OpenHours:  "08:00-22:00",
	},
	{
		ID:         5,
		Name:       "龙岗大运中心",
		Address:    "深圳市龙岗区龙翔大道2999号",
		District:   "龙岗",
		PriceRange: "110-160元/小时",
		Phone:      "0755-28989988",
		Lat:        22.6788,
		Lng:        114.2466,
		Facilities: []string{"更衣室", "淋浴", "停车场", "灯光", "wifi"},
		Status:     "available",
		OpenHours:  "08:00-22:00",
	},
	{
		ID:         6,
		Name:       "龙华文化体育中心",
		Address:    "深圳市龙华区观澜街道高尔夫大道",
		District:   "龙华",
		PriceRange: "70-120元/小时",
		Phone:      "0755-28108888",
		Lat:        22.6812,
		Lng:        114.0566,
		Facilities: []string{"更衣室", "淋浴", "灯光"},
		Status:     "busy",
		OpenHours:  "09:00-21:00",
	},
	{
		ID:         7,
		Name:       "盐田体育中心",
		Address:    "深圳市盐田区深盐路2100号",
		District:   "盐田",
		PriceRange: "60-100元/小时",
		Phone:      "0755-25259988",
		Lat:        22.5573,
		Lng:        114.2378,
		Facilities: []string{"更衣室", "停车场"},
		Status:     "available",
		OpenHours:  "08:00-20:00",
	},
	{
		ID:         8,
		Name:       "光明足球基地",
		Address:    "深圳市光明区观光路与科林路交叉口",
		District:   "光明",
		PriceRange: "50-90元/小时",
		Phone:      "0755-27418888",
		Lat:        22.7812,
		Lng:        113.9366,
		Facilities: []string{"更衣室", "淋浴", "停车场", "灯光"},
		Status:     "available",
		OpenHours:  "08:00-21:00",
	},
	{
		ID:         9,
		Name:       "坪山足球公园",
		Address:    "深圳市坪山区体育五路与体育一路交叉口",
		District:   "坪山",
		PriceRange: "50-80元/小时",
		Phone:      "0755-89628888",
		Lat:        22.6878,
		Lng:        114.3856,
		Facilities: []string{"停车场", "灯光"},
		Status:     "available",
		OpenHours:  "09:00-20:00",
	},
	{
		ID:         10,
		Name:       "大鹏足球训练基地",
		Address:    "深圳市大鹏新区葵鹏路88号",
		District:   "大鹏",
		PriceRange: "60-100元/小时",
		Phone:      "0755-84308888",
		Lat:        22.5888,
		Lng:        114.4722,
		Facilities: []string{"更衣室", "淋浴", "停车场", "灯光"},
		Status:     "available",
		OpenHours:  "08:00-21:00",
	},
	{
		ID:         11,
		Name:       "福田上梅林足球场",
		Address:    "深圳市福田区梅林路与梅中路交叉口",
		District:   "福田",
		PriceRange: "100-140元/小时",
		Phone:      "0755-83128888",
		Lat:        22.5612,
		Lng:        114.0401,
		Facilities: []string{"更衣室", "淋浴", "灯光"},
		Status:     "busy",
		OpenHours:  "08:00-22:00",
	},
	{
		ID:         12,
		Name:       "南山科技园足球场",
		Address:    "深圳市南山区科技中四路与高新中四道交叉口",
		District:   "南山",
		PriceRange: "130-180元/小时",
		Phone:      "0755-26778888",
		Lat:        22.5256,
		Lng:        113.9266,
		Facilities: []string{"更衣室", "淋浴", "停车场", "灯光", "wifi"},
		Status:     "available",
		OpenHours:  "07:00-23:00",
	},
}

func main() {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.GET("/api/fields", func(c *gin.Context) {
		district := c.Query("district")
		status := c.Query("status")

		result := fields
		if district != "" {
			var filtered []Field
			for _, f := range fields {
				if f.District == district {
					filtered = append(filtered, f)
				}
			}
			result = filtered
		}
		if status != "" {
			var filtered []Field
			for _, f := range result {
				if f.Status == status {
					filtered = append(filtered, f)
				}
			}
			result = filtered
		}
		c.JSON(http.StatusOK, gin.H{"fields": result})
	})

	r.GET("/api/fields/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		for _, f := range fields {
			if f.ID == id {
				c.JSON(http.StatusOK, f)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "field not found"})
	})

	r.GET("/api/districts", func(c *gin.Context) {
		districts := []string{"福田", "南山", "罗湖", "宝安", "龙岗", "龙华", "盐田", "光明", "坪山", "大鹏"}
		c.JSON(http.StatusOK, gin.H{"districts": districts})
	})

	r.GET("/api/statuses", func(c *gin.Context) {
		statuses := []map[string]string{
			{"value": "available", "label": "空闲"},
			{"value": "busy", "label": "忙碌"},
			{"value": "full", "label": "已满"},
		}
		c.JSON(http.StatusOK, gin.H{"statuses": statuses})
	})

	r.Run(":8080")
}
