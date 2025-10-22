# Sistem Informasi Manajemen Dokumen Kepolisian (SIMDOKPOL)

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue.svg)
![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey.svg)

**SIMDOKPOL** adalah aplikasi desktop _standalone_ yang dirancang untuk membantu unit kepolisian dalam manajemen dan penerbitan surat keterangan secara efisien, cepat, dan aman. Dibangun untuk berjalan 100% offline, aplikasi ini menyediakan alur kerja modern dalam format yang mudah didistribusikan dan diinstal.

---

## ✨ Fitur Utama

-   **Aplikasi Desktop Standalone**: Berjalan sebagai aplikasi mandiri dengan ikon di _system tray_ dan notifikasi _native_, memberikan pengalaman pengguna yang terintegrasi.
-   **Virtual Host Otomatis**: Setup domain lokal (`simdokpol.local`) secara otomatis pada first-launch untuk pengalaman akses yang lebih profesional.
-   **Installer Profesional**: Didistribusikan sebagai file installer (`.exe`) untuk Windows dan AppImage untuk Linux, membuat proses instalasi menjadi mudah bagi pengguna akhir.
-   **Cross-Platform Icon Support**: Icon systray yang optimal untuk Windows (.ico), Linux & macOS (.png) dengan fallback mechanism.
-   **Alur Setup Awal Terpandu**: Konfigurasi pertama kali yang mudah untuk mengatur detail instansi (KOP surat, nama kantor) dan membuat akun Super Admin.
-   **Manajemen Dokumen Lengkap (CRUD)**: Sistem penuh untuk Membuat, Membaca, Memperbarui, dan Menghapus surat keterangan, termasuk fitur **Buat Ulang (Duplikat)** untuk efisiensi.
-   **Manajemen Pengguna Berbasis Peran**: Dua tingkat hak akses (Super Admin & Operator) dengan fitur untuk menonaktifkan dan mengaktifkan kembali akun pengguna.
-   **Dasbor Analitik Real-Time**: Tampilan ringkasan data dengan kartu statistik dan grafik interaktif untuk memonitor aktivitas operasional.
-   **Formulir Cerdas & Dinamis**: Input tanggal yang konsisten, data barang hilang yang interaktif, dan sistem rekomendasi petugas otomatis berdasarkan regu.
-   **Fitur Backup & Restore**: Super Admin dapat dengan mudah mencadangkan dan memulihkan seluruh database aplikasi.
-   **Modul Audit Log Komprehensif**: Setiap aksi penting (pembuatan/pembaruan/penghapusan data) dicatat secara otomatis untuk akuntabilitas.
-   **Auto-Generated Secure JWT Secret**: Secret key yang aman dibuat otomatis menggunakan cryptographically secure random generator.
-   **Pratinjau Cetak Presisi Tinggi**: Halaman pratinjau cetak yang dirancang agar 100% cocok dengan format fisik surat resmi.
-   **100% Offline**: Semua aset (font, CSS, JavaScript) dan fungsionalitas dirancang untuk berjalan tanpa koneksi internet.

---

## 🛠️ Tumpukan Teknologi (Technology Stack)

-   **Backend**: Go (Golang)
-   **Web Framework**: Gin
-   **ORM & Migrasi**: GORM & golang-migrate
-   **Database**: SQLite
-   **Desktop UI**: Go HTML Templates, CSS, JavaScript (jQuery, Bootstrap, Chart.js)
-   **Integrasi Desktop**:
    -   System Tray: `getlantern/systray`
    -   Notifikasi: `gen2brain/beeep`
-   **Build & Packaging**:
    -   Installer Windows: NSIS
    -   Paket Linux: AppImageKit

---

## 🚀 Memulai (Getting Started)

### Untuk Pengguna Akhir

