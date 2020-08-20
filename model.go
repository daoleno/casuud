package main

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name  string `json:"name"`
	Cards []Card `json:"cards"`
}

type Card struct {
	gorm.Model
	Front string `json:"front"`
	Back  string `json:"back"`
}
