package main

import "testing"

func TestAllocateInvestment(t *testing.T) {
	btcAmt, ethAmt := allocateInvestment(150)
	if btcAmt != 105 && ethAmt != 45 {
		t.Errorf("Expected 105 and 45 but got %v and %v", btcAmt, ethAmt)
	}
}

func TestUsdToCrypto(t *testing.T) {
	cryptoAmt := usdToCrypto(100, .005)
	if cryptoAmt != .5 {
		t.Errorf("Expected .5 but got %v", cryptoAmt)
	}
}
