package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
)

var handler *Handler

func Test(t *testing.T) {
	links := make([]string, 0, 1000)
	baseUrl := "http:/testing.com"
	for i := 0; i < 100; i++ {
		url := fmt.Sprintf("%v/%v", baseUrl, uuid.NewString())
		links = append(links, url)
	}

	for _, l := range links {
		t.Log(l)
	}
}

// 18s
func TestSingleInsert(t *testing.T) {
	baseUrl := "http:/testing.com"
	for i := 0; i < 10000; i++ {
		url := fmt.Sprintf("%v/%v", baseUrl, uuid.NewString())
		err := handler.SingleInsertperRow(url)
		if err != nil {
			t.Fatalf("single insert failed: %v", err)
		}

	}
}

// 16.998s, 21.788s (N = 10,000)
func TestPreparedStatementsExecute(t *testing.T) {
	n := 10000
	links := make([]string, 0, n)
	baseUrl := "http:/testing.com"
	for i := 0; i < n; i++ {
		url := fmt.Sprintf("%v/%v", baseUrl, uuid.NewString())
		links = append(links, url)
	}
	if err := handler.PreparedStatementsExecute(links); err != nil {
		t.Fatal(err.Error())
	}
}

//	0.855s, 0.769s, 0.791s (N = 1000)
// 1.159s, 1.152s, 1.242s (N = 10,000)
func TestBatchInsert(t *testing.T) {
	N := 1000
	links := make([]string, 0, N)
	baseUrl := "http:/testing.com"
	for i := 0; i < N; i++ {
		url := fmt.Sprintf("%v/%v", baseUrl, uuid.NewString())
		links = append(links, url)
	}
	if err := handler.BatchInsert(links); err != nil {
		t.Fatal(err.Error())
	}
}

//	0.937s, 0.889s, 0.921s (N = 1000)
// 2.741s, 2.625s, 2.626s  (N = 1,0000)
func TestTransactionInserts(t *testing.T) {
	N := 1000
	links := make([]string, 0, N)
	baseUrl := "http:/testing.com"
	for i := 0; i < N; i++ {
		url := fmt.Sprintf("%v/%v", baseUrl, uuid.NewString())
		links = append(links, url)
	}
	if err := handler.TransactionInserts(links); err != nil {
		t.Fatal(err.Error())
	}
}

// 0.840s, 0.785s, 0.768s (n = 1000)
// 1.558s, 1.230s, 1.221s (N = 1,0000)
func TestTransactionBatchInserts(t *testing.T) {
	N := 1000
	links := make([]string, 0, N)
	baseUrl := "http:/testing.com"
	for i := 0; i < N; i++ {
		url := fmt.Sprintf("%v/%v", baseUrl, uuid.NewString())
		links = append(links, url)
	}
	if err := handler.TransactionBatchInserts(links); err != nil {
		t.Fatal(err.Error())
	}
}

func TestMain(m *testing.M) {
	db, err := connectDb()
	if err != nil {
		panic(err)
	}

	handler = &Handler{db: db}

	code := m.Run()

	db.Close()
	os.Exit(code)
}
