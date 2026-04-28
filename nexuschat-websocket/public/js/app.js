const socket = io();

const chatBox = document.getElementById('chat-box');
const alertBox = document.getElementById('alert-box');
const recentUsersDiv = document.getElementById('recent-users'); // Elemen UI baru

// --- FITUR DAFTAR USER AKTIF ---
const activeUsers = new Set(); // Menyimpan nama user tanpa duplikat

function renderActiveUsers() {
    recentUsersDiv.innerHTML = '<span class="help-text">Balas Cepat: </span>';
    activeUsers.forEach(user => {
        recentUsersDiv.innerHTML += `<div class="user-pill" onclick="setTarget('${user}')">${user}</div>`;
    });
}

function setTarget(user) {
    document.getElementById('targetId').value = user;
    document.getElementById('msgContent').focus(); // Auto-focus ke kotak ngetik
}
// -------------------------------

const ctx = document.getElementById('trafficChart').getContext('2d');
let messageCount = 0;
const trafficChart = new Chart(ctx, {
    type: 'line',
    data: { labels: ['0s', '5s', '10s', '15s', '20s', '25s'], datasets: [{ label: 'Pesan gRPC/sec', data: [0, 0, 0, 0, 0, 0], borderColor: '#f9e2af', tension: 0.4 }] },
    options: { animation: false, scales: { y: { beginAtZero: true, max: 10 } } }
});

setInterval(() => {
    trafficChart.data.datasets[0].data.shift();
    trafficChart.data.datasets[0].data.push(messageCount);
    trafficChart.update();
    messageCount = 0;
}, 5000);

function getTimeStr() {
    return new Date().toLocaleTimeString('id-ID', { hour: '2-digit', minute: '2-digit' });
}

// 2. SOCKET EVENTS
socket.on('chat_message', (msg) => {
    messageCount++;
    const time = getTimeStr();
    chatBox.innerHTML += `<p><span style="color:#6c7086; font-size:12px;">[${time}]</span> <b>[${msg.sender_id}]</b> ke [${msg.receiver_id}]: ${msg.content}</p>`;
    chatBox.scrollTop = chatBox.scrollHeight;

    // OTOMATIS TAMBAHKAN KE DAFTAR USER AKTIF
    if (msg.sender_id !== 'AdminWeb' && msg.sender_id !== 'SERVER') {
        activeUsers.add(msg.sender_id);
        renderActiveUsers();
    }
});

socket.on('system_alert', (msg) => {
    alertBox.innerText = msg;
});

socket.on('disconnect', () => {
    document.getElementById('online-status').style.backgroundColor = '#f38ba8';
    document.getElementById('status-text').innerText = 'Disconnected';
});

// 3. UI ACTIONS
function sendMessage() {
    const sender = document.getElementById('senderId').value;
    const target = document.getElementById('targetId').value;
    const content = document.getElementById('msgContent').value;
    
    if(!content || !target) return; // Jangan kirim kalau target kosong
    
    socket.emit('send_message', { senderId: sender, receiverId: target, content: content });
    
    const time = getTimeStr();
    chatBox.innerHTML += `<p style="color:#a6e3a1;"><span style="color:#6c7086; font-size:12px;">[${time}]</span> <b>[Me]</b> ke [${target}]: ${content}</p>`;
    document.getElementById('msgContent').value = '';
    chatBox.scrollTop = chatBox.scrollHeight;
}

function kickUser() {
    const target = document.getElementById('kickTarget').value;
    if(target) socket.emit('kick_user', target);
}