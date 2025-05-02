package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/db"
	"gitlab.computing.dcu.ie/fonseca3/2025-csc1097-fonseca3-dagohos2/internal/service"
	"log"
	"strconv"
	"strings"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	parts := strings.Split(value, ",")
	*s = append(*s, parts...)
	return nil
}

type intSlice []int64

func (s *intSlice) String() string {
	strVals := make([]string, len(*s))
	for i, val := range *s {
		strVals[i] = strconv.Itoa(int(val))
	}
	return strings.Join(strVals, ",")
}

func (s *intSlice) Set(value string) error {
	parts := strings.Split(value, ",")
	for _, part := range parts {
		if num, err := strconv.Atoi(strings.TrimSpace(part)); err == nil {
			*s = append(*s, int64(num))
		} else {
			log.Fatal("Invalid mod id(s) provided")
		}
	}
	return nil
}

func getListManager(db db.IDb) service.IListManager {
	return service.NewListManager(db)
}

func getUserManager(db db.IDb) service.IUserManager {
	return service.NewUserManager(db)
}

func parseArgs(cmd *flag.FlagSet, args []string) {
	if err := cmd.Parse(args); err != nil {
		log.Fatal(err)
	}
}

func convertId(id string) int {
	idInt, err := strconv.Atoi(id)
	if err != nil {
		log.Fatalf("Invalid id '%s'\n", id)
	}
	return idInt
}

func printObject(object interface{}) {
	marshal, err := json.MarshalIndent(object, "", "\t")
	if err != nil {
		log.Fatal("Error printing payload:", err)
	}
	fmt.Println(string(marshal))
}
