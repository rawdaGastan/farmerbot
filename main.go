/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/

//go:generate mockgen -source internal/rmb.go -destination mock_rmb.go -package mocks github.com/rawdaGastan/farmerbot FarmHandler

package main

import "github.com/rawdaGastan/farmerbot/cmd"

func main() {
	cmd.Execute()
}
