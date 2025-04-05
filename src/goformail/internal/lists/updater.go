package lists

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func createList(list string, configs map[string]string) error {
	list = list + "@" + configs["EMAIL_DOMAIN"]
	file, err := os.OpenFile(configs["MAP_LOCATION"], os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		if err := file.Close(); err != nil {
			fmt.Println(err)
		}
	}(file)

	// Check user doesn't already exist
	allUsers, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if strings.Contains(string(allUsers), list) {
		fmt.Println(list)
		return errors.New("list " + list + " already exists")
	}

	// Update list file
	newEntry := fmt.Sprintf("%s lmtp:inet:%s:%s\n", list, configs["OWN_ADDRESS"], configs["LMTP_PORT"])
	if !strings.HasSuffix(string(allUsers), "\n") {
		newEntry = "\n" + newEntry
	}
	if _, err = file.WriteString(newEntry); err != nil {
		return err
	}

	return nil
}
