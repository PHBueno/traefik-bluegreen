package redis

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type DataType byte

const (
	array DataType = '*'
)

func getRedisRESP(rd *bufio.Reader) ([]byte, error) {
	resp, err := rd.ReadBytes('\n')

	if err != nil {
		fmt.Fprintln(os.Stdout, "erro na resposta do redis: ", err)
		return nil, err
	}

	return resp[:len(resp)-2], nil // *6

}

func deserializeArray(rd *bufio.Reader) {
	d, _ := getRedisRESP(rd)
	fmt.Fprintf(os.Stdout, "RETORNO:  %s\n", d)

}

func Deserializer(rd io.Reader) error {
	reader := bufio.NewReader(rd)

	respType, err := reader.ReadByte()

	if err != nil {
		return err
	}

	switch DataType(respType) {
	case array:
		deserializeArray(reader)
	}
	return nil

}
