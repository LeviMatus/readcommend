package main

import (
	"github.com/LeviMatus/readcommend/service/cmd"
	_ "github.com/lib/pq"
)

func main() {
	cmd.Execute()
}
