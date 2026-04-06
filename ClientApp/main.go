package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pbAuth "nexuschat/proto/auth"
	pbChat "nexuschat/proto/coremessaging"
	pbMedia "nexuschat/proto/medianotif" // Import servis ke-3

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("❌ Cara pakai: go run main.go [NamaKamu] [NamaLawanChat]")
	}
	myName := os.Args[1]
	friendName := os.Args[2]

	log.Printf("👤 User: %s | 🎯 Target: %s", myName, friendName)

	// --- 1. KONEKSI KE SEMUA SERVICE ---
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	
	connAuth, _ := grpc.NewClient("localhost:50011", opts)
	connChat, _ := grpc.NewClient("localhost:50012", opts)
	connMedia, _ := grpc.NewClient("localhost:50013", opts) // Port 50013
	defer connAuth.Close()
	defer connChat.Close()
	defer connMedia.Close()

	authClient := pbAuth.NewAuthServiceClient(connAuth)
	chatClient := pbChat.NewMessagingServiceClient(connChat)
	mediaClient := pbMedia.NewMediaNotifServiceClient(connMedia)

	// --- 2. AUTHENTICATION ---
	_, err := authClient.Login(context.Background(), &pbAuth.AuthRequest{Username: myName, Password: "123"})
	if err != nil {
		log.Fatalf("❌ Login gagal: %v", err)
	}

	// --- 3. FITUR BARU: SUBSCRIBE NOTIFIKASI (Server Streaming) ---
	go func() {
		streamNotif, _ := mediaClient.SubscribeNotif(context.Background(), &pbMedia.SubscribeRequest{ClientId: myName})
		for {
			res, err := streamNotif.Recv()
			if err != nil { break }
			fmt.Printf("\r📢 [NOTIF SYSTEM]: %s\n> ", res.GetSystemMessage())
		}
	}()

	// --- 4. FITUR BARU: UPLOAD FILE DUMMY (Client Streaming) ---
	go func() {
		// Kita coba upload file kecil setiap 30 detik sebagai simulasi
		time.Sleep(5 * time.Second)
		streamUp, _ := mediaClient.UploadFile(context.Background())
		chunks := [][]byte{[]byte("isi file part 1 "), []byte("isi file part 2")}
		for _, c := range chunks {
			streamUp.Send(&pbMedia.FileChunk{ChunkData: c, Filename: "foto_profile.png"})
		}
		res, _ := streamUp.CloseAndRecv()
		fmt.Printf("\r📁 [Upload Info]: %s\n> ", res.GetFileUrl())
	}()

	// --- 5. CHAT STREAMING (Bi-directional) ---
	stream, _ := chatClient.ChatStream(context.Background())

	// HANDSHAKE: Agar siapa pun bisa mulai chat duluan
	stream.Send(&pbChat.ChatMessage{
		SenderId:   myName,
		ReceiverId: "SERVER",
		Content:    "INITIAL_CONN",
	})

	// THREAD PENDENGAR PESAN
	go func() {
		for {
			in, err := stream.Recv()
			if err != nil { return }
			if in.GetContent() != "INITIAL_CONN" {
				fmt.Printf("\r📩 [%s]: %s\n> ", in.GetSenderId(), in.GetContent())
				// Kirim ACK
				_, _ = chatClient.AckMessage(context.Background(), &pbChat.AckRequest{
					MessageId: in.GetMessageId(), ReceiverId: myName,
				})
			}
		}
	}()

	// THREAD PENGIRIM (Keyboard)
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("\n💬 Chat Ready! Ketik pesan atau 'exit' untuk keluar.\n> ")
	
	counter := 1
	for scanner.Scan() {
		teks := scanner.Text()
		if teks == "exit" { break }

		stream.Send(&pbChat.ChatMessage{
			MessageId: fmt.Sprintf("MSG-%d", counter),
			SenderId: myName, ReceiverId: friendName,
			Content: teks, Timestamp: time.Now().Unix(),
		})
		counter++
		fmt.Print("> ")
	}
	stream.CloseSend()
}