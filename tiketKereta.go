package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type penumpang struct {
	nDpn, nBlkg, tDuduk, status string
	kodePembayaran, kodeBook    int
	WaktuReservasi              time.Time
}

const kapasitas = 40

type gerbong struct {
	nGerbong string
	kursi    [kapasitas]bool
	harga    int
}

var gerbongA = gerbong{nGerbong: "A", harga: 500000}
var gerbongB = gerbong{nGerbong: "B", harga: 300000}
var reservasiList [kapasitas * 2]penumpang
var counterBayar int = 0
var bookingCounter int = 0

func menu() {
	fmt.Println("============================")
	fmt.Println("|  Reservasi Tiket Kereta  |")
	fmt.Println("============================")
	fmt.Println("| 1. Pesan Tiket           |")
	fmt.Println("| 2. Cek Status            |")
	fmt.Println("| 3. Pembayaran            |")
	fmt.Println("| 4. Daftar Penumpang      |")
	fmt.Println("| 5. Exit                  |")
	fmt.Println("============================")
	fmt.Print("Masukan Perintah Anda: ")
}

func pesanTiket() {
	var p penumpang
	var a, b, tiket int
	var tValid, pValid bool

	fmt.Println("\n============================")
	fmt.Println("|        Pesan Tiket       |")
	fmt.Println("============================")
	fmt.Print("Jumlah tiket yang ingin dipesan (max 4): ")

	for !tValid {
		fmt.Scan(&tiket)
		if tiket < 1 || tiket > 4 {
			fmt.Print("Perintah tidak valid, coba lagi: ")
		} else {
			tValid = true
		}
	}

	fmt.Println("============================")
	fmt.Print("Jumlah penumpang Gerbong A: ")
	for !pValid {
		fmt.Scan(&a)
		if a < 0 || a > tiket {
			fmt.Print("Perintah tidak valid, coba lagi: ")
		} else {
			pValid = true
		}
	}
	b = tiket - a
	fmt.Println("Jumlah penumpang Gerbong B:", b)
	fmt.Println("============================")
	counterBayar++
	for i := 0; i < a; i++ {
		p.kodePembayaran = counterBayar
		inputPenumpang(&p, "A", 1)
		fmt.Println("============================")
	}

	for i := 0; i < b; i++ {
		p.kodePembayaran = counterBayar
		inputPenumpang(&p, "B", 1)
		fmt.Println("============================")
	}
	fmt.Printf("Reservasi berhasil. Kode Pembayaran Anda: %d\n", p.kodePembayaran)
	fmt.Printf("Waktu Reservasi: %s\n", p.WaktuReservasi.Format("15:04:05"))
	fmt.Println()
}

func inputPenumpang(p *penumpang, gerbong string, jmlhTiket int) {
	var seatIndex int
	var sukses bool

	fmt.Print("Masukan Nama Depan   : ")
	fmt.Scan(&p.nDpn)

	fmt.Print("Masukan Nama Belakang: ")
	fmt.Scan(&p.nBlkg)

	if gerbong == "A" {
		allocateSeats(&gerbongA, jmlhTiket, &seatIndex, &sukses)
	} else {
		allocateSeats(&gerbongB, jmlhTiket, &seatIndex, &sukses)
	}

	if !sukses {
		fmt.Println("Tidak ada kursi berdekatan yang tersedia untuk jumlah tiket yang diminta di gerbong", gerbong)
		return
	}
	bookingCounter++

	p.kodeBook = bookingCounter

	baris := seatIndex / 4
	kolom := seatIndex % 4
	p.tDuduk = fmt.Sprintf("%s-%d-%d", gerbong, baris+1, kolom+1)
	p.status = "Reserved"
	p.WaktuReservasi = time.Now()
	reservasiList[p.kodeBook] = *p

	fmt.Printf("Kode Booking         : %d\nKursi                : %s\n", p.kodeBook, p.tDuduk)
	fmt.Println()
}

func allocateSeats(g *gerbong, jumlahTiket int, seatIndex *int, success *bool) {
	*seatIndex = -1
	*success = false
	for i := 0; i <= len(g.kursi)-jumlahTiket; i++ {
		allAvailable := true
		for j := 0; j < jumlahTiket; j++ {
			if g.kursi[i+j] {
				allAvailable = false
				break
			}
		}
		if allAvailable {
			for j := 0; j < jumlahTiket; j++ {
				g.kursi[i+j] = true
			}
			*seatIndex = i
			*success = true
			break
		}
	}
	if !*success {
		fmt.Println("Tidak Ada Kursi yang Tersedia")
	}
}

