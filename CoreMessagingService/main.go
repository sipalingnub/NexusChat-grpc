package main

import (
	"context"
	"io"
	"log"
	"net"
	"sync"

	pb "nexuschat/proto/coremessaging"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMessagingServiceServer
	
	mu             sync.RWMutex
	activeStreams  map[string]pb.MessagingService_ChatStreamServer
	messageStatus  map[string]string
}

func (s *server) ChatStream(stream pb.MessagingService_ChatStreamServer) error {
	log.Println("🔌 Klien baru membuka koneksi ChatStream!")
	var currentUserID string // Menyimpan ID user yang terhubung di stream ini

	// Pastikan saat stream ini putus, user dihapus dari memori server (CLEANUP ZOMBIE CONNECTION)
	defer func() {
		if currentUserID != "" {
			s.mu.Lock()
			delete(s.activeStreams, currentUserID)
			s.mu.Unlock()
			log.Printf("🗑️ [Cleanup] Klien %s disconnect, dihapus dari memori.", currentUserID)
		}
	}()

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("❌ Error menerima stream: %v", err)
			return err
		}

		currentUserID = msg.GetSenderId() // Set ID user

		// 1. STATE MANAGEMENT: Daftarkan/Update stream klien ini
		s.mu.Lock()
		s.activeStreams[currentUserID] = stream
		s.messageStatus[msg.GetMessageId()] = "SENT"
		s.mu.Unlock()

		log.Printf("📥 [Pesan Masuk] ID: %s | Dari: %s | Ke: %s", msg.GetMessageId(), currentUserID, msg.GetReceiverId())

		// 2. ROUTING PESAN KE PENERIMA
		s.mu.RLock()
		receiverStream, isOnline := s.activeStreams[msg.GetReceiverId()]
		s.mu.RUnlock()

		if isOnline {
			log.Printf("   -> 🚀 Meneruskan pesan ke %s...", msg.GetReceiverId())
			if err := receiverStream.Send(msg); err != nil {
				log.Printf("❌ Gagal meneruskan pesan: %v", err)
			}
		} else {
			log.Printf("   -> ⚠️ Target %s sedang offline. Pesan ditahan di server.", msg.GetReceiverId())
		}
	}
}

func (s *server) AckMessage(ctx context.Context, req *pb.AckRequest) (*pb.AckResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.messageStatus[req.GetMessageId()] = "DELIVERED"
	log.Printf("✅ [ACK] Status pesan %s di-update menjadi DELIVERED (Diterima oleh %s)", req.GetMessageId(), req.GetReceiverId())
	
	return &pb.AckResponse{
		Success: true,
		Status:  "DELIVERED",
	}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50012")
	if err != nil {
		log.Fatalf("Gagal listen di port 50012: %v", err)
	}

	s := grpc.NewServer()
	
	srv := &server{
		activeStreams: make(map[string]pb.MessagingService_ChatStreamServer),
		messageStatus: make(map[string]string),
	}
	
	pb.RegisterMessagingServiceServer(s, srv)

	log.Printf("🔥 CoreMessagingService V3 (With Auto Cleanup) menyala di port 50012...")
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}