1.  Unduh versi terbaru dari halaman **[Releases](https://github.com/USERNAME/REPO/releases)**.
2.  **Di Windows**:
    -   Jalankan file `SIMDOKPOL-Setup-vX.X.X.exe` **sebagai Administrator** (untuk first-time setup virtual host).
    -   Ikuti wizard instalasi.
    -   Jalankan aplikasi dari shortcut di Desktop atau Start Menu.
    -   Akses melalui: `http://simdokpol.local:8080` atau `http://localhost:8080`
3.  **Di Linux**:
    -   Unduh file `.AppImage`.
    -   Buat file tersebut dapat dieksekusi: `chmod +x SIMDOKPOL-x86_64.AppImage`.
    -   Jalankan pertama kali dengan sudo untuk setup virtual host: `sudo ./SIMDOKPOL-x86_64.AppImage`
    -   Setelah setup, bisa dijalankan normal: `./SIMDOKPOL-x86_64.AppImage`
    -   Akses melalui: `http://simdokpol.local:8080` atau `http://localhost:8080`
4.  **Di macOS**:
    -   Unduh dan ekstrak file `.dmg` atau `.app`.
    -   Jalankan pertama kali dengan sudo untuk setup virtual host: `sudo ./SIMDOKPOL.app`
    -   Setelah setup, bisa dijalankan normal dari Applications.
    -   Akses melalui: `http://simdokpol.local:8080` atau `http://localhost:8080`

> **Catatan Virtual Host**: Aplikasi akan otomatis mencoba mengkonfigurasi domain lokal `simdokpol.local` pada first launch. Jika gagal, instruksi manual akan ditampilkan di console log. Lihat [README_VHOST.md](README_VHOST.md) untuk detail lengkap.

### Untuk Pengembang

#### Prasyarat

-   Go (versi 1.22+)
-   Git
-   C Compiler (TDM-GCC/Mingw-w64 untuk Windows, `build-essential` untuk Debian/Ubuntu, `base-devel` untuk Manjaro/Arch)
-   Untuk _Cross-Compile_ ke Windows dari Linux: `mingw-w64-gcc`
-   Untuk membuat installer: NSIS (bisa diinstal via Wine atau native di Linux)

#### Instalasi & Menjalankan

1.  **Kloning Repositori**:
    ```bash
    git clone [URL_REPOSITORI_ANDA] simdokpol
    cd simdokpol
    ```

2.  **Konfigurasi Environment**:
    ```bash
    cp .env.example .env
    # File .env akan di-generate otomatis dengan secure JWT secret jika tidak ada
    ```

3.  **Instalasi Dependensi**:
    ```bash
    go mod tidy
    ```

4.  **Menjalankan di Mode Pengembangan (dengan Live Reload)**:
    Untuk pengembangan, disarankan menggunakan [Air](https://github.com/cosmtrek/air).
    ```bash
    air
    ```
    Aplikasi akan berjalan di `http://localhost:8080`.

5.  **Setup Virtual Host (Opsional untuk Development)**:
    ```bash
    # Windows (PowerShell as Administrator)
    .\vhost-manager.bat setup
    
    # Linux/macOS
    sudo ./vhost-manager.sh setup
    ```

#### Membangun Aplikasi (Building)

-   **Build untuk Linux**:
    ```bash
    go build -ldflags="-s -w" -o simdokpol ./cmd/main.go
    ```

-   **Build untuk macOS**:
    ```bash
    go build -ldflags="-s -w" -o simdokpol ./cmd/main.go
    ```

-   **Cross-Compile untuk Windows (dari Linux)**:
    ```bash
    CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
    go build -ldflags="-s -w -H windowsgui -extldflags \"-static\"" \
    -o simdokpol.exe ./cmd/main.go
    ```

-   **Build dengan Icon untuk Windows**:
    ```bash
    # Pastikan rsrc sudah terinstall
    go install github.com/akavel/rsrc@latest
    
    # Generate resource file
    rsrc -ico web/static/img/icon.ico -o rsrc.syso
    
    # Build dengan icon
    CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
    go build -ldflags="-s -w -H windowsgui -extldflags \"-static\"" \
    -o simdokpol.exe ./cmd/main.go
    ```

---

## 📂 Struktur Proyek

```
simdokpol/
├── cmd/
│   └── main.go              # Entrypoint aplikasi (systray, vhost, server)
├── internal/
│   ├── config/              # Konfigurasi aplikasi
│   ├── controllers/         # HTTP handlers
│   ├── middleware/          # Middleware (auth, logging, etc.)
│   ├── models/              # Model data
│   ├── repositories/        # Data access layer
│   ├── services/            # Business logic
│   └── utils/               # Utility functions
│       └── vhost_setup.go   # Virtual host setup utility
├── migrations/              # File migrasi skema database
├── web/
│   ├── static/
│   │   ├── css/            # Stylesheet
│   │   ├── js/             # JavaScript files
│   │   └── img/
│   │       ├── icon.ico    # Windows tray icon
│   │       └── icon.png    # Linux/macOS tray icon
│   └── templates/          # HTML templates
├── .air.toml               # Konfigurasi live-reload Air
├── .env.example            # Contoh file environment
├── go.mod                  # Manajemen dependensi Go
├── go.sum                  # Checksum dependensi
├── setup.nsi               # Skrip installer NSIS (Windows)
├── README.md               # Dokumentasi utama
├── README_VHOST.md         # Dokumentasi virtual host setup
└── vhost-manager.sh/bat    # Skrip helper konfigurasi vhost
```

---

## 🔧 Konfigurasi

### Environment Variables

Aplikasi menggunakan file `.env` untuk konfigurasi. File ini akan di-generate otomatis dengan nilai aman jika belum ada.

```env
# JWT Secret Key (auto-generated secara aman)
JWT_SECRET_KEY=<auto-generated-64-char-secure-string>

# Database Configuration
DB_DSN=simdokpol.db?_foreign_keys=on

# Server Port
PORT=8080
```

### Virtual Host Configuration

Domain default: `simdokpol.local`

Untuk mengubah domain, edit `internal/utils/vhost_setup.go`:

```go
const (
    LocalDomain = "nama-domain-anda.local"
    LocalIP     = "127.0.0.1"
)
```

---

## 🎯 Platform-Specific Features

### Windows
- ✅ System tray dengan .ico format
- ✅ Auto-elevate untuk modifikasi hosts file
- ✅ NSIS installer dengan uninstaller
- ✅ Start menu dan desktop shortcuts
- ✅ Windows notification support

### Linux
- ✅ System tray dengan .png format
- ✅ AppImage portable package
- ✅ systemd-resolved DNS flush support
- ✅ Desktop entry untuk application menu
- ✅ Native notification via libnotify

### macOS
- ✅ System tray dengan .png format
- ✅ .app bundle atau .dmg distribution
- ✅ macOS notification center support
- ✅ DNS cache flush dengan dscacheutil
- ✅ Gatekeeper compatible signing (optional)

---

## 📚 Dokumentasi Tambahan

- **[Virtual Host Setup Guide](README_VHOST.md)** - Panduan lengkap setup domain lokal
- **[API Documentation](http://localhost:8080/swagger/)** - Swagger API docs (saat aplikasi running)
- **[Development Guide](DEVELOPMENT.md)** - Panduan untuk developer (coming soon)
- **[Deployment Guide](DEPLOYMENT.md)** - Panduan deployment & distribution (coming soon)

---

## 🛣️ Rencana Pengembangan (Roadmap)

-   [x] **Arsitektur Backend & Frontend yang Bersih**
-   [x] **Alur Kerja Surat Keterangan Hilang (CRUD Lengkap)**
-   [x] **Otentikasi & Otorisasi Berbasis Peran**
-   [x] **Manajemen Pengguna (CRUD, Aktivasi/Deaktivasi)**
-   [x] **Dasbor Analitik dengan Grafik Dinamis**
-   [x] **Alur Setup Awal Terpandu**
-   [x] **Modul Audit Log & Fitur Backup/Restore**
-   [x] **Aplikasi Desktop Standalone (via Systray & Beeep)**
-   [x] **Virtual Host Setup Otomatis (Multiplatform)**
-   [x] **Cross-Platform Icon Support dengan Fallback**
-   [x] **Auto-Generated Secure JWT Secret**
-   [x] **Installer Profesional (NSIS untuk Windows, AppImage untuk Linux)**
-   [x] **Penambahan Unit Test & Integration Test (Blueprint sudah ada)**
-   [x] **Dokumentasi API yang lebih mendalam**
-   [ ] **Peningkatan Fitur Pencetakan (kustomisasi template, dll)**
-   [ ] **Support untuk Multiple Instances di Network**
-   [ ] **Export Data ke Format PDF/Excel**

---

## 🐛 Troubleshooting

### Virtual Host Issues

**Problem**: Domain `simdokpol.local` tidak bisa diakses

**Solution**:
```bash
# Windows (Command Prompt as Admin)
ipconfig /flushdns

# Linux
sudo systemctl restart systemd-resolved

# macOS
sudo dscacheutil -flushcache
sudo killall -HUP mDNSResponder
```

**Problem**: "Permission denied" saat setup virtual host

**Solution**:
- Windows: Jalankan aplikasi sebagai Administrator
- Linux/macOS: Jalankan dengan `sudo` untuk first-time setup

### Icon Not Showing in System Tray

**Solution**:
- Pastikan file `web/static/img/icon.ico` (Windows) atau `icon.png` (Linux/macOS) ada
- Check file permissions: `chmod 644 web/static/img/icon.*`
- Restart aplikasi

### Database Migration Errors

**Solution**:
```bash
# Backup database terlebih dahulu
cp simdokpol.db simdokpol.db.backup

# Reset migrations
rm -rf migrations/*.up.sql
go run cmd/main.go
```

Untuk troubleshooting lebih lanjut, lihat [Issues](https://github.com/USERNAME/REPO/issues) atau buat issue baru.

---

## 🤝 Kontribusi

Kontribusi sangat diterima! Silakan baca [CONTRIBUTING.md](CONTRIBUTING.md) untuk detail tentang code of conduct dan proses submit pull requests.

### Development Workflow

1. Fork repositori
2. Buat feature branch: `git checkout -b feature/AmazingFeature`
3. Commit changes: `git commit -m 'Add some AmazingFeature'`
4. Push ke branch: `git push origin feature/AmazingFeature`
5. Buka Pull Request

---

## 📄 Lisensi

Proyek ini dilisensikan di bawah **MIT License** - lihat file [LICENSE](LICENSE) untuk detail.

---

## 👥 Tim Pengembang

- **Lead Developer** - [MUHAMMAD YUSUF ABDURROHMAN](https://github.com/muhammad1505)
- **Contributors** - [List Contributors](https://github.com/USERNAME/REPO/contributors)

---

## 📞 Kontak & Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/USERNAME/REPO/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/USERNAME/REPO/discussions)
- 📧 **Email**: emailbaruku50@gmail.com

---

## 🙏 Acknowledgments

- Thanks to all contributors who have helped shape SIMDOKPOL
- Built with ❤️ using Go and modern web technologies
- Special thanks to the open-source community

---

**Made with ❤️ for Indonesian Police Units**