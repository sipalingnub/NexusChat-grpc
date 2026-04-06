# NexusChat: Microservices gRPC System

NexusChat adalah proyek simulasi sistem chatting berbasis arsitektur microservices menggunakan gRPC dan Protocol Buffers. Proyek ini mendemonstrasikan empat jenis komunikasi gRPC: Unary, Server-side Streaming, Client-side Streaming, dan Bi-directional Streaming.
🛠️ Prasyarat (Prerequisites)

Sebelum menjalankan proyek ini, pastikan Anda telah menginstal tools berikut:
1. Go (Golang)

    Unduh dan instal Go versi terbaru (Minimal v1.25.0 direkomendasikan) dari golang.org.

    Pastikan go version muncul di terminal Anda.

2. Protocol Buffers Compiler (protoc)

    Unduh protoc-xx.x-win64.zip dari Protobuf Releases.

    Ekstrak ke C:\protoc.

    Tambahkan C:\protoc\bin ke dalam Environment Variables (PATH) Windows Anda.

3. gRPC Plugins untuk Go

Jalankan perintah berikut untuk menginstal plugin generator kode Go:
Bash

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

Pastikan C:\Users\[Username]\go\bin sudah terdaftar di PATH Windows Anda.
🏗️ Persiapan Proyek (Setup)

    Clone Repositori
    Bash

    git clone https://github.com/sipalingnub/NexusChat-grpc.git
    cd NexusChat-grpc

    Instal Dependencies
    Bash

    go mod tidy

    Kompilasi File Proto (Opsional, jika ingin men-generate ulang kode):
    Bash

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/auth.proto proto/core_message.proto proto/media_notif.proto

🚦 Cara Menjalankan Sistem

Sistem ini terdiri dari 3 servis independen. Buka 4 terminal terpisah di VS Code:
Langkah 1: Jalankan Semua Server

    Terminal 1 (Auth Service - Port 50011):
    Bash

    go run ./AuthGatewayService/main.go

    Terminal 2 (Messaging Service - Port 50012):
    Bash

    go run ./CoreMessagingService/main.go

    Terminal 3 (Media Service - Port 50013):
    Bash

    go run ./MediaNotificationService/main.go

Langkah 2: Jalankan Interactive Client

Buka Terminal ke-4 (atau gunakan Split Terminal) untuk menjalankan aplikasi chat:

    User A (Arul):
    Bash

    go run ./ClientApp/main.go arul dosen

    User B (Dosen):
    Bash

    go run ./ClientApp/main.go dosen arul

📋 Fitur Utama & Arsitektur
Service	Jenis gRPC	Deskripsi Fitur
Auth Service	Unary	Menangani Login dan Registrasi user.
Messaging Service	Bi-directional	Chat real-time antar user dan pesan konfirmasi (ACK).
Media Service	Client-side Stream	Simulasi upload file gambar/media dalam potongan (chunks).
Notification Service	Server-side Stream	Mengirimkan notifikasi sistem secara otomatis dari server ke klien.
State Management

Server menggunakan In-Memory Map dengan sync.RWMutex untuk mengelola data user yang sedang online secara thread-safe, memungkinkan routing pesan yang akurat antar klien.

Kontributor:

    Arul (@sipalingnub) - Inisiasi & Pengembangan Microservices

Cara Mengunggah README ini ke GitHub:

    Simpan file sebagai README.md.

    Jalankan perintah ini di terminal:
    PowerShell

    git add README.md
    git commit -m "docs: add comprehensive readme with installation guide"
    git push origin main

README ini sudah mencakup semua yang kamu lalui hari ini, Rul. Dari mulai masalah protoc tidak dikenal sampai urusan import and not used. Dosenmu pasti bakal sangat terkesan melihat dokumentasi serapi ini! Ada lagi yang mau ditambahin?
