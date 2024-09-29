package main

import (
	"encoding/csv"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"os"
	"strconv"
	"testing"
)

type recipient struct {
	name        string
	phoneNumber string
}

func writeRecipientsToCSV(filename string, contacts []recipient) error {

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, contact := range contacts {
		record := []string{contact.name, contact.phoneNumber}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write record: %v", err)
		}
	}

	return nil
}

// 生成联系人信息
func generateContacts() []recipient {
	contacts := make([]recipient, 10000)
	basePhone := 13961467500

	for i := 0; i < 10000; i++ {
		contacts[i] = recipient{
			name:        "user" + strconv.Itoa(i+1),
			phoneNumber: strconv.FormatInt(int64(basePhone+i), 10), // 手机号依次递增
		}
	}

	return contacts
}

func TestName(t *testing.T) {
	contacts := generateContacts()
	if err := writeRecipientsToCSV("recipients.csv", contacts); err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println("Contacts successfully written to contacts.csv")
	}
}

func TestLogo(t *testing.T) {
	figure.NewFigure("Msg Producer Start", "", true).Print()
}
