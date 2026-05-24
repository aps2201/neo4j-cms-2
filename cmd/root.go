package cmd

import (
	"context"
	"fmt"
	"os"

	"aps.web.id/neo4j-cms-2/secrets"
	"github.com/neo4j/neo4j-go-driver/v6/neo4j"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "neo4j-cms",
	Short: "A CLI neo4j CMS tool.",
	Long:  "This is a CLI neo4j CMS tool for my website. Its a nice little tool I've been using.",
}

var neoDriver neo4j.Driver

func GetDriver() neo4j.Driver {
	if neoDriver != nil {
		return neoDriver
	}
	driver, err := neo4j.NewDriver(secrets.Neo4j_cred["NEO4J_URI"], neo4j.BasicAuth(secrets.Neo4j_cred["NEO4J_USERNAME"], secrets.Neo4j_cred["NEO4J_PASSWORD"], ""))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	if err := driver.VerifyConnectivity(ctx); err != nil {
		panic(err)
	}
	//fmt.Println("neo4j connection good")
	neoDriver = driver
	return neoDriver
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
