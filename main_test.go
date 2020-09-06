package main

import (
	"os"
	"testing"
)

type Parent struct {
	FirstName string `gorm:"primaryKey"`
	LastName  string `gorm:"primaryKey"`
}

type Child struct {
	ParentFirstName string
	ParentLastName  string
	Parent          Parent `gorm:"ForeignKey: ParentFirstName,ParentLastName; References: FirstName, LastName"`
}

func TestMain(m *testing.M) {

	DB.AutoMigrate(&Parent{}, &Child{})
	DB.Create(&Child{Parent: Parent{"John", "Black"}})
	DB.Create(&Child{Parent: Parent{"Tom", "White"}})

	result := m.Run()

	DB.Debug().Exec("DELETE FROM children")
	DB.Debug().Exec("DELETE FROM parents")
	os.Exit(result)
}

func TestPreloadWithCompositeFK(t *testing.T) {

	var children []Child
	if err := DB.Preload("Parent").Find(&children).Error; err != nil {
		t.Errorf("Failed, got error: %v", err)
	}
}

func TestWhereInClauseWithCompositeFK(t *testing.T) {

	var parents []Parent
	var params = [][]string{
		[]string{"John", "Black"},
		[]string{"Tom", "White"},
	}
	if err := DB.Where("(first_name, last_name) IN ?", params).Find(&parents).Error; err != nil {
		t.Errorf("Failed, got error: %v", err)
	}
}
