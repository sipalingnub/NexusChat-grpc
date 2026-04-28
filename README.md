# 🚀 NexusChat: Microservices gRPC & Web Dashboard

NexusChat adalah proyek simulasi sistem *chatting* berskala *enterprise* menggunakan arsitektur **Microservices (gRPC)** di *backend* dan **WebSocket (Node.js)** sebagai *BFF (Backend for Frontend)* untuk Web Dashboard. Proyek ini mendemonstrasikan 4 jenis komunikasi gRPC sekaligus serta integrasinya dengan aplikasi Web *real-time*.

## 🏗️ Arsitektur Sistem

1. **Go Microservices (Backend):**
   - **Auth Service** (Port 50011): Unary gRPC untuk Login & Registrasi.
   - **Core Messaging Service** (Port 50012): Bi-directional Streaming untuk lalu lintas pesan antar *user*.
   - **Media Notification Service** (Port 50013): Client Streaming (Upload) & Server Streaming (Notifikasi Push Sistem).
2. **Node.js Middle-Server (API Gateway):** Berjalan di Port 3000. Bertugas sebagai "Penerjemah" komunikasi gRPC dari *Backend Go* menjadi WebSocket agar bisa diakses secara mulus oleh *Browser*.
3. **Frontend Dashboard:** Web UI interaktif berbasis HTML/JS untuk memantau lalu lintas pesan dan notifikasi *server*.

## 🛠️ Prasyarat (Prerequisites)

Pastikan Anda sudah menginstal alat-alat berikut:
1. **Go (Golang)** v1.20 atau yang lebih baru.
2. **Node.js & npm** v18 atau yang lebih baru.
3. **Protocol Buffers Compiler (protoc)** & gRPC Go Plugins (`protoc-gen-go`, `protoc-gen-go-grpc`).

## 🚦 Cara Menjalankan Sistem

Sistem ini membutuhkan beberapa terminal yang berjalan secara bersamaan. Ikuti urutan berikut:

### Langkah 1: Nyalakan Semua Go Microservices
Buka 3 terminal terpisah di direktori *root* proyek, jalankan masing-masing perintah:
```bash
go run ./AuthGatewayService/main.go
go run ./CoreMessagingService/main.go
go run ./MediaNotificationService/main.go
```
## Langkah 2: Nyalakan Node.js Web Dashboard

Buka terminal ke-4, masuk ke dalam folder middle-server, lalu jalankan:
```bash
cd nexuschat-websocket/middle-server
npm install
node server.js
```

👉 Buka browser dan kunjungi: http://localhost:3000

## Langkah 3: Jalankan Klien Terminal (Go CLI)

Untuk menyimulasikan percakapan dengan Web Dashboard, buka terminal baru dan jalankan klien dengan format go run ./ClientApp/main.go [NamaKamu] AdminWeb:
```bash
go run ./ClientApp/main.go arul AdminWeb
go run ./ClientApp/main.go "david bowie" AdminWeb
```

## Fitur Unggulan:

- **Live gRPC Traffic**: Pemantauan beban lalu lintas pesan secara real-time menggunakan Chart.js.
- **Event-Driven System Alerts**: Menerima ping dan notifikasi langsung dari Server Go (Port 50013) tanpa henti (infinite stream).
- **Smart Target Selection**: Nama klien yang baru mengirim pesan akan otomatis muncul sebagai pill button. Admin dapat membalas pesan cukup dengan mengeklik nama tersebut.
- **State Management**: Backend dioptimalkan dengan sync.RWMutex untuk mencegah data race saat banyak klien terhubung.