func cekStatus() {
	fmt.Println("\n============================")
	fmt.Println("|        Cek Status        |")
	fmt.Println("============================")
	fmt.Print("Masukkan kode booking: ")
	var kode int
	fmt.Scan(&kode)

	// Filter out empty entries and create a slice of valid reservations
	var validReservations []penumpang
	for i := 0; i < len(reservasiList); i++ {
		p := reservasiList[i]
		if p.nDpn != "" {
			validReservations = append(validReservations, p)
		}
	}

	kiri := 0
	kanan := len(validReservations) - 1
	found := false

	for kiri <= kanan {
		tengah := (kiri + kanan) / 2
		if validReservations[tengah].kodeBook > kode {
			kanan = tengah - 1
		} else if validReservations[tengah].kodeBook < kode {
			kiri = tengah + 1
		} else {
			p := validReservations[tengah]
			fmt.Printf("Nama           : %s, %s\n", p.nBlkg, p.nDpn)
			fmt.Printf("Kode Booking   : %d\n", p.kodeBook)
			fmt.Printf("Kode Pembayaran: %d\n", p.kodePembayaran)
			fmt.Printf("Tempat Duduk   : %s\n", p.tDuduk)
			fmt.Printf("Status         : %s\n", p.status)
			fmt.Printf("Waktu Reservasi: %s\n", p.WaktuReservasi.Format("15:04:05"))
			found = true
			break
		}
	}

	if !found {
		fmt.Println("Kode booking tidak ditemukan.")
	}
	fmt.Println()
}

func pembayaran() {
	fmt.Println("\n===========================")
	fmt.Println("|        Pembayaran       |")
	fmt.Println("===========================")
	bacaTxt()
	fmt.Print("Masukkan Kode Pembayaran: ")
	var kodePembayaran int
	fmt.Scan(&kodePembayaran)

	var totalHarga int
	var pesananDitemukan bool
	var sudahKadaluarsa bool

	for i := 0; i < len(reservasiList); i++ {
		p := &reservasiList[i]
		if p.kodePembayaran == kodePembayaran {
			pesananDitemukan = true

			if time.Since(p.WaktuReservasi).Minutes() > 1 && p.status == "Reserved" { // assuming expiration time is 60 minutes
				p.status = "Expired"
				fmt.Printf("Reservasi dengan Kode Booking %d telah kadaluarsa.\n", p.kodeBook)
				sudahKadaluarsa = true
			}

			if p.status != "Expired" {
				if p.tDuduk[0] == 'A' {
					totalHarga += gerbongA.harga
				} else {
					totalHarga += gerbongB.harga
				}
			}
		}
	}

	if !pesananDitemukan {
		fmt.Println("Reservasi tidak ditemukan.")
		return
	}

	if sudahKadaluarsa {
		fmt.Println("Reservasi telah kadaluarsa dan tidak dapat dibayar.")
		return
	}

	fmt.Printf("Harga total yang harus dibayar untuk Kode Pembayaran %d adalah %d.\n", kodePembayaran, totalHarga)
	fmt.Print("Konfirmasi pembayaran? (y/n): ")
	var konfirmasi string
	fmt.Scan(&konfirmasi)

	if konfirmasi == "y" {
		for i := 0; i < len(reservasiList); i++ {
			if reservasiList[i].kodePembayaran == kodePembayaran {
				reservasiList[i].status = "Paid"
			}
		}
		saveReservasiData()
		fmt.Println("Pembayaran berhasil.")
	} else {
		fmt.Println("Pembayaran dibatalkan.")
	}
	fmt.Println()
}

