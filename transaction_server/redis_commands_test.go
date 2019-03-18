package main

import (
	"testing"
)

/*
/ These tests assume that the functions
/ dialRedis() and flushRedis() are
/ working properly. If they are not,
/ no tests in this test suite will
/ pass.
*/ 
func TestGetBalance(t *testing.T) {
	client := dialRedis()
	flushRedis(client)

	username := "user"
	balance := 1200.00

	client.Cmd("HSET", username, "Balance", balance)
	result := getBalance(client, username)
	if result != balance {
		t.Errorf("getBalance was incorrect, got: %f, want: %f.", result, balance)
	}
}

func TestAddBalance(t *testing.T) {
	client := dialRedis()
	flushRedis(client)

	username := "user"
	add := 300.00

	addBalance(client, username, add)
	result := getBalance(client, username)
	if result != add {
		t.Errorf("addBalance was incorrect, got: %f, want %f.", result, add)
	}
}

func TestStockOwned(t *testing.T) {
	client := dialRedis()
	flushRedis(client)

	username := "user"
	stock := "ABC"

	result := stockOwned(client, username, stock)
	if result != 0 {
		t.Errorf("stockOwned was incorrect, got %d, want %d.", result, 0)
	}

	amount := 31

	client.Cmd("HSET", username, stock, amount)
	result2 := stockOwned(client, username, stock)
	if result2 != amount {
		t.Errorf("stockOwned was incorrect, got %d, want %d.", result, amount)
	}
}

func TestsExists(t *testing.T) {
	client := dialRedis()
	flushRedis(client)

	username := "user"

	result := exists(client, username)
	if result != false {
		t.Errorf("exists was incorrect, got %t, want %t.", result, false)
	}

	client.Cmd("HMSET", username, "Balance", 0.00)

	result = exists(client, username)
	if result != true {
		t.Errorf("exists was incorrect, got %t, want %t.", result, true)
	}
}

func TestQExists(t *testing.T) {
	client := dialRedis()
	flushRedis(client)

	stock := "ABC"

	result := qExists(client, stock)
	if result != false {
		t.Errorf("qExists was incorrect, got %t, want %t.", result, false)
	}

	client.Cmd("SET", stock, 123.00)

	result = qExists(client, stock)
	if result != true {
		t.Errorf("qExists was incorrect, got %t, want %t.", result, true)
	}
}

func TestSaveTransaction(t *testing.T) {
	client := dialRedis()
	flushRedis(client)

	username := "user"

	result, _ := client.Cmd("ZCOUNT", "HISTORY:"+username, "-inf", "+inf").Int()

	if result != 0 {
		t.Errorf("saveTransaction was incorrect, got %d, want %d.", result, 0)
	}

	command := "ADD"
	amount := "300.00"
	newBalance := "300.00"

	saveTransaction(client, username, command, amount, newBalance)

	result, _ = client.Cmd("ZCOUNT", "HISTORY:"+username, "-inf", "+inf").Int()

	if result != 1 {
		t.Errorf("saveTransaction was incorrect, got %d, want %d.", result, 1)
	}
}

func TestRedisADD(t *testing.T) {
	client := dialRedis()
	flushRedis(client)

	username := "user"
	amount := 123.00

	redisADD(client, username, amount)

	result := exists(client, username)
	if result != true{
		t.Errorf("redisADD was incorrect, got %t, want %t.", result, true)
	}

	newBalance := getBalance(client, username)
	if newBalance != amount {
		t.Errorf("redisADD was incorrect, got %f, want %f.", newBalance, amount)
	}

	// check that the transaction was saved in HISTORY:username
	count, _ := client.Cmd("ZCOUNT", "HISTORY:"+username, "-inf", "+inf").Int()

	if count != 1 {
		t.Errorf("saveTransaction was incorrect, got %d, want %d.", count, 1)
	}
}
