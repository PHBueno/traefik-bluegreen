package redis

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/PHBueno/traefik-bluegreen/pkg/redis/models"
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

func readBulkString(rd *bufio.Reader) (string, error) {
	_, err := rd.ReadString('\n')

	if err != nil {
		fmt.Fprintln(os.Stdout, "erro na leitura de valores vindos do Redis: ", err)
		return "", err
	}

	value, err := rd.ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stdout, "erro na leitura de valores vindos do Redis: ", err)
		return "", err
	}

	return strings.TrimSpace(value), nil

}

func readPair(rd *bufio.Reader) (string, string, error) {
	field, err := readBulkString(rd)
	if err != nil {
		return "", "", err
	}

	value, err := readBulkString(rd)
	if err != nil {
		return "", "", err
	}

	return field, value, nil
}

func deserializeArray(rd *bufio.Reader) (*models.TenantSlot, error) {
	returnBytes, _ := getRedisRESP(rd)
	returnBytesToInt, err := strconv.Atoi(string(returnBytes))

	if err != nil {
		fmt.Fprintln(os.Stdout, "erro na conversão de tipo: ", err)
		return nil, err
	}

	tenantMap := make(map[string]string)

	for i := 0; i < returnBytesToInt/2; i++ {
		field, value, err := readPair(rd)
		if err != nil {
			fmt.Fprintln(os.Stdout, "erro na leitura de valores vindos do Redis: ", err)
			return nil, err
		}
		// Ignora tamanho do Campo
		// rd.ReadString('\n')

		// field, err := rd.ReadString('\n')
		// if err != nil {
		// 	fmt.Fprintln(os.Stdout, "erro na leitura de valores vindos do Redis: ", err)
		// 	return nil, err
		// }

		// // Ignora tamanho do valor
		// rd.ReadString('\n')

		// value, err := rd.ReadString('\n')
		// if err != nil {
		// 	fmt.Fprintln(os.Stdout, "erro na leitura de valores vindos do Redis: ", err)
		// 	return nil, err
		// }

		tenantMap[field] = value
	}

	fmt.Fprintln(os.Stdout, tenantMap)

	return &models.TenantSlot{
		TenantID: tenantMap["tenant"],
		AppName:  tenantMap["app"],
		Slot:     tenantMap["slot"],
	}, nil

}

func Deserializer(rd io.Reader) (*models.TenantSlot, error) {
	reader := bufio.NewReader(rd)

	respType, err := reader.ReadByte()

	if err != nil {
		return nil, err
	}

	switch DataType(respType) {
	case array:
		return deserializeArray(reader)
	// Implementar retorno default
	default:
		return nil, nil
	}

}
