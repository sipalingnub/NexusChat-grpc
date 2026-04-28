package main

import (
	"io"
	"log"
	"net"
	"time" // Import time sekarang sudah benar di sini

	pb "nexuschat/proto/medianotif"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMediaNotifServiceServer
}

// 1. Client Streaming: Menerima potongan file (Chunks)
func (s *server) UploadFile(stream pb.MediaNotifService_UploadFileServer) error {
	var totalSize int
	var fileName string

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			log.Printf("📂 File %s berhasil diunggah. Total ukuran: %d bytes", fileName, totalSize)
			return stream.SendAndClose(&pb.UploadResponse{
				Success: true,
				FileUrl: "http://nexus-storage.com/" + fileName,
			})
		}
		if err != nil {
			return err
		}
		fileName = chunk.GetFilename()
		totalSize += len(chunk.GetChunkData())
	}
}

// 2. Server Streaming: Dorong notifikasi ke semua klien
func (s *server) SubscribeNotif(req *pb.SubscribeRequest, stream pb.MediaNotifService_SubscribeNotifServer) error {
	log.Printf("🔔 User %s berlangganan notifikasi sistem", req.GetClientId())
	
	notifs := []string{
		"Selamat datang di NexusChat!",
		"Pengumuman: Maintenance jam 12 malam nanti.",
		"Tips: Gunakan password yang kuat!",
	}

	// 1. Kirim 3 notif awal
	for _, n := range notifs {
		if err := stream.Send(&pb.Notification{SystemMessage: n, Type: "INFO"}); err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
	}

	// 2. INFINITE LOOP: Agar stream tetap hidup selamanya
	for {
		time.Sleep(15 * time.Second) // Jeda 15 detik
		pesanPing := "Sistem beroperasi normal (Ping)"
		if err := stream.Send(&pb.Notification{SystemMessage: pesanPing, Type: "SYSTEM"}); err != nil {
			log.Printf("📴 Klien %s putus dari stream notif.", req.GetClientId())
			return err // Berhenti kalau Node.js dimatikan
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":50013")
	if err != nil {
		log.Fatalf("Gagal listen port 50013: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterMediaNotifServiceServer(s, &server{})

	log.Printf("📢 MediaNotificationService menyala di port 50013...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Gagal: %v", err)
	}
}