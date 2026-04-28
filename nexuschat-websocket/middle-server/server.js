const express = require('express');
const http = require('http');
const { Server } = require('socket.io');
const grpc = require('@grpc/grpc-js');
const protoLoader = require('@grpc/proto-loader');
const path = require('path');

const app = express();
const server = http.createServer(app);
const io = new Server(server);

app.use(express.static(path.join(__dirname, '../public')));

// --- 1. LOAD PROTO FILES ---
const loaderOptions = { keepCase: true, longs: String, enums: String, defaults: true, oneofs: true };

// Load Core Messaging Proto
const CORE_PROTO = path.join(__dirname, 'proto/core_message.proto');
const corePkg = protoLoader.loadSync(CORE_PROTO, loaderOptions);
const coreMessaging = grpc.loadPackageDefinition(corePkg).coremessaging;

// Load Media Notif Proto (UNTUK UPGRADE 1.0)
const MEDIA_PROTO = path.join(__dirname, 'proto/media_notif.proto');
const mediaPkg = protoLoader.loadSync(MEDIA_PROTO, loaderOptions);
const mediaNotif = grpc.loadPackageDefinition(mediaPkg).medianotif;

// --- 2. KONEKSI KE SERVER GO ---
// Konek ke Port 50012 (Chat)
const chatClient = new coreMessaging.MessagingService('localhost:50012', grpc.credentials.createInsecure());
// Konek ke Port 50013 (Notifikasi)
const notifClient = new mediaNotif.MediaNotifService('localhost:50013', grpc.credentials.createInsecure());

// --- 3. BUKA JALUR NOTIFIKASI DARI GO KE WEB UI ---
const notifStream = notifClient.SubscribeNotif({ client_id: "AdminWebDashboard" });
notifStream.on('data', (notif) => {
    // Jika Go ngirim notif, langsung tembak ke semua browser yang buka Web UI
    console.log(`[Notif System] ${notif.system_message}`);
    io.emit('system_alert', `📢 [${notif.type}]: ${notif.system_message}`);
});
notifStream.on('error', (err) => console.error("⚠️ Stream Notif Terputus."));

// --- 4. WEBSOCKET HANDLER ---
io.on('connection', (socket) => {
    console.log('🔌 Klien Web terhubung:', socket.id);

    const call = chatClient.ChatStream();

    call.write({
        message_id: 'INIT',
        sender_id: 'AdminWeb',
        receiver_id: 'SERVER',
        content: 'INITIAL_CONN',
        timestamp: 0
    });

    socket.on('send_message', (data) => {
        call.write({
            message_id: 'WEB-' + Date.now(),
            sender_id: data.senderId,
            receiver_id: data.receiverId,
            content: data.content,
            timestamp: Math.floor(Date.now() / 1000)
        });
    });

    call.on('data', (msg) => {
        socket.emit('chat_message', msg);
        chatClient.AckMessage({ message_id: msg.message_id, receiver_id: msg.receiver_id }, (err) => {});
    });

    socket.on('kick_user', (userId) => {
        // Broadcast ke web UI lain (simulasi)
        io.emit('system_alert', `🚨 COMMAND: Kick requested for ${userId}`);
    });

    socket.on('disconnect', () => {
        console.log('👋 Klien Web terputus');
        call.end();
    });
});

const PORT = 3000;
server.listen(PORT, () => {
    console.log(`🚀 Middle Server (V1.0) menyala di http://localhost:${PORT}`);
});