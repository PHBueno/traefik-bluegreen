package redis

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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

func deserializeArray(rd *bufio.Reader) error {
	returnBytes, _ := getRedisRESP(rd)
	returnBytesToInt, err := strconv.Atoi(string(returnBytes))

	if err != nil {
		fmt.Fprintln(os.Stdout, "erro na conversão de tipo: ", err)
		return err
	}

	for i := 0; i <= returnBytesToInt/2; i++ {
		// Ignora tamanho do Campo
		rd.ReadString('\n')

		field, _ := rd.ReadString('\n')

		// Ignora tamanho do valor
		rd.ReadString('\n')

		value, _ := rd.ReadString('\n')

		fmt.Fprintf(os.Stdout, "[%d] %s => %s", i, field, value)
	}

	return nil

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