func cetakListPenumpang() {
	fmt.Println("\n============================")
	fmt.Println("|     Daftar Penumpang     |")
	fmt.Println("============================")

	// Menggunakan map untuk melacak penumpang yang telah ditambahkan
	seen := make(map[int]bool)
	var validPenumpang []penumpang

	// Mengumpulkan penumpang berdasarkan kode booking
	for i := 0; i < len(reservasiList); i++ {
		p := reservasiList[i]
		if p.nDpn != "" && !seen[p.kodeBook] {
			validPenumpang = append(validPenumpang, p)
			seen[p.kodeBook] = true
		}
	}

	for pass := 1; pass < len(validPenumpang); pass++ {
		temp := validPenumpang[pass]
		i := pass - 1
		for i >= 0 && validPenumpang[i].nBlkg > temp.nBlkg {
			validPenumpang[i+1] = validPenumpang[i]
			i--
		}
		validPenumpang[i+1] = temp
	}

	fmt.Printf("%15s %7s %10s\n", "Nama Belakang", "Kursi", "Status")
	for i := 0; i < len(validPenumpang); i++ {
		p := validPenumpang[i]
		fmt.Printf("%10s %12s %10s\n", p.nBlkg, p.tDuduk, p.status)
	}
	fmt.Println()
}

func exit() {
	fmt.Println("\n============================")
	fmt.Println("| -= Anda Telah Keluar! =- |")
	fmt.Println("============================\n")
	os.Exit(0)
}

func saveReservasiData() {
	existingData := make(map[string]bool)

	// Membaca data yang sudah ada dalam file
	file, err := os.OpenFile("data.txt", os.O_RDONLY, 0644)
	if err != nil {
		if !os.IsNotExist(err) {
			panic(err)
		}
	} else {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.Split(line, ", ")
			if len(parts) >= 4 {
				kodeBook := strings.TrimPrefix(parts[3], "Kode Booking: ")
				existingData[kodeBook] = true
			}
		}
		file.Close()
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}

	// Membuka file dalam mode append
	file, err = os.OpenFile("data.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Menyimpan data baru
	for _, p := range reservasiList {
		if p.nDpn != "" {
			kodeBook := strconv.Itoa(p.kodeBook)
			if !existingData[kodeBook] {
				_, err := file.WriteString(fmt.Sprintf("Nama Depan: %s, Nama Belakang: %s, Kursi: %s, Kode Booking: %d, Kode Pembayaran: %d, Status: %s, Waktu Reservasi: %s\n",
					p.nDpn, p.nBlkg, p.tDuduk, p.kodeBook, p.kodePembayaran, p.status, p.WaktuReservasi.Format("2006-01-02 15:04:05")))
				if err != nil {
					panic(err)
				}
				existingData[kodeBook] = true
			}
		}
	}
}

func bacaTxt() {
	file, err := os.Open("data.txt")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ", ")
		if len(parts) != 7 {
			continue
		}
		nDpn := strings.TrimPrefix(parts[0], "Nama Depan: ")
		nBlkg := strings.TrimPrefix(parts[1], "Nama Belakang: ")
		tDuduk := strings.TrimPrefix(parts[2], "Kursi: ")
		kodeBook, _ := strconv.Atoi(strings.TrimPrefix(parts[3], "Kode Booking: "))
		kodePembayaran, _ := strconv.Atoi(strings.TrimPrefix(parts[4], "Kode Pembayaran: "))
		status := strings.TrimPrefix(parts[5], "Status: ")
		WaktuReservasi, _ := time.Parse("2006-01-02 15:04:05", strings.TrimPrefix(parts[6], "Waktu Reservasi: "))

		p := penumpang{
			nDpn:           nDpn,
			nBlkg:          nBlkg,
			tDuduk:         tDuduk,
			kodeBook:       kodeBook,
			kodePembayaran: kodePembayaran,
			status:         status,
			WaktuReservasi: WaktuReservasi,
		}

		var gerbong *gerbong
		if p.tDuduk[0] == 'A' {
			gerbong = &gerbongA
		} else {
			gerbong = &gerbongB
		}

		seatParts := strings.Split(p.tDuduk, "-")
		if len(seatParts) == 3 {
			row, _ := strconv.Atoi(seatParts[1])
			col, _ := strconv.Atoi(seatParts[2])
			gerbong.kursi[(row-1)*4+(col-1)] = true
		}

		for i := 0; i < len(reservasiList); i++ {
			if reservasiList[i].nDpn == "" {
				reservasiList[i] = p
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func main() {
	bacaTxt()

	var x int
	for {
		menu()
		fmt.Scan(&x)

		switch x {
		case 1:
			pesanTiket()
		case 2:
			cekStatus()
		case 3:
			pembayaran()
		case 4:
			cetakListPenumpang()
		case 5:
			exit()
		default:
			fmt.Print("Perintah tidak valid, coba lagi: ")
		}
	}
}
