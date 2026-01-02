GUARDPOT
# SSH Terminal - Web-Based SSH Client 

Tarayıcı üzerinden çalışan, gerçek SSH bağlantıları kurabilen full-stack web uygulaması.

## 🏗️ Mimari Özet

```
┌─────────────────────────────────────────────────────────────┐
│                      Frontend (React)                       │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │   Login/    │  │  Dashboard  │  │    Terminal View    │  │
│  │  Register   │  │ (SSH List)  │  │    (xterm.js)       │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
│         │               │                    │              │
│         └───────────────┴────────────────────┘              │
│                         │ HTTP/WebSocket                    │
└─────────────────────────┼───────────────────────────────────┘
                          │
┌─────────────────────────┼───────────────────────────────────┐
│                   Backend (Go)                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Auth API   │  │  SSH CRUD   │  │  WebSocket Handler  │  │
│  │ (JWT/OAuth) │  │     API     │  │   (SSH Bridge)      │  │
│  └─────────────┘  └─────────────┘  └─────────────────────┘  │
│         │               │                    │              │
│         └───────────────┴────────────────────┘              │
│                         │                                   │
│              ┌──────────┴──────────┐                        │
│              │  SQLite Database    │                        │
│              │  (Users, SSH Conn)  │                        │
│              └─────────────────────┘                        │
└─────────────────────────────────────────────────────────────┘
                          │
                          │ TCP/SSH
                          ▼
              ┌─────────────────────┐
              │   Remote Servers    │
              │   (Real SSH)        │
              └─────────────────────┘
```

### Teknoloji Stack'i

| Katman | Teknoloji |
|--------|-----------|
| Frontend | React 18, Vite, xterm.js |
| Backend | Go 1.21+, Gorilla Mux |
| Database | SQLite (GORM ORM) |
| Auth | JWT, Google OAuth 2.0 |
| SSH | golang.org/x/crypto/ssh |
| Real-time | WebSocket (gorilla/websocket) |

---

## 🔐 SSH Bağlantısı Nasıl Çalışır?

1. **Bağlantı Ekleme**: Kullanıcı dashboard'dan SSH bilgilerini girer (Host, Port, Username, Password/Key)
2. **Şifreleme**: Kimlik bilgileri AES-256-GCM ile şifrelenerek SQLite'a kaydedilir
3. **Terminal Açma**: Kullanıcı "Connect" butonuna tıkladığında:
   - Frontend xterm.js terminal oluşturur
   - Backend'e WebSocket bağlantısı açılır
   - Backend `golang.org/x/crypto/ssh` ile gerçek SSH bağlantısı kurar
   - PTY (pseudo-terminal) oturumu başlatılır
4. **Gerçek Zamanlı I/O**: 
   - Kullanıcı girişi → WebSocket → SSH stdin → Uzak sunucu
   - Uzak sunucu çıktısı → SSH stdout → WebSocket → xterm.js ekran

```
Browser ←──WebSocket──→ Go Backend ←──TCP/SSH──→ Remote Server
(xterm.js)              (Bridge)                (Real Shell)
```

---

## 🤖 Geliştirme Süreci ve İşbirliği Metodolojisi

Bu projenin teknik vizyonu, sistem mimarisi ve çekirdek yapı taşları tamamen geliştirici tarafından kurgulanmıştır. Geliştirme sürecinde **Google DeepMind Antigravity (Gemini)** ile "Pair Programming" (Eşli Programlama) metodolojisi izlenmiştir.

### 👨‍💻 Geliştiricinin Rolü (Mimar ve Lider)
- **Mimari Kararlar**: Frontend-Backend iletişim protokolünün (HTTP + WebSocket) tasarlanması.
- **Teknoloji Seçimi**: Go, React, SQLite ve xterm.js gibi kilit teknolojilerin belirlenmesi.
- **Güvenlik Stratejisi**: SSH anahtarlarının şifrelenmesi ve JWT tabanlı kimlik doğrulama yapısının kurgulanması.
- **Kod İnceleme (Code Review)**: Üretilen her kod bloğunun performans ve güvenlik standartlarına göre denetlenmesi.

### 🤖 AI Asistanın Rolü (DeepMind Antigravity)
- **Boilerplate Kod Üretimi**: Tekrar eden veri tabanı modelleri ve API handler taslaklarının hızlıca oluşturulması.
- **Frontend Bileşenleri**: Modern UI/UX pratiklerine uygun React bileşenlerinin (SSH Terminal kartları, Form yapıları) iskeletinin hazırlanması.
- **Hata Analizi**: Go rutinlerindeki olası "race condition" durumlarının tespiti ve çözüm önerileri.
- **Dokümantasyon**: Teknik dokümanların ve kullanıcı kılavuzlarının taslaklarının hazırlanması.



## 🚀 Kurulum ve Çalıştırma

### Gereksinimler
- Go 1.21+
- Node.js 18+
- npm veya yarn

### Backend
```bash
# Proje dizininde
go mod tidy
go run main.go
# Server: http://localhost:8080
```

### Frontend
```bash
cd frontend
npm install
npm run dev
# Dev server: http://localhost:5173
```

### Environment Variables (Opsiyonel)
```bash
# .env dosyası veya sistem ortam değişkenleri
JWT_SECRET=your-secret-key
ENCRYPTION_KEY=32-byte-encryption-key-here!!
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
```

### Production Build
```bash
cd frontend
npm run build
# Build output: frontend/dist/

# Backend production
go build -o server main.go
./server
```

---

## 📁 Proje Yapısı

```
Go Backend/
├── main.go                 # Uygulama giriş noktası
├── go.mod                  # Go modül tanımları
├── internal/
│   ├── config/             # Yapılandırma
│   ├── database/           # SQLite bağlantısı
│   ├── handlers/           # HTTP/WebSocket handlers
│   │   ├── auth.go         # Kimlik doğrulama
│   │   ├── ssh.go          # SSH CRUD
│   │   └── terminal.go     # WebSocket SSH bridge
│   ├── middleware/         # Auth, CORS middleware
│   ├── models/             # User, SSHConnection modelleri
│   └── utils/              # Şifreleme, hash yardımcıları
└── frontend/
    ├── src/
    │   ├── components/     # React bileşenleri
    │   ├── context/        # Auth context
    │   ├── pages/          # Sayfa bileşenleri
    │   ├── services/       # API servisleri
    │   └── index.css       # Global stiller
    └── package.json
```

---

## 🔒 Güvenlik Özellikleri

- **JWT Token**: 7 gün geçerli, HttpOnly opsiyonu
- **Bcrypt**: Şifre hash'leme (cost: 10)
- **AES-256-GCM**: SSH kimlik bilgisi şifreleme
- **CORS**: Cross-origin istek kontrolü
- **Protected Routes**: Auth middleware ile API koruması

---

## ✨ Özellikler

- ✅ Kullanıcı kayıt ve giriş
- ✅ Google OAuth entegrasyonu
- ✅ SSH bağlantı yönetimi (CRUD)
- ✅ Gerçek SSH terminal (simülasyon değil!)
- ✅ Password ve Private Key authentication
- ✅ Terminal resize desteği
- ✅ Modern dark theme UI
- ✅ Responsive tasarım

---

## 📝 Lisans

MIT License
