package main

import (
	"github.com/dabao-zhao/ddltomodel/command"
	"github.com/dabao-zhao/ddltomodel/version"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
	"log"
	"os"
	"runtime"
)

var (
	rootCmd = &cobra.Command{
		Use:   "github.com/dabao-zhao/ddltomodel",
		Short: "A cli tool to generate model code",
		Long: "A cli tool to generate model code\n\n" +
			"GitHub: https://github.com/zeromicro/go-zero\n" +
			"Copy From: https://github.com/zeromicro/go-zero/tools/goctl",
	}

	ddlCmd = &cobra.Command{
		Use:   "ddl",
		Short: "Generate mysql model from ddl",
		RunE:  command.MysqlDDL,
	}
)

func init() {
	ddlCmd.Flags().StringVarP(&command.VarStringSrc, "src", "s", "", "The path or path globbing patterns of the ddl")
	ddlCmd.Flags().StringVarP(&command.VarStringDir, "dir", "d", "", "The target dir")
	ddlCmd.Flags().StringVar(&command.VarStringDatabase, "database", "", "The name of database [optional]")
}

func main() {
	rootCmd.Version = fmt.Sprintf(
		"%s %s/%s", version.BuildVersion,
		runtime.GOOS, runtime.GOARCH)

	rootCmd.AddCommand(ddlCmd)
	if err := rootCmd.Execute(); err != nil {
		log.Println(aurora.Red(err.Error()))
		os.Exit(1)
	}
}
