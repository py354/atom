package main

import (
	"atom/internal/config"
	"atom/pkg/api"
	"bufio"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"
	"strconv"
	"strings"
)

const HELP_MSG = `
Команды:
help - получить данное сообщение
get-all
get <номер>
delete <номер>
add-mileage <номер> <пробег>
get-cost <номер>
register <номер> <модель> <назначение> <год выпуска> <пробег>
назначение должно быть SHARING, TAXI или DELIVERY
`

func main() {
	conf, err := config.GetConfig()
	if err != nil {
		log.Fatalln(err)
	}

	addr := fmt.Sprintf("%s:%s", conf.Grpc.Host, conf.Grpc.Port)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	ctx := context.Background()
	client := api.NewFleetClient(conn)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Для вывода команд, испроьзуй команду help")
	commands := []Command{
		helpCommand, getAllCommand, getCar, deleteCar, addMileage, getCost, registerCar,
	}

	for {
		text, _ := reader.ReadString('\n')
		args := strings.Split(strings.TrimRight(text, "\n"), " ")

		for _, cmd := range commands {
			if cmd(ctx, client, args) {
				fmt.Println("-------------------------------------")
				continue
			}
		}

	}
}

type Command func(ctx context.Context, client api.FleetClient, args []string) bool

func helpCommand(ctx context.Context, client api.FleetClient, args []string) bool {
	if len(args) != 1 || args[0] != "help" {
		return false
	}

	fmt.Printf(HELP_MSG)
	return true
}

func getAllCommand(ctx context.Context, client api.FleetClient, args []string) bool {
	if len(args) != 1 || args[0] != "get-all" {
		return false
	}

	cars, err := client.GetAllCars(ctx, &api.Empty{})
	fmt.Println("Ошибка:", err)
	if err == nil {
		for _, car := range cars.GetCars() {
			fmt.Println(parseCar(car))
		}
	}
	return true
}

func getCar(ctx context.Context, client api.FleetClient, args []string) bool {
	if len(args) != 2 || args[0] != "get" {
		return false
	}

	car, err := client.GetCar(ctx, &api.RP{RegistrationPlate: args[1]})
	fmt.Println("Ошибка:", err)
	fmt.Println(parseCar(car))
	return true
}

func getCost(ctx context.Context, client api.FleetClient, args []string) bool {
	if len(args) != 2 || args[0] != "get-cost" {
		return false
	}

	car, err := client.GetEstimatedCost(ctx, &api.RP{RegistrationPlate: args[1]})
	fmt.Println("Ошибка:", err)
	fmt.Println(car.GetCost())
	return true
}

func deleteCar(ctx context.Context, client api.FleetClient, args []string) bool {
	if len(args) != 2 || args[0] != "delete" {
		return false
	}

	_, err := client.DeleteCar(ctx, &api.RP{RegistrationPlate: args[1]})
	fmt.Println("Ошибка:", err)
	return true
}

func addMileage(ctx context.Context, client api.FleetClient, args []string) bool {
	if len(args) != 3 || args[0] != "add-mileage" {
		return false
	}

	mileage, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println("Ошибка:", err)
		return true
	}

	_, err = client.AddMileage(ctx, &api.AddMileageRequest{RegistrationPlate: args[1], Mileage: int32(mileage)})
	fmt.Println("Ошибка:", err)
	return true
}

func registerCar(ctx context.Context, client api.FleetClient, args []string) bool {
	if len(args) != 6 || args[0] != "register" {
		return false
	}

	//register <номер> <модель> <назначение> <год выпуска> <пробег>
	year, err := strconv.Atoi(args[4])
	if err != nil {
		fmt.Println("Ошибка:", err)
		return true
	}

	mileage, err := strconv.Atoi(args[5])
	if err != nil {
		fmt.Println("Ошибка:", err)
		return true
	}

	_, err = client.RegisterCar(ctx, &api.Car{
		RegistrationPlate: args[1],
		Model:             args[2],
		Purpose:           api.Purpose(api.Purpose_value[args[3]]),
		ManufactureYear:   int32(year),
		Mileage:           int32(mileage),
	})
	fmt.Println("Ошибка:", err)
	return true
}

func parseCar(car *api.Car) string {
	template := "[%s] moodal:%s, purpose:%s, year:%d, mmileage:%d"
	return fmt.Sprintf(template, car.GetRegistrationPlate(), car.GetModel(),
		car.Purpose.String(), car.ManufactureYear, car.Mileage)
}
