package main

import (
	"context"
	"log"
	"net"

	// Import file hasil generate kita tadi
	pb "nexuschat/proto/auth" 
	
	"google.golang.org/grpc"
)

// server struct digunakan untuk mengimplementasikan AuthService yang ada di proto
type server struct {
	pb.UnimplementedAuthServiceServer
}

// Implementasi fungsi Register
func (s *server) Register(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	log.Printf("Menerima request Register untuk user: %s", req.GetUsername())
	
	// Di sini nanti kita bisa tambahkan logika simpan ke map/database
	return &pb.AuthResponse{
		Success: true, 
		Message: "Registrasi berhasil, silakan login!",
	}, nil
}

// Implementasi fungsi Login
func (s *server) Login(ctx context.Context, req *pb.AuthRequest) (*pb.AuthResponse, error) {
	log.Printf("Menerima request Login untuk user: %s", req.GetUsername())
	
	// Di sini nanti kita tambahkan validasi password dan pembuatan JWT sungguhan
	return &pb.AuthResponse{
		Success: true, 
		Message: "Login sukses!", 
		Token:   "dummy-jwt-token-12345", // Token pura-pura sementara
	}, nil
}

func main() {
	// 1. Tentukan port untuk service ini (Sesuai rancangan presentasi: 50011)
	lis, err := net.Listen("tcp", ":50011")
	if err != nil {
		log.Fatalf("Gagal listen di port 50011: %v", err)
	}

	// 2. Buat instance gRPC server baru
	s := grpc.NewServer()

	// 3. Daftarkan service Auth kita ke dalam gRPC server
	pb.RegisterAuthServiceServer(s, &server{})

	log.Printf("🚀 AuthGatewayService sudah menyala dan mendengarkan di port 50011...")
	
	// 4. Jalankan server!
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}