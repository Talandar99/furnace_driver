package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"os"
	pb "talandar.dev/furnace_driver/proto"
)

func main() {
	report := &pb.SensorReport{
		Id: 47,
		Ver: &pb.Version{
			Major: 1,
			Minor: 3,
			Patch: 2,
		},
		Status: "ERROR",
		Logs: "0001 INF Sensor init OK\n" +
			"0042 DBG Calibration done\n" +
			"0187 WRN Watchdog reset\n" +
			"0203 ERR Adc overflow CH:2\n" +
			"0519 INF Firmware update 1.3.2",
		Data: []*pb.DataPoint{
			{
				Time:   "2026-10-04 08:00:00",
				Temp:   22.3,
				Humi:   float64Ptr(45.1),
				Enable: true,
			},
			{
				Time:   "2026-10-04 08:30:00",
				Temp:   21.8,
				Humi:   float64Ptr(47.6),
				Enable: false,
			},
			{
				Time:   "2026-10-04 09:00:00",
				Temp:   22.1,
				Enable: true,
				// Humi brak => odpowiednik null
			},
		},
	}

	data, err := proto.Marshal(report)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("report.bin", data, 0o644)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Zapisano %d bajtów do report.bin\n\n", len(data))

	raw, err := os.ReadFile("report.bin")
	if err != nil {
		panic(err)
	}

	var decoded pb.SensorReport
	err = proto.Unmarshal(raw, &decoded)
	if err != nil {
		panic(err)
	}

	printReport(&decoded)
}

func float64Ptr(v float64) *float64 {
	return &v
}

func printReport(r *pb.SensorReport) {
	fmt.Println("ID:", r.Id)

	if r.Ver != nil {
		fmt.Printf("Version: %d.%d.%d\n", r.Ver.Major, r.Ver.Minor, r.Ver.Patch)
	}

	fmt.Println("Status:", r.Status)
	fmt.Println("Logs:")
	fmt.Println(r.Logs)
	fmt.Println()

	for i, dp := range r.Data {
		fmt.Printf("DataPoint #%d\n", i+1)
		fmt.Println("  Time  :", dp.Time)
		fmt.Println("  Temp  :", dp.Temp)
		fmt.Println("  Enable:", dp.Enable)

		if dp.Humi != nil {
			fmt.Println("  Humi  :", *dp.Humi)
		} else {
			fmt.Println("  Humi  : <null / brak>")
		}

		fmt.Println()
	}
}
