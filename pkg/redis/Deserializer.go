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

func deserializeArray(rd *bufio.Reader) (*models.TenantSlot, error) {
	returnBytes, _ := getRedisRESP(rd)
	returnBytesToInt, err := strconv.Atoi(string(returnBytes))

	if err != nil {
		fmt.Fprintln(os.Stdout, "erro na conversão de tipo: ", err)
		return nil, err
	}

	tenantMap := make(map[string]string)

	for i := 0; i < returnBytesToInt/2; i++ {
		// Ignora tamanho do Campo
		rd.ReadString('\n')

		field, err := rd.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stdout, "erro na leitura de valores vindos do Redis: ", err)
			return nil, err
		}

		// Ignora tamanho do valor
		rd.ReadString('\n')

		value, err := rd.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stdout, "erro na leitura de valores vindos do Redis: ", err)
			return nil, err
		}

		tenantMap[strings.TrimSpace(field)] = strings.TrimSpace(value)
	}

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
